package plugin

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
	plugin_exec "github.com/thegeeklab/wp-plugin-go/v5/exec"
	plugin_file "github.com/thegeeklab/wp-plugin-go/v5/file"
	plugin_types "github.com/thegeeklab/wp-plugin-go/v5/types"
	plugin_util "github.com/thegeeklab/wp-plugin-go/v5/util"
)

const (
	daemonBackoffMaxRetries      = 3
	daemonBackoffInitialInterval = 2 * time.Second
	daemonBackoffMultiplier      = 3.5
)

var (
	ErrGitCloneDestintionNotValid = errors.New("destination not valid")
	ErrTypeAssertionFailed        = errors.New("type assertion failed")
)

func (p *Plugin) run(ctx context.Context) error {
	if err := p.FlagsFromContext(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := p.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := p.Execute(ctx); err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	return nil
}

// Validate handles the settings validation of the plugin.
func (p *Plugin) Validate() error {
	var err error

	// This default cannot be set in the cli flag, as the CI_* environment variables
	// can be set empty, resulting in an empty default value.
	if p.Settings.Repo.Branch == "" {
		p.Settings.Repo.Branch = "main"
	}

	if p.Settings.Repo.WorkDir == "" {
		p.Settings.Repo.WorkDir, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}
	}

	if p.Settings.Partial {
		p.Settings.Repo.Depth = 1
		p.Settings.Repo.Filter = "tree:0"
	}

	return nil
}

// Execute provides the implementation of the plugin.
func (p *Plugin) Execute(ctx context.Context) error {
	var err error

	homeDir := plugin_util.GetUserHomeDir()
	batchCmd := make([]*plugin_exec.Cmd, 0)

	fmt.Println(p.Settings.Repo.WorkDir)

	// Handle repo initialization.
	if err := os.MkdirAll(p.Settings.Repo.WorkDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create working directory: %w", err)
	}

	p.Settings.Repo.IsEmpty, err = plugin_file.IsDirEmpty(p.Settings.Repo.WorkDir)
	if err != nil {
		return fmt.Errorf("failed to check working directory: %w", err)
	}

	isDir, err := plugin_file.IsDir(filepath.Join(p.Settings.Repo.WorkDir, ".git"))
	if err != nil {
		return fmt.Errorf("failed to check working directory: %w", err)
	}

	if !isDir {
		batchCmd = append(batchCmd, p.Settings.Repo.Init())
		batchCmd = append(batchCmd, p.Settings.Repo.RemoteAdd())

		if p.Settings.SSHKey != "" {
			batchCmd = append(batchCmd, p.Settings.Repo.ConfigSSHCommand(p.Settings.SSHKey))
		}
	}

	batchCmd = append(batchCmd, p.Settings.Repo.ConfigSSLVerify(p.Network.InsecureSkipVerify))

	netrc := p.Settings.Netrc
	if err := WriteNetrc(homeDir, netrc.Machine, netrc.Login, netrc.Password); err != nil {
		return err
	}

	// Handle clone
	if p.Settings.Repo.CommitSha == "" {
		// fetch and checkout by ref
		log.Info().Msg("no commit information: using head checkout")

		batchCmd = append(batchCmd, p.Settings.Repo.FetchSource(p.Settings.Repo.CommitRef))
		batchCmd = append(batchCmd, p.Settings.Repo.CheckoutHead())
	} else {
		batchCmd = append(batchCmd, p.Settings.Repo.FetchSource(p.Settings.Repo.CommitSha))
		batchCmd = append(batchCmd, p.Settings.Repo.CheckoutSha())
	}

	if p.Settings.Tags {
		batchCmd = append(batchCmd, p.Settings.Repo.FetchTags())
	}

	for name, submoduleURL := range p.Settings.Repo.Submodules {
		batchCmd = append(batchCmd, p.Settings.Repo.ConfigRemapSubmodule(name, submoduleURL))
	}

	if p.Settings.Recursive {
		batchCmd = append(batchCmd, p.Settings.Repo.SubmoduleUpdate())
	}

	if p.Settings.Lfs {
		batchCmd = append(batchCmd, p.Settings.Repo.FetchLFS())
		batchCmd = append(batchCmd, p.Settings.Repo.CheckoutLFS())
	}

	for _, cmd := range batchCmd {
		if cmd == nil {
			continue
		}

		buf := new(bytes.Buffer)

		// Don' set GIT_TERMINAL_PROMPT=0 as it prevents git from loading .netrc
		defaultEnvVars := []string{
			"GIT_LFS_SKIP_SMUDGE=1", // prevents git-lfs from retrieving any LFS files
		}

		if p.Settings.Home != "" {
			if _, err := os.Stat(p.Settings.Home); !os.IsNotExist(err) {
				defaultEnvVars = append(defaultEnvVars, fmt.Sprintf("HOME=%s", p.Settings.Home))
			}
		}

		cmd.Env = append(os.Environ(), defaultEnvVars...)
		cmd.Stdout = io.MultiWriter(os.Stdout, buf)
		cmd.Stderr = io.MultiWriter(os.Stderr, buf)
		cmd.Dir = p.Settings.Repo.WorkDir

		err := cmd.Run()

		switch {
		case err != nil && shouldRetry(buf.String()):
			return retryCmd(ctx, cmd)
		case err != nil:
			return err
		}
	}

	return nil
}

func (p *Plugin) FlagsFromContext() error {
	submodules, ok := p.Context.Generic("submodule-override").(*plugin_types.MapFlag)
	if !ok {
		return fmt.Errorf("%w: failed to read submodule-override input", ErrTypeAssertionFailed)
	}

	p.Settings.Repo.Submodules = submodules.Get()

	return nil
}

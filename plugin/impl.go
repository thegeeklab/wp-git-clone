// Copyright (c) 2023, Robert Kaussow <mail@thegeeklab.de>

// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/thegeeklab/wp-git-clone/git"
	"github.com/thegeeklab/wp-plugin-go/types"
	"golang.org/x/sys/execabs"
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

//nolint:revive
func (p *Plugin) run(ctx context.Context) error {
	if err := p.FlagsFromContext(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := p.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := p.Execute(); err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	return nil
}

// Validate handles the settings validation of the plugin.
func (p *Plugin) Validate() error {
	if p.Settings.WorkDir == "" {
		var err error
		if p.Settings.WorkDir, err = os.Getwd(); err != nil {
			return err
		}
	}

	if p.Settings.Partial {
		p.Settings.Depth = 1
		p.Settings.Filter = "tree:0"
	}

	return nil
}

// Execute provides the implementation of the plugin.
func (p *Plugin) Execute() error {
	cmds := make([]*execabs.Cmd, 0)

	// Handle init
	initPath := filepath.Join(p.Settings.WorkDir, ".git")

	if err := os.MkdirAll(p.Settings.WorkDir, os.ModePerm); err != nil {
		return err
	}

	//nolint:nestif
	if _, err := os.Stat(initPath); os.IsNotExist(err) {
		cmds = append(cmds, git.ConfigSafeDirectory(p.Settings.Repo))

		if err := p.execCmd(git.Init(p.Settings.Repo), new(bytes.Buffer)); err != nil {
			return err
		}

		if p.Settings.UseSSH {
			cmds = append(cmds, git.RemoteAdd(p.Settings.Repo.RemoteSSH))
			if p.Settings.SSHKey != "" {
				cmds = append(cmds, git.ConfigSSHCommand(p.Settings.SSHKey))
			}
		} else {
			cmds = append(cmds, git.RemoteAdd(p.Settings.Repo.RemoteURL))
		}
	}

	if p.Settings.Repo.InsecureSkipSSLVerify {
		cmds = append(cmds, git.ConfigSSLVerify(p.Settings.Repo))
	}

	if err := git.WriteNetrc(p.Settings.Netrc.Machine, p.Settings.Netrc.Login, p.Settings.Netrc.Password); err != nil {
		return err
	}

	// Handle clone

	if p.Settings.Repo.CommitSha == "" {
		// fetch and checkout by ref
		log.Info().Msg("no commit information: using head checkout")

		cmds = append(cmds, git.FetchSource(p.Settings.Repo.CommitRef, p.Settings.Depth, p.Settings.Filter))
		cmds = append(cmds, git.CheckoutHead())
	} else {
		cmds = append(cmds, git.FetchSource(p.Settings.Repo.CommitSha, p.Settings.Depth, p.Settings.Filter))
		cmds = append(cmds, git.CheckoutSha(p.Settings.Repo))
	}

	if p.Settings.Tags {
		cmds = append(cmds, git.FetchTags())
	}

	for name, submoduleURL := range p.Settings.Repo.Submodules {
		cmds = append(cmds, git.ConfigRemapSubmodule(name, submoduleURL))
	}

	if p.Settings.Recursive {
		cmds = append(cmds, git.SubmoduleUpdate(p.Settings.Repo))
	}

	if p.Settings.Lfs {
		cmds = append(cmds, git.FetchLFS())
		cmds = append(cmds, git.CheckoutLFS())
	}

	for _, cmd := range cmds {
		log.Debug().Msgf("+ %s", strings.Join(cmd.Args, " "))

		buf := new(bytes.Buffer)
		err := p.execCmd(cmd, buf)

		switch {
		case err != nil && shouldRetry(buf.String()):
			return retryCmd(cmd)
		case err != nil:
			return err
		}
	}

	return nil
}

func (p *Plugin) FlagsFromContext() error {
	submodules, ok := p.Context.Generic("submodule-override").(*types.MapFlag)
	if !ok {
		return fmt.Errorf("%w: failed to read submodule-override input", ErrTypeAssertionFailed)
	}

	p.Settings.Repo.Submodules = submodules.Get()

	return nil
}

func (p *Plugin) execCmd(cmd *execabs.Cmd, buf *bytes.Buffer) error {
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
	cmd.Dir = p.Settings.WorkDir

	trace(cmd)

	return cmd.Run()
}

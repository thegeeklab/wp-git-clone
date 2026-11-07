package plugin

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v7"
	"github.com/rs/zerolog/log"
	plugin_exec "github.com/thegeeklab/wp-plugin-go/v8/exec"
)

const (
	netrcFile = `machine %s
login %s
password %s
`
	sshConfigFile = `Host *
    IdentityFile %s
    IdentitiesOnly yes
`
	sshConfigFileName = "wp-git-clone.conf"
	strictFilePerm    = 0o600
	sshDirPerm        = 0o700
)

// shouldRetry returns true if the command should be re-executed. Currently
// this only returns true if the remote ref does not exist.
func shouldRetry(s string) bool {
	return strings.Contains(s, "find remote ref")
}

func retryCmd(ctx context.Context, cmd *plugin_exec.Cmd) error {
	bf := backoff.NewExponentialBackOff()
	bf.InitialInterval = daemonBackoffInitialInterval
	bf.Multiplier = daemonBackoffMultiplier

	bfo := func() (any, error) {
		// copy the original command
		retry := plugin_exec.Command(cmd.Path, cmd.Args...)
		retry.Env = cmd.Env
		retry.Stdout = cmd.Stdout
		retry.Stderr = cmd.Stderr
		retry.Dir = cmd.Dir

		return nil, retry.Run()
	}

	bfn := func(err error, delay time.Duration) {
		log.Error().Msgf("failed to find remote ref: %v: retry in %s", err, delay.Truncate(time.Second))
	}

	_, err := backoff.Retry(ctx, bfo,
		backoff.WithBackOff(bf),
		backoff.WithMaxTries(daemonBackoffMaxRetries),
		backoff.WithNotify(bfn))

	return err
}

// WriteNetrc writes the netrc file.
func WriteNetrc(path, machine, login, password string) error {
	netrcPath := filepath.Join(path, ".netrc")
	netrcContent := fmt.Sprintf(netrcFile, machine, login, password)

	if err := os.WriteFile(netrcPath, []byte(netrcContent), strictFilePerm); err != nil {
		return fmt.Errorf("failed to create .netrc file: %w", err)
	}

	return nil
}

// sshConfigQuote wraps s in double quotes and escapes embedded double quotes
// and backslashes so ssh parses it as a single argument even when it contains
// whitespace.
func sshConfigQuote(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)

	return `"` + s + `"`
}

// WriteSSHConfig writes an OpenSSH client configuration that points the
// default ssh invocation at the given identity file. It avoids assembling a
// shell command (core.sshCommand) entirely: git execs ssh directly, and ssh
// parses this file without a shell.
//
// The identity is written to a dedicated config.d file and wired in through an
// Include directive, so any pre-existing ~/.ssh/config is left untouched.
func WriteSSHConfig(path, keyPath string) error {
	sshDir := filepath.Join(path, ".ssh")

	if err := os.MkdirAll(sshDir, sshDirPerm); err != nil {
		return fmt.Errorf("failed to create .ssh directory: %w", err)
	}

	configDir := filepath.Join(sshDir, "config.d")

	if err := os.MkdirAll(configDir, sshDirPerm); err != nil {
		return fmt.Errorf("failed to create ssh config.d directory: %w", err)
	}

	configPath := filepath.Join(configDir, sshConfigFileName)
	configContent := fmt.Sprintf(sshConfigFile, sshConfigQuote(keyPath))

	if err := os.WriteFile(configPath, []byte(configContent), strictFilePerm); err != nil {
		return fmt.Errorf("failed to create ssh config file: %w", err)
	}

	if err := ensureSSHInclude(sshDir, configPath); err != nil {
		return err
	}

	return nil
}

// ensureSSHInclude adds an Include directive for the plugin's config.d file to
// the main ssh config, preserving any existing content and staying idempotent.
// The directive uses an absolute path so its resolution never depends on $HOME.
func ensureSSHInclude(sshDir, configPath string) error {
	mainConfig := filepath.Join(sshDir, "config")
	includeLine := "Include " + sshConfigQuote(configPath)

	existing, err := os.ReadFile(mainConfig)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read ssh config: %w", err)
	}

	for _, line := range strings.Split(string(existing), "\n") {
		if strings.TrimSpace(line) == includeLine {
			return nil
		}
	}

	content := string(existing)

	if content != "" && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	content += includeLine + "\n"

	// #nosec G703 -- mainConfig is the fixed "config" file inside the user's own ~/.ssh directory.
	if err := os.WriteFile(mainConfig, []byte(content), strictFilePerm); err != nil {
		return fmt.Errorf("failed to update ssh config: %w", err)
	}

	return nil
}

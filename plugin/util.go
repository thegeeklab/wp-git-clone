package plugin

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/rs/zerolog/log"
	plugin_exec "github.com/thegeeklab/wp-plugin-go/v4/exec"
)

const (
	netrcFile = `machine %s
login %s
password %s
`
	strictFilePerm = 0o600
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
		retry := plugin_exec.Command(cmd.Cmd.Path, cmd.Cmd.Args...)
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

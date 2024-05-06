package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/rs/zerolog/log"
	"github.com/thegeeklab/wp-plugin-go/v2/types"
	"golang.org/x/sys/execabs"
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

func newBackoff(maxRetries uint64) backoff.BackOff {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = daemonBackoffInitialInterval
	b.Multiplier = daemonBackoffMultiplier

	return backoff.WithMaxRetries(b, maxRetries)
}

func retryCmd(cmd *types.Cmd) error {
	backoffOps := func() error {
		// copy the original command
		//nolint:gosec
		retry := &types.Cmd{
			Cmd: execabs.Command(cmd.Cmd.Path, cmd.Cmd.Args...),
		}
		retry.Env = cmd.Env
		retry.Stdout = cmd.Stdout
		retry.Stderr = cmd.Stderr
		retry.Dir = cmd.Dir

		return retry.Run()
	}
	backoffLog := func(err error, delay time.Duration) {
		log.Error().Msgf("failed to find remote ref: %v: retry in %s", err, delay.Truncate(time.Second))
	}

	return backoff.RetryNotify(backoffOps, newBackoff(daemonBackoffMaxRetries), backoffLog)
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

package plugin

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/rs/zerolog/log"
	"golang.org/x/sys/execabs"
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

func trace(cmd *execabs.Cmd) {
	fmt.Fprintf(os.Stdout, "+ %s\n", strings.Join(cmd.Args, " "))
}

func retryCmd(cmd *execabs.Cmd) error {
	backoffOps := func() error {
		// copy the original command
		//nolint:gosec
		retry := execabs.Command(cmd.Args[0], cmd.Args[1:]...)
		retry.Dir = cmd.Dir
		retry.Env = cmd.Env
		retry.Stdout = os.Stdout
		retry.Stderr = os.Stderr

		trace(cmd)

		return cmd.Run()
	}
	backoffLog := func(err error, delay time.Duration) {
		log.Error().Msgf("failed to find remote ref: %v: retry in %s", err, delay.Truncate(time.Second))
	}

	return backoff.RetryNotify(backoffOps, newBackoff(daemonBackoffMaxRetries), backoffLog)
}

package plugin

import (
	"fmt"
	"os"
	"strings"

	"github.com/cenkalti/backoff/v4"
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

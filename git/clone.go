package git

import (
	"fmt"

	"golang.org/x/sys/execabs"
)

// FetchSource fetches the source from remote.
func FetchSource(ref string, depth int, filter string) *execabs.Cmd {
	args := []string{
		"fetch",
	}

	if depth != 0 {
		args = append(args, fmt.Sprintf("--depth=%d", depth))
	}

	if filter != "" {
		args = append(args, "--filter="+filter)
	}

	args = append(args, "origin")
	args = append(args, fmt.Sprintf("+%s:", ref))

	return execabs.Command(
		gitBin,
		args...,
	)
}

// FetchTags fetches the source from remote.
func FetchTags() *execabs.Cmd {
	args := []string{
		"fetch",
		"--tags",
		"--quiet",
		"origin",
	}

	return execabs.Command(
		gitBin,
		args...,
	)
}

// FetchLFS fetches lfs.
func FetchLFS() *execabs.Cmd {
	args := []string{
		"lfs",
		"fetch",
	}

	return execabs.Command(
		gitBin,
		args...,
	)
}

// CheckoutHead handles head checkout.
func CheckoutHead() *execabs.Cmd {
	args := []string{
		"checkout",
		"--force",
		"--quiet",
		"FETCH_HEAD",
	}

	return execabs.Command(
		gitBin,
		args...,
	)
}

// CheckoutSha handles commit checkout.
func CheckoutSha(repo Repository) *execabs.Cmd {
	args := []string{
		"reset",
		"--hard",
		"--quiet",
		repo.CommitSha,
	}

	return execabs.Command(
		gitBin,
		args...,
	)
}

// CheckoutLFS handles commit checkout.
func CheckoutLFS() *execabs.Cmd {
	args := []string{
		"lfs",
		"checkout",
	}

	return execabs.Command(
		gitBin,
		args...,
	)
}

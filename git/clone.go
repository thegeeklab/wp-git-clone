package git

import (
	"fmt"

	"github.com/thegeeklab/wp-plugin-go/v3/types"
	"golang.org/x/sys/execabs"
)

// FetchSource fetches the source from remote.
func (r *Repository) FetchSource(ref string) *types.Cmd {
	args := []string{
		"fetch",
	}

	if r.Depth != 0 {
		args = append(args, fmt.Sprintf("--depth=%d", r.Depth))
	}

	if r.Filter != "" {
		args = append(args, "--filter", r.Filter)
	}

	args = append(args, "origin")
	args = append(args, fmt.Sprintf("+%s:", ref))

	return &types.Cmd{
		Cmd: execabs.Command(gitBin, args...),
	}
}

// FetchTags fetches the source from remote.
func (r *Repository) FetchTags() *types.Cmd {
	args := []string{
		"fetch",
		"--tags",
		"--quiet",
		"origin",
	}

	return &types.Cmd{
		Cmd: execabs.Command(gitBin, args...),
	}
}

// FetchLFS fetches lfs.
func (r *Repository) FetchLFS() *types.Cmd {
	args := []string{
		"lfs",
		"fetch",
	}

	return &types.Cmd{
		Cmd: execabs.Command(gitBin, args...),
	}
}

// CheckoutHead handles head checkout.
func (r *Repository) CheckoutHead() *types.Cmd {
	args := []string{
		"checkout",
		"--force",
		"--quiet",
		"FETCH_HEAD",
	}

	return &types.Cmd{
		Cmd: execabs.Command(gitBin, args...),
	}
}

// CheckoutSha handles commit checkout.
func (r *Repository) CheckoutSha() *types.Cmd {
	args := []string{
		"reset",
		"--hard",
		"--quiet",
		r.CommitSha,
	}

	return &types.Cmd{
		Cmd: execabs.Command(gitBin, args...),
	}
}

// CheckoutLFS handles commit checkout.
func (r *Repository) CheckoutLFS() *types.Cmd {
	args := []string{
		"lfs",
		"checkout",
	}

	return &types.Cmd{
		Cmd: execabs.Command(gitBin, args...),
	}
}

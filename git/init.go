package git

import (
	"github.com/thegeeklab/wp-plugin-go/v2/types"
	"golang.org/x/sys/execabs"
)

// RemoteRemove drops the defined remote from a git repo.
func (r *Repository) Init() *types.Cmd {
	args := []string{
		"init",
		"-b",
		r.Branch,
	}

	return &types.Cmd{
		Cmd: execabs.Command(gitBin, args...),
	}
}

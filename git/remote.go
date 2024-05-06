package git

import (
	"github.com/thegeeklab/wp-plugin-go/v2/types"
	"golang.org/x/sys/execabs"
)

// RemoteAdd adds an additional remote to a git repo.
func (r *Repository) RemoteAdd() *types.Cmd {
	args := []string{
		"remote",
		"add",
		"origin",
		r.RemoteURL,
	}

	return &types.Cmd{
		Cmd: execabs.Command(gitBin, args...),
	}
}

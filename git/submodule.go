package git

import (
	"github.com/thegeeklab/wp-plugin-go/v3/types"
	"golang.org/x/sys/execabs"
)

// SubmoduleUpdate recursively initializes and updates submodules.
func (r *Repository) SubmoduleUpdate() *types.Cmd {
	args := []string{
		"submodule",
		"update",
		"--init",
		"--recursive",
	}

	if r.SubmodulePartial {
		args = append(args, "--depth=1", "--recommend-shallow")
	}

	if r.SubmoduleRemote {
		args = append(args, "--remote")
	}

	return &types.Cmd{
		Cmd: execabs.Command(gitBin, args...),
	}
}

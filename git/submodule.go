package git

import (
	"golang.org/x/sys/execabs"
)

// SubmoduleUpdate recursively initializes and updates submodules.
func SubmoduleUpdate(repo Repository) *execabs.Cmd {
	args := []string{
		"submodule",
		"update",
		"--init",
		"--recursive",
	}

	if repo.SubmodulePartial {
		args = append(args, "--depth=1", "--recommend-shallow")
	}

	cmd := execabs.Command(
		gitBin,
		args...,
	)

	if repo.SubmoduleRemote {
		cmd.Args = append(cmd.Args, "--remote")
	}

	return cmd
}

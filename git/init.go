package git

import (
	"golang.org/x/sys/execabs"
)

// RemoteRemove drops the defined remote from a git repo.
func Init(repo Repository) *execabs.Cmd {
	args := []string{
		"init",
	}

	if repo.Branch != "" {
		args = []string{
			"init",
			"-b",
			repo.Branch,
		}
	}

	cmd := execabs.Command(
		gitBin,
		args...,
	)

	return cmd
}

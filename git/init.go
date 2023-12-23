package git

import (
	"golang.org/x/sys/execabs"
)

// RemoteRemove drops the defined remote from a git repo.
func Init(repo Repository) *execabs.Cmd {
	args := []string{
		"init",
		"-b",
		repo.Branch,
	}

	return execabs.Command(
		gitBin,
		args...,
	)
}

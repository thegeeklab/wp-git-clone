package git

import (
	"golang.org/x/sys/execabs"
)

// RemoteAdd adds an additional remote to a git repo.
func RemoteAdd(url string) *execabs.Cmd {
	args := []string{
		"remote",
		"add",
		"origin",
		url,
	}

	return execabs.Command(
		gitBin,
		args...,
	)
}

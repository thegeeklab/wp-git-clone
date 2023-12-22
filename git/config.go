package git

import (
	"fmt"
	"strconv"

	"golang.org/x/sys/execabs"
)

// ConfigSSLVerify disables globally the git ssl verification.
func ConfigSSLVerify(repo Repository) *execabs.Cmd {
	args := []string{
		"config",
		"--local",
		"http.sslVerify",
		strconv.FormatBool(repo.InsecureSSLVerify),
	}

	return execabs.Command(
		gitBin,
		args...,
	)
}

// ConfigSafeDirectory disables globally the git ssl verification.
func ConfigSafeDirectory(repo Repository) *execabs.Cmd {
	args := []string{
		"config",
		"--local",
		"--replace-all",
		"safe.directory",
		repo.SafeDirectory,
	}

	return execabs.Command(
		gitBin,
		args...,
	)
}

// ConfigRemapSubmodule returns a git command that, when executed configures git to
// remap submodule urls.
func ConfigRemapSubmodule(name, url string) *execabs.Cmd {
	args := []string{
		"config",
		"--local",
		fmt.Sprintf("submodule.%s.url", name),
		url,
	}

	return execabs.Command(
		gitBin,
		args...,
	)
}

// ConfigSSHCommand sets custom SSH key.
func ConfigSSHCommand(sshKey string) *execabs.Cmd {
	args := []string{
		"config",
		"--local",
		"core.sshCommand",
		"ssh -i " + sshKey,
	}

	return execabs.Command(
		gitBin,
		args...,
	)
}

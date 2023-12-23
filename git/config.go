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
		"--global",
		"http.sslVerify",
		strconv.FormatBool(!repo.InsecureSkipSSLVerify),
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
		"--global",
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
		"--global",
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
		"--global",
		"core.sshCommand",
		"ssh -i " + sshKey,
	}

	return execabs.Command(
		gitBin,
		args...,
	)
}

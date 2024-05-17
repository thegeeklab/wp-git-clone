package git

import (
	"fmt"
	"strconv"

	"github.com/thegeeklab/wp-plugin-go/v3/types"
	"golang.org/x/sys/execabs"
)

// ConfigSSLVerify disables globally the git ssl verification.
func (r *Repository) ConfigSSLVerify(skipVerify bool) *types.Cmd {
	args := []string{
		"config",
		"--global",
		"http.sslVerify",
		strconv.FormatBool(!skipVerify),
	}

	return &types.Cmd{
		Cmd: execabs.Command(gitBin, args...),
	}
}

// ConfigSafeDirectory disables globally the git ssl verification.
func (r *Repository) ConfigSafeDirectory() *types.Cmd {
	args := []string{
		"config",
		"--global",
		"--replace-all",
		"safe.directory",
		r.SafeDirectory,
	}

	return &types.Cmd{
		Cmd: execabs.Command(gitBin, args...),
	}
}

// ConfigRemapSubmodule returns a git command that, when executed configures git to
// remap submodule urls.
func (r *Repository) ConfigRemapSubmodule(name, url string) *types.Cmd {
	args := []string{
		"config",
		"--global",
		fmt.Sprintf("submodule.%s.url", name),
		url,
	}

	return &types.Cmd{
		Cmd: execabs.Command(gitBin, args...),
	}
}

// ConfigSSHCommand sets custom SSH key.
func (r *Repository) ConfigSSHCommand(sshKey string) *types.Cmd {
	args := []string{
		"config",
		"--global",
		"core.sshCommand",
		"ssh -i " + sshKey,
	}

	cmd := &types.Cmd{
		Cmd: execabs.Command(gitBin, args...),
	}
	cmd.SetTrace(false)

	return cmd
}

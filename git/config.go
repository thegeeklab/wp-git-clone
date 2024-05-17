package git

import (
	"fmt"
	"strconv"

	plugin_exec "github.com/thegeeklab/wp-plugin-go/v3/exec"
)

// ConfigSSLVerify disables globally the git ssl verification.
func (r *Repository) ConfigSSLVerify(skipVerify bool) *plugin_exec.Cmd {
	args := []string{
		"config",
		"--global",
		"http.sslVerify",
		strconv.FormatBool(!skipVerify),
	}

	cmd := plugin_exec.Command(gitBin, args...)
	cmd.Trace = false

	return cmd
}

// ConfigSafeDirectory disables globally the git ssl verification.
func (r *Repository) ConfigSafeDirectory() *plugin_exec.Cmd {
	args := []string{
		"config",
		"--global",
		"--replace-all",
		"safe.directory",
		r.SafeDirectory,
	}

	cmd := plugin_exec.Command(gitBin, args...)
	cmd.Trace = false

	return cmd
}

// ConfigRemapSubmodule returns a git command that, when executed configures git to
// remap submodule urls.
func (r *Repository) ConfigRemapSubmodule(name, url string) *plugin_exec.Cmd {
	args := []string{
		"config",
		"--global",
		fmt.Sprintf("submodule.%s.url", name),
		url,
	}

	cmd := plugin_exec.Command(gitBin, args...)
	cmd.Trace = false

	return cmd
}

// ConfigSSHCommand sets custom SSH key.
func (r *Repository) ConfigSSHCommand(sshKey string) *plugin_exec.Cmd {
	args := []string{
		"config",
		"--global",
		"core.sshCommand",
		"ssh -i " + sshKey,
	}

	cmd := plugin_exec.Command(gitBin, args...)
	cmd.Trace = false

	return cmd
}

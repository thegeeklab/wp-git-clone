package git

import (
	"os"

	plugin_exec "github.com/thegeeklab/wp-plugin-go/v4/exec"
)

// SubmoduleUpdate recursively initializes and updates submodules.
func (r *Repository) SubmoduleUpdate() *plugin_exec.Cmd {
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

	cmd := plugin_exec.Command(gitBin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

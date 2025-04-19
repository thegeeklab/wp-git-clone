package git

import (
	"os"

	plugin_exec "github.com/thegeeklab/wp-plugin-go/v6/exec"
)

// RemoteRemove drops the defined remote from a git repo.
func (r *Repository) Init() *plugin_exec.Cmd {
	args := []string{
		"init",
		"-b",
		r.Branch,
	}

	cmd := plugin_exec.Command(gitBin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

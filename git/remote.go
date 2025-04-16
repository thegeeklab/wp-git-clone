package git

import (
	"os"

	plugin_exec "github.com/thegeeklab/wp-plugin-go/v5/exec"
)

// RemoteAdd adds an additional remote to a git repo.
func (r *Repository) RemoteAdd() *plugin_exec.Cmd {
	args := []string{
		"remote",
		"add",
		"origin",
		r.RemoteURL,
	}

	cmd := plugin_exec.Command(gitBin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

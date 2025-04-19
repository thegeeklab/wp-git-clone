package git

import (
	"fmt"
	"os"

	plugin_exec "github.com/thegeeklab/wp-plugin-go/v6/exec"
)

// FetchSource fetches the source from remote.
func (r *Repository) FetchSource(ref string) *plugin_exec.Cmd {
	args := []string{
		"fetch",
	}

	if r.Depth != 0 {
		args = append(args, fmt.Sprintf("--depth=%d", r.Depth))
	}

	if r.Filter != "" {
		args = append(args, "--filter", r.Filter)
	}

	args = append(args, "origin")
	args = append(args, fmt.Sprintf("+%s:", ref))

	cmd := plugin_exec.Command(gitBin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

// FetchTags fetches the source from remote.
func (r *Repository) FetchTags() *plugin_exec.Cmd {
	args := []string{
		"fetch",
		"--tags",
		"--quiet",
		"origin",
	}

	cmd := plugin_exec.Command(gitBin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

// FetchLFS fetches lfs.
func (r *Repository) FetchLFS() *plugin_exec.Cmd {
	args := []string{
		"lfs",
		"fetch",
	}

	cmd := plugin_exec.Command(gitBin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

// CheckoutHead handles head checkout.
func (r *Repository) CheckoutHead() *plugin_exec.Cmd {
	args := []string{
		"checkout",
		"--force",
		"--quiet",
		"FETCH_HEAD",
	}

	cmd := plugin_exec.Command(gitBin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

// CheckoutSha handles commit checkout.
func (r *Repository) CheckoutSha() *plugin_exec.Cmd {
	args := []string{
		"reset",
		"--hard",
		"--quiet",
		r.CommitSha,
	}

	cmd := plugin_exec.Command(gitBin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

// CheckoutLFS handles commit checkout.
func (r *Repository) CheckoutLFS() *plugin_exec.Cmd {
	args := []string{
		"lfs",
		"checkout",
	}

	cmd := plugin_exec.Command(gitBin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

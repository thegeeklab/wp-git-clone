package main

import (
	"github.com/thegeeklab/wp-git-clone/plugin"
	"github.com/thegeeklab/wp-plugin-go/types"
	"github.com/urfave/cli/v2"
)

// settingsFlags has the cli.Flags for the plugin.Settings.
//
//go:generate go run docs.go flags.go
func settingsFlags(settings *plugin.Settings, category string) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "remote",
			Usage:       "git remote HTTP clone url",
			EnvVars:     []string{"PLUGIN_REMOTE", "CI_REPO_CLONE_URL"},
			Destination: &settings.Repo.RemoteURL,
			DefaultText: "$CI_REPO_CLONE_URL",
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "remote-ssh",
			Usage:       "git remote SSH clone url",
			EnvVars:     []string{"PLUGIN_REMOTE_SSH", "CI_REPO_CLONE_SSH_URL"},
			Destination: &settings.Repo.RemoteSSH,
			DefaultText: "$CI_REPO_CLONE_SSH_URL",
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "workdir",
			Usage:       "path to clone git repository",
			EnvVars:     []string{"PLUGIN_WORKDIR", "CI_WORKSPACE"},
			Destination: &settings.WorkDir,
			DefaultText: "$CI_WORKSPACE",
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "sha",
			Usage:       "git commit sha",
			EnvVars:     []string{"PLUGIN_COMMIT_SHA", "CI_COMMIT_SHA"},
			Destination: &settings.Repo.CommitSha,
			DefaultText: "$CI_COMMIT_SHA",
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "ref",
			Usage:       "git commit ref",
			EnvVars:     []string{"PLUGIN_COMMIT_REF", "CI_COMMIT_REF"},
			Value:       "refs/heads/main",
			Destination: &settings.Repo.CommitRef,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "netrc.machine",
			Usage:       "netrc machine",
			EnvVars:     []string{"CI_NETRC_MACHINE"},
			Destination: &settings.Netrc.Machine,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "netrc.username",
			Usage:       "netrc username",
			EnvVars:     []string{"CI_NETRC_USERNAME"},
			Destination: &settings.Netrc.Password,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "netrc.password",
			Usage:       "netrc password",
			EnvVars:     []string{"CI_NETRC_PASSWORD"},
			Destination: &settings.Netrc.Password,
			Category:    category,
		},
		&cli.IntFlag{
			Name:        "depth",
			Usage:       "clone depth",
			EnvVars:     []string{"PLUGIN_DEPTH"},
			Destination: &settings.Depth,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "recursive",
			Usage:       "clone submodules",
			EnvVars:     []string{"PLUGIN_RECURSIVE"},
			Value:       true,
			Destination: &settings.Recursive,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "tags",
			Usage:       "fetch git tags during clone",
			EnvVars:     []string{"PLUGIN_TAGS"},
			Value:       true,
			Destination: &settings.Tags,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "insecure-skip-ssl-verify",
			Usage:       "skip SSL verification of the remote machine",
			EnvVars:     []string{"PLUGIN_INSECURE_SKIP_SSL_VERIFY"},
			Destination: &settings.Repo.InsecureSkipSSLVerify,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "submodule-update-remote",
			Usage:       "update remote submodules",
			EnvVars:     []string{"PLUGIN_SUBMODULES_UPDATE_REMOTE", "PLUGIN_SUBMODULE_UPDATE_REMOTE"},
			Destination: &settings.Repo.SubmoduleRemote,
			Category:    category,
		},
		&cli.GenericFlag{
			Name:     "submodule-override",
			Usage:    "JSON map of submodule overrides",
			EnvVars:  []string{"PLUGIN_SUBMODULE_OVERRIDE"},
			Value:    &types.MapFlag{},
			Category: category,
		},
		&cli.BoolFlag{
			Name:        "submodule-partial",
			Usage:       "update submodules via partial clone (`depth=1`)",
			EnvVars:     []string{"PLUGIN_SUBMODULES_PARTIAL", "PLUGIN_SUBMODULE_PARTIAL"},
			Value:       true,
			Destination: &settings.Repo.SubmodulePartial,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "lfs",
			Usage:       "whether to retrieve LFS content if available",
			EnvVars:     []string{"PLUGIN_LFS"},
			Value:       true,
			Destination: &settings.Lfs,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "branch",
			Usage:       "change branch name",
			EnvVars:     []string{"PLUGIN_BRANCH", "CI_COMMIT_BRANCH", "CI_REPO_DEFAULT_BRANCH"},
			Value:       "main",
			Destination: &settings.Repo.Branch,
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "partial",
			Usage:       "enable/disable partial clone",
			EnvVars:     []string{"PLUGIN_PARTIAL"},
			Destination: &settings.Partial,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "safe-directory",
			Usage:       "define/replace safe directories",
			EnvVars:     []string{"PLUGIN_SAFE_DIRECTORY", "CI_WORKSPACE"},
			Destination: &settings.Repo.SafeDirectory,
			DefaultText: "$CI_WORKSPACE",
			Category:    category,
		},
		&cli.BoolFlag{
			Name:        "use-ssh",
			Usage:       "using SSH for git clone",
			EnvVars:     []string{"PLUGIN_USE_SSH"},
			Destination: &settings.UseSSH,
			Category:    category,
		},
		&cli.StringFlag{
			Name:        "ssh-key",
			Usage:       "Private key for SSH clone",
			EnvVars:     []string{"PLUGIN_SSH_KEY"},
			Destination: &settings.SSHKey,
			Category:    category,
		},
	}
}

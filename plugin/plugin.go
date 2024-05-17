package plugin

import (
	"fmt"

	"github.com/thegeeklab/wp-git-clone/git"
	plugin_base "github.com/thegeeklab/wp-plugin-go/v3/plugin"
	plugin_types "github.com/thegeeklab/wp-plugin-go/v3/types"
	"github.com/urfave/cli/v2"
)

//go:generate go run ../internal/doc/main.go -output=../docs/data/data-raw.yaml

// Plugin implements provide the plugin.
type Plugin struct {
	*plugin_base.Plugin
	Settings *Settings
}

type Netrc struct {
	Machine  string
	Login    string
	Password string
}

// Settings for the plugin.
type Settings struct {
	Recursive bool
	Tags      bool
	Lfs       bool
	Partial   bool
	Home      string
	SSHKey    string

	Netrc Netrc
	Repo  git.Repository
}

func New(e plugin_base.ExecuteFunc, build ...string) *Plugin {
	p := &Plugin{
		Settings: &Settings{},
	}

	options := plugin_base.Options{
		Name:                "wp-git-clone",
		Description:         "Clone git repository",
		Flags:               Flags(p.Settings, plugin_base.FlagsPluginCategory),
		Execute:             p.run,
		HideWoodpeckerFlags: true,
	}

	if len(build) > 0 {
		options.Version = build[0]
	}

	if len(build) > 1 {
		options.VersionMetadata = fmt.Sprintf("date=%s", build[1])
	}

	if e != nil {
		options.Execute = e
	}

	p.Plugin = plugin_base.New(options)

	return p
}

// Flags returns a slice of CLI flags for the plugin.
func Flags(settings *Settings, category string) []cli.Flag {
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
			Destination: &settings.Repo.WorkDir,
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
			Destination: &settings.Repo.Depth,
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
			Value:    &plugin_types.MapFlag{},
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
		&cli.StringFlag{
			Name:        "ssh-key",
			Usage:       "private key for SSH clone",
			EnvVars:     []string{"PLUGIN_SSH_KEY"},
			Destination: &settings.SSHKey,
			Category:    category,
		},
	}
}

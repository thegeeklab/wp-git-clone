package plugin

import (
	"github.com/thegeeklab/wp-git-clone/git"
	wp "github.com/thegeeklab/wp-plugin-go/plugin"
)

// Plugin implements provide the plugin.
type Plugin struct {
	*wp.Plugin
	Settings *Settings
}

type Netrc struct {
	Machine  string
	Login    string
	Password string
}

// Settings for the plugin.
type Settings struct {
	Depth     int
	Recursive bool
	Tags      bool
	Lfs       bool
	Partial   bool
	Filter    string
	UseSSH    bool
	SSHKey    string
	Home      string
	WorkDir   string

	Netrc Netrc
	Repo  git.Repository
}

func New(options wp.Options, settings *Settings) *Plugin {
	p := &Plugin{}

	if options.Execute == nil {
		options.Execute = p.run
	}

	p.Plugin = wp.New(options)
	p.Settings = settings

	return p
}

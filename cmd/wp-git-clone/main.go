package main

import (
	"fmt"

	"github.com/thegeeklab/wp-git-clone/plugin"

	wp "github.com/thegeeklab/wp-plugin-go/v2/plugin"
)

//nolint:gochecknoglobals
var (
	BuildVersion = "devel"
	BuildDate    = "00000000"
)

func main() {
	settings := &plugin.Settings{}
	options := wp.Options{
		Name:            "wp-git-clone",
		Description:     "Clone git repository",
		Version:         BuildVersion,
		VersionMetadata: fmt.Sprintf("date=%s", BuildDate),
		Flags:           settingsFlags(settings, wp.FlagsPluginCategory),
	}

	plugin.New(options, settings).Run()
}

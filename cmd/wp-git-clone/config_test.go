package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thegeeklab/wp-git-clone/plugin"
	wp "github.com/thegeeklab/wp-plugin-go/plugin"
)

func Test_pluginOptions(t *testing.T) {
	tests := []struct {
		name string
		envs map[string]string
		want string
	}{
		{
			name: "ensure default branch is set",
			envs: map[string]string{
				"CI_COMMIT_BRANCH": "",
			},
			want: "main",
		},
	}

	for _, tt := range tests {
		for key, value := range tt.envs {
			t.Setenv(key, value)
		}

		settings := &plugin.Settings{}
		options := wp.Options{
			Name:    "wp-git-clone",
			Flags:   settingsFlags(settings, wp.FlagsPluginCategory),
			Execute: func(_ context.Context) error { return nil },
		}

		got := plugin.New(options, settings)

		_ = got.App.Run([]string{"wp-git-clone"})
		_ = got.Validate()
		_ = got.FlagsFromContext()

		assert.EqualValues(t, tt.want, got.Settings.Repo.Branch)
	}
}

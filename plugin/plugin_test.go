package plugin

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

func setupPluginTest(t *testing.T) (*Plugin, error) {
	t.Helper()

	cli.HelpPrinter = func(_ io.Writer, _ string, _ interface{}) {}
	got := New(func(_ context.Context) error { return nil })
	err := got.App.Run(t.Context(), []string{"wp-s3-action"})

	return got, err
}

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

		got, _ := setupPluginTest(t)

		assert.EqualValues(t, tt.want, got.Settings.Repo.Branch)
	}
}

package plugin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
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

		got := New(func(_ context.Context) error { return nil })

		_ = got.App.Run([]string{"wp-git-clone"})
		_ = got.Validate()
		_ = got.FlagsFromContext()

		assert.EqualValues(t, tt.want, got.Settings.Repo.Branch)
	}
}

package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoteAdd(t *testing.T) {
	tests := []struct {
		name string
		repo Repository
		want []string
	}{
		{
			name: "add remote with valid inputs",
			repo: Repository{
				RemoteURL: "https://example.com/repo.git",
			},
			want: []string{gitBin, "remote", "add", "origin", "https://example.com/repo.git"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.repo.RemoteAdd()
			assert.Equal(t, tt.want, cmd.Cmd.Args)
			assert.Equal(t, tt.repo.WorkDir, cmd.Cmd.Dir)
		})
	}
}

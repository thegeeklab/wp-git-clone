package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUpdateSubmodules tests if the arguments to `git submodule update`
// are constructed properly.
func TestUpdateSubmodules(t *testing.T) {
	tests := []struct {
		name string
		repo *Repository
		want []string
	}{
		{
			name: "full submodule update",
			repo: &Repository{
				SubmodulePartial: false,
			},
			want: []string{
				gitBin,
				"submodule",
				"update",
				"--init",
				"--recursive",
			},
		},
		{
			name: "partial submodule update",
			repo: &Repository{
				SubmodulePartial: true,
			},
			want: []string{
				gitBin,
				"submodule",
				"update",
				"--init",
				"--recursive",
				"--depth=1",
				"--recommend-shallow",
			},
		},
		{
			name: "submodule update with remote",
			repo: &Repository{
				SubmoduleRemote: true,
			},
			want: []string{
				gitBin,
				"submodule",
				"update",
				"--init",
				"--recursive",
				"--remote",
			},
		},
		{
			name: "submodule update with remote and partial",
			repo: &Repository{
				SubmoduleRemote:  true,
				SubmodulePartial: true,
			},
			want: []string{
				gitBin,
				"submodule",
				"update",
				"--init",
				"--recursive",
				"--depth=1",
				"--recommend-shallow",
				"--remote",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.repo.SubmoduleUpdate()
			assert.Equal(t, tt.want, cmd.Args)
		})
	}
}

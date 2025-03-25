package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name     string
		repo     Repository
		expected []string
	}{
		{
			name: "init repo",
			repo: Repository{
				Branch: "main",
			},
			expected: []string{gitBin, "init", "-b", "main"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.repo.Init()
			assert.Equal(t, tt.expected, cmd.Args)
		})
	}
}

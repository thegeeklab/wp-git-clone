package git

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetch(t *testing.T) {
	testdata := []struct {
		name  string
		repo  *Repository
		tags  bool
		depth int
		want  []string
	}{
		{
			name: "fetch main without tags",
			repo: &Repository{
				CommitRef: "refs/heads/main",
			},
			tags:  false,
			depth: 0,
			want: []string{
				gitBin,
				"fetch",
				"origin",
				"+refs/heads/main:",
			},
		},
		{
			name: "fetch main without tags with depth",
			repo: &Repository{
				CommitRef: "refs/heads/main",
				Depth:     50,
			},
			tags:  false,
			depth: 50,
			want: []string{
				gitBin,
				"fetch",
				"--depth=50",
				"origin",
				"+refs/heads/main:",
			},
		},
		{
			name: "fetch main with tags and depth",
			repo: &Repository{
				CommitRef: "refs/heads/main",
				Depth:     100,
			},
			tags:  true,
			depth: 100,
			want: []string{
				gitBin,
				"fetch",
				"--depth=100",
				"origin",
				"+refs/heads/main:",
			},
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.repo.FetchSource(tt.repo.CommitRef)
			require.Equal(t, tt.want, cmd.Args)
		})
	}
}

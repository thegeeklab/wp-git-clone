package git

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetchSource(t *testing.T) {
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

func TestFetchTags(t *testing.T) {
	testdata := []struct {
		name  string
		repo  *Repository
		tags  bool
		depth int
		want  []string
	}{
		{
			name: "fetch tags",
			repo: &Repository{},
			tags: true,
			want: []string{
				gitBin,
				"fetch",
				"--tags",
				"--quiet",
				"origin",
			},
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.repo.FetchTags()
			require.Equal(t, tt.want, cmd.Args)
		})
	}
}

func TestFetchLFS(t *testing.T) {
	testdata := []struct {
		name string
		repo *Repository
		want []string
	}{
		{
			name: "fetch LFS",
			repo: &Repository{},
			want: []string{
				gitBin,
				"lfs",
				"fetch",
			},
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.repo.FetchLFS()
			require.Equal(t, tt.want, cmd.Args)
		})
	}
}

func TestCheckoutHead(t *testing.T) {
	testdata := []struct {
		name string
		repo *Repository
		want []string
	}{
		{
			name: "checkout head",
			repo: &Repository{},
			want: []string{
				gitBin,
				"checkout",
				"--force",
				"--quiet",
				"FETCH_HEAD",
			},
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.repo.CheckoutHead()
			require.Equal(t, tt.want, cmd.Args)
		})
	}
}

func TestCheckoutSha(t *testing.T) {
	testdata := []struct {
		name string
		repo *Repository
		want []string
	}{
		{
			name: "checkout sha",
			repo: &Repository{
				CommitSha: "abcd1234",
			},
			want: []string{
				gitBin,
				"reset",
				"--hard",
				"--quiet",
				"abcd1234",
			},
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.repo.CheckoutSha()
			require.Equal(t, tt.want, cmd.Args)
		})
	}
}

func TestCheckoutLFS(t *testing.T) {
	testdata := []struct {
		name string
		repo *Repository
		want []string
	}{
		{
			name: "checkout LFS with no arguments",
			repo: &Repository{},
			want: []string{
				gitBin,
				"lfs",
				"checkout",
			},
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.repo.CheckoutLFS()
			require.Equal(t, tt.want, cmd.Args)
		})
	}
}

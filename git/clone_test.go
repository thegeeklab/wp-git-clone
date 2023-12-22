package git

import (
	"testing"
)

// TestFetch tests if the arguments to `git fetch` are constructed properly.
func TestFetch(t *testing.T) {
	testdata := []struct {
		ref   string
		tags  bool
		depth int
		exp   []string
	}{
		{
			"refs/heads/master",
			false,
			0,
			[]string{
				"/usr/bin/git",
				"fetch",
				"--no-tags",
				"origin",
				"+refs/heads/master:",
			},
		},
		{
			"refs/heads/master",
			false,
			50,
			[]string{
				"/usr/bin/git",
				"fetch",
				"--no-tags",
				"--depth=50",
				"origin",
				"+refs/heads/master:",
			},
		},
		{
			"refs/heads/master",
			true,
			100,
			[]string{
				"/usr/bin/git",
				"fetch",
				"--tags",
				"--depth=100",
				"origin",
				"+refs/heads/master:",
			},
		},
	}
	for _, td := range testdata {
		c := FetchSource(td.ref, td.tags, td.depth, "")
		if len(c.Args) != len(td.exp) {
			t.Errorf("Expected: %s, got %s", td.exp, c.Args)
		}

		for i := range c.Args {
			if c.Args[i] != td.exp[i] {
				t.Errorf("Expected: %s, got %s", td.exp, c.Args)
			}
		}
	}
}

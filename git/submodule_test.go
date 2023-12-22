package git

import (
	"testing"
)

// TestUpdateSubmodules tests if the arguments to `git submodule update`
// are constructed properly.
func TestUpdateSubmodules(t *testing.T) {
	tests := []struct {
		partial bool
		exp     []string
	}{
		{
			false,
			[]string{
				"/usr/bin/git",
				"submodule",
				"update",
				"--init",
				"--recursive",
			},
		},
		{
			true,
			[]string{
				"/usr/bin/git",
				"submodule",
				"update",
				"--init",
				"--recursive",
				"--depth=1",
				"--recommend-shallow",
			},
		},
	}
	for _, tt := range tests {
		repo := Repository{
			SubmoduleRemote:  false,
			SubmodulePartial: tt.partial,
		}

		c := SubmoduleUpdate(repo)
		if len(c.Args) != len(tt.exp) {
			t.Errorf("Expected: %s, got %s", tt.exp, c.Args)
		}

		for i := range c.Args {
			if c.Args[i] != tt.exp[i] {
				t.Errorf("Expected: %s, got %s", tt.exp, c.Args)
			}
		}
	}
}

// TestUpdateSubmodules tests if the arguments to `git submodule update`
// are constructed properly.
func TestUpdateSubmodulesRemote(t *testing.T) {
	tests := []struct {
		exp []string
	}{
		{
			[]string{
				"/usr/bin/git",
				"submodule",
				"update",
				"--init",
				"--recursive",
				"--remote",
			},
		},
		{
			[]string{
				"/usr/bin/git",
				"submodule",
				"update",
				"--init",
				"--recursive",
				"--remote",
			},
		},
	}
	for _, tt := range tests {
		repo := Repository{
			SubmoduleRemote:  true,
			SubmodulePartial: false,
		}

		c := SubmoduleUpdate(repo)
		if len(c.Args) != len(tt.exp) {
			t.Errorf("Expected: %s, got %s", tt.exp, c.Args)
		}

		for i := range c.Args {
			if c.Args[i] != tt.exp[i] {
				t.Errorf("Expected: %s, got %s", tt.exp, c.Args)
			}
		}
	}
}

package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/thegeeklab/wp-git-clone/git"
)

type testCommit struct {
	name      string
	path      string
	clone     string
	event     string
	commit    string
	ref       string
	file      string
	data      string
	dataSize  int64
	recursive bool
	lfs       bool
}

// TestClone tests the ability to clone a specific commit into
// a fresh, empty directory every time.
func TestClone(t *testing.T) {
	for _, tt := range getCommits() {
		dir := setup()
		defer teardown(dir)

		plugin := Plugin{
			Settings: &Settings{
				Repo: git.Repository{
					RemoteURL: tt.clone,
					CommitRef: tt.ref,
					CommitSha: tt.commit,
					Branch:    "main",
				},
				Pipeline: Pipeline{
					Event: tt.event,
				},
				WorkDir:   filepath.Join(dir, tt.path),
				Recursive: tt.recursive,
				Lfs:       tt.lfs,
			},
		}

		if err := plugin.Execute(); err != nil {
			t.Errorf("Expected successful clone. Got error. %s.", err)
		}

		if tt.data != "" {
			data := readFile(plugin.Settings.WorkDir, tt.file)
			if data != tt.data {
				t.Errorf("Expected %s to contain [%s]. Got [%s].", tt.file, tt.data, data)
			}
		}

		if tt.dataSize != 0 {
			size := getFileSize(plugin.Settings.WorkDir, tt.file)
			if size != tt.dataSize {
				t.Errorf("Expected %s size to be [%d]. Got [%d].", tt.file, tt.dataSize, size)
			}
		}
	}
}

// TestCloneNonEmpty tests the ability to clone a specific commit into
// a non-empty directory. This is useful if the git workspace is cached
// and re-stored for every workflow.
func TestCloneNonEmpty(t *testing.T) {
	dir := setup()
	defer teardown(dir)

	for _, tt := range getCommits() {
		plugin := Plugin{
			Settings: &Settings{
				Repo: git.Repository{
					RemoteURL: tt.clone,
					CommitRef: tt.ref,
					CommitSha: tt.commit,
					Branch:    "main",
				},
				Pipeline: Pipeline{
					Event: tt.event,
				},
				WorkDir:   filepath.Join(dir, tt.path),
				Recursive: tt.recursive,
				Lfs:       tt.lfs,
			},
		}

		fmt.Println(plugin.Settings.Repo.CommitSha, tt.commit, fmt.Sprintf("%q", tt.data))

		if err := plugin.Execute(); err != nil {
			t.Errorf("Expected successful clone. Got error. %s.", err)
		}

		if tt.data != "" {
			data := readFile(plugin.Settings.WorkDir, tt.file)
			if data != tt.data {
				t.Errorf("Expected %s to contain [%q]. Got [%q].", tt.file, tt.data, data)

				break
			}
		}

		if tt.dataSize != 0 {
			size := getFileSize(plugin.Settings.WorkDir, tt.file)
			if size != tt.dataSize {
				t.Errorf("Expected %s size to be [%d]. Got [%d].", tt.file, tt.dataSize, size)
			}
		}
	}
}

// helper function that will setup a temporary workspace
// to which we can clone the repositroy.
func setup() string {
	dir, _ := os.MkdirTemp("/tmp", "plugin_git_test_")
	_ = os.Mkdir(dir, os.ModePerm)

	return dir
}

// helper function to delete the temporary workspace.
func teardown(dir string) {
	os.RemoveAll(dir)
}

// helper function to read a file in the temporary worskapce.
func readFile(dir, file string) string {
	filename := filepath.Join(dir, file)
	fmt.Println(filename)
	data, _ := os.ReadFile(filename)

	return string(data)
}

func getFileSize(dir, file string) int64 {
	filename := filepath.Join(dir, file)
	fi, _ := os.Stat(filename)

	return fi.Size()
}

func getCommits() []testCommit {
	return []testCommit{
		{
			name:   "first commit",
			path:   "octocat/Hello-World",
			clone:  "https://github.com/octocat/Hello-World.git",
			event:  "push",
			commit: "553c2077f0edc3d5dc5d17262f6aa498e69d6f8e",
			ref:    "refs/heads/master",
			file:   "README",
			data:   "Hello World!",
		},
		{
			name:   "head commit",
			path:   "octocat/Hello-World",
			clone:  "https://github.com/octocat/Hello-World.git",
			event:  "push",
			commit: "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			ref:    "refs/heads/master",
			file:   "README",
			data:   "Hello World!\n",
		},
		{
			name:   "pull request commit",
			path:   "octocat/Hello-World",
			clone:  "https://github.com/octocat/Hello-World.git",
			event:  "pull_request",
			commit: "762941318ee16e59dabbacb1b4049eec22f0d303",
			ref:    "refs/pull/6/merge",
			file:   "README",
			data:   "Hello World!\n",
		},
		{
			name:   "branch",
			path:   "octocat/Hello-World",
			clone:  "https://github.com/octocat/Hello-World.git",
			event:  "push",
			commit: "b3cbd5bbd7e81436d2eee04537ea2b4c0cad4cdf",
			ref:    "refs/heads/test",
			file:   "CONTRIBUTING.md",
			data:   "## Contributing\n",
		},
		{
			name:   "tags",
			path:   "github/mime-types",
			clone:  "https://github.com/github/mime-types.git",
			event:  "tag",
			commit: "bf68d60215a167c935bc5976b7d06a7ffb290926",
			ref:    "refs/tags/v1.17",
			file:   ".gitignore",
			data:   "*.swp\n*~\n.rake_tasks~\nhtml\ndoc\npkg\npublish\ncoverage\n",
		},
		{
			name:      "submodules",
			path:      "test-assets/woodpecker-git-test-submodule",
			clone:     "https://github.com/test-assets/woodpecker-git-test-submodule.git",
			event:     "push",
			commit:    "cc020eb6aaa601c13ca7b0d5db9d1ca694e7a003",
			ref:       "refs/heads/main",
			file:      "Hello-World/README",
			data:      "Hello World!\n",
			recursive: true,
		},
		{
			name:  "checkout with ref only",
			path:  "octocat/Hello-World",
			clone: "https://github.com/octocat/Hello-World.git",
			event: "push",
			// commit: "a11fb45a696bf1d696fc9ab2c733f8f123aa4cf5",
			ref:  "pull/2403/head",
			file: "README",
			data: "Hello World!\n\nsomething is changed!\n",
		},
		// test lfs, please do not change order, otherwise TestCloneNonEmpty will fail ###
		{
			name:     "checkout with lfs skip",
			path:     "test-assets/woodpecker-git-test-lfs",
			clone:    "https://github.com/test-assets/woodpecker-git-test-lfs.git",
			event:    "push",
			commit:   "69d4dadb4c2899efb73c0095bb58a6454d133cef",
			ref:      "refs/heads/main",
			file:     "4M.bin",
			dataSize: 132,
		},
		{
			name:     "checkout with lfs",
			path:     "test-assets/woodpecker-git-test-lfs",
			clone:    "https://github.com/test-assets/woodpecker-git-test-lfs.git",
			event:    "push",
			commit:   "69d4dadb4c2899efb73c0095bb58a6454d133cef",
			ref:      "refs/heads/main",
			file:     "4M.bin",
			dataSize: 4194304,
			lfs:      true,
		},
	}
}

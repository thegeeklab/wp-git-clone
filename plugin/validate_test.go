package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thegeeklab/wp-git-clone/git"
)

func TestValidateSSHKey(t *testing.T) {
	tests := []struct {
		name string
		key  string
		err  bool
	}{
		{name: "empty", key: "", err: false},
		{name: "absolute path", key: "/tmp/ssh/key", err: false},
		{name: "relative path", key: "relative/key", err: true},
		{name: "shell injection", key: "; cat ~/.netrc; #", err: true},
		{name: "command substitution", key: "/tmp/$(id)", err: true},
		{name: "path with spaces", key: "/tmp/my key", err: false},
		{name: "tab", key: "/tmp/my\tkey", err: true},
		{name: "backtick", key: "/tmp/`id`", err: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSSHKey(tt.key)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateBranch(t *testing.T) {
	tests := []struct {
		name   string
		branch string
		err    bool
	}{
		{name: "empty", branch: "", err: false},
		{name: "valid", branch: "main", err: false},
		{name: "valid with slash", branch: "feature/foo", err: false},
		{name: "leading dash", branch: "--upload-pack=evil", err: true},
		{name: "whitespace", branch: "main; rm -rf /", err: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateBranch(tt.branch)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateCommitSha(t *testing.T) {
	tests := []struct {
		name string
		sha  string
		err  bool
	}{
		{name: "empty", sha: "", err: false},
		{name: "short sha", sha: "abcd1234", err: false},
		{name: "full sha", sha: "553c2077f0edc3d5dc5d17262f6aa498e69d6f8e", err: false},
		{name: "leading dash", sha: "--upload-pack=evil", err: true},
		{name: "non hex", sha: "refs/heads/main", err: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCommitSha(tt.sha)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateCommitRef(t *testing.T) {
	tests := []struct {
		name string
		ref  string
		err  bool
	}{
		{name: "empty", ref: "", err: false},
		{name: "full ref", ref: "refs/heads/main", err: false},
		{name: "short ref", ref: "pull/2403/head", err: false},
		{name: "refspec injection", ref: "refs/heads/*:refs/heads/*", err: true},
		{name: "leading dash", ref: "--upload-pack=evil", err: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCommitRef(tt.ref)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateSubmodule(t *testing.T) {
	tests := []struct {
		name    string
		subName string
		subURL  string
		err     bool
	}{
		{name: "valid", subName: "mysubmodule", subURL: "https://example.com/mysubmodule.git", err: false},
		{name: "newline in name", subName: "evil\nsafe.directory", subURL: "https://example.com/x.git", err: true},
		{name: "leading dash url", subName: "mysubmodule", subURL: "--upload-pack=evil", err: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSubmodule(tt.subName, tt.subURL)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateNetrc(t *testing.T) {
	tests := []struct {
		name  string
		netrc Netrc
		err   bool
	}{
		{name: "valid", netrc: Netrc{Machine: "example.com", Login: "user", Password: "pass"}, err: false},
		{name: "newline in machine", netrc: Netrc{Machine: "example.com\nevil.com"}, err: true},
		{name: "newline in login", netrc: Netrc{Login: "user\nevil"}, err: true},
		{name: "newline in password", netrc: Netrc{Password: "pass\nmachine evil.com\nlogin x\npassword y"}, err: true},
		{name: "carriage return", netrc: Netrc{Password: "pass\r"}, err: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateNetrc(tt.netrc)
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateSettings(t *testing.T) {
	tests := []struct {
		name     string
		settings *Settings
		err      bool
	}{
		{
			name: "valid settings",
			settings: &Settings{
				Repo: git.Repository{
					RemoteURL: "https://github.com/octocat/Hello-World.git",
					Branch:    "main",
					CommitRef: "refs/heads/main",
					CommitSha: "553c2077f0edc3d5dc5d17262f6aa498e69d6f8e",
					Filter:    "tree:0",
				},
			},
			err: false,
		},
		{
			name: "malicious ssh key",
			settings: &Settings{
				SSHKey: "; cat ~/.netrc; #",
				Repo: git.Repository{
					Branch: "main",
				},
			},
			err: true,
		},
		{
			name: "malicious commit sha",
			settings: &Settings{
				Repo: git.Repository{
					Branch:    "main",
					CommitSha: "--upload-pack=evil",
				},
			},
			err: true,
		},
		{
			name: "malicious netrc password",
			settings: &Settings{
				Netrc: Netrc{Password: "pass\nmachine evil.com\nlogin x\npassword y"},
				Repo: git.Repository{
					Branch: "main",
				},
			},
			err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Plugin{Settings: tt.settings}
			err := p.Validate()

			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigSSLVerify(t *testing.T) {
	tests := []struct {
		name       string
		repo       Repository
		skipVerify bool
		want       []string
	}{
		{
			name:       "enable SSL verification",
			repo:       Repository{},
			skipVerify: false,
			want:       []string{gitBin, "config", "--global", "http.sslVerify", "true"},
		},
		{
			name:       "disable SSL verification",
			repo:       Repository{},
			skipVerify: true,
			want:       []string{gitBin, "config", "--global", "http.sslVerify", "false"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.repo.ConfigSSLVerify(tt.skipVerify)
			assert.Equal(t, tt.want, cmd.Args)
		})
	}
}

func TestConfigSafeDirectory(t *testing.T) {
	tests := []struct {
		name    string
		repo    Repository
		safeDir string
		want    []string
	}{
		{
			name: "set safe directory",
			repo: Repository{
				SafeDirectory: "/path/to/safe/dir",
			},
			want: []string{gitBin, "config", "--global", "--replace-all", "safe.directory", "/path/to/safe/dir"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.repo.ConfigSafeDirectory()
			assert.Equal(t, tt.want, cmd.Args)
		})
	}
}

func TestConfigRemapSubmodule(t *testing.T) {
	tests := []struct {
		name    string
		repo    Repository
		subName string
		subURL  string
		want    []string
	}{
		{
			name:    "remap submodule URL",
			repo:    Repository{},
			subName: "mysubmodule",
			subURL:  "https://example.com/mysubmodule.git",
			want: []string{
				gitBin, "config", "--global", "submodule.mysubmodule.url",
				"https://example.com/mysubmodule.git",
			},
		},
		{
			name:    "remap submodule URL with spaces",
			repo:    Repository{},
			subName: "my submodule",
			subURL:  "https://example.com/my submodule.git",
			want: []string{
				gitBin, "config", "--global", "submodule.my submodule.url",
				"https://example.com/my submodule.git",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.repo.ConfigRemapSubmodule(tt.subName, tt.subURL)
			assert.Equal(t, tt.want, cmd.Args)
		})
	}
}

func TestConfigSSHCommand(t *testing.T) {
	tests := []struct {
		name   string
		repo   Repository
		sshKey string
		want   []string
	}{
		{
			name:   "set SSH command with key",
			repo:   Repository{},
			sshKey: "/path/to/ssh/key",
			want:   []string{gitBin, "config", "--global", "core.sshCommand", "ssh -i /path/to/ssh/key"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.repo.ConfigSSHCommand(tt.sshKey)
			assert.Equal(t, tt.want, cmd.Args)
		})
	}
}

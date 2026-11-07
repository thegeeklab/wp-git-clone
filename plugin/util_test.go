package plugin

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sshConfPath(dir string) string {
	return filepath.Join(dir, ".ssh", "config.d", "wp-git-clone.conf")
}

func TestWriteSSHConfig(t *testing.T) {
	dir := t.TempDir()

	err := WriteSSHConfig(dir, "/path/to/ssh/key")
	require.NoError(t, err)

	confPath := sshConfPath(dir)
	mainConfig := filepath.Join(dir, ".ssh", "config")

	confContent, err := os.ReadFile(confPath)
	require.NoError(t, err)
	assert.Equal(t, "Host *\n    IdentityFile \"/path/to/ssh/key\"\n    IdentitiesOnly yes\n", string(confContent))

	mainContent, err := os.ReadFile(mainConfig)
	require.NoError(t, err)
	assert.Equal(t, "Include \""+confPath+"\"\n", string(mainContent))

	sshDirInfo, err := os.Stat(filepath.Join(dir, ".ssh"))
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o700), sshDirInfo.Mode().Perm())

	confInfo, err := os.Stat(confPath)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), confInfo.Mode().Perm())
}

func TestWriteSSHConfigWithSpaces(t *testing.T) {
	dir := t.TempDir()

	err := WriteSSHConfig(dir, "/path with space/key")
	require.NoError(t, err)

	content, err := os.ReadFile(sshConfPath(dir))
	require.NoError(t, err)

	assert.Equal(t, "Host *\n    IdentityFile \"/path with space/key\"\n    IdentitiesOnly yes\n", string(content))
}

func TestWriteSSHConfigPreservesExistingConfig(t *testing.T) {
	dir := t.TempDir()
	sshDir := filepath.Join(dir, ".ssh")

	require.NoError(t, os.MkdirAll(sshDir, 0o700))

	existing := "Host *\n    StrictHostKeyChecking no\n"
	require.NoError(t, os.WriteFile(filepath.Join(sshDir, "config"), []byte(existing), 0o600))

	err := WriteSSHConfig(dir, "/path/to/ssh/key")
	require.NoError(t, err)

	includeLine := "Include \"" + sshConfPath(dir) + "\""

	content, err := os.ReadFile(filepath.Join(sshDir, "config"))
	require.NoError(t, err)

	assert.Contains(t, string(content), existing)
	assert.Contains(t, string(content), includeLine)

	// Idempotent: a second run must not duplicate the Include line.
	err = WriteSSHConfig(dir, "/path/to/ssh/key")
	require.NoError(t, err)

	content, err = os.ReadFile(filepath.Join(sshDir, "config"))
	require.NoError(t, err)
	assert.Equal(t, 1, strings.Count(string(content), includeLine))
}

func TestWriteSSHConfigParseableBySSH(t *testing.T) {
	ssh, err := exec.LookPath("ssh")
	if err != nil {
		t.Skip("ssh binary not available")
	}

	home := t.TempDir()
	sshDir := filepath.Join(home, ".ssh")
	require.NoError(t, os.MkdirAll(sshDir, 0o700))

	sshConfig := []byte("Host *\n    StrictHostKeyChecking no\n")
	require.NoError(t, os.WriteFile(filepath.Join(sshDir, "config"), sshConfig, 0o600))

	err = WriteSSHConfig(home, "/tmp/path with space/key")
	require.NoError(t, err)

	cmd := exec.CommandContext(t.Context(), ssh, "-F", filepath.Join(sshDir, "config"), "-G", "example.com")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err)

	assert.Contains(t, string(out), "identityfile /tmp/path with space/key")
}

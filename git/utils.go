package git

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

const (
	netrcFile = `
machine %s
login %s
password %s
`
)

const (
	strictFilePerm = 0o600
)

// WriteNetrc writes the netrc file.
func WriteNetrc(machine, login, password string) error {
	netrcContent := fmt.Sprintf(
		netrcFile,
		machine,
		login,
		password,
	)

	home := "/root"

	if currentUser, err := user.Current(); err == nil {
		home = currentUser.HomeDir
	}

	netpath := filepath.Join(
		home,
		".netrc",
	)

	return os.WriteFile(
		netpath,
		[]byte(netrcContent),
		strictFilePerm,
	)
}

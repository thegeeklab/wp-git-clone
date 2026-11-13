package plugin

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	shaPattern        = regexp.MustCompile(`^[0-9a-fA-F]{7,64}$`)
	identifierPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9/._-]*$`)
	filterPattern     = regexp.MustCompile(`^[A-Za-z0-9:._=-]+$`)
	sshKeyPattern     = regexp.MustCompile(`^/[A-Za-z0-9._ /-]+$`)
)

var (
	errInvalidValue         = errors.New("invalid value")
	errInvalidBranch        = errors.New("invalid branch name")
	errInvalidCommitRef     = errors.New("invalid commit ref")
	errInvalidCommitSha     = errors.New("invalid commit sha")
	errInvalidFilter        = errors.New("invalid clone filter")
	errInvalidSSHKey        = errors.New("invalid ssh key path")
	errInvalidNetrcMachine  = errors.New("invalid netrc machine")
	errInvalidNetrcLogin    = errors.New("invalid netrc login")
	errInvalidNetrcPassword = errors.New("invalid netrc password")
	errEmptySubmoduleName   = errors.New("empty submodule name")
	errInvalidSubmoduleName = errors.New("invalid submodule name")
	errEmptySubmoduleURL    = errors.New("empty submodule url")
)

// validateNoLeadingDash rejects values that begin with a dash, which git would
// otherwise interpret as a command-line option (option injection).
func validateNoLeadingDash(value, label string) error {
	if strings.HasPrefix(value, "-") {
		return fmt.Errorf("%w: %s %q", errInvalidValue, label, value)
	}

	return nil
}

func validateBranch(branch string) error {
	if branch == "" {
		return nil
	}

	if !identifierPattern.MatchString(branch) {
		return fmt.Errorf("%w: %q", errInvalidBranch, branch)
	}

	return nil
}

func validateCommitRef(ref string) error {
	if ref == "" {
		return nil
	}

	if !identifierPattern.MatchString(ref) {
		return fmt.Errorf("%w: %q", errInvalidCommitRef, ref)
	}

	return nil
}

func validateCommitSha(sha string) error {
	if sha == "" {
		return nil
	}

	if !shaPattern.MatchString(sha) {
		return fmt.Errorf("%w: %q", errInvalidCommitSha, sha)
	}

	return nil
}

func validateFilter(filter string) error {
	if filter == "" {
		return nil
	}

	if !filterPattern.MatchString(filter) {
		return fmt.Errorf("%w: %q", errInvalidFilter, filter)
	}

	return nil
}

func validateSSHKey(sshKey string) error {
	if sshKey == "" {
		return nil
	}

	if !sshKeyPattern.MatchString(sshKey) {
		return fmt.Errorf("%w %q: must be an absolute path without shell metacharacters", errInvalidSSHKey, sshKey)
	}

	return nil
}

func validateNetrc(netrc Netrc) error {
	if strings.ContainsAny(netrc.Machine, "\n\r") {
		return fmt.Errorf("%w: %q", errInvalidNetrcMachine, netrc.Machine)
	}

	if strings.ContainsAny(netrc.Login, "\n\r") {
		return fmt.Errorf("%w: %q", errInvalidNetrcLogin, netrc.Login)
	}

	if strings.ContainsAny(netrc.Password, "\n\r") {
		return fmt.Errorf("%w: %q", errInvalidNetrcPassword, netrc.Password)
	}

	return nil
}

func validateRemoteURL(url string) error {
	if url == "" {
		return nil
	}

	return validateNoLeadingDash(url, "remote url")
}

func validateSafeDirectory(dir string) error {
	if dir == "" {
		return nil
	}

	return validateNoLeadingDash(dir, "safe directory")
}

func validateSubmodule(name, url string) error {
	if name == "" {
		return errEmptySubmoduleName
	}

	if strings.ContainsAny(name, "\n\r") {
		return fmt.Errorf("%w: %q", errInvalidSubmoduleName, name)
	}

	if url == "" {
		return fmt.Errorf("%w: %q", errEmptySubmoduleURL, name)
	}

	if err := validateNoLeadingDash(url, fmt.Sprintf("submodule url for %q", name)); err != nil {
		return err
	}

	return nil
}

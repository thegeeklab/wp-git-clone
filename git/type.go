package git

const gitBin = "/usr/bin/git"

type Repository struct {
	RemoteURL        string
	RemoteSSH        string
	Branch           string
	CommitSha        string
	CommitRef        string
	Submodules       map[string]string
	SubmoduleRemote  bool
	SubmodulePartial bool

	InsecureSkipSSLVerify bool
	SafeDirectory         string
	InitExists            bool
}

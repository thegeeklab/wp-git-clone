package git

type Repository struct {
	RemoteURL        string
	RemoteSSH        string
	Branch           string
	CommitSha        string
	CommitRef        string
	Submodules       map[string]string
	SubmoduleRemote  bool
	SubmodulePartial bool

	SafeDirectory string
	WorkDir       string
	IsEmpty       bool
	Filter        string
	Depth         int
}

const gitBin = "/usr/bin/git"

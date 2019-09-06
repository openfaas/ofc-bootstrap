package version

var (
	//GitCommit Git Commit SHA
	GitCommit string
	//Version version of the CLI
	Version string
)

//GetVersion get lastest version
func GetVersion() string {
	if len(Version) == 0 {
		return "dev"
	}
	return Version
}

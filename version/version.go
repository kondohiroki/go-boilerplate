package version

var (
	// GitCommit is the current HEAD set using ldflags.
	GitCommit string

	// Version is the built version.
	Version string = AppSemVer
)

func init() {
	if GitCommit != "" {
		Version += "-" + GitCommit
	}
}

const (
	// AppSemVer is app version.
	AppSemVer = "0.0.1"
)

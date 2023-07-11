package version

var (
	BaseVersion = "0.0.1-dev"
	GitVersion  string
	GoVersion   string
	BuildTime   string
	Version     = BaseVersion + " build on " + BuildTime + " with Git Commit " + GitVersion
)

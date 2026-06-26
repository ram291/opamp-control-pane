package version

// These variables are set at build time via ldflags.
var (
	// Version is the current version of the supervisor.
	Version = "dev"

	// CommitHash is the git commit hash at build time.
	CommitHash = "unknown"

	// BuildDate is the timestamp of the build.
	BuildDate = "unknown"
)

// Info returns a formatted version string.
func Info() string {
	return "opamp-control-pane version=" + Version + " commit=" + CommitHash + " built=" + BuildDate
}
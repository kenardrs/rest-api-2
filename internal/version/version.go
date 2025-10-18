package version

// Application version information
const (
	Version   = "1.1.0"
	BuildDate = "2025-10-18"
)

// Build information (can be overridden during build)
var (
	GitCommit = "dev"
	GoVersion = "1.25.1"
)

// Info returns formatted version information
func Info() string {
	return Version
}

// Full returns complete version details
func Full() map[string]string {
	return map[string]string{
		"version":    Version,
		"build_date": BuildDate,
		"git_commit": GitCommit,
		"go_version": GoVersion,
	}
}

package core


import (
	"fmt"
	"runtime"
)

// The following fields are populated at build time using -ldflags -X.
// Note that DATE is omitted for reproducible builds
var (
	buildVersion     = "unknown"
	buildGitRevision = "unknown"
	buildUser        = "unknown"
	buildHost        = "unknown"
	buildStatus      = "unknown"
	buildTime        = "unknown"
)

// BuildInfo describes version information about the binary build.
type BuildInfo struct {
	Version       string `json:"version"`
	GitRevision   string `json:"revision"`
	User          string `json:"user"`
	Host          string `json:"host"`
	GolangVersion string `json:"golang_version"`
	BuildStatus   string `json:"status"`
	BuildTime     string `json:"time"`
}

var (
	// Info exports the build version information.
	Info BuildInfo
)

// String produces a single-line version info
//
// This looks like:
//
// ```
// user@host-<version>-<git revision>-<build status>
// ```
func (b BuildInfo) String() string {
	return fmt.Sprintf("%v@%v-%v-%v-%v-%v",
		b.User,
		b.Host,
		b.Version,
		b.GitRevision,
		b.BuildStatus,
		b.BuildTime)
}

// LongForm returns a multi-line version information
//
// This looks like:
//
// ```
// Version: <version>
// GitRevision: <git revision>
// User: user@host
// GolangVersion: go1.10.2
// BuildStatus: <build status>
// ```
func (b BuildInfo) LongForm() string {
	return fmt.Sprintf(`Version: %v
GitRevision: %v
User: %v@%v
GolangVersion: %v
BuildStatus: %v
BuildTime: %v
`,
		b.Version,
		b.GitRevision,
		b.User,
		b.Host,
		b.GolangVersion,
		b.BuildStatus,
		b.BuildTime)
}

func init() {
	Info = BuildInfo{
		Version:       buildVersion,
		GitRevision:   buildGitRevision,
		User:          buildUser,
		Host:          buildHost,
		GolangVersion: runtime.Version(),
		BuildStatus:   buildStatus,
		BuildTime:     buildTime,
	}
}

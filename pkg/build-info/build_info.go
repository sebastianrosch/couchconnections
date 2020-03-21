package buildinfo

import (
	"fmt"
	"runtime"
)

// Build information. Populated at build-time.
var (
	Version   string
	Revision  string
	Branch    string
	BuildUser string
	BuildDate string
	GoVersion = runtime.Version()
)

// BuildInfo contains information added at build time such as version
type BuildInfo struct {
	Version   string
	Revision  string
	Branch    string
	BuildUser string
	BuildDate string
	GoVersion string
}

// NewDefaultBuildInfo returns the build information obtained from ldflags
func NewDefaultBuildInfo() BuildInfo {
	return NewBuildInfo(Version, Branch, BuildDate, BuildUser, GoVersion, Revision)
}

// NewBuildInfo returns an object with the build information
func NewBuildInfo(version, branch, buildDate, buildUser, goVersion, revision string) BuildInfo {
	return BuildInfo{
		Version:   version,
		Branch:    branch,
		BuildDate: buildDate,
		BuildUser: buildUser,
		GoVersion: goVersion,
		Revision:  revision,
	}
}

// Info returns version, branch and revision information.
func (bi *BuildInfo) Info() string {
	return fmt.Sprintf("(version=%s, branch=%s, revision=%s)", bi.Version, bi.Branch, bi.Revision)
}

// BuildContext returns goVersion, buildUser and buildDate information.
func (bi *BuildInfo) BuildContext() string {
	return fmt.Sprintf("(go=%s, user=%s, date=%s)", bi.GoVersion, bi.BuildUser, bi.BuildDate)
}

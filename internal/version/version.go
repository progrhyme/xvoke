// Package version tells the version of the running command.
package version

import (
	"runtime/debug"
	"strings"
)

const devVersion = "dev"

// Get returns the module version recorded at build time, without the leading "v".
// It returns "dev" when the build records no module version, as with "go run" or "-buildvcs=false".
// A build from a working tree records a pseudo-version.
func Get() string {
	info, ok := debug.ReadBuildInfo()
	return fromBuildInfo(info, ok)
}

func fromBuildInfo(info *debug.BuildInfo, ok bool) string {
	if !ok || info.Main.Version == "" || info.Main.Version == "(devel)" {
		return devVersion
	}
	return strings.TrimPrefix(info.Main.Version, "v")
}

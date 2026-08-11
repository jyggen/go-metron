// Package version reports go-metron's own module version, as seen by an
// importing binary's build info.
package version

import (
	"runtime/debug"
	"strings"
)

// modulePath locates this module's own entry in a consumer's build info.
const modulePath = "github.com/jyggen/go-metron"

// Version reports this module's released version, with the "v" prefix
// trimmed (a Go module convention, not an HTTP one), or "devel" when there is
// no released version to read: non-module builds, filesystem replaces, and
// this module's own test binaries (golang/go#68045).
func Version() string {
	bi, ok := debug.ReadBuildInfo()
	if ok && bi.Main.Path != modulePath {
		for _, dep := range bi.Deps {
			if dep.Path != modulePath {
				continue
			}

			v := dep.Version
			if dep.Replace != nil {
				v = dep.Replace.Version
			}

			// A filesystem replace reports "(devel)" rather than a version.
			if v != "" && v != "(devel)" {
				return strings.TrimPrefix(v, "v")
			}

			break
		}
	}

	return "devel"
}

package metron

import "runtime/debug"

// modulePath locates this module's own entry in a consumer's build info.
const modulePath = "github.com/jyggen/go-metron"

// defaultUserAgent tracks the released version, so it needs no manual bump.
var defaultUserAgent = "go-metron/" + moduleVersion()

// moduleVersion reports this module's version from the importing binary, or
// "devel" when there is no released version to read: non-module builds,
// filesystem replaces, and our own test binaries (golang/go#68045).
func moduleVersion() string {
	bi, ok := debug.ReadBuildInfo()
	if ok && bi.Main.Path != modulePath {
		for _, dep := range bi.Deps {
			if dep.Path != modulePath {
				continue
			}

			version := dep.Version
			if dep.Replace != nil {
				version = dep.Replace.Version
			}

			// A filesystem replace reports "(devel)" rather than a version.
			if version != "" && version != "(devel)" {
				return version
			}

			break
		}
	}

	return "devel"
}

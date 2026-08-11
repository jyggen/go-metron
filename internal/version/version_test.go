package version_test

import (
	"testing"

	"github.com/jyggen/go-metron/internal/version"
	"github.com/stretchr/testify/require"
)

// TestVersion pins the "devel" fallback: this module's own test binary has
// empty Deps (golang/go#68045), so it can never observe a released version.
func TestVersion(t *testing.T) {
	t.Parallel()

	require.Equal(t, "devel", version.Version())
}

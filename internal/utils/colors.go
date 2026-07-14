// SHARED UTILITY: This package is shared between the CLI interface (internal/cli)
// and the Web server (internal/web). Modifying public APIs will impact both contexts.

package utils

// Color constants matching the original shell scripts and Makefile styles.
const (
	Red    = "\033[0;31m"
	Green  = "\033[0;32m"
	Yellow = "\033[0;33m"
	Cyan   = "\033[0;36m"
	Reset  = "\033[0m"
)

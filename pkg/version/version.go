package version

import (
	"fmt"
	"runtime"
)

// Build information. Populated at build-time.
var (
	Version   string
	Revision  string
	BuildDate string
	GoVersion = runtime.Version()
)

// Info returns version, branch and revision information.
func Info() string {
	return fmt.Sprintf("am-kafka-project, version: %s, revision: %s, build date: %s, %s", Version, Revision, BuildDate, GoVersion)
}

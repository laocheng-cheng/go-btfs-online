package env

import (
	"os"
)

// In the env package to resolve circular dependencies
// Since config package needs to log and log package needs to config
var (
	LogFile = "status_server.log"
)

func init() {
	if lf := os.Getenv("SS_LOG_FILE"); lf != "" {
		LogFile = lf
	}
}

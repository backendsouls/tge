package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	devLog    *log.Logger
	systemTpl string = "[Ding! %s]"
)

func init() {
	devLog = log.New(os.Stderr, "[DEV] ", log.LstdFlags)
}

// SetSystemTemplate overrides the default system log format string.
// It should contain exactly one %s for the message.
func SetSystemTemplate(template string) {
	if template != "" {
		systemTpl = template
	}
}

// Dev logs information intended for developers.
func Dev(format string, args ...any) {
	devLog.Printf(format, args...)
}

// System logs events intended for the "novel author" view.
func System(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	// Print directly to stdout for the author to see
	fmt.Printf(systemTpl+"\n", msg)
}

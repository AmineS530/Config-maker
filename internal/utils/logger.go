package utils

import (
	"fmt"
	"io"
)

// Logger provides stylized output methods for CLI terminal formatting.
type Logger struct {
	Out io.Writer
}

// Info prints a cyan info message.
func (l *Logger) Info(format string, a ...interface{}) {
	fmt.Fprintf(l.Out, "%s[INFO] %s%s\n", Cyan, fmt.Sprintf(format, a...), Reset)
}

// Success prints a green success message.
func (l *Logger) Success(format string, a ...interface{}) {
	fmt.Fprintf(l.Out, "%s[SUCCESS] %s%s\n", Green, fmt.Sprintf(format, a...), Reset)
}

// Warning prints a yellow warning message.
func (l *Logger) Warning(format string, a ...interface{}) {
	fmt.Fprintf(l.Out, "%s[WARNING] %s%s\n", Yellow, fmt.Sprintf(format, a...), Reset)
}

// Error prints a red error message.
func (l *Logger) Error(format string, a ...interface{}) {
	fmt.Fprintf(l.Out, "%s[ERROR] %s%s\n", Red, fmt.Sprintf(format, a...), Reset)
}

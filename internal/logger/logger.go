package logger

import "fmt"

const (
	reset      = "\033[0m"
	bold       = "\033[1m"
	red        = "\033[31m"
	green      = "\033[32m"
	yellow     = "\033[33m"
	blue       = "\033[34m"
	magenta    = "\033[35m"
	grey       = "\033[37m"
	cyan       = "\033[96m"
	lightGreen = "\033[92m"
)

type Logger struct{}

func New() *Logger {
	return &Logger{}
}

// cyan for Information message
func (l *Logger) Info(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s%s==> %s%s\n", cyan, bold, msg, reset)
}

// yellow for Warning
func (l *Logger) Warn(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s%s\n", yellow, msg, reset)
}

// red for error
func (l *Logger) Error(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s%s\n", red, msg, reset)
}

// greeen for messages
func (l *Logger) Message(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s%s\n", green, msg, reset)
}

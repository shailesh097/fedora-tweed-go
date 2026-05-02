package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	cyan   = "\033[96m"
)

type Logger struct {
	termOut io.Writer // colored op to terminal
	termErr io.Writer // colored err t terminal stderr
	fileOut io.Writer // normal op to log file
}

func New() (*Logger, func(), error) {
	logFile, err := createLogFile()
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to create log file: %w", err)
	}

	closefile := func() { logFile.Close() }

	return &Logger{
		termOut: os.Stdout,
		termErr: os.Stderr,
		fileOut: logFile,
	}, closefile, nil
}

func createLogFile() (*os.File, error) {
	//get home dir
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(home, "personal_logfiles", "fedora_tweed_go")

	// create dirs
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	timestamp := time.Now().Format("2006-01-02_15_04_05")
	filename := fmt.Sprintf("setup_%s.log", timestamp)
	fullPath := filepath.Join(dir, filename)
	return os.Create(fullPath) // create file in personal log dir
}

func (l *Logger) writeBoth(termWriter io.Writer, color, prefix, msg string) {
	fmt.Fprintf(termWriter, "%s%s%s %s%s\n", color, bold, prefix, msg, reset) // terminal
	fmt.Fprintf(l.fileOut, "%s %s\n", prefix, msg)                            // logfile
}

// cyan for Information message
func (l *Logger) Info(format string, args ...any) {
	l.writeBoth(l.termOut, cyan, "==>", fmt.Sprintf(format, args...))
}

// yellow for Warning
func (l *Logger) Warn(format string, args ...any) {
	l.writeBoth(l.termOut, yellow, "", fmt.Sprintf(format, args...))
}

// red for error
func (l *Logger) Error(format string, args ...any) {
	l.writeBoth(l.termErr, red, "", fmt.Sprintf(format, args...))
}

// greeen for messages
func (l *Logger) Message(format string, args ...any) {
	l.writeBoth(l.termOut, green, "", fmt.Sprintf(format, args...))
}

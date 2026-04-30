package logger

import "os"

func errWriter() *os.File {
	return os.Stderr
}

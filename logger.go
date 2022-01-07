package p2p

import (
	"fmt"
	"os"
)

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type stdLogger struct{}

func NewStdLogger() (l *stdLogger) {
	return &stdLogger{}
}

func (l *stdLogger) Info(msg string) {
	_, _ = fmt.Fprintln(os.Stderr, msg)
}

func (l *stdLogger) Warn(msg string) {
	_, _ = fmt.Fprintln(os.Stderr, msg)
}

func (l *stdLogger) Error(msg string) {
	_, _ = fmt.Fprintln(os.Stderr, msg)
}

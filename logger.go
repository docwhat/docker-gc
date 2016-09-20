package main

import (
	"fmt"
	"os"
)

// Logger is a logger that understands the application's configuration.
type Logger struct {
	config appConfig
}

// NewLogger initializes and returns a new Logger instance.
func NewLogger(config appConfig) Logger {
	return Logger{config: config}
}

// Info emits an informational message.
func (l Logger) Info(msg string, a ...interface{}) {
	if !l.config.Quiet {
		fmt.Printf(fmt.Sprintf("%s\n", msg), a...)
	}
}

// Fatal emits a message and aborts.
func (l Logger) Fatal(msg string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, fmt.Sprintf("FATAL: %s\n", msg), a...)
	os.Exit(1)
}

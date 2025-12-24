package logger

import (
	"fmt"
	"log"
	"os"
)

// Logger is a simple structured logger
type Logger struct {
	debug   bool
	env     string
	infoLog *log.Logger
	errorLog *log.Logger
	debugLog *log.Logger
}

// New creates a new logger
func New(debug bool, env string) *Logger {
	return &Logger{
		debug:   debug,
		env:     env,
		infoLog: log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lmsgprefix),
		errorLog: log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lmsgprefix),
		debugLog: log.New(os.Stdout, "[DEBUG] ", log.LstdFlags|log.Lmsgprefix),
	}
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.infoLog.Printf(format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.errorLog.Printf(format, args...)
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.debug {
		l.debugLog.Printf(format, args...)
	}
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.errorLog.Fatalf(format, args...)
}

// WithFields logs with fields (simplified)
func (l *Logger) WithFields(fields map[string]interface{}) *FieldLogger {
	return &FieldLogger{
		logger: l,
		fields: fields,
	}
}

// FieldLogger is a logger with predefined fields
type FieldLogger struct {
	logger *Logger
	fields map[string]interface{}
}

// Info logs an info message with fields
func (fl *FieldLogger) Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fl.logger.Info("%s %v", msg, fl.fields)
}

// Error logs an error message with fields
func (fl *FieldLogger) Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fl.logger.Error("%s %v", msg, fl.fields)
}

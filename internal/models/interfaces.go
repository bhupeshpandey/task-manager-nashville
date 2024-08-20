package models

type LogLevel string

const (
	ErrorLevel LogLevel = "error"
	WarnLevel  LogLevel = "warn"
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
)

type Logger interface {
	Log(logLevel LogLevel, message string, log ...interface{})
}

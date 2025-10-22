package logging

import "sync"

type LogLevel int

const (
	LevelNone LogLevel = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
)

type Logger interface {
	SetLevel(level LogLevel)
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

var (
	mu           sync.RWMutex
	globalLogger Logger = &noopLogger{}
)

func SetGlobalLogger(l Logger) {
	mu.Lock()
	defer mu.Unlock()
	globalLogger = l
}

func SetLevel(level LogLevel) {
	mu.RLock()
	defer mu.RUnlock()
	globalLogger.SetLevel(level)
}

func Debug(msg string, args ...interface{}) { globalLogger.Debug(msg, args...) }
func Info(msg string, args ...interface{})  { globalLogger.Info(msg, args...) }
func Warn(msg string, args ...interface{})  { globalLogger.Warn(msg, args...) }
func Error(msg string, args ...interface{}) { globalLogger.Error(msg, args...) }

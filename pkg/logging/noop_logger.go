//go:build !logging
// +build !logging

package logging

type noopLogger struct{}

func init() {
	SetGlobalLogger(&noopLogger{})
}

func (n *noopLogger) SetLevel(level LogLevel)               {}
func (n *noopLogger) Debug(msg string, args ...interface{}) {}
func (n *noopLogger) Info(msg string, args ...interface{})  {}
func (n *noopLogger) Warn(msg string, args ...interface{})  {}
func (n *noopLogger) Error(msg string, args ...interface{}) {}

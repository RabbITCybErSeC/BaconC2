//go:build logging
// +build logging

package logging

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type stdLogger struct {
	mu    sync.RWMutex
	level LogLevel
	l     *log.Logger
}

func init() {
	SetGlobalLogger(NewStdLogger())
}

func NewStdLogger() *stdLogger {
	return &stdLogger{
		level: LevelInfo,
		l:     log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (s *stdLogger) SetLevel(level LogLevel) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.level = level
}

func (s *stdLogger) shouldLog(lvl LogLevel) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return lvl <= s.level
}

func (s *stdLogger) Debug(msg string, args ...interface{}) {
	if s.shouldLog(LevelDebug) {
		s.l.Printf("[DEBUG] " + fmt.Sprintf(msg, args...))
	}
}

func (s *stdLogger) Info(msg string, args ...interface{}) {
	if s.shouldLog(LevelInfo) {
		s.l.Printf("[INFO] " + fmt.Sprintf(msg, args...))
	}
}

func (s *stdLogger) Warn(msg string, args ...interface{}) {
	if s.shouldLog(LevelWarn) {
		s.l.Printf("[WARN] " + fmt.Sprintf(msg, args...))
	}
}

func (s *stdLogger) Error(msg string, args ...interface{}) {
	if s.shouldLog(LevelError) {
		s.l.Printf("[ERROR] " + fmt.Sprintf(msg, args...))
	}
}

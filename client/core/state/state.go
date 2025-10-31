package state

import (
	"os"
	"sync"
)

type AgentState struct {
	mu               sync.RWMutex
	workingDirectory string
	environment      map[string]string
}

func NewAgentState() *AgentState {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "/"
	}

	return &AgentState{
		workingDirectory: cwd,
		environment:      make(map[string]string),
	}
}

func (s *AgentState) GetWorkingDirectory() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.workingDirectory
}

func (s *AgentState) SetWorkingDirectory(path string) error {
	// Validate the directory exists
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return os.ErrInvalid
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.workingDirectory = path
	return nil
}

func (s *AgentState) GetEnv(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.environment[key]
	return val, ok
}

func (s *AgentState) SetEnv(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.environment[key] = value
}

func (s *AgentState) GetAllEnv() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	envCopy := make(map[string]string, len(s.environment))
	for k, v := range s.environment {
		envCopy[k] = v
	}
	return envCopy
}

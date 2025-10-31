package filesystem

import (
	"os"
	"path/filepath"
	"testing"

	command_handler "github.com/RabbITCybErSeC/BaconC2/pkg/commands/handlers"
	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
)

type mockAgentState struct {
	workingDir  string
	environment map[string]string
}

func newMockAgentState() *mockAgentState {
	cwd, _ := os.Getwd()
	return &mockAgentState{
		workingDir:  cwd,
		environment: make(map[string]string),
	}
}

func (m *mockAgentState) GetWorkingDirectory() string {
	return m.workingDir
}

func (m *mockAgentState) SetWorkingDirectory(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return os.ErrInvalid
	}
	m.workingDir = path
	return nil
}

func (m *mockAgentState) GetEnv(key string) (string, bool) {
	val, ok := m.environment[key]
	return val, ok
}

func (m *mockAgentState) SetEnv(key, value string) {
	m.environment[key] = value
}

func (m *mockAgentState) GetAllEnv() map[string]string {
	return m.environment
}

func TestPwdHandler(t *testing.T) {
	state := newMockAgentState()
	cmd := models.Command{
		ID:      "test-1",
		Command: "pwd",
		Type:    models.CommandTypeInternal,
	}

	ctx := command_handler.NewCommandContext(cmd, state)
	result := PwdHandler(ctx)

	if result.Status != models.CommandStatusCompleted {
		t.Errorf("Expected status completed, got %s", result.Status)
	}

	if result.Output == "" {
		t.Error("Expected output to contain working directory")
	}
}

func TestCdHandler(t *testing.T) {
	state := newMockAgentState()

	tempDir := os.TempDir()

	cmd := models.Command{
		ID:      "test-2",
		Command: "cd",
		Args:    []string{tempDir},
		Type:    models.CommandTypeInternal,
	}

	ctx := command_handler.NewCommandContext(cmd, state)
	result := CdHandler(ctx)

	if result.Status != models.CommandStatusCompleted {
		t.Errorf("Expected status completed, got %s: %s", result.Status, result.Output)
	}

	if state.GetWorkingDirectory() != filepath.Clean(tempDir) {
		t.Errorf("Expected working directory to be %s, got %s", tempDir, state.GetWorkingDirectory())
	}
}

func TestCdHandlerInvalidDirectory(t *testing.T) {
	state := newMockAgentState()

	cmd := models.Command{
		ID:      "test-3",
		Command: "cd",
		Args:    []string{"/nonexistent/directory/path"},
		Type:    models.CommandTypeInternal,
	}

	ctx := command_handler.NewCommandContext(cmd, state)
	result := CdHandler(ctx)

	if result.Status != models.CommandStatusFailed {
		t.Errorf("Expected status failed for invalid directory, got %s", result.Status)
	}
}

func TestCdHandlerRelativePath(t *testing.T) {
	state := newMockAgentState()

	cwd, _ := os.Getwd()
	state.SetWorkingDirectory(cwd)

	cmd := models.Command{
		ID:      "test-4",
		Command: "cd",
		Args:    []string{".."},
		Type:    models.CommandTypeInternal,
	}

	ctx := command_handler.NewCommandContext(cmd, state)
	result := CdHandler(ctx)

	if result.Status != models.CommandStatusCompleted {
		t.Errorf("Expected status completed, got %s: %s", result.Status, result.Output)
	}

	expectedDir := filepath.Clean(filepath.Join(cwd, ".."))
	if state.GetWorkingDirectory() != expectedDir {
		t.Errorf("Expected working directory to be %s, got %s", expectedDir, state.GetWorkingDirectory())
	}
}

func TestLsHandler(t *testing.T) {
	state := newMockAgentState()

	tempDir := os.TempDir()
	state.SetWorkingDirectory(tempDir)

	cmd := models.Command{
		ID:      "test-5",
		Command: "ls",
		Type:    models.CommandTypeInternal,
	}

	ctx := command_handler.NewCommandContext(cmd, state)
	result := LsHandler(ctx)

	if result.Status != models.CommandStatusCompleted {
		t.Errorf("Expected status completed, got %s: %s", result.Status, result.Output)
	}

	if result.Output == "" {
		t.Error("Expected output to contain directory listing")
	}
}

func TestLsHandlerWithPath(t *testing.T) {
	state := newMockAgentState()

	tempDir := os.TempDir()

	cmd := models.Command{
		ID:      "test-6",
		Command: "ls",
		Args:    []string{tempDir},
		Type:    models.CommandTypeInternal,
	}

	ctx := command_handler.NewCommandContext(cmd, state)
	result := LsHandler(ctx)

	if result.Status != models.CommandStatusCompleted {
		t.Errorf("Expected status completed, got %s: %s", result.Status, result.Output)
	}
}

func TestLsHandlerInvalidPath(t *testing.T) {
	state := newMockAgentState()

	cmd := models.Command{
		ID:      "test-7",
		Command: "ls",
		Args:    []string{"/nonexistent/path"},
		Type:    models.CommandTypeInternal,
	}

	ctx := command_handler.NewCommandContext(cmd, state)
	result := LsHandler(ctx)

	if result.Status != models.CommandStatusFailed {
		t.Errorf("Expected status failed for invalid path, got %s", result.Status)
	}
}

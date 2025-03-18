package store

import (
	"fmt"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/server/db"
	"github.com/RabbITCybErSeC/BaconC2/server/models"
	"gorm.io/gorm"
)

type AgentStoreInterface interface {
	Register(agent *models.Agent) error
	Get(id string) (*models.Agent, error)
	List() ([]*models.Agent, error)
	UpdateLastSeen(id string) error
	UpdateAgentCommands(id string, cmd models.Command) error
}

var ErrAgentNotFound = fmt.Errorf("agent not found")

type AgentStore struct {
	db db.AgentRepositoryInterface
}

func NewAgentStore(repo db.AgentRepositoryInterface) *AgentStore {
	return &AgentStore{
		db: repo,
	}
}

func (s *AgentStore) Register(agent *models.Agent) error {
	agent.LastSeen = time.Now()
	agent.IsActive = true
	agent.Commands = []models.Command{}
	return s.db.Save(agent)
}

func (s *AgentStore) Get(id string) (*models.Agent, error) {
	agent, err := s.db.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrAgentNotFound
		}
		return nil, err
	}
	return agent, nil
}

func (s *AgentStore) List() ([]*models.Agent, error) {
	agents, err := s.db.GetAll()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	agentPtrs := make([]*models.Agent, len(agents))

	for i := range agents {
		if now.Sub(agents[i].LastSeen) > 5*time.Minute {
			agents[i].IsActive = false
		}
		agentPtrs[i] = &agents[i]
	}

	return agentPtrs, nil
}

func (s *AgentStore) UpdateLastSeen(id string) error {
	agent, err := s.Get(id)
	if err != nil {
		return err
	}

	agent.LastSeen = time.Now()
	agent.IsActive = true
	return s.db.Save(agent)
}

func (s *AgentStore) UpdateAgentCommands(id string, cmd models.Command) error {
	agent, err := s.Get(id)
	if err != nil {
		return err
	}

	agent.CommandMu.Lock()
	defer agent.CommandMu.Unlock()

	found := false
	for i, existingCmd := range agent.Commands {
		if existingCmd.ID == cmd.ID {
			agent.Commands[i] = cmd
			found = true
			break
		}
	}

	if !found {
		agent.Commands = append(agent.Commands, cmd)
	}

	return s.db.Save(agent)
}

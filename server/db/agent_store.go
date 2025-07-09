package db

import (
	"fmt"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
	local_models "github.com/RabbITCybErSeC/BaconC2/server/models"
	"gorm.io/gorm"
)

type IAgentRepository interface {
	Save(agent *local_models.ServerAgentModel) error
	Get(id string) (*local_models.ServerAgentModel, error)
	GetAll() ([]local_models.ServerAgentModel, error)
	Delete(id string) error

	GetWithExtendedInfo(id string) (*local_models.ServerAgentModel, error)
	GetWithAllRelations(id string) (*local_models.ServerAgentModel, error)
	GetActiveAgents() ([]local_models.ServerAgentModel, error)

	UpdateExtendedInfo(agentID string, info *models.ExtendedAgentInfo) error
	GetExtendedInfo(agentID string) (*models.ExtendedAgentInfo, error)

	CreateSession(session *local_models.AgentSession) error
	EndSession(sessionID string) error
	GetActiveSessions() ([]local_models.AgentSession, error)
	GetAgentSessions(agentID string) ([]local_models.AgentSession, error)

	AddCommand(command *local_models.AgentCommand) error
	GetPendingCommands(agentID string) ([]local_models.AgentCommand, error)
	GetCommands(agentID string, limit int) ([]local_models.AgentCommand, error)
	UpdateCommandStatus(commandID string, status string) error
}

type AgentRepository struct {
	db *gorm.DB
}

func NewAgentRepository(db *gorm.DB) *AgentRepository {
	return &AgentRepository{db: db}
}

// AutoMigrate creates all tables
func (s *AgentRepository) AutoMigrate() error {
	return s.db.AutoMigrate(
		&local_models.ServerAgentModel{},
		&local_models.AgentCommand{},
		&local_models.AgentSession{},
	)
}

// Basic CRUD operations
func (s *AgentRepository) Save(agent *local_models.ServerAgentModel) error {
	return s.db.Save(agent).Error
}

func (s *AgentRepository) Get(id string) (*local_models.ServerAgentModel, error) {
	var agent local_models.ServerAgentModel
	err := s.db.First(&agent, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func (s *AgentRepository) GetAll() ([]local_models.ServerAgentModel, error) {
	var agents []local_models.ServerAgentModel
	err := s.db.Find(&agents).Error
	if err != nil {
		return nil, err
	}
	return agents, nil
}

func (s *AgentRepository) Delete(id string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Delete related records first
		if err := tx.Where("agent_id = ?", id).Delete(&local_models.AgentCommand{}).Error; err != nil {
			return err
		}
		if err := tx.Where("agent_id = ?", id).Delete(&local_models.AgentSession{}).Error; err != nil {
			return err
		}
		// Delete the agent
		return tx.Delete(&local_models.ServerAgentModel{}, "id = ?", id).Error
	})
}

// Extended operations with relationships
func (s *AgentRepository) GetWithExtendedInfo(id string) (*local_models.ServerAgentModel, error) {
	var agent local_models.ServerAgentModel
	err := s.db.Preload("ExtendedInfo").First(&agent, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func (s *AgentRepository) GetWithAllRelations(id string) (*local_models.ServerAgentModel, error) {
	var agent local_models.ServerAgentModel
	err := s.db.
		Preload("Commands").
		Preload("ExtendedInfo").
		Preload("Sessions").
		First(&agent, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func (s *AgentRepository) GetActiveAgents() ([]local_models.ServerAgentModel, error) {
	var agents []local_models.ServerAgentModel
	err := s.db.Where("is_active = ?", true).
		Preload("ExtendedInfo").
		Find(&agents).Error
	if err != nil {
		return nil, err
	}
	return agents, nil
}

// Extended info operations
func (s *AgentRepository) UpdateExtendedInfo(agentID string, info *models.ExtendedAgentInfo) error {
	info.AgentID = agentID
	return s.db.Save(info).Error
}

func (s *AgentRepository) GetExtendedInfo(agentID string) (*models.ExtendedAgentInfo, error) {
	var info models.ExtendedAgentInfo
	err := s.db.Where("agent_id = ?", agentID).First(&info).Error
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// Session operations
func (s *AgentRepository) CreateSession(session *local_models.AgentSession) error {
	return s.db.Create(session).Error
}

func (s *AgentRepository) EndSession(sessionID string) error {
	now := time.Now()
	return s.db.Model(&local_models.AgentSession{}).
		Where("session_id = ?", sessionID).
		Updates(map[string]interface{}{
			"end_time":  &now,
			"is_active": false,
		}).Error
}

func (s *AgentRepository) GetActiveSessions() ([]local_models.AgentSession, error) {
	var sessions []local_models.AgentSession
	err := s.db.Where("is_active = ?", true).Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (s *AgentRepository) GetAgentSessions(agentID string) ([]local_models.AgentSession, error) {
	var sessions []local_models.AgentSession
	err := s.db.Where("agent_id = ?", agentID).
		Order("created_at DESC").
		Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

// Command operations
func (s *AgentRepository) AddCommand(command *local_models.AgentCommand) error {
	return s.db.Create(command).Error
}

func (s *AgentRepository) GetPendingCommands(agentID string) ([]local_models.AgentCommand, error) {
	var commands []local_models.AgentCommand
	err := s.db.Where("agent_id = ? AND status = ?", agentID, "pending").
		Order("created_at ASC").
		Find(&commands).Error
	if err != nil {
		return nil, err
	}
	return commands, nil
}

func (s *AgentRepository) GetCommands(agentID string, limit int) ([]local_models.AgentCommand, error) {
	var commands []local_models.AgentCommand
	query := s.db.Where("agent_id = ?", agentID).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&commands).Error
	if err != nil {
		return nil, err
	}
	return commands, nil
}

func (s *AgentRepository) UpdateCommandStatus(commandID string, status string) error {
	return s.db.Model(&local_models.AgentCommand{}).
		Where("id = ?", commandID).
		Update("status", status).Error
}

// Utility methods
func (s *AgentRepository) GetAgentStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Count total agents
	var totalAgents int64
	if err := s.db.Model(&local_models.ServerAgentModel{}).Count(&totalAgents).Error; err != nil {
		return nil, err
	}
	stats["total_agents"] = totalAgents

	// Count active agents
	var activeAgents int64
	if err := s.db.Model(&local_models.ServerAgentModel{}).Where("is_active = ?", true).Count(&activeAgents).Error; err != nil {
		return nil, err
	}
	stats["active_agents"] = activeAgents

	// Count active sessions
	var activeSessions int64
	if err := s.db.Model(&local_models.AgentSession{}).Where("is_active = ?", true).Count(&activeSessions).Error; err != nil {
		return nil, err
	}
	stats["active_sessions"] = activeSessions

	// Count total commands
	var totalCommands int64
	if err := s.db.Model(&local_models.AgentCommand{}).Count(&totalCommands).Error; err != nil {
		return nil, err
	}
	stats["total_commands"] = totalCommands

	return stats, nil
}

// Search and filter methods
func (s *AgentRepository) SearchAgents(query string, limit int) ([]local_models.ServerAgentModel, error) {
	var agents []local_models.ServerAgentModel

	searchQuery := fmt.Sprintf("%%%s%%", query)
	err := s.db.Where("id LIKE ? OR hostname LIKE ?", searchQuery, searchQuery).
		Limit(limit).
		Preload("ExtendedInfo").
		Find(&agents).Error
	if err != nil {
		return nil, err
	}
	return agents, nil
}

// Get agents that haven't been seen recently
func (s *AgentRepository) GetStaleAgents(threshold time.Duration) ([]local_models.ServerAgentModel, error) {
	var agents []local_models.ServerAgentModel
	cutoff := time.Now().Add(-threshold)

	err := s.db.Where("last_seen < ? AND is_active = ?", cutoff, true).
		Find(&agents).Error
	if err != nil {
		return nil, err
	}
	return agents, nil
}

// Update agent last seen timestamp
func (s *AgentRepository) UpdateLastSeen(agentID string) error {
	return s.db.Model(&local_models.ServerAgentModel{}).
		Where("id = ?", agentID).
		Update("last_seen", time.Now()).Error
}

// Get recent commands for an agent
func (s *AgentRepository) GetRecentCommands(agentID string, since time.Time) ([]local_models.AgentCommand, error) {
	var commands []local_models.AgentCommand
	err := s.db.Where("agent_id = ? AND created_at > ?", agentID, since).
		Order("created_at DESC").
		Find(&commands).Error
	if err != nil {
		return nil, err
	}
	return commands, nil
}

// Batch operations
func (s *AgentRepository) MarkAgentsInactive(agentIDs []string) error {
	return s.db.Model(&local_models.ServerAgentModel{}).
		Where("id IN ?", agentIDs).
		Updates(map[string]any{
			"is_active": false,
			"last_seen": time.Now(),
		}).Error
}

// Get active session for an agent
func (s *AgentRepository) GetActiveSessionForAgent(agentID string) (*local_models.AgentSession, error) {
	var session local_models.AgentSession
	err := s.db.Where("agent_id = ? AND is_active = ?", agentID, true).
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

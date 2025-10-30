package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
	local_models "github.com/RabbITCybErSeC/BaconC2/server/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrNotFound = errors.New("record not found")

type IAgentRepository interface {
	SaveAgent(ctx context.Context, agent *local_models.ServerAgentModel) error
	GetAgent(ctx context.Context, id string) (*local_models.ServerAgentModel, error)
	GetAllAgents(ctx context.Context) ([]local_models.ServerAgentModel, error)
	DeleteAgent(ctx context.Context, id string) error

	GetWithExtendedInfo(ctx context.Context, id string) (*local_models.ServerAgentModel, error)
	GetWithAllRelations(ctx context.Context, id string) (*local_models.ServerAgentModel, error)
	GetActiveAgents(ctx context.Context) ([]local_models.ServerAgentModel, error)

	UpdateExtendedInfo(ctx context.Context, agentID string, info *models.ExtendedAgentInfo) error
	UpdateLastSeen(ctx context.Context, agentID string) error

	GetExtendedInfo(ctx context.Context, agentID string) (*models.ExtendedAgentInfo, error)

	CreateSession(ctx context.Context, session *local_models.AgentSession) error
	EndSession(ctx context.Context, sessionID string) error
	GetActiveSessions(ctx context.Context) ([]local_models.AgentSession, error)
	GetAgentSessions(ctx context.Context, agentID string) ([]local_models.AgentSession, error)

	SaveCommand(ctx context.Context, command *local_models.AgentCommand) error
	GetCommandsByStatus(ctx context.Context, agentID string, status models.CommandStatus) ([]local_models.AgentCommand, error)
	GetCommands(ctx context.Context, agentID string, limit int) ([]local_models.AgentCommand, error)
	UpdateCommandStatus(ctx context.Context, commandID string, status models.CommandStatus) error

	SaveCommandResult(ctx context.Context, result *local_models.AgentCommandResult) error
	UpdateCommandStatusWithResult(ctx context.Context, agentID, commandID string, status models.CommandStatus, output any) error
	GetCommandResult(ctx context.Context, commandID string) (*models.CommandResult, error)
}

type AgentRepository struct {
	db *gorm.DB
}

func NewAgentRepository(db *gorm.DB) *AgentRepository {
	return &AgentRepository{db: db}
}

func (s *AgentRepository) SaveAgent(ctx context.Context, agent *local_models.ServerAgentModel) error {
	if err := s.db.WithContext(ctx).Save(agent).Error; err != nil {
		return fmt.Errorf("could not save agent %s: %w", agent.ID, err)
	}
	return nil
}

func (s *AgentRepository) GetAgent(ctx context.Context, id string) (*local_models.ServerAgentModel, error) {
	var agent local_models.ServerAgentModel
	err := s.db.WithContext(ctx).First(&agent, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("could not get agent %s: %w", id, err)
	}
	return &agent, nil
}

func (s *AgentRepository) GetAllAgents(ctx context.Context) ([]local_models.ServerAgentModel, error) {
	var agents []local_models.ServerAgentModel
	if err := s.db.WithContext(ctx).Find(&agents).Error; err != nil {
		return nil, fmt.Errorf("could not get all agents: %w", err)
	}
	return agents, nil
}

func (s *AgentRepository) DeleteAgent(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Note: Depending on DB engine and schema, foreign key constraints with cascading deletes
		// might be a cleaner solution. This manual approach is DB-agnostic.
		if err := tx.Where("agent_id = ?", id).Delete(&local_models.AgentCommand{}).Error; err != nil {
			return fmt.Errorf("could not delete commands for agent %s: %w", id, err)
		}
		if err := tx.Where("agent_id = ?", id).Delete(&local_models.AgentSession{}).Error; err != nil {
			return fmt.Errorf("could not delete sessions for agent %s: %w", id, err)
		}
		if err := tx.Delete(&local_models.ServerAgentModel{}, "id = ?", id).Error; err != nil {
			return fmt.Errorf("could not delete agent %s: %w", id, err)
		}
		return nil
	})
}

func (s *AgentRepository) GetWithExtendedInfo(ctx context.Context, id string) (*local_models.ServerAgentModel, error) {
	var agent local_models.ServerAgentModel
	err := s.db.WithContext(ctx).Preload("ExtendedInfo").First(&agent, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("could not get agent %s with extended info: %w", id, err)
	}
	return &agent, nil
}

func (s *AgentRepository) GetWithAllRelations(ctx context.Context, id string) (*local_models.ServerAgentModel, error) {
	var agent local_models.ServerAgentModel
	err := s.db.WithContext(ctx).
		Preload("Commands").
		Preload("ExtendedInfo").
		Preload("Sessions").
		First(&agent, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("could not get agent %s with all relations: %w", id, err)
	}
	return &agent, nil
}

func (s *AgentRepository) GetActiveAgents(ctx context.Context) ([]local_models.ServerAgentModel, error) {
	var agents []local_models.ServerAgentModel
	err := s.db.WithContext(ctx).Where("is_active = ?", true).
		Preload("ExtendedInfo").
		Find(&agents).Error
	if err != nil {
		return nil, fmt.Errorf("could not get active agents: %w", err)
	}
	return agents, nil
}

func (s *AgentRepository) UpdateExtendedInfo(ctx context.Context, agentID string, info *models.ExtendedAgentInfo) error {
	info.AgentID = agentID
	if err := s.db.WithContext(ctx).Save(info).Error; err != nil {
		return fmt.Errorf("could not update extended info for agent %s: %w", agentID, err)
	}
	return nil
}

func (s *AgentRepository) UpdateLastSeen(ctx context.Context, agentID string) error {
	err := s.db.WithContext(ctx).Model(&local_models.ServerAgentModel{}).
		Where("id = ?", agentID).
		UpdateColumn("last_seen", time.Now()).Error
	if err != nil {
		return fmt.Errorf("could not update last_seen for agent %s: %w", agentID, err)
	}
	return nil
}

func (s *AgentRepository) GetExtendedInfo(ctx context.Context, agentID string) (*models.ExtendedAgentInfo, error) {
	var info models.ExtendedAgentInfo
	err := s.db.WithContext(ctx).Where("agent_id = ?", agentID).First(&info).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("could not get extended info for agent %s: %w", agentID, err)
	}
	return &info, nil
}

func (s *AgentRepository) CreateSession(ctx context.Context, session *local_models.AgentSession) error {
	if err := s.db.WithContext(ctx).Create(session).Error; err != nil {
		return fmt.Errorf("could not create session for agent %s: %w", session.AgentID, err)
	}
	return nil
}

func (s *AgentRepository) EndSession(ctx context.Context, sessionID string) error {
	updates := map[string]interface{}{
		"is_active": false,
		"end_time":  time.Now(),
	}
	err := s.db.WithContext(ctx).Model(&local_models.AgentSession{}).
		Where("session_id = ?", sessionID).
		Updates(updates).Error
	if err != nil {
		return fmt.Errorf("could not end session %s: %w", sessionID, err)
	}
	return nil
}

func (s *AgentRepository) GetActiveSessions(ctx context.Context) ([]local_models.AgentSession, error) {
	var sessions []local_models.AgentSession
	err := s.db.WithContext(ctx).Where("is_active = ?", true).Find(&sessions).Error
	if err != nil {
		return nil, fmt.Errorf("could not get active sessions: %w", err)
	}
	return sessions, nil
}

func (s *AgentRepository) GetAgentSessions(ctx context.Context, agentID string) ([]local_models.AgentSession, error) {
	var sessions []local_models.AgentSession
	err := s.db.WithContext(ctx).Where("agent_id = ?", agentID).
		Order("created_at DESC").
		Find(&sessions).Error
	if err != nil {
		return nil, fmt.Errorf("could not get sessions for agent %s: %w", agentID, err)
	}
	return sessions, nil
}

// Command operations

func (s *AgentRepository) SaveCommand(ctx context.Context, command *local_models.AgentCommand) error {
	if err := s.db.WithContext(ctx).Create(command).Error; err != nil {
		return fmt.Errorf("could not add command %s for agent %s: %w", command.ID, command.AgentID, err)
	}
	return nil
}

func (s *AgentRepository) GetCommandsByStatus(ctx context.Context, agentID string, status models.CommandStatus) ([]local_models.AgentCommand, error) {
	var commands []local_models.AgentCommand
	err := s.db.WithContext(ctx).Where("agent_id = ? AND status = ?", agentID, status).
		Order("created_at ASC").
		Find(&commands).Error
	if err != nil {
		return nil, fmt.Errorf("could not get commands with status %d for agent %s: %w", status, agentID, err)
	}
	return commands, nil
}

// updateCommandStatus is a private helper for use in transactions
func (s *AgentRepository) updateCommandStatus(db *gorm.DB, ctx context.Context, commandID string, status models.CommandStatus) error {
	result := db.WithContext(ctx).Model(&local_models.AgentCommand{}).
		Where("id = ?", commandID).
		Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("could not update status for command %s: %w", commandID, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no command found with ID %s to update", commandID)
	}
	return nil
}

func (s *AgentRepository) UpdateCommandStatus(ctx context.Context, commandID string, status models.CommandStatus) error {
	return s.updateCommandStatus(s.db, ctx, commandID, status)
}

func (s *AgentRepository) GetCommands(ctx context.Context, agentID string, limit int) ([]local_models.AgentCommand, error) {
	var commands []local_models.AgentCommand
	query := s.db.WithContext(ctx).Where("agent_id = ?", agentID).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&commands).Error; err != nil {
		return nil, fmt.Errorf("could not get commands for agent %s: %w", agentID, err)
	}
	return commands, nil
}

// saveCommandResult is a private helper for use in transactions
func (s *AgentRepository) saveCommandResult(db *gorm.DB, ctx context.Context, result *local_models.AgentCommandResult) error {
	err := db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"output", "status"}),
	}).Create(result).Error

	if err != nil {
		return fmt.Errorf("could not save command result for command %s: %w", result.ID, err)
	}
	return nil
}

func (s *AgentRepository) SaveCommandResult(ctx context.Context, result *local_models.AgentCommandResult) error {
	return s.saveCommandResult(s.db, ctx, result)
}

func (s *AgentRepository) UpdateCommandStatusWithResult(ctx context.Context, agentID, commandID string, status models.CommandStatus, output any) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := s.updateCommandStatus(tx, ctx, commandID, status); err != nil {
			return err
		}

		outputBytes, err := json.Marshal(output)
		if err != nil {
			return fmt.Errorf("failed to marshal command output for command %s: %w", commandID, err)
		}

		commandResult := local_models.AgentCommandResult{
			AgentID: agentID,
			CommandResult: models.CommandResult{
				ID:     commandID,
				Status: status,
				Output: string(outputBytes),
			},
		}

		if err := s.saveCommandResult(tx, ctx, &commandResult); err != nil {
			return err
		}

		return nil
	})
}

func (s *AgentRepository) GetCommandResult(ctx context.Context, commandID string) (*models.CommandResult, error) {
	var result local_models.AgentCommandResult
	err := s.db.WithContext(ctx).Where("id = ?", commandID).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("could not get result for command %s: %w", commandID, err)
	}
	return &result.CommandResult, nil
}

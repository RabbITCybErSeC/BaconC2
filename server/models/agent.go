package models

import (
	"sync"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
)

type AgentCommand struct {
	AgentID        string `json:"agent_id" gorm:"index"`
	models.Command `gorm:"embedded"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type AgentCommandResult struct {
	AgentID              string `json:"agent_id" gorm:"index"`
	models.CommandResult `gorm:"embedded"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type ServerAgentModel struct {
	models.Agent `gorm:"embedded"`
	LastSeen     time.Time                `json:"last_seen" gorm:"column:last_seen"`
	IsActive     bool                     `json:"is_active" gorm:"column:is_active"`
	Commands     []AgentCommand           `json:"-" gorm:"foreignKey:AgentID"`
	ExtendedInfo models.ExtendedAgentInfo `json:"extended_info" gorm:"foreignKey:ID;references:AgentID"`
	Sessions     []AgentSession           `json:"sessions,omitempty" gorm:"foreignKey:AgentID;references:ID"`
	CommandMu    sync.Mutex               `json:"-" gorm:"-"`
}

type AgentSession struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	AgentID   string     `json:"agent_id" gorm:"index"`
	SessionID string     `json:"session_id" gorm:"uniqueIndex;size:255"`
	StartTime time.Time  `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
	IPAddress string     `json:"ip_address" gorm:"size:45"`
	UserAgent string     `json:"user_agent" gorm:"size:500"`
	IsActive  bool       `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func NewAgentCommand(agentID string, cmd models.Command) *AgentCommand {
	now := time.Now()
	return &AgentCommand{
		AgentID:   agentID,
		Command:   cmd,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func NewAgentCommandResult(agentID string, result models.CommandResult) *AgentCommandResult {
	now := time.Now()
	return &AgentCommandResult{
		AgentID:       agentID,
		CommandResult: result,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func NewServerAgentModel(agent models.Agent) *ServerAgentModel {
	return &ServerAgentModel{
		Agent:    agent,
		LastSeen: time.Now(),
		IsActive: true,
		Commands: make([]AgentCommand, 0),
		Sessions: make([]AgentSession, 0),
	}
}

func NewAgentSession(agentID, sessionID, ipAddress, userAgent string) *AgentSession {
	now := time.Now()
	return &AgentSession{
		AgentID:   agentID,
		SessionID: sessionID,
		StartTime: now,
		EndTime:   nil,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

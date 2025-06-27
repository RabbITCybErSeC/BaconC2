package models

import (
	"sync"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
)

type AgentCommand struct {
	ID      string `json:"id" gorm:"primaryKey"`
	Command models.Command
	AgentID string `json:"agent_id" gorm:"index"`
}

type ServerAgentModel struct {
	BaseAgentModel models.Agent
	LastSeen       time.Time      `json:"last_seen" gorm:"column:last_seen"`
	IsActive       bool           `json:"is_active" gorm:"column:is_active"`
	Commands       []AgentCommand `json:"-" gorm:"foreignKey:AgentID"`
	CommandMu      sync.Mutex     `json:"-" gorm:"-"`
}

package models

import (
	"sync"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
)

type AgentCommand struct {
	Command models.Command `gorm:"embedded"`
	AgentID string         `json:"agent_id" gorm:"index"`
}
type ServerAgentModel struct {
	BaseAgentModel models.Agent   `json:"-" gorm:"embedded"`
	LastSeen       time.Time      `json:"last_seen" gorm:"column:last_seen"`
	IsActive       bool           `json:"is_active" gorm:"column:is_active"`
	Commands       []AgentCommand `json:"-" gorm:"foreignKey:AgentID"`
	CommandMu      sync.Mutex     `json:"-" gorm:"-"`
}

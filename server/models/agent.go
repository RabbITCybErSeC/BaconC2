package models

import (
	"sync"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
)

type AgentCommand struct {
	AgentID string         `json:"agent_id" gorm:"index"`
	Command models.Command `gorm:"embedded"`
}
type ServerAgentModel struct {
	BaseAgentModel models.Agent             `json:"-" gorm:"embedded"`
	LastSeen       time.Time                `json:"last_seen" gorm:"column:last_seen"`
	IsActive       bool                     `json:"is_active" gorm:"column:is_active"`
	Commands       []AgentCommand           `json:"-" gorm:"foreignKey:AgentID"`
	CommandMu      sync.Mutex               `json:"-" gorm:"-"`
	ExtendedInfo   models.ExtendedAgentInfo `json:"extended_info" gorm:"foreignKey:ID;references:AgentID"`
}

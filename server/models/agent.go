package models

import (
	"sync"
	"time"

	"github.com/RabbITCybErSeC/BaconC2/pkg/models"
)

type AgentCommand struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	AgentID   string         `json:"agent_id" gorm:"index"`
	Command   models.Command `gorm:"embedded"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type ServerAgentModel struct {
	BaseAgentModel models.Agent             `json:"-" gorm:"embedded"`
	LastSeen       time.Time                `json:"last_seen" gorm:"column:last_seen"`
	IsActive       bool                     `json:"is_active" gorm:"column:is_active"`
	Commands       []AgentCommand           `json:"-" gorm:"foreignKey:AgentID"`
	ExtendedInfo   models.ExtendedAgentInfo `json:"extended_info" gorm:"foreignKey:ID;references:AgentID"`
	Sessions       []AgentSession           `json:"sessions,omitempty" gorm:"foreignKey:AgentID;references:ID"`
	CommandMu      sync.Mutex               `json:"-" gorm:"-"`
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

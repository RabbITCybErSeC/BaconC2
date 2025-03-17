package models

import (
	"sync"
	"time"
)

type Command struct {
	ID      string `json:"id" gorm:"primaryKey"`
	Command string `json:"command"`
	Status  string `json:"status"`
	Output  string `json:"output,omitempty"`
	AgentID string `json:"agent_id" gorm:"index"`
}

type Agent struct {
	ID        string     `json:"id" gorm:"primaryKey"`
	Hostname  string     `json:"hostname"`
	IP        string     `json:"ip"`
	LastSeen  time.Time  `json:"lastSeen"`
	OS        string     `json:"os"`
	IsActive  bool       `json:"isActive"`
	Protocol  string     `json:"protocol"`
	Commands  []Command  `json:"-" gorm:"foreignKey:AgentID"`
	CommandMu sync.Mutex `gorm:"-"`
}

package db

import (
	"github.com/RabbITCybErSeC/BaconC2/server/models"
	"gorm.io/gorm"
)

type AgentRepositoryInterface interface {
	Save(agent *models.Agent) error
	Get(id string) (*models.Agent, error)
	GetAll() ([]models.Agent, error)
}

type AgentRepository struct {
	db *gorm.DB
}

func NewAgentRepository(db *gorm.DB) *AgentRepository {
	return &AgentRepository{db: db}
}

func (r *AgentRepository) Save(agent *models.Agent) error {
	return r.db.Save(agent).Error
}

func (r *AgentRepository) Get(id string) (*models.Agent, error) {
	var agent models.Agent
	err := r.db.First(&agent, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func (r *AgentRepository) GetAll() ([]models.Agent, error) {
	var agents []models.Agent
	err := r.db.Find(&agents).Error
	if err != nil {
		return nil, err
	}
	return agents, nil
}

package db

import (
	local_models "github.com/RabbITCybErSeC/BaconC2/server/models"
	"gorm.io/gorm"
)

type AgentRepositoryInterface interface {
	Save(agent *local_models.ServerAgentModel) error
	Get(id string) (*local_models.ServerAgentModel, error)
	GetAll() ([]local_models.ServerAgentModel, error)
}

type AgentRepository struct {
	db *gorm.DB
}

func NewAgentRepository(db *gorm.DB) *AgentRepository {
	return &AgentRepository{db: db}
}

func (r *AgentRepository) Save(agent *local_models.ServerAgentModel) error {
	return r.db.Save(agent).Error
}

func (r *AgentRepository) Get(id string) (*local_models.ServerAgentModel, error) {
	var agent local_models.ServerAgentModel
	err := r.db.First(&agent, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func (r *AgentRepository) GetAll() ([]local_models.ServerAgentModel, error) {
	var agents []local_models.ServerAgentModel
	err := r.db.Find(&agents).Error
	if err != nil {
		return nil, err
	}
	return agents, nil
}

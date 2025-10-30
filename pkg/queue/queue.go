package queue

import "github.com/RabbITCybErSeC/BaconC2/pkg/models"

type ICommandExecutor interface {
	Execute(cmd models.Command) models.CommandResult
	ProcessCommandQueue()
}

type ICommandQueue interface {
	Add(cmd models.Command) error
	Get() (models.Command, bool)
	List() ([]models.Command, error)
	RemoveFirst() (models.Command, error)
	RemoveAt(index int) (models.Command, error)
	RemoveItem(cmd models.Command) error
	Size() int
}

type IResultQueue interface {
	Add(cmd models.CommandResult) error
	Get() (models.CommandResult, bool)
	List() ([]models.CommandResult, error)
	RemoveFirst() (models.CommandResult, error)
	RemoveAt(index int) (models.CommandResult, error)
	RemoveItem(cmd models.CommandResult) error
	Size() int
}

type IServerCommandQueue interface {
	MultiQueue[models.Command]
}

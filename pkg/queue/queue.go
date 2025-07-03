package queue

import "github.com/RabbITCybErSeC/BaconC2/pkg/models"

type ICommandExecutor interface {
	Execute(cmd models.Command) models.CommandResult
	ProcessCommandQueue()
}

type ICommandQueue interface {
	GenericQueue[models.Command]
}

type IResultQueue interface {
	GenericQueue[models.CommandResult]
}

type IServerCommandQueue interface {
	MultiQueue[models.Command]
}

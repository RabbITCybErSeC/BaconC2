package command_handler

var commandRegistry *CommandHandlerRegistry

func init() {
	commandRegistry = NewCommandHandlerRegistry()
}

func GetGlobalCommandRegistry() *CommandHandlerRegistry {
	return commandRegistry
}

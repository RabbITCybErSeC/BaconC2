package command_handler

var commandRegistry *CommandHandlerRegistry

func init() {
	commandRegistry = NewCommandHandlerRegistry()
	// setup default built-in handlers
	commandRegistry.RegisterHandler(*NewGetRegisteredClientHandlers())
}

func GetGlobalCommandRegistry() *CommandHandlerRegistry {
	return commandRegistry
}

package command_handler

var CommandRegistry *CommandHandlerRegistry

func init() {
	CommandRegistry = NewCommandHandlerRegistry()
}

// type RegisterFunc func(registry *command_handler.CommandHandlerRegistry)

// type CommandHandlerRegistry struct {
// 	handlers map[string]RegisterFunc
// }

// var registeredHandlers []RegisterFunc

// func AddRegisterFunc(fn RegisterFunc) {
// 	registeredHandlers = append(registeredHandlers, fn)
// }

// func NewCommandHandlerRegistry() *CommandHandlerRegistry {
// 	return &CommandHandlerRegistry{
// 		handlers: make(map[string]interface{}), // Adjust as needed
// 	}
// }

// func RegisterAllHandlers(registry *command_handler.CommandHandlerRegistry) {
// 	for _, fn := range registeredHandlers {
// 		fn(registry)
// 	}
// }

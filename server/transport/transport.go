package transport

type ITransportProtocol interface {
	Start() error
	Stop() error
	Name() string
}

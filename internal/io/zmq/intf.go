package zmq

type IO interface {
	Events() chan *Event
	Commands() chan *Command
	Input() chan *Message
	Output() chan *Message
	Connected(endpoint string) bool
	ConnectedList() []string
}

type EventType uint8

const (
	OnConnected EventType = iota
	OnDisconnected
	ErrorConnection
	ErrorSocket
)

type CommandType uint8

const (
	TryConnect CommandType = iota
	TryDisconnect
)

type Event struct {
	EventType EventType
	Endpoint  string
}

type Command struct {
	CommandType CommandType
	Endpoint    string
}

type Message struct {
	Message  []byte
	Endpoint string
}

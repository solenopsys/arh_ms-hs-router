package ws

type IO interface {
	GetEvents() chan *Event
	GetCommands() chan *Command
	GetInput() chan *Message
	GetOutput() chan *Message
}

type Auth interface {
	GetAuth(token string) (uint16, error)
}

type EventType uint8

const (
	OnConnected EventType = iota
	OnDisconnected
	ErrorAuth
	ErrorSocket
	UpgradeError
)

type CommandType uint8

const (
	WsTryConnect CommandType = iota
	WsTryDisconnect
)

type Event struct {
	EventType EventType
	Key       string
}

type Command struct {
	EventType CommandType
	Endpoint  string
}

type Message struct {
	Message       []byte
	ConnectionKey string
	User          uint16
	From          string
}

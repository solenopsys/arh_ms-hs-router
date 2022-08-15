package core

type CommandType uint8

const (
	ZmqUpdateServices CommandType = iota
	ZmqTryConnect
	ZmqTryDisconnect
	ZmqOnConnected
	ZmqOnDisconnected // ? нет обработчика
	WsOnConnected
	WsUpgradeError // ? нет обработчика
	WsTryDisconnect
	WsOnDisconnected // ? нет обработчика
	UpdateConfigMap
)

func (me CommandType) String() string {
	return [...]string{
		"Zmq UpdateServices",
		"Zmq TryConnect",
		"Zmq TryDisconnect",
		"Zmq OnConnected",
		"Zmq OnInitialized",
		"Zmq OnDisconnected",
		"Ws OnConnected",
		"Ws OnInitialized",
		"Ws UpgradeError",
		"Ws TryDisconnect",
		"Ws OnDisconnected",
	}[me]
}

type Command struct {
	MessageType CommandType
	Key         string
}

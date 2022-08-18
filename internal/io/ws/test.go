package ws

func NewWsHub() *Hub {
	return &Hub{
		Connections: make(map[string]*Connection),
	}
}

//
//state.IncStat("WsUpgradeError")
//state.IncStat("NewWsConnection")
// state.IncStat("EndpointDisconnected")

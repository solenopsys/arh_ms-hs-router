package ws

func NewWsHub() *Hub {
	return &Hub{
		connections: make(map[string]*Connection),
	}
}

//
//state.IncStat("WsUpgradeError")
//state.IncStat("NewWsConnection")
// state.IncStat("EndpointDisconnected")

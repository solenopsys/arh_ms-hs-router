package main

import (
	"github.com/gorilla/websocket"
	"k8s.io/klog/v2"
	"net/http"
)

type WSConnect struct {
	connection *websocket.Conn
	userId     uint16
}

func newWsHub() *WsHub {
	return &WsHub{
		clients: make(map[string]*WSConnect),
	}
}

type WsHub struct {
	clients map[string]*WSConnect
}

func getAuth(token string) (uint16, error) {
	return 1, nil
}

func (wsHub *WsHub) tryConnectionProcessing(state *State, w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token != "" {
		userId, err := getAuth(token)
		if err == nil {
			setCorsHeaders(w)
			connection, err := upgrader.Upgrade(w, r, nil)
			var id = connection.RemoteAddr().String()
			if err != nil {
				klog.Error("Wpgrade ws:", err)
				state.command <- &Command{key: id, messageType: WsUpgradeError}

				state.incStat("WsUpgradeError")
			}

			state.incStat("NewWsConnection")

			wsHub.clients[id] = &WSConnect{userId: userId, connection: connection}
			state.command <- &Command{key: id, messageType: WsOnConnected}
		} else {
			w.WriteHeader(403)
		}

	}

}

func (wsHub *WsHub) onConnectProcessing(state *State, endpoint string) {
	go wsHub.wsInMessageProcessing(state, endpoint)
	state.command <- &Command{key: endpoint, messageType: WsOnInitialized}
}

func (wsHub *WsHub) tryDisconnectProcessing(state *State, key string) {
	if client, ok := wsHub.clients[key]; ok {
		delete(wsHub.clients, key)
		err := client.connection.Close()
		if err != nil {
			klog.Error("Disconnect error:", err)
		} else {
			state.command <- &Command{key: key, messageType: WsOnDisconnected}
		}
	}
}

func (wsHub *WsHub) wsInMessageProcessing(state *State, key string) { //todo обработать отключение через контекст.
	conn := wsHub.clients[key]
	for {
		_, message, err := conn.connection.ReadMessage()
		if err != nil {
			klog.Error("Read message error:", err)
			state.command <- &Command{key: key, messageType: WsOnDisconnected}
			break
		}
		klog.Infof("Ws in мessage:", string(message))
		state.messages <- &MessageWrapper{body: message, userId: conn.userId, key: key, messageType: WsInput}
	}
}

func (wsHub *WsHub) sendToWsMessageProcessing(key string, body []byte) {
	conn := wsHub.clients[key]
	if conn != nil {
		err := conn.connection.WriteMessage(websocket.BinaryMessage, body)
		if err != nil {
			klog.Error("read:", err)
		}
	}

}

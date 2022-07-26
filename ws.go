package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
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

//ok
func (wsHub *WsHub) tryConnectionProcessing(state *State, w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token != "" {
		userId, err := getAuth(token)
		if err == nil {
			setCorsHeaders(w)
			//responseHeader http.Header https://pkg.go.dev/github.com/gorilla/websocket#hdr-Origin_Considerations
			connection, err := upgrader.Upgrade(w, r, nil)
			var id = connection.RemoteAddr().String()
			if err != nil {
				log.Print("upgrade:", err) //todo обработать ошибку
				state.command <- &Command{key: id, messageType: WsUpgradeError}
			}

			wsHub.clients[id] = &WSConnect{userId: userId, connection: connection}
			state.command <- &Command{key: id, messageType: WsOnConnected}
		} else {
			w.WriteHeader(403)
		}

	}

}

//ok
func (wsHub *WsHub) onConnectProcessing(state *State, endpoint string) {
	go wsHub.wsInMessageProcessing(state, endpoint)
	state.command <- &Command{key: endpoint, messageType: WsOnInitialized}
}

//ok
func (wsHub *WsHub) tryDisconnectProcessing(state *State, key string) {
	if client, ok := wsHub.clients[key]; ok {
		delete(wsHub.clients, key)
		err := client.connection.Close()
		if err != nil {
			log.Println("read:", err) //todo обработать
		} else {
			state.command <- &Command{key: key, messageType: WsOnDisconnected}
		}
	}
}

//ok
func (wsHub *WsHub) wsInMessageProcessing(state *State, key string) { //todo обработать отключение через контекст.
	conn := wsHub.clients[key]
	for {
		mt, message, err := conn.connection.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			state.command <- &Command{key: key, messageType: WsOnDisconnected}
			break
		}
		log.Println("mt:", mt) //todo убрать мусор
		fmt.Println(string(message))
		state.messages <- &MessageWrapper{body: message, userId: conn.userId, key: key, messageType: WsInput}
	}
}

func (wsHub *WsHub) sendToWsMessageProcessing(key string, body []byte) {
	conn := wsHub.clients[key]
	if conn != nil {
		err := conn.connection.WriteMessage(websocket.BinaryMessage, body)
		if err != nil {
			log.Println("read:", err) //todo обработать
		}
	}

}

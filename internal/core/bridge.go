package core

import (
	"solenopsys.org/zmq_router/internal/io/ws"
	"solenopsys.org/zmq_router/internal/io/zmq"
	"solenopsys.org/zmq_router/pkg/utils"
)

type MessagesBridge struct {
	WsHub     ws.IO
	ZmqHub    zmq.IO
	Router    *Routing
	Statistic *utils.Statistic
}

func (cd MessagesBridge) Init() {
	go cd.wsToZmqConveyor()
	go cd.zmqToWsConveyor()
}

func (cd MessagesBridge) stat(message string) {
	if cd.Statistic != nil {
		cd.Statistic.Increment <- message
	}
}

func (cd MessagesBridge) wsToZmqConveyor() {
	for {
		srcMessage := <-cd.WsHub.GetInput()
		pack := MessagePackage{Raw: srcMessage.Message, UserId: srcMessage.User, Key: srcMessage.ConnectionKey}
		if pack.IsFirst() {
			stream, result := cd.Router.WsRouteNew(&pack)
			if result == Ok {
				cd.stat("FirstWsMessageOk")
				cd.ZmqHub.GetOutput() <- &zmq.Message{Message: pack.UserInjectedBody(), Endpoint: stream.ZmqEndpoint}
			} else {
				cd.stat("FirstWsMessageErr-" + result.String())
				cd.WsHub.GetOutput() <- &ws.Message{Message: pack.ErrorResponseBody(result.String()), ConnectionKey: srcMessage.ConnectionKey}
			}
		} else {
			stream, result := cd.Router.WsRoute(&pack)
			if result == Ok {
				cd.stat("SecondWsMessageOk")
				cd.ZmqHub.GetOutput() <- &zmq.Message{Message: pack.Raw, Endpoint: stream.ZmqEndpoint}
			} else {
				cd.stat("SecondWsMessageErr-" + result.String())
				cd.WsHub.GetOutput() <- &ws.Message{Message: pack.ErrorResponseBody(result.String()), ConnectionKey: srcMessage.ConnectionKey}
			}
		}
	}
}

func (cd MessagesBridge) zmqToWsConveyor() {
	for {
		srcMessage := <-cd.ZmqHub.GetInput()
		pack := MessagePackage{Raw: srcMessage.Message, Key: srcMessage.Endpoint}

		stream, result := cd.Router.ZqmRoute(&pack)
		if result == Ok {
			cd.stat("ZmqMessageOk")
			cd.WsHub.GetOutput() <- &ws.Message{Message: pack.Raw, ConnectionKey: stream.WsEndpoint}
		} else {
			cd.stat("ZmqMessageErr-" + result.String())
			cd.ZmqHub.GetOutput() <- &zmq.Message{Message: pack.ErrorResponseBody(result.String()), Endpoint: srcMessage.Endpoint}
		}
	}
}

package core

import (
	"solenopsys.org/zmq_router/internal/io/ws"
	"solenopsys.org/zmq_router/internal/io/zmq"
	"solenopsys.org/zmq_router/pkg/utils"
)

type MessagesBridge struct {
	wsHub     ws.IO
	zmqHub    zmq.IO
	router    *Routing
	statistic *utils.Statistic
}

func (cd MessagesBridge) init() {
	go cd.wsToZmqConveyor()
	go cd.zmqToWsConveyor()
}

func (cd MessagesBridge) wsToZmqConveyor() {
	for {
		srcMessage := <-cd.wsHub.Input()
		pack := MessagePackage{raw: srcMessage.Message, UserId: srcMessage.User, ConnectionKey: srcMessage.ConnectionKey}
		if pack.IsFirst() {
			stream, result := cd.router.wsRouteNew(&pack)
			if result == Ok {
				cd.statistic.Increment <- "FirstWsMessageOk"
				cd.zmqHub.Output() <- &zmq.Message{Message: pack.userInjectedBody(), Endpoint: stream.zmqEndpoint}
			} else {
				cd.statistic.Increment <- "FirstWsMessageErr"
				cd.statistic.Increment <- result.String()
				cd.wsHub.Output() <- &ws.Message{Message: pack.errorResponseBody(result.String()), ConnectionKey: srcMessage.ConnectionKey}
			}
		} else {
			stream, result := cd.router.wsRoute(&pack)
			if result == Ok {
				cd.statistic.Increment <- "SecondWsMessageOk"
				cd.zmqHub.Output() <- &zmq.Message{Message: pack.raw, Endpoint: stream.zmqEndpoint}
			} else {
				cd.statistic.Increment <- "SecondWsMessageErr"
				cd.statistic.Increment <- result.String()
				cd.wsHub.Output() <- &ws.Message{Message: pack.errorResponseBody(result.String()), ConnectionKey: srcMessage.ConnectionKey}
			}
		}
	}
}

func (cd MessagesBridge) zmqToWsConveyor() {
	for {
		srcMessage := <-cd.zmqHub.Input()
		pack := MessagePackage{raw: srcMessage.Message, ConnectionKey: srcMessage.Endpoint}

		stream, result := cd.router.zqmRoute(&pack)
		if result == Ok {
			cd.statistic.Increment <- "ZmqMessageOk"
			cd.wsHub.Output() <- &ws.Message{Message: pack.userInjectedBody(), ConnectionKey: stream.wsEndpoint}
		} else {
			cd.statistic.Increment <- "ZmqMessageErr"
			cd.statistic.Increment <- result.String()
			cd.zmqHub.Output() <- &zmq.Message{Message: pack.errorResponseBody(result.String()), Endpoint: srcMessage.Endpoint}
		}
	}
}

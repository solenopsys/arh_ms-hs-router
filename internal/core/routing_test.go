package core

import (
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"solenopsys.org/zmq_router/pkg/utils"
)

type MessageMock struct {
	stream  uint32
	service uint16
	key     string
	first   bool
}

func (m MessageMock) Stream() uint32 {
	return m.stream
}
func (m MessageMock) Service() uint16 {
	return m.service
}
func (m MessageMock) ConnectionKey() string {
	return m.key
}
func (m MessageMock) IsFirst() bool {
	return m.first
}

var _ = Describe("Router", func() {
	var rounting *Routing
	var endpoint string = ""

	BeforeEach(func() {
		endpoint = ""
		rounting = &Routing{
			Routes: make(map[uint32]*StreamConfig),

			RandomService: func(serviceId uint16) (string, error) {
				if endpoint == "" {
					return endpoint, errors.New("not found")
				} else {
					return endpoint, nil
				}

			},
		}
	})

	Context("Add routes", func() {

		It("should be empty", func() {
			Expect(len(rounting.Routes)).To(Equal(0))
		})

		It("should be right route", func() {
			endpoint = "zmq1"
			routeNew, result := rounting.WsRouteNew(&MessageMock{12132, 2123, "ws1", true})
			Expect(result).To(Equal(Ok))
			Expect(len(rounting.Routes)).To(Equal(1))
			_, has := rounting.Routes[12132]
			Expect(has).To(Equal(true))

			Expect(routeNew.WsEndpoint).To(Equal("ws1"))
			Expect(routeNew.ZmqEndpoint).To(Equal("zmq1"))
		})

		It("should be first flag", func() {
			endpoint = "zmq1"
			_, result1 := rounting.WsRouteNew(&MessageMock{12122, 2123, "ws1", true})
			Expect(result1).To(Equal(Ok))
			_, result2 := rounting.WsRouteNew(&MessageMock{12132, 2123, "ws1", false})
			Expect(result2).To(Equal(FirstFlagNotFound))
		})

		It("should be streams collision", func() {
			endpoint = "zmq1"
			st := &MessageMock{12132, 2123, "ws1", true}
			rounting.WsRouteNew(st)
			_, result := rounting.WsRouteNew(st)
			Expect(result).To(Equal(StreamCollision))
		})

		It("should be stream not found", func() {
			endpoint = "zmq1"
			_, result := rounting.WsRoute(&MessageMock{12132, 2123, "ws1", true})
			Expect(result).To(Equal(StreamNotFound))
		})

		It("should be endpoint not found", func() {
			_, result := rounting.WsRouteNew(&MessageMock{12132, 2123, "ws1", true})
			Expect(result).To(Equal(EndpointNotFound))
		})

		It("should be stream not found", func() {
			_, result := rounting.ZqmRoute(&MessageMock{12132, 2123, "ws1", true})
			Expect(result).To(Equal(StreamNotFound))
		})

		It("should be stream not found", func() {
			endpoint = "zmq1"
			rounting.WsRouteNew(&MessageMock{12132, 2123, "ws1", true})
			rounting.WsRouteNew(&MessageMock{12133, 2123, "ws2", true})
			_, result1 := rounting.ZqmRoute(&MessageMock{12132, 2123, "ws1", true})
			Expect(result1).To(Equal(Ok))
			_, result2 := rounting.ZqmRoute(&MessageMock{12132, 2123, "ws2", true})
			Expect(result2).To(Equal(AccessDenied))
		})
	})

	Context("Find and remove routes", func() {

		It("should be find streams by ws and zmq", func() {
			endpoint = "zmq1"
			rounting.WsRouteNew(&MessageMock{1, 2123, "ws1", true})
			rounting.WsRouteNew(&MessageMock{2, 2123, "ws2", true})

			endpoint = "zmq2"
			rounting.WsRouteNew(&MessageMock{3, 2123, "ws2", true})
			rounting.WsRouteNew(&MessageMock{4, 2123, "ws4", true})

			ws2 := rounting.SelectRoutesByWs("ws2")
			utils.SortUint32(ws2)
			Expect(ws2).To(Equal([]uint32{2, 3}))
			ws4 := rounting.SelectRoutesByWs("ws4")
			utils.SortUint32(ws4)
			Expect(ws4).To(Equal([]uint32{4}))
			zmq1 := rounting.SelectRoutesByZmq("zmq1")
			utils.SortUint32(zmq1)
			Expect(zmq1).To(Equal([]uint32{1, 2}))
			zmq2 := rounting.SelectRoutesByZmq("zmq2")
			utils.SortUint32(zmq2)
			Expect(zmq2).To(Equal([]uint32{3, 4}))
		})

		It("should be find streams by ws and zmq", func() {
			endpoint = "zmq1"
			rounting.WsRouteNew(&MessageMock{1, 2123, "ws1", true})
			rounting.WsRouteNew(&MessageMock{2, 2123, "ws2", true})

			Expect(len(rounting.Routes)).To(Equal(2))
			rounting.DeleteRoute(2)
			Expect(len(rounting.Routes)).To(Equal(1))

		})
	})
})

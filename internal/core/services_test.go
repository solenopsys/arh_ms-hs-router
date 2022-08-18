package core

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type EndpointsMock struct {
	endpoints map[string]string
}

func (e EndpointsMock) UpdateEndpoints() (map[string]string, error) {
	return e.endpoints, nil
}

type MappingMock struct {
	mapping map[string]uint16
}

func (m MappingMock) UpdateMapping() (map[string]uint16, error) {
	return m.mapping, nil
}

var _ = Describe("ServicesMap", func() {
	var services *ServicesController

	BeforeEach(func() {

		services = &ServicesController{
			EndpointsApi: EndpointsMock{endpoints: map[string]string{
				"tcp://192.168.122.1:9999": "serv1",
				"tcp://192.168.122.2:9999": "serv2",
				"tcp://192.168.122.3:9999": "serv2",
				"tcp://192.168.122.4:9999": "serv3",
				"tcp://192.168.122.5:9999": "serv3",
				"tcp://192.168.122.6:9999": "serv3",
				"tcp://192.168.122.7:9999": "serv4",
				"tcp://192.168.122.8:9999": "serv4",
			}},
			MappingApi:   MappingMock{mapping: map[string]uint16{"serv1": 1, "serv2": 2, "serv3": 3, "serv4": 4}},
			ServicesMap:  make(map[string]uint16),
			EndpointsMap: make(map[string]uint16),
			Groups:       make(map[uint16][]string),
		}

	})

	Context("Check load values", func() {
		It("check services map", func() {
			services.SyncEndpoints()
			Expect(services.Services()).To(Equal(map[string]uint16{"serv1": 1, "serv2": 2, "serv3": 3, "serv4": 4}))
		})

		It("check endpoints map", func() {
			services.SyncEndpoints()
			Expect(services.Endpoints()).To(Equal(map[string]uint16{
				"tcp://192.168.122.1:9999": 1,
				"tcp://192.168.122.2:9999": 2,
				"tcp://192.168.122.3:9999": 2,
				"tcp://192.168.122.4:9999": 3,
				"tcp://192.168.122.5:9999": 3,
				"tcp://192.168.122.6:9999": 3,
				"tcp://192.168.122.7:9999": 4,
				"tcp://192.168.122.8:9999": 4,
			}))
		})

		It("check groups", func() {
			services.SyncEndpoints()
			Expect(services.Groups[1]).To(ConsistOf([]string{"tcp://192.168.122.1:9999"}))
			Expect(services.Groups[2]).To(ConsistOf([]string{"tcp://192.168.122.2:9999", "tcp://192.168.122.3:9999"}))
			Expect(services.Groups[3]).To(ConsistOf([]string{"tcp://192.168.122.6:9999", "tcp://192.168.122.4:9999", "tcp://192.168.122.5:9999"}))
			Expect(services.Groups[4]).To(ConsistOf([]string{"tcp://192.168.122.7:9999", "tcp://192.168.122.8:9999"}))
		})

		It("check rand", func() {
			services.SyncEndpoints()

			Expect(services.RandomEndpointByType(1)).To(BeElementOf([]string{"tcp://192.168.122.1:9999"}))
			Expect(services.RandomEndpointByType(2)).To(BeElementOf([]string{"tcp://192.168.122.2:9999", "tcp://192.168.122.3:9999"}))
			Expect(services.RandomEndpointByType(3)).To(BeElementOf([]string{"tcp://192.168.122.6:9999", "tcp://192.168.122.4:9999", "tcp://192.168.122.5:9999"}))
			Expect(services.RandomEndpointByType(4)).To(BeElementOf([]string{"tcp://192.168.122.7:9999", "tcp://192.168.122.8:9999"}))
		})
	})

})

package core

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func pack(b []byte) MessagePackage {
	return MessagePackage{Raw: b, UserId: 100, Key: "bla"}
}

var _ = Describe("Message", func() {

	Describe("Parsing tests", func() {
		Context("First Flag 4 byte", func() {
			It("should be true", func() {
				Expect(pack([]byte{33, 43, 44, 34, 15, 4}).IsFirst()).To(Equal(true))
			})
			It("should be false", func() {
				Expect(pack([]byte{33, 43, 44, 34, 14, 4}).IsFirst()).To(Equal(false))
			})
		})

		Context("Stream Parse 0:4 bytes", func() {
			It("should be 34", func() {
				Expect(pack([]byte{0, 0, 0, 34, 15, 4}).Stream()).To(Equal(uint32(34)))
			})
			It("should be 290", func() {
				Expect(pack([]byte{0, 0, 1, 34, 15, 4}).Stream()).To(Equal(uint32(290)))
			})
			It("should be 371261730", func() {
				Expect(pack([]byte{22, 33, 1, 34, 15, 4}).Stream()).To(Equal(uint32(371261730)))
			})
		})

		Context("Service Parse 6:8 bytes", func() {
			It("should be 10", func() {
				Expect(pack([]byte{0, 0, 0, 34, 15, 4, 0, 10}).Service()).To(Equal(uint16(10)))
			})
			It("should be 10", func() {
				Expect(pack([]byte{0, 0, 0, 34, 15, 4, 1, 10}).Service()).To(Equal(uint16(266)))
			})
		})

		Context("Function Parse 5 bytes", func() {
			It("should be 4", func() {
				Expect(pack([]byte{0, 0, 0, 34, 15, 4, 0, 10}).Function()).To(Equal(uint8(4)))
			})
			It("should be 3", func() {
				Expect(pack([]byte{0, 0, 0, 34, 15, 3, 1, 10}).Function()).To(Equal(uint8(3)))
			})
		})

		Context("State Parse 4 byte", func() {
			It("should be 11", func() {
				Expect(pack([]byte{0, 0, 0, 34, 11, 4, 0, 10}).State()).To(Equal(uint8(11)))
			})
			It("should be 33", func() {
				Expect(pack([]byte{0, 0, 0, 34, 33, 4, 1, 10}).State()).To(Equal(uint8(33)))
			})
		})

		Context("Body Parse", func() {
			It("first message 2 bytes body", func() {
				Expect(pack([]byte{0, 0, 0, 34, 15, 4, 0, 10, 32, 34}).Body()).To(Equal([]byte{32, 34}))
			})
			It("first message 4 bytes body", func() {
				Expect(pack([]byte{0, 0, 0, 34, 14, 4, 0, 10, 32, 34}).Body()).To(Equal([]byte{0, 10, 32, 34}))
			})
		})

	})

	Describe("Message clone tests", func() {
		Context("Clone width user id", func() {
			It("User 100", func() {
				item := MessagePackage{Raw: []byte{0, 0, 0, 34, 15, 4, 0, 0, 32, 34}, UserId: 100, Key: "bla"}
				Expect(item.Raw).To(Equal([]byte{0, 0, 0, 34, 15, 4, 0, 0, 32, 34}))
				Expect(item.UserInjectedBody()).To(Equal([]byte{0, 0, 0, 34, 15, 4, 0, 100, 32, 34}))
			})
		})

		Context("Clone to error", func() {
			It("Clone width error message", func() {
				item := MessagePackage{Raw: []byte{1, 0, 0, 34, 15, 4, 0, 0, 32, 34}, UserId: 100, Key: "bla"}
				Expect(item.Raw).To(Equal([]byte{1, 0, 0, 34, 15, 4, 0, 0, 32, 34}))
				bytes := []byte("ERR1")
				header := []byte{1, 0, 0, 34, 13, 4}
				body := append(header, bytes[:]...)
				Expect(item.ErrorResponseBody("ERR1")).To(Equal(body))
			})
		})
	})
})

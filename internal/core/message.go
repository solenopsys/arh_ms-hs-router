package core

import "encoding/binary"

type Direction uint8

const FirstFrame = uint8(15)
const OrdinaryFrame = uint8(0)

type MessagePackage struct {
	raw           []byte
	UserId        uint16
	ConnectionKey string
}

func (m MessagePackage) Stream() uint32 {
	return binary.BigEndian.Uint32(m.raw[:4])
}

func (m MessagePackage) Service() uint16 {
	return binary.BigEndian.Uint16(m.raw[6:8])
}

func (m MessagePackage) State() uint8 {
	return m.raw[4]
}

func (m MessagePackage) Function() uint8 {
	return m.raw[5]
}

func (m MessagePackage) IsFirst() bool {
	return m.State() == FirstFrame
}

func (m MessagePackage) Body() []byte {
	if m.IsFirst() {
		return m.raw[8:]
	} else {
		return m.raw[6:]
	}
}

func (m MessagePackage) userInjectedBody() []byte {
	newBody := make([]byte, len(m.raw))
	copy(newBody, m.raw)
	binary.BigEndian.PutUint16(newBody[6:8], m.UserId)
	return newBody
}

func (m MessagePackage) errorResponseBody(errorKey string) []byte {
	header := make([]byte, 6)
	binary.BigEndian.PutUint32(header[0:4], m.Stream())
	header[4] = StreamError
	header[5] = m.Function()
	body := append(header, []byte(errorKey)...)
	return body
}

package core

import "encoding/binary"

type Direction uint8

const FirstFrame = uint8(15)
const OrdinaryFrame = uint8(0)

type MessageInt interface {
	Stream() uint32
	Service() uint16
	State() uint8
	Function() uint8
	Body() []byte
	IsFirst() bool
	ConnectionKey() string
}

type MessageRoute interface {
	Stream() uint32
	Service() uint16
	ConnectionKey() string
	IsFirst() bool
}

type MessagePackage struct {
	Raw    []byte
	UserId uint16
	Key    string
}

func (m MessagePackage) Stream() uint32 {
	return binary.BigEndian.Uint32(m.Raw[:4])
}

func (m MessagePackage) Service() uint16 {
	return binary.BigEndian.Uint16(m.Raw[6:8])
}

func (m MessagePackage) ConnectionKey() string {
	return m.Key
}

func (m MessagePackage) State() uint8 {
	return m.Raw[4]
}

func (m MessagePackage) Function() uint8 {
	return m.Raw[5]
}

func (m MessagePackage) IsFirst() bool {
	return m.State() == FirstFrame
}

func (m MessagePackage) Body() []byte {
	if m.IsFirst() {
		return m.Raw[8:]
	} else {
		return m.Raw[6:]
	}
}

func (m MessagePackage) UserInjectedBody() []byte {
	newBody := make([]byte, len(m.Raw))
	copy(newBody, m.Raw)
	binary.BigEndian.PutUint16(newBody[6:8], m.UserId)
	return newBody
}

func (m MessagePackage) ErrorResponseBody(errorKey string) []byte {
	header := make([]byte, 6)
	binary.BigEndian.PutUint32(header[0:4], m.Stream())
	header[4] = StreamError
	header[5] = m.Function()
	body := append(header, []byte(errorKey)...)
	return body
}

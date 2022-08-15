package testio

//func procFunc(router *zmq4.Socket, messageHandler MessageHandler, message []byte, fromAddress []byte) {
//var header []byte = message[:8]
//body := message[8:]
//
//var streamId uint32 = binary.BigEndian.Uint32(header[:4])
//var serviceId uint16 = binary.BigEndian.Uint16(header[4:6])
//var functionId uint16 = binary.BigEndian.Uint16(header[6:8])
//
//
//
//
//
//}
//
//func send(){
//	resultBody := messageHandler(body, streamId, serviceId, functionId)
//	ars := [][]byte{header, resultBody}
//	result := bytes.Join(ars, []byte{})
//}

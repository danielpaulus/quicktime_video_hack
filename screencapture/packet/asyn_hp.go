package packet

import "encoding/binary"

//NewAsynHPD0 creates the bytes needed for stopping video streaming
func NewAsynHPD0() []byte {
	length := 20
	data := make([]byte, length)
	binary.LittleEndian.PutUint32(data, uint32(length))
	binary.LittleEndian.PutUint32(data[4:], AsynPacketMagic)
	binary.LittleEndian.PutUint64(data[8:], EmptyCFType)
	binary.LittleEndian.PutUint32(data[16:], HPD0)
	return data
}

//NewAsynHPA0 creates the bytes needed for stopping audio streaming
func NewAsynHPA0(clockRef uint64) []byte {
	length := 20
	data := make([]byte, length)
	binary.LittleEndian.PutUint32(data, uint32(length))
	binary.LittleEndian.PutUint32(data[4:], AsynPacketMagic)
	binary.LittleEndian.PutUint64(data[8:], clockRef)
	binary.LittleEndian.PutUint32(data[16:], HPA0)
	return data
}

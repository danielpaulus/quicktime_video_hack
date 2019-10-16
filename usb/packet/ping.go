package packet

import "encoding/binary"

const (
	PingPacketMagic uint32 = 0x70696E67
	PingLength      uint32 = 16
	PingHeader      uint64 = 0x0000000100000000
)

//NewPingPacketAsBytes generates a new default Ping packet
func NewPingPacketAsBytes() []byte {
	packetBytes := make([]byte, 16)
	binary.LittleEndian.PutUint32(packetBytes, PingLength)
	binary.LittleEndian.PutUint32(packetBytes[4:], PingPacketMagic)
	binary.LittleEndian.PutUint64(packetBytes[8:], PingHeader)
	return packetBytes
}

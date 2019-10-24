package packet

import (
	"encoding/binary"
	"fmt"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
)

//SyncCvrpPacket contains all info from a CVRP packet sent by the device
type SyncCvrpPacket struct {
	ClockRef       CFTypeID
	CorrelationID  uint64
	DeviceClockRef CFTypeID
	Payload        coremedia.StringKeyDict
}

//NewSyncCvrpPacketFromBytes parses a SyncCvrpPacket from a []byte
func NewSyncCvrpPacketFromBytes(data []byte) (SyncCvrpPacket, error) {
	remainingBytes, clockRef, correlationID, err := ParseSyncHeader(data, CVRP)
	if err != nil {
		return SyncCvrpPacket{}, err
	}
	packet := SyncCvrpPacket{ClockRef: clockRef, CorrelationID: correlationID}
	packet.ClockRef = binary.LittleEndian.Uint64(data[4:])
	if packet.ClockRef != EmptyCFType {
		return packet, fmt.Errorf("CVRP packet should have empty CFTypeID for ClockRef but has:%x", packet.ClockRef)
	}

	packet.DeviceClockRef = binary.LittleEndian.Uint64(remainingBytes)

	payloadDict, err := coremedia.NewStringDictFromBytes(remainingBytes[8:])
	if err != nil {
		return packet, err
	}
	packet.Payload = payloadDict
	return packet, nil
}

//NewReply creates a RPLY packet containing the given clockRef and serializes it to a []byte
func (sp SyncCvrpPacket) NewReply(clockRef CFTypeID) []byte {
	return clockRefReply(clockRef, sp.CorrelationID)
}

func (sp SyncCvrpPacket) String() string {
	return fmt.Sprintf("SYNC_CVRP{ClockRef:%x, CorrelationID:%x, DeviceClockRef:%x, Payload:%s}", sp.ClockRef, sp.CorrelationID, sp.DeviceClockRef, sp.Payload.String())
}

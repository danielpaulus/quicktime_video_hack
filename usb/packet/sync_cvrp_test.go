package packet_test

import (
	"github.com/danielpaulus/quicktime_video_hack/usb/packet"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

func TestCvrp(t *testing.T) {
	dat, err := ioutil.ReadFile("fixtures/cvrp-request")
	if err != nil {
		log.Fatal(err)
	}
	cvrp, err := packet.NewSyncCvrpPacketFromBytes(dat[4:])
	if assert.NoError(t, err) {
		assert.Equal(t, 3, len(cvrp.Payload.Entries))
		assert.Equal(t, packet.EmptyCFType, cvrp.ClockRef)
		assert.Equal(t, packet.SyncPacketMagic, cvrp.SyncMagic)
		assert.Equal(t, packet.CVRP, cvrp.MessageType)
		assert.Equal(t, uint64(0x113538da0), cvrp.DeviceClockRef)
		assert.Equal(t, uint64(0x1135659d0), cvrp.CorrelationID)
	}
	testSerializationOfCvrpReply(cvrp, t)
}

func testSerializationOfCvrpReply(cvrp packet.SyncCvrpPacket, t *testing.T) {
	var clockRef packet.CFTypeID = 0x00007FA66CD10250
	replyBytes := cvrp.NewReply(clockRef)
	expectedReplyBytes, err := ioutil.ReadFile("fixtures/cvrp-reply")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, expectedReplyBytes, replyBytes)
}

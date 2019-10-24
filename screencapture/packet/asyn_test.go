package packet_test

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/packet"
	"github.com/stretchr/testify/assert"
)

type asynTestCase struct {
	actualBytes   []byte
	expectedBytes []byte
	description   string
}

func TestAsynPacket(t *testing.T) {
	cases := []asynTestCase{
		{
			actualBytes:   packet.AsynNeedPacketBytes(0x0000000102c16ca0),
			expectedBytes: loadFromFile("asyn-need"),
			description:   "Expect Asyn Need to be correctly serialized",
		},
		{
			actualBytes:   packet.NewAsynHpd1Packet(packet.CreateHpd1DeviceInfoDict()),
			expectedBytes: loadFromFile("asyn-hpd1"),
			description:   "Expect Asyn HPD1 to be correctly serialized",
		},
		{
			actualBytes:   packet.NewAsynHpa1Packet(packet.CreateHpa1DeviceInfoDict(), 0x00000001145392F0),
			expectedBytes: loadFromFile("asyn-hpa1"),
			description:   "Expect Asyn HPA1 to be correctly serialized",
		},
		{
			actualBytes:   packet.NewAsynHPA0(0x0000000102C5FC10),
			expectedBytes: loadFromFile("asyn-hpa0"),
			description:   "Expect Asyn HPA0 to be correctly serialized",
		},
		{
			actualBytes:   packet.NewAsynHPD0(),
			expectedBytes: loadFromFile("asyn-hpd0"),
			description:   "Expect Asyn HPA0 to be correctly serialized",
		},
	}
	for _, testCase := range cases {
		assert.Equal(t, testCase.expectedBytes, testCase.actualBytes, testCase.description)
	}
}

func loadFromFile(name string) []byte {
	dat, err := ioutil.ReadFile("fixtures/" + name)
	if err != nil {
		log.Fatal(err)
	}
	return dat
}

package screencapture_test

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/danielpaulus/quicktime_video_hack/screencapture"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/packet"
	"github.com/stretchr/testify/assert"
)

type UsbTestDummy struct {
	dataReceiver        chan []byte
	cmSampleBufConsumer chan coremedia.CMSampleBuffer
}

func (u UsbTestDummy) Consume(buf coremedia.CMSampleBuffer) error {
	u.cmSampleBufConsumer <- buf
	return nil
}

func (u UsbTestDummy) WriteDataToUsb(data []byte) {
	u.dataReceiver <- data
}

func TestMessageProcessorStopsOnUnknownPacket(t *testing.T) {
	usbDummy := UsbTestDummy{}
	stopChannel := make(chan interface{})
	mp := screencapture.NewMessageProcessor(usbDummy, stopChannel, usbDummy)
	go func() { mp.ReceiveData(make([]byte, 4)) }()
	<-stopChannel
}

type syncTestCase struct {
	receivedData  []byte
	expectedReply [][]byte
	description   string
}

func TestMessageProcessorRespondsCorrectlyToSyncMessages(t *testing.T) {
	clokRequest := loadFromFile("clok-request")[4:]
	parsedClokRequest, _ := packet.NewSyncClokPacketFromBytes(clokRequest)

	cvrpRequest := loadFromFile("cvrp-request")[4:]
	parsedCvrpRequest, _ := packet.NewSyncCvrpPacketFromBytes(cvrpRequest)

	cwpaRequest := loadFromFile("cwpa-request1")[4:]
	parsedCwpaRequest, _ := packet.NewSyncCwpaPacketFromBytes(cwpaRequest)

	cases := []syncTestCase{
		{
			receivedData:  packet.NewPingPacketAsBytes()[4:],
			expectedReply: [][]byte{packet.NewPingPacketAsBytes()},
			description:   "Expect Ping as a response to a ping packet",
		},
		{
			receivedData:  loadFromFile("afmt-request")[4:],
			expectedReply: [][]byte{loadFromFile("afmt-reply")},
			description:   "Expect correct reply for afmt",
		},
		{
			receivedData:  clokRequest,
			expectedReply: [][]byte{parsedClokRequest.NewReply(parsedClokRequest.ClockRef + 0x10000)},
			description:   "Expect correct reply for clok",
		},
		{
			receivedData:  cvrpRequest,
			expectedReply: [][]byte{packet.AsynNeedPacketBytes(parsedCvrpRequest.DeviceClockRef), parsedCvrpRequest.NewReply(parsedCvrpRequest.DeviceClockRef + 0x1000AF)},
			description:   "Expect correct reply for cvrp",
		},
		{
			receivedData: cwpaRequest,
			expectedReply: [][]byte{packet.NewAsynHpd1Packet(packet.CreateHpd1DeviceInfoDict()),
				parsedCwpaRequest.NewReply(parsedCwpaRequest.DeviceClockRef + 1000),
				packet.NewAsynHpd1Packet(packet.CreateHpd1DeviceInfoDict()),
				packet.NewAsynHpa1Packet(packet.CreateHpa1DeviceInfoDict(), parsedCwpaRequest.DeviceClockRef)},
			description: "Expect correct reply for cwpa",
		},
		{
			receivedData:  loadFromFile("og-request")[4:],
			expectedReply: [][]byte{loadFromFile("og-reply")},
			description:   "Expect correct reply for og",
		},
		{
			receivedData:  loadFromFile("stop-request")[4:],
			expectedReply: [][]byte{loadFromFile("stop-reply")},
			description:   "Expect correct reply for stop",
		},
	}

	usbDummy := UsbTestDummy{dataReceiver: make(chan []byte)}
	stopChannel := make(chan interface{})
	mp := screencapture.NewMessageProcessorWithClockBuilder(usbDummy, stopChannel, usbDummy,
		func(ID uint64) coremedia.CMClock { return coremedia.NewCMClockWithHostTime(5) })

	for _, testCase := range cases {
		go func() { mp.ReceiveData(testCase.receivedData) }()
		for _, expectedResponse := range testCase.expectedReply {
			response := <-usbDummy.dataReceiver
			assert.Equal(t, expectedResponse, response, testCase.description)
		}
	}

}

func TestMessageProcessorRespondsCorrectlyToTimeSyncMessages(t *testing.T) {
	timeBytes := loadFromFile("time-request1")[4:]
	timeRequest, err := packet.NewSyncTimePacketFromBytes(timeBytes)
	if err != nil {
		log.Fatal(err)
	}
	testCases := map[string]struct {
		receivedData []byte
		timeRequest  packet.SyncTimePacket
	}{
		"check on time request it sends a reply valid CMTime and correlationID": {timeBytes, timeRequest},
	}

	usbDummy := UsbTestDummy{dataReceiver: make(chan []byte)}
	stopChannel := make(chan interface{})
	mp := screencapture.NewMessageProcessorWithClockBuilder(usbDummy, stopChannel, usbDummy,
		func(ID uint64) coremedia.CMClock { return coremedia.NewCMClockWithHostTime(5) })

	for k, testCase := range testCases {
		go func() { mp.ReceiveData(testCase.receivedData) }()
		response := <-usbDummy.dataReceiver
		fmt.Printf("%x", response)
		assert.Equal(t, uint32(len(response)), binary.LittleEndian.Uint32(response), k)
		assert.Equal(t, packet.ReplyPacketMagic, binary.LittleEndian.Uint32(response[4:]), k)
		assert.Equal(t, testCase.timeRequest.CorrelationID, binary.LittleEndian.Uint64(response[8:]), k)
		_, err := coremedia.NewCMTimeFromBytes(response[16:])
		assert.NoError(t, err)
	}
}

func TestMessageProcessorForwardsFeed(t *testing.T) {
	dat, err := ioutil.ReadFile("packet/fixtures/asyn-feed")
	if err != nil {
		log.Fatal(err)
	}

	usbDummy := UsbTestDummy{dataReceiver: make(chan []byte), cmSampleBufConsumer: make(chan coremedia.CMSampleBuffer)}
	stopChannel := make(chan interface{})
	mp := screencapture.NewMessageProcessor(usbDummy, stopChannel, usbDummy)
	go func() { mp.ReceiveData(dat[4:]) }()
	response := <-usbDummy.cmSampleBufConsumer
	expected := "{OutputPresentationTS:CMTime{95911997690984/1000000000, flags:KCMTimeFlagsHasBeenRounded, epoch:0}, NumSamples:1, Nalus:[{len:30 type:SEI},{len:90712 type:IDR},], fdsc:fdsc:{MediaType:Video, VideoDimension:(1126x2436), Codec:AVC-1, PPS:27640033ac5680470133e69e6e04040404, SPS:28ee3cb0, Extensions:IndexKeyDict:[{49 : IndexKeyDict:[{105 : 0x01640033ffe1001127640033ac5680470133e69e6e0404040401000428ee3cb0fdf8f800},]},{52 : H.264},]}, attach:IndexKeyDict:[{28 : IndexKeyDict:[{46 : Float64[2436.000000]},{47 : Float64[2436.000000]},]},{29 : Int32[0]},{26 : IndexKeyDict:[{46 : Float64[1126.000000]},{47 : Float64[2436.000000]},{45 : Float64[0.000000]},{44 : Float64[0.000000]},]},{27 : IndexKeyDict:[{46 : Float64[1126.000000]},{47 : Float64[2436.000000]},{45 : Float64[0.000000]},{44 : Float64[0.000000]},]},], sary:IndexKeyDict:[{4 : %!s(bool=false)},], SampleTimingInfoArray:{Duration:CMTime{1/60, flags:KCMTimeFlagsHasBeenRounded, epoch:0}, PresentationTS:CMTime{95911997690984/1000000000, flags:KCMTimeFlagsHasBeenRounded, epoch:0}, DecodeTS:CMTime{0/0, flags:KCMTimeFlagsValid, epoch:0}}}"

	assert.Equal(t, expected, response.String())
}

func TestMessageProcessorShutdownMessagesAreCorrect(t *testing.T) {
	usbDummy := UsbTestDummy{dataReceiver: make(chan []byte), cmSampleBufConsumer: make(chan coremedia.CMSampleBuffer)}
	stopChannel := make(chan interface{})
	mp := screencapture.NewMessageProcessor(usbDummy, stopChannel, usbDummy)
	waitCloseSessionChannel := make(chan interface{})

	go func() {
		mp.CloseSession()
		var signal interface{}
		waitCloseSessionChannel <- signal
	}()
	expectedHPA0 := packet.NewAsynHPA0(0x0)
	expectedHPD0 := packet.NewAsynHPD0()
	hpa0 := <-usbDummy.dataReceiver
	hpd0 := <-usbDummy.dataReceiver

	assert.Equal(t, expectedHPA0, hpa0)
	assert.Equal(t, expectedHPD0, hpd0)

	go func() {
		mp.ReceiveData(loadFromFile("asyn-rels")[4:])
		mp.ReceiveData(loadFromFile("asyn-rels")[4:])
	}()

	assert.Equal(t, expectedHPD0, <-usbDummy.dataReceiver)
	<-waitCloseSessionChannel
}

func loadFromFile(name string) []byte {
	dat, err := ioutil.ReadFile("packet/fixtures/" + name)
	if err != nil {
		log.Fatal(err)
	}
	return dat
}

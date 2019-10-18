package rtpsupport

import (
	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	"github.com/pion/rtp"
	"github.com/pion/rtp/codecs"
	"net"
)

type Rtpserver struct {
	packetizer rtp.Packetizer
	clientConn net.Conn
}

func NewRtpServer() Rtpserver {
	//payload type https://docs.microsoft.com/en-us/openspecs/office_protocols/ms-rtp/3b8dc3c6-34b8-4827-9b38-3b00154f471c
	payloader := codecs.H264Payloader{}
	packetizer := rtp.NewPacketizer(1500, 0x60, 5, &payloader, rtp.NewRandomSequencer(), 90000)
	server := Rtpserver{packetizer: packetizer}

	return server
}

func (srv *Rtpserver) StartServerSocket() {
	conn, err := net.Dial("udp", "127.0.0.1:5000")
	if err != nil {
		panic(err)
	}
	srv.clientConn = conn
	// Handle connections in a new goroutine.

}

func (srv Rtpserver) Consume(buf coremedia.CMSampleBuffer) error {
	packets := srv.packetizer.Packetize(buf.SampleData, 1)
	for _, packet := range packets {
		packet.Timestamp = uint32(float64(buf.OutputPresentationTimestamp.CMTimeValue) * 0.00009)
		println(packet.Timestamp)
		data, _ := packet.Marshal()
		srv.clientConn.Write(data)
	}

	return nil
}

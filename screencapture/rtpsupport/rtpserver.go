package rtpsupport

import (
	"encoding/binary"
	"fmt"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	"github.com/pion/rtp"
	"github.com/pion/rtp/codecs"
	"log"
	"net"
)

type Rtpserver struct {
	packetizer rtp.Packetizer
	clientConn net.Conn
	host       string
	port       int
}

func NewRtpServer(host string, port int) *Rtpserver {
	//payload type https://docs.microsoft.com/en-us/openspecs/office_protocols/ms-rtp/3b8dc3c6-34b8-4827-9b38-3b00154f471c
	payloader := codecs.H264Payloader{}
	packetizer := rtp.NewPacketizer(30000, 0x60, 5, &payloader, rtp.NewRandomSequencer(), 90000)
	server := Rtpserver{packetizer: packetizer, host: host, port: port}

	return &server
}

func (srv *Rtpserver) StartServerSocket() {
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", srv.host, srv.port))
	if err != nil {
		panic(err)
	}
	srv.clientConn = conn
	// Handle connections in a new goroutine.
}

func (srv Rtpserver) Consume(buf coremedia.CMSampleBuffer) error {
	if buf.HasFormatDescription {
		err := srv.writeNalu(buf.FormatDescription.PPS, buf)
		if err != nil {
			return err
		}
		err = srv.writeNalu(buf.FormatDescription.SPS, buf)
		if err != nil {
			return err
		}
	}
	srv.writeNalus(buf)

	return nil
}

func (nfw Rtpserver) writeNalus(bytes coremedia.CMSampleBuffer) error {
	slice := bytes.SampleData
	for len(slice) > 0 {
		length := binary.BigEndian.Uint32(slice)
		err := nfw.writeNalu(slice[4:length+4], bytes)
		if err != nil {
			return err
		}
		slice = slice[length+4:]
	}
	return nil
}

func (srv Rtpserver) writeNalu(naluBytes []byte, buf coremedia.CMSampleBuffer) error {
	packets := srv.packetizer.Packetize(naluBytes, 1)
	for _, packet := range packets {
		packet.Timestamp = uint32(float64(buf.OutputPresentationTimestamp.CMTimeValue) * 0.00009)
		//println(packet.Timestamp)
		data, _ := packet.Marshal()
		_,err := srv.clientConn.Write(data)
		if err!=nil{
			log.Fatal("write failed", err)
		}
		//log.Printf("written:%d",n)
	}
	return nil
}

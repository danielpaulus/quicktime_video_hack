package rtpsupport

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	"github.com/pion/rtp"
	"github.com/pion/rtp/codecs"
)

//https://developer.ridgerun.com/wiki/index.php?title=Streaming_RAW_Video_with_GStreamer#Build_udpsrc_for_IMX6
type Rtpserver struct {
	packetizer      rtp.Packetizer
	audioPacketizer rtp.Packetizer
	clientConn      net.Conn
	audioRtpSocket  net.Conn
	host            string
	port            int
	dumpfile        *os.File
}

func NewRtpServer(host string, port int) *Rtpserver {
	file, _ := os.Create("/Users/danielpaulus/tmp/dump-be.bin")
	//payload type https://docs.microsoft.com/en-us/openspecs/office_protocols/ms-rtp/3b8dc3c6-34b8-4827-9b38-3b00154f471c
	payloader := codecs.H264Payloader{}
	packetizer := rtp.NewPacketizer(60000, 0x60, 5, &payloader, rtp.NewRandomSequencer(), 90000)
	audioPayloader := codecs.G711Payloader{}
	server := Rtpserver{packetizer: packetizer, host: host, port: port, dumpfile: file,
		audioPacketizer: rtp.NewPacketizer(8000, 0x96, 6, &audioPayloader, rtp.NewRandomSequencer(), 48000),
	}

	return &server
}

func (srv *Rtpserver) StartServerSocket() {
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", srv.host, srv.port))
	if err != nil {
		log.Warn("Failed connecting to video UDP socket")
	}
	srv.clientConn = conn

	audioConn, err := net.Dial("udp", fmt.Sprintf("%s:%d", srv.host, srv.port+1))
	if err != nil {
		panic(err)
	}
	srv.audioRtpSocket = audioConn
}

func (srv Rtpserver) Consume(buf coremedia.CMSampleBuffer) error {
	if buf.MediaType == coremedia.MediaTypeSound {
		return srv.sendAudioSample(buf)
	}

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

func (srv Rtpserver) sendAudioSample(buf coremedia.CMSampleBuffer) error {
	packets := srv.audioPacketizer.Packetize(buf.SampleData, uint32(buf.NumSamples))
	for _, packet := range packets {
		//packet.Timestamp = uint32(float64(buf.OutputPresentationTimestamp.CMTimeValue) * 0.00009)
		packet.Timestamp = packet.Timestamp
		println(packet.Timestamp)
		packet.PayloadType = 96
		//println(packet.Timestamp)
		data, _ := packet.Marshal()

		srv.dumpfile.Write(data)
		_, err := srv.audioRtpSocket.Write(data)
		if err != nil {
			log.Warn("write failed", err)
		}
		//log.Printf("written:%d",n)
	}
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

		//println(packet.Timestamp)
		data, _ := packet.Marshal()
		_, err := srv.clientConn.Write(data)
		if err != nil {
			log.Fatal("write failed", err)
		}
		//log.Printf("written:%d",n)
	}
	return nil
}

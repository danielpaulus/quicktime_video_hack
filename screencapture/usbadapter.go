package screencapture

import (
	"github.com/google/gousb"
	log "github.com/sirupsen/logrus"
	"time"
)

//UsbAdapter reads and writes from AV Quicktime USB Bulk endpoints
type UsbAdapter struct {
	outEndpoint *gousb.OutEndpoint
}

func (usa UsbAdapter) WriteDataToUsb(bytes []byte) {
	n, err := usa.outEndpoint.Write(bytes)
	if err != nil {
		log.Error("failed sending to usb", err)
	}
	log.Debugf("bytes written:%d", n)
}

//StartReading claims the AV Quicktime USB Bulk endpoints and starts reading until a stopSignal is sent.
//Every received data is added to a frameextractor and when it is complete, sent to the UsbDataReceiver.
func (usa *UsbAdapter) StartReading(device IosDevice, receiver UsbDataReceiver, stopSignal chan interface{}) {
	muxConfig, qtConfig := findConfigurations(device.usbDevice.Desc)
	device.QTConfigIndex = qtConfig
	if qtConfig == -1 {
		log.Error("AV Quicktime USB Config was not enabled")
		return
	}
	device.UsbMuxConfigIndex = muxConfig
	config, err := device.enableQuickTimeConfig()
	defer func() {
		log.Debug("closing Device")
		err := config.Close()
		if err != nil {
			log.Warn("Failed closing device in shutdown", err)
		}
		log.Debug("re-enabling default device config")
		err = device.enableUsbMuxConfig()
		if err != nil {
			log.Error("Failed re-enabling UsbMuxConfig, your device might be broken.", err)
		}
	}()
	if err != nil {
		log.Fatal("Failed enabling Quicktime Device Config. Is Quicktime running on your Machine? If so, close it.", err)
		return
	}

	log.Infof("QT Config is active: %s", config.String())

	//in idx muss sicher der endpoint rein
	duration, _ := time.ParseDuration("20ms")
	device.usbDevice.ControlTimeout = duration
	val, err := device.usbDevice.Control(0x02, 0x01, 0, 0x86, make([]byte, 0))
	if err != nil {
		log.Fatal("failed control", err)
		return
	}
	log.Infof("Clear Feature RC: %d", val)

	val1, err1 := device.usbDevice.Control(0x02, 0x01, 0, 0x05, make([]byte, 0))
	if err1 != nil {
		log.Fatal("failed control", err1)
		return
	}
	log.Infof("Clear Feature RC: %d", val1)

	iface, err := grabQuickTimeInterface(config)
	if err != nil {
		log.Fatal("Couldnt get Quicktime Interface")
		return
	}
	log.Debugf("Got QT iface:%s", iface.String())

	inEndpoint, err := iface.InEndpoint(grabInBulk(iface.Setting))
	if err != nil {
		log.Fatal("couldnt get InEndpoint")
		return
	}
	log.Debugf("Inbound Bulk: %s", inEndpoint.String())

	outEndpoint, err := iface.OutEndpoint(grabOutBulk(iface.Setting))
	if err != nil {
		log.Fatal("couldnt get OutEndpoint")
		return
	}
	log.Debugf("Outbound Bulk: %s", outEndpoint.String())
	usa.outEndpoint = outEndpoint

	stream, err := inEndpoint.NewStream(4096, 5)
	if err != nil {
		log.Fatal("couldnt create stream")
		return
	}
	log.Debug("Endpoint claimed")

	go func() {

		frameExtractor := NewLengthFieldBasedFrameExtractor()
		for {
			buffer := make([]byte, 65536)

			n, err := stream.Read(buffer)
			if err != nil {
				log.Error("couldn't read bytes", err)
				return
			}
			frame, isCompleteFrame := frameExtractor.ExtractFrame(buffer[:n])
			if isCompleteFrame {
				receiver.ReceiveData(frame)
			}
		}
	}()

	<-stopSignal
	receiver.CloseSession()
	log.Info("Closing usb stream")
	device.usbDevice.Reset()
	err = stream.Close()
	if err != nil {
		log.Error("Error closing stream", err)
	}
	log.Info("Closing usb interface")
	iface.Close()

}

func grabOutBulk(setting gousb.InterfaceSetting) int {
	for _, v := range setting.Endpoints {
		if v.Direction == gousb.EndpointDirectionOut {
			return v.Number
		}
	}
	//TODO: error
	return -1
}

func grabInBulk(setting gousb.InterfaceSetting) int {
	for _, v := range setting.Endpoints {
		if v.Direction == gousb.EndpointDirectionIn {
			return v.Number
		}
	}
	//TODO: error
	return -1
}

func grabQuickTimeInterface(config *gousb.Config) (*gousb.Interface, error) {
	_, ifaceIndex := findInterfaceForSubclass(config.Desc, QuicktimeSubclass)
	return config.Interface(ifaceIndex, 0)
}

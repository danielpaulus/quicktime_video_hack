package screencapture

import (
	"errors"

	"github.com/google/gousb"
	log "github.com/sirupsen/logrus"
)

//UsbAdapter reads and writes from AV Quicktime USB Bulk endpoints
type UsbAdapter struct {
	outEndpoint *gousb.OutEndpoint
}

//WriteDataToUsb implements the UsbWriter interface and sends the byte array to the usb bulk endpoint.
func (usa UsbAdapter) WriteDataToUsb(bytes []byte) {
	_, err := usa.outEndpoint.Write(bytes)
	if err != nil {
		log.Error("failed sending to usb", err)
	}
}

//StartReading claims the AV Quicktime USB Bulk endpoints and starts reading until a stopSignal is sent.
//Every received data is added to a frameextractor and when it is complete, sent to the UsbDataReceiver.
func (usa *UsbAdapter) StartReading(device IosDevice, receiver UsbDataReceiver, stopSignal chan interface{}) error {
	ctx = gousb.NewContext()
	defer func() { ctx.Close() }()
	usbDevice, err := ctx.OpenDeviceWithVIDPID(device.VID, device.PID)
	if err != nil {
		return err
	}
	if !device.IsActivated() {
		return errors.New("device not activated for screen mirroring")
	}
	confignum, _ := usbDevice.ActiveConfigNum()

	log.Debugf("Config is active: %d", confignum)

	config, err := usbDevice.Config(confignum)
	if err != nil {
		return errors.New("Could not retrieve config")
	}

	err = sendQTConfigControlRequest(usbDevice)

	if err != nil {
		log.Error("Error disabling config", err)
	}

	log.Infof("QT Config is active: %s", config.String())

	val, err := usbDevice.Control(0x02, 0x01, 0, 0x86, make([]byte, 0))
	if err != nil {
		log.Warn("failed control", err)
	}
	log.Infof("Clear Feature RC: %d", val)

	val1, err1 := usbDevice.Control(0x02, 0x01, 0, 0x05, make([]byte, 0))
	if err1 != nil {
		log.Warn("failed control", err1)
	}
	log.Infof("Clear Feature RC: %d", val1)

	iface, err := grabQuickTimeInterface(config)
	if err != nil {
		log.Error("Couldnt get Quicktime Interface")
		return err
	}
	log.Debugf("Got QT iface:%s", iface.String())

	inEndpoint, err := iface.InEndpoint(grabInBulk(iface.Setting))
	if err != nil {
		log.Error("couldnt get InEndpoint")
		return err
	}
	log.Debugf("Inbound Bulk: %s", inEndpoint.String())

	outEndpoint, err := iface.OutEndpoint(grabOutBulk(iface.Setting))
	if err != nil {
		log.Error("couldnt get OutEndpoint")
		return err
	}
	log.Debugf("Outbound Bulk: %s", outEndpoint.String())
	usa.outEndpoint = outEndpoint

	stream, err := inEndpoint.NewStream(4096, 5)
	if err != nil {
		log.Fatal("couldnt create stream")
		return err
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

	err = stream.Close()
	if err != nil {
		log.Error("Error closing stream", err)
	}
	log.Info("Closing usb interface")
	iface.Close()

	err = sendQTDisableConfigControlRequest(usbDevice)
	if err != nil {
		log.Error("Error sending disable control request", err)
	}
	return nil
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

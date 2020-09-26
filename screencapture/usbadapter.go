package screencapture

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/pkg/errors"

	"github.com/google/gousb"
	log "github.com/sirupsen/logrus"
)

//UsbAdapter reads and writes from AV Quicktime USB Bulk endpoints
type UsbAdapter struct {
	outEndpoint   *gousb.OutEndpoint
	Dump          bool
	DumpOutWriter io.Writer
	DumpInWriter  io.Writer
}

//WriteDataToUsb implements the UsbWriter interface and sends the byte array to the usb bulk endpoint.
func (usbAdapter *UsbAdapter) WriteDataToUsb(bytes []byte) {
	_, err := usbAdapter.outEndpoint.Write(bytes)
	if err != nil {
		log.Error("failed sending to usb", err)
	}
	if usbAdapter.Dump {
		_, err := usbAdapter.DumpOutWriter.Write(bytes)
		if err != nil {
			log.Fatal("Failed dumping data:%v", err)
		}
	}
}

//StartReading claims the AV Quicktime USB Bulk endpoints and starts reading until a stopSignal is sent.
//Every received data is added to a frameextractor and when it is complete, sent to the UsbDataReceiver.
func (usbAdapter *UsbAdapter) StartReading(device IosDevice, receiver UsbDataReceiver, stopSignal chan interface{}) error {
	ctx, cleanUp := createContext()
	defer cleanUp()

	usbDevice, err := OpenDevice(ctx, device)
	if err != nil {
		return err
	}
	if !device.IsActivated() {
		return errors.New("device not activated for screen mirroring")
	}

	confignum, _ := usbDevice.ActiveConfigNum()
	log.Debugf("Config is active: %d, QT config is: %d", confignum, device.QTConfigIndex)

	config, err := usbDevice.Config(device.QTConfigIndex)
	if err != nil {
		return errors.New("Could not retrieve config")
	}

	log.Debugf("QT Config is active: %s", config.String())

	iface, err := findAndClaimQuickTimeInterface(config)
	if err != nil {
		log.Debug("could not get Quicktime Interface")
		return err
	}
	log.Debugf("Got QT iface:%s", iface.String())

	inboundBulkEndpointIndex, inboundBulkEndpointAddress, err := findBulkEndpoint(iface.Setting, gousb.EndpointDirectionIn)
	if err != nil {
		return err
	}

	outboundBulkEndpointIndex, outboundBulkEndpointAddress, err := findBulkEndpoint(iface.Setting, gousb.EndpointDirectionOut)
	if err != nil {
		return err
	}

	err = clearFeature(usbDevice, inboundBulkEndpointAddress, outboundBulkEndpointAddress)
	if err != nil {
		return err
	}

	inEndpoint, err := iface.InEndpoint(inboundBulkEndpointIndex)
	if err != nil {
		log.Error("couldnt get InEndpoint")
		return err
	}
	log.Debugf("Inbound Bulk: %s", inEndpoint.String())

	outEndpoint, err := iface.OutEndpoint(outboundBulkEndpointIndex)
	if err != nil {
		log.Error("couldnt get OutEndpoint")
		return err
	}
	log.Debugf("Outbound Bulk: %s", outEndpoint.String())
	usbAdapter.outEndpoint = outEndpoint

	stream, err := inEndpoint.NewStream(4096, 5)
	if err != nil {
		log.Fatal("couldnt create stream")
		return err
	}
	log.Debug("Endpoint claimed")
	log.Infof("Device '%s' USB connection ready, waiting for ping..", device.SerialNumber)
	go func() {
		lengthBuffer := make([]byte, 4)
		for {

			n, err := io.ReadFull(stream, lengthBuffer)
			if err != nil {
				log.Errorf("Failed reading 4bytes length with err:%s only received: %d", err, n)
				return
			}
			//the 4 bytes header are included in the length, so we need to subtract them
			//here to know how long the payload will be
			length := binary.LittleEndian.Uint32(lengthBuffer) - 4
			dataBuffer := make([]byte, length)

			n, err = io.ReadFull(stream, dataBuffer)
			if err != nil {
				log.Errorf("Failed reading payload with err:%s only received: %d/%d bytes", err, n, length)
				close(stopSignal)
				return
			}
			if usbAdapter.Dump {
				_, err := usbAdapter.DumpInWriter.Write(dataBuffer)
				if err != nil {
					log.Fatal("Failed dumping data:%v", err)
				}
			}
			receiver.ReceiveData(dataBuffer)
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

	sendQTDisableConfigControlRequest(usbDevice)
	log.Debug("Resetting device config")
	_, err = usbDevice.Config(device.UsbMuxConfigIndex)
	if err != nil {
		log.Warn(err)
	}

	return nil
}

func clearFeature(usbDevice *gousb.Device, inboundBulkEndpointAddress gousb.EndpointAddress, outboundBulkEndpointAddress gousb.EndpointAddress) error {
	val, err := usbDevice.Control(0x02, 0x01, 0, uint16(inboundBulkEndpointAddress), make([]byte, 0))
	if err != nil {
		return errors.Wrap(err, "clear feature failed")
	}
	log.Debugf("Clear Feature RC: %d", val)

	val, err = usbDevice.Control(0x02, 0x01, 0, uint16(outboundBulkEndpointAddress), make([]byte, 0))
	log.Debugf("Clear Feature RC: %d", val)
	return errors.Wrap(err, "clear feature failed")
}

func findBulkEndpoint(setting gousb.InterfaceSetting, direction gousb.EndpointDirection) (int, gousb.EndpointAddress, error) {
	for _, v := range setting.Endpoints {
		if v.Direction == direction {
			return v.Number, v.Address, nil

		}
	}
	return 0, 0, errors.New("Inbound Bulkendpoint not found")
}

func findAndClaimQuickTimeInterface(config *gousb.Config) (*gousb.Interface, error) {
	log.Debug("Looking for quicktime interface..")
	found, ifaceIndex := findInterfaceForSubclass(config.Desc, QuicktimeSubclass)
	if !found {
		return nil, fmt.Errorf("did not find interface %v", config)
	}
	log.Debugf("Found Quicktimeinterface: %d", ifaceIndex)
	return config.Interface(ifaceIndex, 0)
}

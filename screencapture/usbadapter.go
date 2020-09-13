package screencapture

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

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

	val, err := usbDevice.Control(0x02, 0x01, 0, 0x86, make([]byte, 0))
	if err != nil {
		log.Debug("failed control", err)
	}
	log.Debugf("Clear Feature RC: %d", val)

	val, err = usbDevice.Control(0x02, 0x01, 0, 0x05, make([]byte, 0))
	if err != nil {
		log.Debug("failed control", err)
	}
	log.Debugf("Clear Feature RC: %d", val)

	iface, err := grabQuickTimeInterface(config)
	if err != nil {
		log.Debug("could not get Quicktime Interface")
		return err
	}
	log.Debugf("Got QT iface:%s", iface.String())

	inboundBulkEndpointIndex, err := grabInBulk(iface.Setting)
	if err != nil {
		return err
	}
	inEndpoint, err := iface.InEndpoint(inboundBulkEndpointIndex)
	if err != nil {
		log.Error("couldnt get InEndpoint")
		return err
	}
	log.Debugf("Inbound Bulk: %s", inEndpoint.String())

	outboundBulkEndpointIndex, err := grabOutBulk(iface.Setting)
	if err != nil {
		return err
	}
	outEndpoint, err := iface.OutEndpoint(outboundBulkEndpointIndex)
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
	log.Infof("Device '%s' USB connection ready, waiting for ping..", device.SerialNumber)
	go func() {
		for {
			buffer := make([]byte, 4)

			n, err := io.ReadFull(stream, buffer)
			if err != nil {
				log.Errorf("Failed reading 4bytes length with err:%s only received: %d", err, n)
				return
			}
			//the 4 bytes header are included in the length, so we need to subtract them
			//here to know how long the payload will be
			length := binary.LittleEndian.Uint32(buffer) - 4
			dataBuffer := make([]byte, length)

			n, err = io.ReadFull(stream, dataBuffer)
			if err != nil {
				log.Errorf("Failed reading payload with err:%s only received: %d/%d bytes", err, n, length)
				return
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

func grabOutBulk(setting gousb.InterfaceSetting) (int, error) {
	for _, v := range setting.Endpoints {
		if v.Direction == gousb.EndpointDirectionOut {
			return v.Number, nil
		}
	}
	return 0, errors.New("Outbound Bulkendpoint not found")
}

func grabInBulk(setting gousb.InterfaceSetting) (int, error) {
	for _, v := range setting.Endpoints {
		if v.Direction == gousb.EndpointDirectionIn {
			return v.Number, nil
		}
	}
	return 0, errors.New("Inbound Bulkendpoint not found")
}

func grabQuickTimeInterface(config *gousb.Config) (*gousb.Interface, error) {
	log.Debug("Looking for quicktime interface..")
	found, ifaceIndex := findInterfaceForSubclass(config.Desc, QuicktimeSubclass)
	if !found {
		return nil, fmt.Errorf("did not find interface %v", config)
	}
	log.Debugf("Found Quicktimeinterface: %d", ifaceIndex)
	return config.Interface(ifaceIndex, 0)
}

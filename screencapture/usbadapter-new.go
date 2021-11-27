package screencapture

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"

	"github.com/google/gousb"
	log "github.com/sirupsen/logrus"
)

//UsbAdapterNew reads and writes from AV Quicktime USB Bulk endpoints
type UsbAdapterNew struct {
	outEndpoint   *gousb.OutEndpoint
	Dump          bool
	DumpOutWriter io.Writer
	DumpInWriter  io.Writer
	stream        *gousb.ReadStream
	usbDevice     *gousb.Device
	contextClose  func()
	iface         *gousb.Interface
	iosDevice     IosDevice
}

//WriteDataToUsb implements the UsbWriter interface and sends the byte array to the usb bulk endpoint.
func (usbAdapter *UsbAdapterNew) WriteDataToUsb(bytes []byte) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
	_, err := usbAdapter.outEndpoint.WriteContext(ctx, bytes)
	if err != nil {
		return err
	}
	if usbAdapter.Dump {
		_, err := usbAdapter.DumpOutWriter.Write(bytes)
		if err != nil {
			log.Fatalf("Failed dumping data:%v", err)
		}
	}
	return nil
}

func (usbAdapter *UsbAdapterNew) InitializeUSB(device IosDevice) error {
	ctx, cleanUp := createContext()
	usbAdapter.contextClose = cleanUp
	usbAdapter.outEndpoint = &gousb.OutEndpoint{}
	usbAdapter.stream = &gousb.ReadStream{}
	usbAdapter.usbDevice = &gousb.Device{}
	usbAdapter.iface = &gousb.Interface{}
	usbAdapter.iosDevice = device

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

	stream, err := inEndpoint.NewStream(4096, 5)
	if err != nil {
		log.Fatal("couldnt create stream")
		return err
	}
	log.Debug("Endpoint claimed")
	log.Debugf("Outbound Bulk: %s", outEndpoint.String())
	usbAdapter.outEndpoint = outEndpoint
	usbAdapter.stream = stream
	usbAdapter.usbDevice = usbDevice
	usbAdapter.iface = iface

	return nil
}

func (usbAdapter *UsbAdapterNew) Close() error {
	log.Info("Closing usb stream")

	err := usbAdapter.stream.Close()
	if err != nil {
		log.Error("Error closing stream", err)
	}
	log.Info("Closing usb interface")
	usbAdapter.iface.Close()

	sendQTDisableConfigControlRequest(usbAdapter.usbDevice)
	log.Debug("Resetting device config")
	_, err = usbAdapter.usbDevice.Config(usbAdapter.iosDevice.UsbMuxConfigIndex)
	if err != nil {
		log.Warn(err)
	}
	usbAdapter.contextClose()
	return nil
}

func (usbAdapter *UsbAdapterNew) ReadFrame() ([]byte, error) {
	lengthBuffer := make([]byte, 4)
	for {
		n, err := io.ReadFull(usbAdapter.stream, lengthBuffer)
		if err != nil {
			return []byte{}, fmt.Errorf("failed reading 4bytes length with err:%s only received: %d", err, n)
		}
		//the 4 bytes header are included in the length, so we need to subtract them
		//here to know how long the payload will be
		length := binary.LittleEndian.Uint32(lengthBuffer) - 4
		dataBuffer := make([]byte, length)

		n, err = io.ReadFull(usbAdapter.stream, dataBuffer)
		if err != nil {
			return []byte{}, err
		}
		if usbAdapter.Dump {
			_, err := usbAdapter.DumpInWriter.Write(dataBuffer)
			if err != nil {
				log.Fatalf("Failed dumping data:%v", err)
			}
		}
		return dataBuffer, nil
	}
}

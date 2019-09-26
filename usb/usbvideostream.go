package usb

import (
	"encoding/hex"
	"github.com/google/gousb"
	log "github.com/sirupsen/logrus"
	"time"
)

// EnableQTConfig enables the hidden QuickTime Device configuration that will expose two new bulk endpoints.
// We will send a control transfer to the device via USB which will cause the device to disconnect and then
// re-connect with a new device configuration. Usually the usbmuxd will automatically enable that new config
// as it will detect it as the device's preferredConfig.
func EnableQTConfig(devices []IosDevice, attachedDevicesChannel chan string) error {
	for _, device := range devices {
		if isValidIosDeviceWithActiveQTConfig(device.usbDevice.Desc) {
			log.Debugf("Skipping %s because it already has an active QT config", device.SerialNumber)
			continue
		}

		var err error = nil
		err = sendQTConfigControlRequest(device)
		if err != nil {
			return err
		}
		//FIXME: For now we assume just one device on the host
		attachedUdid := <-attachedDevicesChannel
		log.Info("Device '%s' reattached", attachedUdid)
	}
	return nil
}

func sendQTConfigControlRequest(device IosDevice) error {
	response := make([]byte, 0)
	val, err := device.usbDevice.Control(0x40, 0x52, 0x00, 0x02, response)

	if err != nil {
		log.Fatal("Failed sending control transfer for enabling hidden QT config", err)
		return err
	}
	log.Debugf("Enabling QT config RC:%d", val)
	return nil
}

func StartReading(device IosDevice) {
	log.Debug("Enabling Quicktime Config for %s", device.SerialNumber)

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
			log.Fatal("Failed re-enabling UsbMuxConfig, your device might be broken.", err)
		}
	}()
	if err != nil {
		log.Fatal("Failed enabling Quicktime Device Config. Is Quicktime running on your Machine? If so, close it.")
		return
	}

	log.Info("QT Config is active: %s", config.String())

	//in idx muss sicher der endpoint rein
	duration, _ := time.ParseDuration("20ms")
	device.usbDevice.ControlTimeout = duration
	val, err := device.usbDevice.Control(0x02, 0x01, 0, 0x86, make([]byte, 0))
	if err != nil {
		log.Fatal("failed control", err)
		return
	}
	log.Infof("Clear Feature RC: %d", val)

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

	stream, err := inEndpoint.NewStream(512, 1)
	if err != nil {
		log.Fatal("couldnt create stream")
		return
	}
	log.Debugf("Endpoint claimed: %s", stream)

	buffer := make([]byte, 512)
	n, err := stream.Read(buffer)
	if err != nil {
		log.Fatal("coudlnt read bytes", err)
		return
	}
	log.Info("read %d bytes:%s", n, hex.Dump(buffer))
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

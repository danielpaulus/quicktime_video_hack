package screencapture

import (
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
		err := enableQTConfigSingleDevice(device, attachedDevicesChannel)
		if err != nil {
			return err
		}
	}
	return nil
}

func enableQTConfigSingleDevice(device IosDevice, attachedDevicesChannel chan string) error {
	if isValidIosDeviceWithActiveQTConfig(device.usbDevice.Desc) {
		log.Debugf("Skipping %s because it already has an active QT config", device.SerialNumber)
		return nil
	}

	err := sendQTConfigControlRequest(device)
	if err != nil {
		return err
	}

	duratio, _ := time.ParseDuration("2s")
	time.Sleep(duratio)
	ctx.Close()
	ctx = gousb.NewContext()
	device.usbDevice, err = findBySerialNumber(device.SerialNumber)

	return err
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

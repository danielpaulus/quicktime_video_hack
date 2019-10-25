package screencapture

import (
	"fmt"
	"github.com/google/gousb"
	log "github.com/sirupsen/logrus"
	"time"
)

// EnableQTConfig enables the hidden QuickTime Device configuration that will expose two new bulk endpoints.
// We will send a control transfer to the device via USB which will cause the device to disconnect and then
// re-connect with a new device configuration. Usually the usbmuxd will automatically enable that new config
// as it will detect it as the device's preferredConfig.
func EnableQTConfig(devices []IosDevice) error {
	for _, device := range devices {
		err := enableQTConfigSingleDevice(device)
		if err != nil {
			return err
		}
	}
	return nil
}

func enableQTConfigSingleDevice(device IosDevice) error {
	udid := device.SerialNumber
	if isValidIosDeviceWithActiveQTConfig(device.usbDevice.Desc) {
		log.Debugf("Skipping %s because it already has an active QT config", udid)
		return nil
	}

	err := sendQTConfigControlRequest(device)
	if err != nil {
		return err
	}

	var i int
	for {
		log.Infof("Checking for active QT config for %s", udid)
		time.Sleep(500 * time.Millisecond)
		err = ctx.Close()
		if err != nil {
			log.Warn("failed closing context", err)
		}
		log.Debug("Reopening Context")
		ctx = gousb.NewContext()
		device.usbDevice, err = findBySerialNumber(udid)
		if err != nil {
			log.Debugf("device not found:%s", err)
			continue
		}
		i++
		if i > 10 {
			log.Error("Failed activating config")
			return fmt.Errorf("could not activate Quicktime Config for %s", udid)
		}
		break
	}
	log.Infof("QTConfig for %s activated", udid)
	return err
}

func sendQTConfigControlRequest(device IosDevice) error {
	response := make([]byte, 0)
	val, err := device.usbDevice.Control(0x40, 0x52, 0x00, 0x02, response)

	if err != nil {
		log.Warn("Failed sending control transfer for enabling hidden QT config. Seems like this happens sometimes but it still works usually.", err)
	}
	log.Debugf("Enabling QT config RC:%d", val)
	return nil
}

func sendQTDisableConfigControlRequest(device IosDevice) error {
	response := make([]byte, 0)
	val, err := device.usbDevice.Control(0x40, 0x52, 0x00, 0x00, response)

	if err != nil {
		log.Fatal("Failed sending control transfer for disabling hidden QT config", err)
		return err
	}
	log.Debugf("Disabled QT config RC:%d", val)
	return nil
}

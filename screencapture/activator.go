package screencapture

import (
	"fmt"
	"time"

	"github.com/google/gousb"
	log "github.com/sirupsen/logrus"
)

// EnableQTConfig enables the hidden QuickTime Device configuration that will expose two new bulk endpoints.
// We will send a control transfer to the device via USB which will cause the device to disconnect and then
// re-connect with a new device configuration. Usually the usbmuxd will automatically enable that new config
// as it will detect it as the device's preferredConfig.
func EnableQTConfig(device IosDevice) (IosDevice, error) {
	udid := device.SerialNumber
	ctx := gousb.NewContext()
	usbDevice, err := ctx.OpenDeviceWithVIDPID(device.VID, device.PID)
	if err != nil {
		return IosDevice{}, err
	}
	if isValidIosDeviceWithActiveQTConfig(usbDevice.Desc) {
		log.Debugf("Skipping %s because it already has an active QT config", udid)
		return device, nil
	}

	sendQTConfigControlRequest(usbDevice)

	var i int
	for {
		log.Debugf("Checking for active QT config for %s", udid)

		err = ctx.Close()
		if err != nil {
			log.Warn("failed closing context", err)
		}
		time.Sleep(500 * time.Millisecond)
		log.Debug("Reopening Context")
		ctx = gousb.NewContext()
		device, err = device.ReOpen(ctx)
		if err != nil {
			log.Debugf("device not found:%s", err)
			continue
		}
		i++
		if i > 10 {
			log.Debug("Failed activating config")
			return IosDevice{}, fmt.Errorf("could not activate Quicktime Config for %s", udid)
		}
		break
	}
	log.Debugf("QTConfig for %s activated", udid)
	return device, err
}

func sendQTConfigControlRequest(device *gousb.Device) {
	response := make([]byte, 0)
	val, err := device.Control(0x40, 0x52, 0x00, 0x02, response)
	if err != nil {
		log.Warnf("Failed sending control transfer for enabling hidden QT config. Seems like this happens sometimes but it still works usually: %s", err)
	}
	log.Debugf("Enabling QT config RC:%d", val)
}

func sendQTDisableConfigControlRequest(device *gousb.Device) {
	response := make([]byte, 0)
	val, err := device.Control(0x40, 0x52, 0x00, 0x00, response)

	if err != nil {
		log.Warnf("Failed sending control transfer for disabling hidden QT config:%s", err)

	}
	log.Debugf("Disabled QT config RC:%d", val)
}

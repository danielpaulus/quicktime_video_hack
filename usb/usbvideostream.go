package usb

import log "github.com/sirupsen/logrus"

func StartReading(device IosDevice) {
	log.Debug("Enabling Quicktime Config for %s", device.SerialNumber)
	_, err := device.enableQuickTimeConfig()
	defer func() {
		err := device.enableUsbMuxConfig()
		log.Fatal("Failed re-enabling UsbMuxConfig, your device might be broken.", err)
	}()
	if err != nil {
		log.Fatal("Failed enabling Quicktime Device Config. Is Quicktime running on your Machine? If so, close it.")
	}

}

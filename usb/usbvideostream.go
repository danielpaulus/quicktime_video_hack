package usb

import log "github.com/sirupsen/logrus"

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

	log.Info("Config is active: %s", config.String())

}

package usb

import (
	"github.com/google/gousb"
	log "github.com/sirupsen/logrus"
)

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

	iface, err := grabQuickTimeInterface(config)
	if err != nil {
		log.Fatal("Couldnt get Quicktime Interface")
		return
	}
	inEndpoint, err := iface.InEndpoint(grabInBulk(iface.Setting))
	if err!=nil{
		log.Fatal("couldnt get InEndpoint")
		return
	}
	stream, err := inEndpoint.NewStream(8, 3)
	if err!=nil{
		log.Fatal("couldnt create stream")
		return
	}
	buffer := make([]byte, 70000)
	n, err := stream.Read(buffer)
	if err!=nil{
		log.Fatal("coudlnt read bytes")
		return
	}
	log.Info("read %d bytes", n)
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

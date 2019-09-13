package main

import (
	"github.com/danielpaulus/quicktime_video_hack/usb"
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
)

func main() {
	usage := `Q.uickTime V.ideo H.ack or qvh client v0.01

Usage:
  qvh devices 
 
Options:
  -h --help     Show this screen.
  --version     Show version.
  -u=<udid>, --udid     UDID of the device.
  -o=<filepath>, --output
  `
	arguments, _ := docopt.ParseDoc(usage)
	devices, _ := arguments.Bool("devices")
	if devices {
		log.Info("iOS Devices with QT Endpoint:")
		devices, err := usb.FindIosDevices()
		if err != nil {
			log.Fatal("Error finding Devices" + error.Error())
		}
		output := usb.PrintDeviceDetails(devices)
		log.Info(output)
	}

}

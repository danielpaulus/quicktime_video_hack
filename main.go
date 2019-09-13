package main

import (
	"github.com/danielpaulus/quicktime_video_hack/usb"
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
)

func main() {
	usage := `Q.uickTime V.ideo H.ack or qvh client v0.01
		If you do not specify a udid, the first device will be taken by default.

Usage:
  qvh devices
 
Options:
  -h --help     Show this screen.
  --version     Show version.
  -u=<udid>, --udid     UDID of the device.
  -o=<filepath>, --output
  `
	arguments, _ := docopt.ParseDoc(usage)

	udid, _ := arguments.String("--udid")
	//TODO:add device selection here
	log.Info(udid)

	devices, err := usb.FindIosDevices()
	if err != nil {
		log.Fatal("Error finding Devices", err)
	}

	devicesCommand, _ := arguments.Bool("devices")
	if devicesCommand {
		log.Info("iOS Devices with QT Endpoint:")

		output := usb.PrintDeviceDetails(devices)
		log.Info(output)
		return
	}

	rawStreamCommand, _ := arguments.Bool("raw")
	if rawStreamCommand {
		dev := devices[0]
		usb.StartReading(dev)
		return
	}
}

package main

import (
	"bufio"
	"github.com/danielpaulus/quicktime_video_hack/screencapture"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

func main() {
	usage := `Q.uickTime V.ideo H.ack or qvh client v0.01
		If you do not specify a udid, the first device will be taken by default.

Usage:
  qvh devices
  qvh activate
  qvh dumpraw <outfile>
 
Options:
  -h --help     Show this screen.
  --version     Show version.
  -u=<udid>, --udid     UDID of the device.
  -o=<filepath>, --output

The commands work as following:
	devices		lists iOS devices attached to this host and tells you if video streaming was activated for them
	activate	only enables the video streaming config for the given device
	dumpraw		will start video recording and dump it to a raw h264 file playable by VLC. 
				Run like: "qvh dumpraw /home/yourname/out.h264"
  `
	arguments, _ := docopt.ParseDoc(usage)
	//TODO: add verbose switch to conf this
	log.SetLevel(log.DebugLevel)
	udid, _ := arguments.String("--udid")
	//TODO:add device selection here
	log.Info(udid)

	devicesCommand, _ := arguments.Bool("devices")
	if devicesCommand {
		devices()
		return
	}

	activateCommand, _ := arguments.Bool("activate")
	if activateCommand {
		activate()
		return
	}

	rawStreamCommand, _ := arguments.Bool("dumpraw")
	if rawStreamCommand {
		outFilePath, err := arguments.String("<outfile>")
		if err != nil {
			log.Error("Missing outfile parameter. Please specify a valid path like '/home/me/out.h264'")
			return
		}
		dumpraw(outFilePath)
	}
}

func waitForSigInt(stopSignalChannel chan interface{}) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Debugf("Signal received: %s", sig)
			var stopSignal interface{}
			stopSignalChannel <- stopSignal
		}
	}()
}

// Just dump a list of what was discovered to the console
func devices() {
	cleanup := screencapture.Init()
	deviceList, err := screencapture.FindIosDevices()
	defer cleanup()
	log.Infof("(%d) iOS Devices with UsbMux Endpoint:", len(deviceList))

	if err != nil {
		log.Fatal("Error finding iOS Devices", err)
	}
	output := screencapture.PrintDeviceDetails(deviceList)
	log.Info(output)
}

// This command is for testing if we can enable the hidden Quicktime device config
func activate() {
	cleanup := screencapture.Init()
	deviceList, err := screencapture.FindIosDevices()
	defer cleanup()
	if err != nil {
		log.Fatal("Error finding iOS Devices", err)
	}

	log.Info("iOS Devices with UsbMux Endpoint:")

	output := screencapture.PrintDeviceDetails(deviceList)
	log.Info(output)

	err = screencapture.EnableQTConfig(deviceList)
	if err != nil {
		log.Fatal("Error enabling QT config", err)
	}

	qtDevices, err := screencapture.FindIosDevicesWithQTEnabled()
	if err != nil {
		log.Fatal("Error finding QT Devices", err)
	}
	qtOutput := screencapture.PrintDeviceDetails(qtDevices)
	if len(qtDevices) != len(deviceList) {
		log.Warnf("Less qt devices (%d) than plain usbmux devices (%d)", len(qtDevices), len(deviceList))
	}
	log.Info("iOS Devices with QT Endpoint:")
	log.Info(qtOutput)
}

func dumpraw(outFilePath string) {
	cleanup := screencapture.Init()
	deviceList, err := screencapture.FindIosDevices()
	defer cleanup()
	if err != nil {
		log.Fatal("Error finding iOS Devices", err)
	}
	log.Infof("Writing output to:%s", outFilePath)
	dev := deviceList[0]

	file, err := os.Create(outFilePath)
	if err != nil {
		log.Debugf("Error creating file:%s", err)
		log.Errorf("Could not open file '%s'", outFilePath)
	}
	writer := coremedia.NewNaluFileWriter(bufio.NewWriter(file))
	adapter := screencapture.UsbAdapter{}
	stopSignal := make(chan interface{})
	waitForSigInt(stopSignal)
	mp := screencapture.NewMessageProcessor(&adapter, stopSignal, writer)

	adapter.StartReading(dev, &mp, stopSignal)
	return
}

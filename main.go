package main

import (
	"github.com/danielpaulus/go-ios/usbmux"
	"github.com/danielpaulus/quicktime_video_hack/usb"
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	usage := `Q.uickTime V.ideo H.ack or qvh client v0.01
		If you do not specify a udid, the first device will be taken by default.

Usage:
  qvh devices
  qvh activate
  qvh dumpraw
 
Options:
  -h --help     Show this screen.
  --version     Show version.
  -u=<udid>, --udid     UDID of the device.
  -o=<filepath>, --output
  `
	arguments, _ := docopt.ParseDoc(usage)
	//TODO: add verbose switch to conf this
	log.SetLevel(log.DebugLevel)
	udid, _ := arguments.String("--udid")
	//TODO:add device selection here
	log.Info(udid)

	cleanup := usb.Init()
	defer cleanup()
	deviceList, err := usb.FindIosDevices()

	if err != nil {
		log.Fatal("Error finding iOS Devices", err)
	}

	devicesCommand, _ := arguments.Bool("devices")
	if devicesCommand {
		devices(deviceList)
		return
	}

	activateCommand, _ := arguments.Bool("activate")
	if activateCommand {
		activate(deviceList)
		return
	}

	rawStreamCommand, _ := arguments.Bool("dumpraw")
	if rawStreamCommand {
		dev := deviceList[0]
		usb.StartReading(dev)
		return
	}
}

// Just dump a list of what was discovered to the console
func devices(devices []usb.IosDevice) {
	log.Infof("(%d) iOS Devices with UsbMux Endpoint:", len(devices))

	output := usb.PrintDeviceDetails(devices)
	log.Info(output)
}

// This command is for testing if we can enable the hidden Quicktime device config
func activate(devices []usb.IosDevice) {
	log.Info("iOS Devices with UsbMux Endpoint:")

	output := usb.PrintDeviceDetails(devices)
	log.Info(output)

	//This channel will get a UDID string whenever a device is connected
	attachedChannel := make(chan string)
	listenForDeviceChanges(attachedChannel)
	err := usb.EnableQTConfig(devices, attachedChannel)
	if err != nil {
		log.Fatal("Error enabling QT config", err)
	}

	qtDevices, err := usb.FindIosDevicesWithQTEnabled()
	if err != nil {
		log.Fatal("Error finding QT Devices", err)
	}
	qtOutput := usb.PrintDeviceDetails(qtDevices)
	if len(qtDevices) != len(devices) {
		log.Warnf("Less qt devices (%d) than plain usbmux devices (%d)", len(qtDevices), len(devices))
	}
	log.Info("iOS Devices with QT Endpoint:")
	log.Info(qtOutput)
}

func listenForDeviceChanges(attachedChannel chan string) {
	muxConnection := usbmux.NewUsbMuxConnection()

	usbmuxDeviceEventReader, err := muxConnection.Listen()
	if err != nil {
		log.Fatal("Failed issuing LISTEN command", err)
		os.Exit(1)
	}

	//read first message and throw away
	_, err = usbmuxDeviceEventReader()
	if err != nil {
		log.Fatalf("error reading from LISTEN command", err)
		os.Exit(1)
	}

	//keep reading attached messages and publish on channel
	go func() {
		defer muxConnection.Close()
		for {
			//the usbmuxDeviceEventReader blocks until a message is received
			msg, err := usbmuxDeviceEventReader()
			if err != nil {
				log.Error("Stopped listening because of error")
				return
			}
			if msg.DeviceAttached() {
				log.Debugf("Received attached message for %s", msg.Properties.SerialNumber)
				attachedChannel <- msg.Properties.SerialNumber
			}
		}
	}()
}

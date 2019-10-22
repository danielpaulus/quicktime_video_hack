package main

import (
	"bufio"
	"os"
	"os/signal"

	"github.com/danielpaulus/quicktime_video_hack/screencapture"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
)

func main() {
	usage := `Q.uickTime V.ideo H.ack or qvh client v0.01
		If you do not specify a udid, the first device will be taken by default.

Usage:
  qvh devices
  qvh activate
  qvh record <h264file> <wavfile>
 
Options:
  -h --help     Show this screen.
  --version     Show version.
  -u=<udid>, --udid     UDID of the device.
  -o=<filepath>, --output

The commands work as following:
	devices		lists iOS devices attached to this host and tells you if video streaming was activated for them
	activate	only enables the video streaming config for the given device
	record		will start video&audio recording. Video will be saved in a raw h264 file playable by VLC.
				Audio will be saved in a uncompressed wav file.
				Run like: "qvh record /home/yourname/out.h264 /home/yourname/out.wav"
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

	rawStreamCommand, _ := arguments.Bool("record")
	if rawStreamCommand {
		h264FilePath, err := arguments.String("<h264file>")
		if err != nil {
			log.Error("Missing <h264file> parameter. Please specify a valid path like '/home/me/out.h264'")
			return
		}
		waveFilePath, err := arguments.String("<wavfile>")
		if err != nil {
			log.Error("Missing <wavfile> parameter. Please specify a valid path like '/home/me/out.raw'")
			return
		}
		record(h264FilePath, waveFilePath)
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

func record(h264FilePath string, wavFilePath string) {
	activate()
	cleanup := screencapture.Init()
	deviceList, err := screencapture.FindIosDevices()
	defer cleanup()
	if err != nil {
		log.Fatal("Error finding iOS Devices", err)
	}
	log.Infof("Writing video output to:'%s' and audio to: %s", h264FilePath, wavFilePath)
	dev := deviceList[0]

	h264File, err := os.Create(h264FilePath)
	if err != nil {
		log.Debugf("Error creating h264File:%s", err)
		log.Errorf("Could not open h264File '%s'", h264FilePath)
	}
	wavFile, err := os.Create(wavFilePath)
	if err != nil {
		log.Debugf("Error creating wav file:%s", err)
		log.Errorf("Could not open wav file '%s'", wavFilePath)
	}

	writer := coremedia.NewAVFileWriter(bufio.NewWriter(h264File), bufio.NewWriter(wavFile))

	defer func() {
		stat, err := wavFile.Stat()
		if err != nil {
			log.Fatal("Could not get wav file stats", err)
		}
		err = coremedia.WriteWavHeader(int(stat.Size()), wavFile)
		if err != nil {
			log.Fatalf("Error writing wave header %s might be invalid. %s", wavFilePath, err.Error())
		}
		err = wavFile.Close()
		if err != nil {
			log.Fatalf("Error closing wave file. '%s' might be invalid. %s", wavFilePath, err.Error())
		}
		err = h264File.Close()
		if err != nil {
			log.Fatalf("Error closing h264File '%s'. %s", h264FilePath, err.Error())
		}

	}()
	adapter := screencapture.UsbAdapter{}
	stopSignal := make(chan interface{})
	waitForSigInt(stopSignal)
	mp := screencapture.NewMessageProcessor(&adapter, stopSignal, writer)

	adapter.StartReading(dev, &mp, stopSignal)
	return
}

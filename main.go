package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	"github.com/danielpaulus/quicktime_video_hack/screencapture"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/gstadapter"
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
)

const version = "v0.2-beta"

func main() {
	usage := fmt.Sprintf(`Q.uickTime V.ideo H.ack (qvh) %s

Usage:
  qvh devices [-v]
  qvh activate [--udid=<udid>] [-v]
  qvh record <h264file> <wavfile> [-v] [--udid=<udid>]
  qvh gstreamer [-v]
  qvh --version | version


Options:
  -h --help       Show this screen.
  -v              Enable verbose mode (debug logging).
  --version       Show version.
  --udid=<udid>   UDID of the device. If not specified, the first found device will be used automatically.

The commands work as following:
	devices		lists iOS devices attached to this host and tells you if video streaming was activated for them
	activate	enables the video streaming config for the device specified by --udid
	record		will start video&audio recording. Video will be saved in a raw h264 file playable by VLC.
				Audio will be saved in a uncompressed wav file.
				Run like: "qvh record /home/yourname/out.h264 /home/yourname/out.wav"
	gstreamer   qvh will open a new window and push AV data to gstreamer.
  `, version)
	arguments, _ := docopt.ParseDoc(usage)
	log.SetFormatter(&log.JSONFormatter{})

	verboseLoggingEnabled, _ := arguments.Bool("-v")
	if verboseLoggingEnabled {
		log.Info("Set Debug mode")
		log.SetLevel(log.DebugLevel)
	}
	shouldPrintVersionNoDashes, _ := arguments.Bool("version")
	shouldPrintVersion, _ := arguments.Bool("--version")
	if shouldPrintVersionNoDashes || shouldPrintVersion {
		printVersion()
		return
	}

	devicesCommand, _ := arguments.Bool("devices")
	if devicesCommand {
		devices()
		return
	}

	udid, _ := arguments.String("--udid")
	log.Debugf("requested udid:'%s'", udid)

	activateCommand, _ := arguments.Bool("activate")
	if activateCommand {
		activate(udid)
		return
	}

	recordCommand, _ := arguments.Bool("record")
	if recordCommand {
		h264FilePath, err := arguments.String("<h264file>")
		if err != nil {
			printErrJSON(err, "Missing <h264file> parameter. Please specify a valid path like '/home/me/out.h264'")
			return
		}
		waveFilePath, err := arguments.String("<wavfile>")
		if err != nil {
			printErrJSON(err, "Missing <wavfile> parameter. Please specify a valid path like '/home/me/out.raw'")
			return
		}
		record(h264FilePath, waveFilePath, udid)
	}
	gstreamerCommand, _ := arguments.Bool("gstreamer")
	if gstreamerCommand {
		startGStreamer(udid)
	}
}

func printVersion() {
	versionMap := map[string]interface{}{
		"version": version,
	}
	printJSON(versionMap)
}

func startGStreamer(udid string) {
	log.Debug("Starting Gstreamer")
	gStreamer := gstadapter.New()
	startWithConsumer(gStreamer, udid)
}

// Just dump a list of what was discovered to the console
func devices() {
	deviceList, err := screencapture.FindIosDevices()
	if err != nil {
		printErrJSON(err, "Error finding iOS Devices")
	}
	log.Debugf("Found (%d) iOS Devices with UsbMux Endpoint", len(deviceList))

	if err != nil {
		printErrJSON(err, "Error finding iOS Devices")
	}
	output := screencapture.PrintDeviceDetails(deviceList)

	printJSON(map[string]interface{}{"devices": output})
}

// This command is for testing if we can enable the hidden Quicktime device config
func activate(udid string) {
	device, err := screencapture.FindIosDevice(udid)
	if err != nil {
		printErrJSON(err, "no device found to activate")
		return
	}

	log.Debugf("Enabling device: %v", device)
	device, err = screencapture.EnableQTConfig(device)
	if err != nil {
		printErrJSON(err, "Error enabling QT config")
		return
	}

	printJSON(map[string]interface{}{
		"device_activated": device.DetailsMap(),
	})
}

func record(h264FilePath string, wavFilePath string, udid string) {
	log.Debugf("Writing video output to:'%s' and audio to: %s", h264FilePath, wavFilePath)

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
	startWithConsumer(writer, udid)
}

func startWithConsumer(consumer screencapture.CmSampleBufConsumer, udid string) {
	device, err := screencapture.FindIosDevice(udid)
	if err != nil {
		printErrJSON(err, "no device found to activate")
		return
	}

	device, err = screencapture.EnableQTConfig(device)
	if err != nil {
		printErrJSON(err, "Error enabling QT config")
		return
	}

	adapter := screencapture.UsbAdapter{}
	stopSignal := make(chan interface{})
	waitForSigInt(stopSignal)
	mp := screencapture.NewMessageProcessor(&adapter, stopSignal, consumer)

	err = adapter.StartReading(device, &mp, stopSignal)
	if err != nil {
		printErrJSON(err, "failed connecting to usb")
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

func printErrJSON(err error, msg string) {
	printJSON(map[string]interface{}{
		"original_error": err.Error(),
		"error_message":  msg,
	})
}
func printJSON(output map[string]interface{}) {
	text, err := json.Marshal(output)
	if err != nil {
		log.Fatalf("Broken json serialization, error: %s", err)
	}
	println(string(text))
}

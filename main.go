package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	stdlog "log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/danielpaulus/go-ios/ios"
	"github.com/danielpaulus/quicktime_video_hack/screencapture"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/diagnostics"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/gstadapter"
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
)

const version = "v0.6-beta"

func main() {
	usage := fmt.Sprintf(`Q.uickTime V.ideo H.ack (qvh) %s

Usage:
  qvh devices [-v]
  qvh activate [--udid=<udid>] [-v]
  qvh deactivate [--udid=<udid>] [-v]
  qvh record <h264file> <wavfile> [--udid=<udid>] [-v]
  qvh audio <outfile> (--mp3 | --ogg | --wav) [--udid=<udid>] [-v]
  qvh gstreamer [--pipeline=<pipeline>] [--examples] [--udid=<udid>] [-v]
  qvh diagnostics <outfile> [--dump=<dumpfile>] [--udid=<udid>]
  qvh --version | version


Options:
  -h --help       Show this screen.
  -v              Enable verbose mode (debug logging).
  --version       Show version.
  --udid=<udid>   UDID of the device. If not specified, the first found device will be used automatically.

The commands work as following:
	devices		lists iOS devices attached to this host and tells you if video streaming was activated for them
	
	activate	enables the video streaming config for the device specified by --udid

	deactivate	disables the video streaming config for the device specified by --udid (in case it is stuck on streaming config)

	record		will start video&audio recording. Video will be saved in a raw h264 file playable by VLC. 
	            	Audio will be saved in a uncompressed wav file. Run like: "qvh record /home/yourname/out.h264 /home/yourname/out.wav"

	audio		Records only audio from the device. It does not change the status bar like the video recording mode does.
			The recorded audio will be saved in <outfile> with the selected format. Currently (--mp3 | --ogg | --wav) are supported.
			Adding more formats is trivial though so create an issue or a PR if you need something :-)

	gstreamer	If no additional param is provided, qvh will open a new window and push AV data to gstreamer.
			If "qvh gstreamer --examples" is provided, qvh will print some common gstreamer pipeline examples.
			If --pipeline is provided, qvh will use the provided gstreamer pipeline instead of 
			displaying audio and video in a window. 

	diagnostics	The diagnostics mode is added for running longterm tests to debug and ensure stability. 
			It will log several metrics and debug logs. Optionally specify a dump file with the --dump option that
			will store raw bytes of all messages. Be aware though, that this file will grow quite large over time. 
  `, version)
	arguments, _ := docopt.ParseDoc(usage)
	log.SetFormatter(&log.JSONFormatter{})

	verboseLoggingEnabled, _ := arguments.Bool("-v")
	if verboseLoggingEnabled {
		log.Info("Set Debug mode")
		log.SetLevel(log.DebugLevel)
	}
	stdlog.SetOutput(new(LogrusWriter))
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
	device, err := findDevice(udid)
	if err != nil {
		printErrJSON(err, "no device found to use")
	}
	checkDeviceIsPaired(device)

	activateCommand, _ := arguments.Bool("activate")
	if activateCommand {
		activate(device)
		return
	}

	deactivateCommand, _ := arguments.Bool("deactivate")
	if deactivateCommand {
		deactivate(device)
		return
	}

	audioCommand, _ := arguments.Bool("audio")
	if audioCommand {
		outfile, err := arguments.String("<outfile>")
		if err != nil {
			printErrJSON(err, "Missing <outfile> parameter. Please specify a valid path like '/home/me/out.h264'")
			return
		}
		log.Infof("Recording audio only to file: %s", outfile)
		mp3, _ := arguments.Bool("--mp3")
		ogg, _ := arguments.Bool("--ogg")
		wav, _ := arguments.Bool("--wav")
		log.Debugf("recording audio only format mp3:%t ogg: %t wav:%t to file: %s", mp3, ogg, wav, outfile)
		if wav {
			recordAudioWav(outfile, device)
			return
		}
		if ogg {
			recordAudioGst(outfile, device, gstadapter.OGG)
			return
		}
		recordAudioGst(outfile, device, gstadapter.MP3)
		return
	}

	diagnostics, _ := arguments.Bool("diagnostics")
	if diagnostics {
		log.SetLevel(log.DebugLevel)
		logfileName := fmt.Sprintf("logfile-%d.log", time.Now().Unix())
		logfile, err := os.OpenFile(logfileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			println("logging to", logfileName, " execute: 'tail -f ", logfileName, "' for logs. Press CTRL+C to stop recording.")
			log.SetOutput(logfile)
		} else {
			log.Info("Failed to log to file, using default stderr")
		}

		outfile, err := arguments.String("<outfile>")
		if err != nil {
			printErrJSON(err, "Missing <outfile> parameter. Please specify a valid path like '/home/me/out.json'")
			return
		}
		dump, _ := arguments.String("--dump")
		runDiagnostics(outfile, dump != "", dump, device)
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
		record(h264FilePath, waveFilePath, device)
	}
	gstreamerCommand, _ := arguments.Bool("gstreamer")
	if gstreamerCommand {
		shouldPrintExamples, _ := arguments.Bool("--examples")
		if shouldPrintExamples {
			printExamples()
			return
		}
		gstPipeline, _ := arguments.String("--pipeline")
		if gstPipeline == "" {
			startGStreamer(device)
			return
		}
		startGStreamerWithCustomPipeline(device, gstPipeline)
	}
}

//findDevice grabs the first device on the host for a empty --udid
//or tries to find the provided device otherwise
func findDevice(udid string) (screencapture.IosDevice, error) {
	if udid == "" {
		return screencapture.FindIosDevice("")
	}
	usbSerial, err := screencapture.ValidateUdid(udid)
	if err != nil {
		return screencapture.IosDevice{}, err
	}
	log.Debugf("requested usb-serial:'%s' from udid:%s", usbSerial, udid)

	return screencapture.FindIosDevice(usbSerial)
}

func printVersion() {
	versionMap := map[string]interface{}{
		"version": version,
	}
	printJSON(versionMap)
}

func printExamples() {

	examples := `Examples:
	
	Writing an MP4 file
	This pipeline will save the recording in video.mp4 with h264 and aac format. The default settings 
	of this pipeline will create a compressed video that takes up way less space than raw h264.
	Note that you need to set "ignore-length" on the wavparse because we are streaming and do not know the length in advance.

	Write MP4 file Mac OSX: 
	vtdec is the hardware accelerated decoder on the mac. 

	qvh gstreamer --pipeline "mp4mux name=mux ! filesink location=video.mp4 \
	queue name=audio_target ! wavparse ignore-length=true ! audioconvert ! faac ! aacparse ! mux. \
	queue name=video_target ! h264parse ! vtdec ! videoconvert ! x264enc  tune=zerolatency !  mux."
	
	Write MP4 file Linux:
    note that I am using software en and decoding, if you have intel VAAPI available, maybe use those. 

	gstreamer --pipeline "mp4mux name=mux ! filesink location=video.mp4 \
    queue name=audio_target ! wavparse ignore-length=true ! audioconvert ! avenc_aac ! aacparse ! mux. \
    queue name=video_target ! h264parse ! avdec_h264 ! videoconvert ! x264enc tune=zerolatency ! mux."
	`
	fmt.Print(examples)
}

func recordAudioGst(outfile string, device screencapture.IosDevice, audiotype string) {
	log.Debug("Starting Gstreamer with audio pipeline")
	gStreamer, err := gstadapter.NewWithAudioPipeline(outfile, audiotype)
	if err != nil {
		printErrJSON(err, "Failed creating custom pipeline")
		return
	}
	startWithConsumer(gStreamer, device, true)
}

func runDiagnostics(outfile string, dump bool, dumpFile string, device screencapture.IosDevice) {
	log.Debugf("diagnostics mode: %s  dump:%t %s device:%s", outfile, dump, dumpFile, device.SerialNumber)
	metricsFile, err := os.Create(outfile)
	if err != nil {
		log.Errorf("Could not open file '%s'", outfile)
	}
	defer metricsFile.Close()
	consumer := diagnostics.NewDiagnosticsConsumer(metricsFile, time.Second*10)
	if dump {
		startWithConsumerDump(consumer, device, dumpFile)
		return
	}
	startWithConsumer(consumer, device, false)
}

func recordAudioWav(outfile string, device screencapture.IosDevice) {
	log.Debug("Starting Gstreamer with audio pipeline")
	wavFile, err := os.Create(outfile)
	if err != nil {
		log.Debugf("Error creating wav file:%s", err)
		log.Errorf("Could not open wav file '%s'", outfile)
	}
	wavFileWriter := coremedia.NewAVFileWriterAudioOnly(wavFile)

	defer func() {
		stat, err := wavFile.Stat()
		if err != nil {
			log.Fatal("Could not get wav file stats", err)
		}
		err = coremedia.WriteWavHeader(int(stat.Size()), wavFile)
		if err != nil {
			log.Fatalf("Error writing wave header %s might be invalid. %s", outfile, err.Error())
		}
		err = wavFile.Close()
		if err != nil {
			log.Fatalf("Error closing wave file. '%s' might be invalid. %s", outfile, err.Error())
		}

	}()
	startWithConsumer(wavFileWriter, device, true)
}

func startGStreamerWithCustomPipeline(device screencapture.IosDevice, pipelineString string) {
	log.Debug("Starting Gstreamer with custom pipeline")
	gStreamer, err := gstadapter.NewWithCustomPipeline(pipelineString)
	if err != nil {
		printErrJSON(err, "Failed creating custom pipeline")
		return
	}
	startWithConsumer(gStreamer, device, false)
}

func startGStreamer(device screencapture.IosDevice) {
	log.Debug("Starting Gstreamer")
	gStreamer := gstadapter.New()
	startWithConsumer(gStreamer, device, false)
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
func activate(device screencapture.IosDevice) {
	log.Debugf("Enabling device: %v", device)
	var err error
	device, err = screencapture.EnableQTConfig(device)
	if err != nil {
		printErrJSON(err, "Error enabling QT config")
		return
	}

	printJSON(map[string]interface{}{
		"device_activated": device.DetailsMap(),
	})
}

func deactivate(device screencapture.IosDevice) {
        log.Debugf("Disabling device: %v", device)
        var err error
        device, err = screencapture.DisableQTConfig(device)
        if err != nil {
                printErrJSON(err, "Error disabling QT config")
                return
        }

        printJSON(map[string]interface{}{
                "device_activated": device.DetailsMap(),
        })
}

func record(h264FilePath string, wavFilePath string, device screencapture.IosDevice) {
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
	startWithConsumer(writer, device, false)
}

func startWithConsumer(consumer screencapture.CmSampleBufConsumer, device screencapture.IosDevice, audioOnly bool) {
	var err error
	device, err = screencapture.EnableQTConfig(device)
	if err != nil {
		printErrJSON(err, "Error enabling QT config")
		return
	}

	adapter := screencapture.UsbAdapter{}
	stopSignal := make(chan interface{})
	waitForSigInt(stopSignal)

	mp := screencapture.NewMessageProcessor(&adapter, stopSignal, consumer, audioOnly)

	err = adapter.StartReading(device, &mp, stopSignal)
	consumer.Stop()
	if err != nil {
		printErrJSON(err, "failed connecting to usb")
	}
}

func startWithConsumerDump(consumer screencapture.CmSampleBufConsumer, device screencapture.IosDevice, dumpPath string) {
	var err error
	device, err = screencapture.EnableQTConfig(device)
	if err != nil {
		printErrJSON(err, "Error enabling QT config")
		return
	}

	inboundMessagesFile, err := os.Create("inbound-" + dumpPath)
	if err != nil {
		log.Fatalf("Could not open file: %v", err)
	}
	defer inboundMessagesFile.Close()
	outboundMessagesFile, err := os.Create("outbound-" + dumpPath)
	if err != nil {
		log.Fatalf("Could not open file: %v", err)
	}
	defer outboundMessagesFile.Close()
	log.Debug("Start dumping all binary transfer")
	adapter := screencapture.UsbAdapter{Dump: true, DumpInWriter: inboundMessagesFile, DumpOutWriter: outboundMessagesFile}
	stopSignal := make(chan interface{})
	waitForSigInt(stopSignal)

	mp := screencapture.NewMessageProcessor(&adapter, stopSignal, consumer, false)

	err = adapter.StartReading(device, &mp, stopSignal)
	consumer.Stop()
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

func checkDeviceIsPaired(device screencapture.IosDevice) {
	dev, err := ios.GetDevice(screencapture.Correct24CharacterSerial(device.SerialNumber))
	if err != nil {
		printErrJSON(err, "device not found, is it still connected?")
		os.Exit(1)
	}
	allValues, err := ios.GetValuesPlist(dev)
	if err != nil {
		printErrJSON(err, "failed getting deviceinfo, you need to pair the device before running qvh")
		os.Exit(1)
	}
	log.Infof("found %s %s for udid %s", allValues["DeviceName"], allValues["ProductVersion"], dev.Properties.SerialNumber)
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

//this is to ban these irritating "2021/04/29 14:27:59 handle_events: error: libusb: interrupted [code -10]" libusb messages
type LogrusWriter int

const interruptedError = "interrupted [code -10]"

func (LogrusWriter) Write(data []byte) (int, error) {
	logmessage := string(data)
	if strings.Contains(logmessage, interruptedError) {
		log.Tracef("gousb_logs:%s", logmessage)
		return len(data), nil
	}
	log.Infof("gousb_logs:%s", logmessage)
	return len(data), nil
}

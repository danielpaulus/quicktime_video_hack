package screencapture

import (
	"bufio"
	"fmt"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

var writerCount = 0

func StartWithConsumer(consumer CmSampleBufConsumer, device IosDevice, audioOnly bool) {
	consumer, vidfile, audiofile := newWriter()
	var err error
	device, err = EnableQTConfig(device)
	if err != nil {
		log.Fatalf("error enabling QT config %v for device %v", err, device)
	}

	usbAdapter := &UsbAdapterNew{}
	err = usbAdapter.InitializeUSB(device)
	if err != nil {
		log.Fatalf("failed initializing usb with error %v for device %v", err, device)
	}

	valeriaInterface := NewValeriaInterface(usbAdapter)
	defer CloseAll(usbAdapter, valeriaInterface)
	go func() {
		err := valeriaInterface.StartReadLoop()
		log.Info("Valeria read loop stopped.")
		if err != nil {
			log.Errorf("Valeria read loop stopped with error %v", err)
		}
	}()
	setupSession(valeriaInterface)

	go func() {
		for {
			buf := valeriaInterface.Local.ReadSampleBuffer()
			go func() {
				err := valeriaInterface.Remote.RequestSampleData()
				if err != nil {
					log.Debug("failed sending need")
					return
				}
			}()
			err = consumer.Consume(buf)
			if err != nil {
				log.Warnf("consumer %v failed to consume buffer %v with error %v", consumer, buf, err)
			}
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGUSR1)
	signal.Notify(c, syscall.SIGUSR2)
	signal.Notify(c, syscall.SIGTERM)
	for {

		sig := <-c
		switch sig {
		case syscall.SIGUSR1:
			//pause
			log.Info("pause")
			valeriaInterface.Remote.StopVideo()
			valeriaInterface.Remote.StopAudio()
            vidfile.Close()
			stat, err := audiofile.Stat()
			if err != nil {
				log.Fatal("Could not get wav file stats", err)
			}
			err = coremedia.WriteWavHeader(int(stat.Size()), audiofile)
			if err != nil {
				log.Fatalf("Error writing wave header %s might be invalid. %s", audiofile, err.Error())
			}
			err = audiofile.Close()
			if err != nil {
				log.Fatalf("Error closing wave file. '%s' might be invalid. %s", audiofile, err.Error())
			}
		case syscall.SIGUSR2:
			//resume
			consumer, vidfile, audiofile = newWriter()
			log.Info("resume")
			valeriaInterface.Remote.EnableAudio()
			valeriaInterface.Remote.EnableVideo()
		default:
			return
		}
	}
}

func StartWithConsumerDump(consumer CmSampleBufConsumer, device IosDevice, dumpPath string) {}

func newWriter() (coremedia.AVFileWriter, *os.File, *os.File) {
	h264FilePath := fmt.Sprintf("video-%03d.h264", writerCount)
	wavFilePath := fmt.Sprintf("audio-%03d.wav", writerCount)
	writerCount++
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

	return coremedia.NewAVFileWriter(bufio.NewWriter(h264File), bufio.NewWriter(wavFile)), h264File, wavFile
}

func setupSession(valeriaInterface ValeriaInterface) {
	err := valeriaInterface.Local.AwaitPing()
	if err != nil {
		log.Errorf("ping timed out failed %v", err)
		return
	}

	log.Info("ping received, responding..")
	err = valeriaInterface.Remote.Ping()
	if err != nil {
		log.Errorf("failed sending Ping %v", err)
		return
	}

	log.Info("handshake complete, awaiting audio clock sync")
	err = valeriaInterface.Local.AwaitAudioClockSync()
	if err != nil {
		log.Errorf("audio clock sync failed %v", err)
		return
	}

	log.Info("audio clock sync ok, enabling video")
	err = valeriaInterface.Remote.EnableVideo()
	if err != nil {
		log.Errorf("failed enabling video %v", err)
		return
	}

	log.Infof("enabling audio")
	err = valeriaInterface.Remote.EnableAudio()
	if err != nil {
		log.Errorf("failed enabling audio %v", err)
		return
	}

	log.Info("awaiting video clock sync")
	err = valeriaInterface.Local.AwaitVideoClockSync()
	if err != nil {
		log.Errorf("failed waiting for video clock sync %v", err)
		return
	}

	log.Info("sending initial sample data request")
	err = valeriaInterface.Remote.RequestSampleData()
	if err != nil {
		log.Errorf("failed requesting sample data %v", err)
		return
	}
}

func CloseAll(usbAdapter *UsbAdapterNew, valeriaInterface ValeriaInterface) {
	log.Info("stopping audio")
	err := valeriaInterface.Remote.StopAudio()
	if err != nil {
		log.Errorf("error stopping audio", err)
	}

	log.Info("stopping video")
	err = valeriaInterface.Remote.StopVideo()
	if err != nil {
		log.Errorf("error stopping video", err)
	}
	log.Info("awaiting audio release")
	err = valeriaInterface.Local.AwaitAudioClockRelease()
	if err != nil {
		log.Errorf("error waiting audio clock release", err)
	}

	log.Info("awaiting video release")
	err = valeriaInterface.Local.AwaitVideoClockRelease()
	if err != nil {
		log.Errorf("error waiting video clock release", err)
	}

	log.Info("shutting down usbadapter")
	err = usbAdapter.Close()
	log.Info("Stream closed successfullly, good bye :-)")
}

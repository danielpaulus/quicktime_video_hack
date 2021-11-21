package screencapture

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

func StartWithConsumer(consumer CmSampleBufConsumer, device IosDevice, audioOnly bool) {
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
			err := valeriaInterface.Remote.RequestSampleData()
			if err != nil {
				log.Debug("failed sending need")
				return
			}
			err = consumer.Consume(buf)
			if err != nil {
				log.Warnf("consumer %v failed to consume buffer %v with error %v", consumer, buf, err)
			}
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
func StartWithConsumerDump(consumer CmSampleBufConsumer, device IosDevice, dumpPath string){}

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

package usb

import "github.com/danielpaulus/quicktime_video_hack/usb/coremedia"

type CmSampleBufConsumer interface {
	Consume(buf coremedia.CMSampleBuffer) error
}

type UsbDataReceiver interface {
	ReceiveData(data []byte)
}

type UsbWriter interface {
	writeDataToUsb(data []byte)
}

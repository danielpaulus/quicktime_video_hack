package screencapture

import "github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"

//CmSampleBufConsumer is a simple interface with one function that consumes CMSampleBuffers
type CmSampleBufConsumer interface {
	Consume(buf coremedia.CMSampleBuffer) error
}

//UsbDataReceiver should receive USB SYNC, ASYN and PING packets with the correct length and with the 4 bytes length removed.
type UsbDataReceiver interface {
	ReceiveData(data []byte)
	CloseSession()
}

//UsbWriter can be used to send data to a USB Endpoint
type UsbWriter interface {
	WriteDataToUsb(data []byte)
}

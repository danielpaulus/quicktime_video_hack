package screencapture_test

import (
	"github.com/danielpaulus/quicktime_video_hack/screencapture"
	"github.com/danielpaulus/quicktime_video_hack/screencapture/coremedia"
	"testing"
)

type UsbTestDummy struct{}

func (u UsbTestDummy) Consume(buf coremedia.CMSampleBuffer) error {
	panic("implement me")
}

func (u UsbTestDummy) WriteDataToUsb(data []byte) {

	panic("implement me")
}

func TestMessageProcessorReturnsWhenStopped(t *testing.T) {
	usbDummy := UsbTestDummy{}
	stopChannel := make(chan interface{})
	mp := screencapture.NewMessageProcessor(usbDummy, stopChannel, usbDummy)
	print(mp)
}

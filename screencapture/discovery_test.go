package screencapture_test

import (
	"testing"

	"github.com/danielpaulus/quicktime_video_hack/screencapture"
	"github.com/stretchr/testify/assert"
)

func TestSerialConvertedToCorrectUdid(t *testing.T) {
	iphoneXrXsStyleSerial := "xxxxxxxxxxxxxxxxxxxxxxxx"
	regularSerial := "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	XrXSStyleDevice := screencapture.IosDevice{SerialNumber: iphoneXrXsStyleSerial}
	regularDevice := screencapture.IosDevice{SerialNumber: regularSerial}
	details := screencapture.PrintDeviceDetails([]screencapture.IosDevice{XrXSStyleDevice, regularDevice})
	assert.Equal(t, 2, len(details))
	assert.Equal(t, "xxxxxxxx-xxxxxxxxxxxxxxxx", details[0]["udid"])
	assert.Equal(t, regularSerial, details[1]["udid"])
}

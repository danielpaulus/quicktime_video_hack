package screencapture_test

import (
	"strings"
	"testing"

	"github.com/danielpaulus/quicktime_video_hack/screencapture"
	"github.com/stretchr/testify/assert"
)

const iphoneXrXsStyleSerial = "xxxxxxxxxxxxxxxxxxxxxxxx\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"
const iphoneXrXsStyleUdid = "xxxxxxxx-xxxxxxxxxxxxxxxx"
const regularSerial = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

func TestSerialConvertedToCorrectUdid(t *testing.T) {
	XrXSStyleDevice := screencapture.IosDevice{SerialNumber: iphoneXrXsStyleSerial}
	regularDevice := screencapture.IosDevice{SerialNumber: regularSerial}
	details := screencapture.PrintDeviceDetails([]screencapture.IosDevice{XrXSStyleDevice, regularDevice})
	assert.Equal(t, 2, len(details))
	assert.Equal(t, iphoneXrXsStyleUdid, details[0]["udid"])
	assert.Equal(t, regularSerial, details[1]["udid"])
}

func TestValidateUdid(t *testing.T) {
	serial, err := screencapture.ValidateUdid(iphoneXrXsStyleUdid)
	assert.NoError(t, err)
	assert.Equal(t, 40, len(serial))
	assert.Equal(t, iphoneXrXsStyleSerial, serial)
	serial, err = screencapture.ValidateUdid(regularSerial)
	assert.NoError(t, err)
	assert.Equal(t, regularSerial, serial)

	_, err = screencapture.ValidateUdid(regularSerial + "toolong")
	assert.Error(t, err)

	_, err = screencapture.ValidateUdid(strings.ReplaceAll(iphoneXrXsStyleUdid, "-", "x"))
	assert.Error(t, err)

}

package packet

import "github.com/danielpaulus/quicktime_video_hack/usb/dict"

type AsynTjmpPacket struct {
	AsyncMagic  uint32
	ClockRef    CFTypeID
	MessageType uint32
	Property    dict.StringKeyEntry
}

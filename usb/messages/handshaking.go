package messages

import (
	"github.com/danielpaulus/quicktime_video_hack/usb/dict"
)

func CreateDeviceInfoDict() dict.StringKeyDict {
	resultDict := dict.StringKeyDict{Entries: make([]dict.StringKeyEntry, 6)}
	resultDict.Entries[0] = dict.StringKeyEntry{
		Key:   "BufferAheadInterval",
		Value: dict.NewNSNumberFromUFloat64(0.07300000000000001),
	}

	resultDict.Entries[1] = dict.StringKeyEntry{
		Key:   "deviceUID",
		Value: "Valeria",
	}

	resultDict.Entries[2] = dict.StringKeyEntry{
		Key:   "ScreenLatency",
		Value: dict.NewNSNumberFromUFloat64(0.04),
	}

	resultDict.Entries[3] = dict.StringKeyEntry{
		Key:   "formats",
		Value: createLpcmInfo(),
	}

	resultDict.Entries[4] = dict.StringKeyEntry{
		Key:   "EDIDAC3Support",
		Value: dict.NewNSNumberFromUInt32(0),
	}

	resultDict.Entries[5] = dict.StringKeyEntry{
		Key:   "deviceName",
		Value: "Valeria",
	}
	return resultDict
}

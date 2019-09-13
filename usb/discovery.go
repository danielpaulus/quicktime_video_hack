package usb

import (
	"errors"
	"github.com/sirupsen/logrus"
)
import "github.com/google/gousb"

type IosDevice struct {
}

const (
	//Interesting, maybe the subclass type Application was chosen intentionally by Apple
	//Because this config enables the basic communication between MacOSX Apps and iOS Devices over USBMuxD
	UsbMuxSubclass gousb.Class = gousb.ClassApplication
	//You can observe this config being activated as soon as you enable iOS Screen Sharing in Quicktime (That's how I
	//found out :-D )
	QuicktimeSubclass gousb.Class = 0x2A
)

func FindIosDevices() ([]IosDevice, error) {
	ctx := gousb.NewContext()
	defer ctx.Close()

	devices, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		// this function is called for every device present.
		// Returning true means the device should be opened.
		b, _, _ := isValidIosDevice(desc)
		return b
	})
	if err != nil {
		return nil, err
	}
	for _, d := range devices {
		serial, err := d.SerialNumber()
		if err != nil {
			logrus.Fatalf("Failed to get Device UDID for '%s'.. what the hell?!", serial)
		}
		product, err:= d.Product()
		if err != nil{
			logrus.Fatalf("Failed to get Device Name for '%s'.. what the hell?!", serial)
		}

		logrus.Infof("'%s'  %s serial: %s",product, d.String(), serial)
	}
	return nil, errors.New("bla")
}

func isValidIosDevice(desc *gousb.DeviceDesc) (bool, int, int) {
	var muxConfigIndex = -1
	var qtConfigIndex = -1

	for k, v := range desc.Configs {
		if isMuxConfig(v) && !isQtConfig(v) {
			muxConfigIndex = k
		}
		if isQtConfig(v) {
			qtConfigIndex = k
		}
	}
	if muxConfigIndex == -1 || qtConfigIndex == -1 {
		return false, 0, 0
	}
	return true, muxConfigIndex, qtConfigIndex
}

func isQtConfig(confDesc gousb.ConfigDesc) bool {
	b, _ := findInterfaceForSubclass(confDesc, QuicktimeSubclass)
	return b
}

func isMuxConfig(confDesc gousb.ConfigDesc) bool {
	b, _ := findInterfaceForSubclass(confDesc, UsbMuxSubclass)
	return b
}

func findInterfaceForSubclass(confDesc gousb.ConfigDesc, subClass gousb.Class) (bool, int) {
	for i := range confDesc.Interfaces {
		//usually those interfaces have only one altsetting
		isVendorClass := confDesc.Interfaces[i].AltSettings[0].Class == gousb.ClassVendorSpec
		isCorrectSubClass := confDesc.Interfaces[i].AltSettings[0].SubClass == subClass
		if isVendorClass && isCorrectSubClass {
			return true, i
		}
	}
	return false, -1
}

func (d *IosDevice) String() string {
	return "IosDevice"
}

func (d *IosDevice) enableUsbMuxConfig() error {
	return nil
}

func (d *IosDevice) enableQuickTimeConfig() error {
	return nil
}

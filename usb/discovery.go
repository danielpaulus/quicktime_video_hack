package usb

import (
	"fmt"
	"strings"
)
import "github.com/google/gousb"

type IosDevice struct {
	usbDevice    *gousb.Device
	SerialNumber string
	ProductName  string
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
	iosDevices, err := mapToIosDevice(devices)
	if err != nil {
		return nil, err
	}

	return iosDevices, nil
}

func mapToIosDevice(devices []*gousb.Device) ([]IosDevice, error) {
	iosDevices := make([]IosDevice, len(devices))
	for i, d := range devices {
		serial, err := d.SerialNumber()
		if err != nil {
			return nil, err
		}
		product, err := d.Product()
		if err != nil {
			return nil, err
		}
		iosDevice := IosDevice{d, serial, product}
		iosDevices[i] = iosDevice
	}
	return iosDevices, nil
}

func PrintDeviceDetails(devices []IosDevice) string {
	var sb strings.Builder
	for _, d := range devices {
		sb.WriteString(fmt.Sprintf("'%s'  %s serial: %s", d.ProductName, d.usbDevice.String(), d.SerialNumber))
	}
	return sb.String()
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
		//usually the interfaces we care about have only one altsetting
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

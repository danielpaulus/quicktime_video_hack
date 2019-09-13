package usb

import (
	"fmt"
	"log"
	"strings"
)
import "github.com/google/gousb"

type IosDevice struct {
	usbDevice         *gousb.Device
	SerialNumber      string
	ProductName       string
	UsbMuxConfigIndex int
	QTConfigIndex     int
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
	defer func() {
		err := ctx.Close()
		if err != nil {
			log.Fatal("Failed while closing usb Context" + err.Error())
		}
	}()

	devices, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		// this function is called for every device present.
		// Returning true means the device should be opened.
		return isValidIosDevice(desc)
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
		muxConfigIndex, qtConfigIndex := findConfigurations(d.Desc)
		iosDevice := IosDevice{d, serial, product, muxConfigIndex, qtConfigIndex}
		iosDevices[i] = iosDevice
	}
	return iosDevices, nil
}

func PrintDeviceDetails(devices []IosDevice) string {
	var sb strings.Builder
	for _, d := range devices {
		sb.WriteString(d.String())
	}
	return sb.String()
}

func isValidIosDevice(desc *gousb.DeviceDesc) bool {
	muxConfigIndex, qtConfigIndex := findConfigurations(desc)
	if muxConfigIndex == -1 || qtConfigIndex == -1 {
		return false
	}
	return true
}

func findConfigurations(desc *gousb.DeviceDesc) (int, int) {
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
	return muxConfigIndex, qtConfigIndex
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
	return fmt.Sprintf("'%s'  %s serial: %s", d.ProductName, d.usbDevice.String(), d.SerialNumber)
}

//Always call this when you're done recording your video to
//put the device back into non-video mode
func (d *IosDevice) enableUsbMuxConfig() error {
	config, err := d.usbDevice.Config(d.UsbMuxConfigIndex)
	if err != nil {
		return err
	}
	return config.Close()
}

//This enables the config needed for grabbing video of the device
//it should open two additional bulk endpoints where video frames
//will be received
func (d *IosDevice) enableQuickTimeConfig() (*gousb.Config, error) {
	config, err := d.usbDevice.Config(d.QTConfigIndex)
	if err != nil {
		return nil, err
	}
	return config, nil
}

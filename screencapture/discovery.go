package screencapture

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)
import "github.com/google/gousb"

//IosDevice contains a gousb.Device pointer for a found device and some additional info like the device udid
type IosDevice struct {
	usbDevice         *gousb.Device
	SerialNumber      string
	ProductName       string
	UsbMuxConfigIndex int
	QTConfigIndex     int
}

const (
	//UsbMuxSubclass is the subclass used for USBMux USB configuration.
	UsbMuxSubclass = gousb.ClassApplication
	//QuicktimeSubclass is the subclass used for the Quicktime USB configuration.
	QuicktimeSubclass gousb.Class = 0x2A
)

// FindIosDevicesWithQTEnabled finds iOS devices connected on USB ports by looking for their
// USBMux compatible Bulk Endpoints and QuickTime Video Stream compatible Bulk Endpoints
func FindIosDevicesWithQTEnabled() ([]IosDevice, error) {
	return findIosDevices(isValidIosDeviceWithActiveQTConfig)
}

// FindIosDevices finds iOS devices connected on USB ports by looking for their
// USBMux compatible Bulk Endpoints
func FindIosDevices() ([]IosDevice, error) {
	return findIosDevices(isValidIosDevice)
}

var ctx *gousb.Context

//Init initializes a new Context and returns a func to close it later.
//Be sure to run it with defer
func Init() func() {
	ctx = gousb.NewContext()
	return func() {
		err := ctx.Close()
		if err != nil {
			log.Fatal("Failed while closing usb Context" + err.Error())
		}
	}
}

func findIosDevices(validDeviceChecker func(desc *gousb.DeviceDesc) bool) ([]IosDevice, error) {
	devices, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		// this function is called for every device present.
		// Returning true means the device should be opened.
		return validDeviceChecker(desc)
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

func findBySerialNumber(udid string) (*gousb.Device, error) {
	devices, err := FindIosDevices()
	if err != nil {
		return nil, err
	}
	for _, d := range devices {
		if d.SerialNumber == udid {
			return d.usbDevice, nil
		}
	}
	return nil, errors.New("not found")
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

//PrintDeviceDetails returns a pretty string for printing device details to the console.
func PrintDeviceDetails(devices []IosDevice) string {
	var sb strings.Builder
	for _, d := range devices {
		sb.WriteString(d.String())
	}
	return sb.String()
}

func isValidIosDevice(desc *gousb.DeviceDesc) bool {
	muxConfigIndex, _ := findConfigurations(desc)
	if muxConfigIndex == -1 {
		return false
	}
	return true
}

func isValidIosDeviceWithActiveQTConfig(desc *gousb.DeviceDesc) bool {
	_, qtConfigIndex := findConfigurations(desc)
	if qtConfigIndex == -1 {
		return false
	}
	return true
}

func findConfigurations(desc *gousb.DeviceDesc) (int, int) {
	var muxConfigIndex = -1
	var qtConfigIndex = -1

	for _, v := range desc.Configs {
		if isMuxConfig(v) && !isQtConfig(v) {
			muxConfigIndex = v.Number
			log.Debugf("Found MuxConfig %d for Device %s", muxConfigIndex, desc.String())
		}
		if isQtConfig(v) {
			qtConfigIndex = v.Number
			log.Debugf("Found QTConfig %d for Device %s", qtConfigIndex, desc.String())
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
			return true, confDesc.Interfaces[i].Number
		}
	}
	return false, -1
}

func (d *IosDevice) String() string {
	return fmt.Sprintf("'%s'  %s serial: %s", d.ProductName, d.usbDevice.String(), d.SerialNumber)
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

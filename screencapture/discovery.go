package screencapture

import (
	"errors"
	"fmt"

	"github.com/google/gousb"
	log "github.com/sirupsen/logrus"
)

//IosDevice contains a gousb.Device pointer for a found device and some additional info like the device udid
type IosDevice struct {
	SerialNumber      string
	ProductName       string
	UsbMuxConfigIndex int
	QTConfigIndex     int
	VID               gousb.ID
	PID               gousb.ID
	UsbInfo           string
}

//ReOpen creates a new Ios device, opening it using VID and PID, using the given context
func (d *IosDevice) ReOpen(ctx *gousb.Context) (IosDevice, error) {
	dev, err := ctx.OpenDeviceWithVIDPID(d.VID, d.PID)
	if err != nil {
		return IosDevice{}, err
	}
	idev, err := mapToIosDevice([]*gousb.Device{dev})
	if err != nil {
		return IosDevice{}, err
	}
	return idev[0], nil
}

const (
	//UsbMuxSubclass is the subclass used for USBMux USB configuration.
	UsbMuxSubclass = gousb.ClassApplication
	//QuicktimeSubclass is the subclass used for the Quicktime USB configuration.
	QuicktimeSubclass gousb.Class = 0x2A
)

// FindIosDevices finds iOS devices connected on USB ports by looking for their
// USBMux compatible Bulk Endpoints
func FindIosDevices() ([]IosDevice, error) {
	ctx, cleanUp := createContext()
	defer cleanUp()
	return findIosDevices(ctx, isValidIosDevice)
}

func createContext() (*gousb.Context, func()) {
	ctx := gousb.NewContext()
	log.Debugf("Opened usbcontext:%v", ctx)
	cleanUp := func() {
		err := ctx.Close()
		if err != nil {
			log.Fatalf("Error closing usb context: %v", ctx)
		}
	}
	return ctx, cleanUp
}

// FindIosDevice finds a iOS device by udid or picks the first one if udid == ""
func FindIosDevice(udid string) (IosDevice, error) {
	ctx, cleanUp := createContext()
	defer cleanUp()
	list, err := findIosDevices(ctx, isValidIosDevice)
	if err != nil {
		return IosDevice{}, err
	}
	if len(list) == 0 {
		return IosDevice{}, errors.New("no iOS devices are connected to this host")
	}
	if udid == "" {
		log.Debugf("no udid specified, using '%s'", list[0].SerialNumber)
		return list[0], nil
	}
	for _, device := range list {
		if udid == device.SerialNumber {
			return device, nil
		}
	}
	return IosDevice{}, fmt.Errorf("device with udid:'%s' not found", udid)
}

func findIosDevices(ctx *gousb.Context, validDeviceChecker func(desc *gousb.DeviceDesc) bool) ([]IosDevice, error) {
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

func mapToIosDevice(devices []*gousb.Device) ([]IosDevice, error) {
	iosDevices := make([]IosDevice, len(devices))
	for i, d := range devices {
		log.Debugf("Getting serial for: %s", d.String())
		serial, err := d.SerialNumber()
		log.Debug("Got serial" + serial)
		if err != nil {
			return nil, err
		}
		product, err := d.Product()
		if err != nil {
			return nil, err
		}

		muxConfigIndex, qtConfigIndex := findConfigurations(d.Desc)
		iosDevice := IosDevice{serial, product, muxConfigIndex, qtConfigIndex, d.Desc.Vendor, d.Desc.Product, d.String()}
		d.Close()
		iosDevices[i] = iosDevice

	}
	return iosDevices, nil
}

//PrintDeviceDetails returns a list of device details ready to be JSON converted.
func PrintDeviceDetails(devices []IosDevice) []map[string]interface{} {
	result := make([]map[string]interface{}, len(devices))
	for k, device := range devices {
		result[k] = device.DetailsMap()
	}
	return result
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

//IsActivated returns a boolean that is true when this device was enabled for screen mirroring and false otherwise.
func (d *IosDevice) IsActivated() bool {
	return d.QTConfigIndex != -1
}

//DetailsMap contains all the info for a device in a map ready to be JSON encoded
func (d *IosDevice) DetailsMap() map[string]interface{} {
	return map[string]interface{}{
		"deviceName":               d.ProductName,
		"usb_device_info":          d.UsbInfo,
		"udid":                     d.SerialNumber,
		"screen_mirroring_enabled": d.IsActivated(),
	}
}

func (d *IosDevice) String() string {
	return fmt.Sprintf("'%s'  %s serial: %s, qt_mode:%t", d.ProductName, d.UsbInfo, d.SerialNumber, d.IsActivated())
}

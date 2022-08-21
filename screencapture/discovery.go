package screencapture

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/gousb"
	log "github.com/sirupsen/logrus"
)

//IosDevice contains a gousb.Device pointer for a found device and some additional info like the device usbSerial
type IosDevice struct {
	SerialNumber      string
	ProductName       string
	UsbMuxConfigIndex int
	QTConfigIndex     int
	VID               gousb.ID
	PID               gousb.ID
	UsbInfo           string
}

//OpenDevice finds a gousb.Device by using the provided iosDevice.SerialNumber. It returns an open device handle.
//Opening using VID and PID is not specific enough, as different iOS devices can have identical VID/PID combinations.
func OpenDevice(ctx *gousb.Context, iosDevice IosDevice) (*gousb.Device, error) {
	deviceList, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return true
	})

	if err != nil {
		log.Warn("Error opening usb devices", err)
	}
	var usbDevice *gousb.Device = nil
	for _, device := range deviceList {
		sn, err := device.SerialNumber()
		if err != nil {
			log.Warn("Error retrieving Serialnumber", err)
		}
		if sn == iosDevice.SerialNumber {
			usbDevice = device
		} else {
			device.Close()
		}
	}

	if usbDevice == nil {
		return nil, fmt.Errorf("Unable to find device:%+v", iosDevice)
	}
	return usbDevice, nil
}

//ReOpen creates a new Ios device, opening it using VID and PID, using the given context
func (d IosDevice) ReOpen(ctx *gousb.Context) (IosDevice, error) {

	dev, err := OpenDevice(ctx, d)
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

// FindIosDevice finds a iOS device by usbSerial or picks the first one if usbSerial == ""
func FindIosDevice(usbSerial string) (IosDevice, error) {
	ctx, cleanUp := createContext()
	defer cleanUp()
	list, err := findIosDevices(ctx, isValidIosDevice)
	if err != nil {
		return IosDevice{}, err
	}
	if len(list) == 0 {
		return IosDevice{}, errors.New("no iOS devices are connected to this host")
	}
	if usbSerial == "" {
		log.Infof("no usbSerial specified, using '%s'", list[0].SerialNumber)
		return list[0], nil
	}
	for _, device := range list {
		if usbSerial == device.SerialNumber {
			return device, nil
		}
	}
	return IosDevice{}, fmt.Errorf("device with usbSerial:'%s' not found", usbSerial)
}

func findIosDevices(ctx *gousb.Context, validDeviceChecker func(desc *gousb.DeviceDesc) bool) ([]IosDevice, error) {
	devices, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		// this function is called for every device present.
		// Returning true means the device should be opened.
		return validDeviceChecker(desc)
	})
	if err != nil {
		log.Warnf("OpenDevices showed some errors, this might be a problem: %v %v", err, devices)
	}
	iosDevices, err := mapToIosDevice(devices)
	if err != nil {
		return nil, fmt.Errorf("mapToIosDevice: %w", err)
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
	log.Infof("descriptor: %+v", desc)
	muxConfigIndex, qtConfig := findConfigurations(desc)
	log.Infof("configs: %d %d", muxConfigIndex, qtConfig)
	if muxConfigIndex == -1 {
		log.Infof("don't open")
		return false
	}
	log.Infof("open")
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
	for _, iface := range confDesc.Interfaces {
		for _, alt := range iface.AltSettings {
			isVendorClass := alt.Class == gousb.ClassVendorSpec
			isCorrectSubClass := alt.SubClass == subClass
			log.Debugf("iface:%v altsettings:%d isvendor:%t isub:%t", iface, len(iface.AltSettings), isVendorClass, isCorrectSubClass)
			if isVendorClass && isCorrectSubClass {
				return true, iface.Number
			}
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
		"udid":                     Correct24CharacterSerial(d.SerialNumber),
		"screen_mirroring_enabled": d.IsActivated(),
	}
}

//Usually iosDevices have a 40 character USB serial which equals the usbSerial used in usbmuxd, Xcode etc.
//There is an exception, some devices like the Xr and Xs have a 24 character USB serial. Usbmux, Xcode etc.
//however insert a dash after the 8th character in this case. To be compatible with other MacOS X and iOS tools,
//we insert the dash here as well.
func Correct24CharacterSerial(usbSerial string) string {
	usbSerial = strings.Trim(usbSerial, "\x00")
	if len(usbSerial) == 24 {
		return fmt.Sprintf("%s-%s", usbSerial[0:8], usbSerial[8:])
	}
	return usbSerial
}

const sixteenTimesZero = "\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"

//ValidateUdid checks if a given udid is 25 or 40 characters long.
//25 character udids must be of format xxxxxxxx-xxxxxxxxxxxxxxxx.
//Serialnumbers on the usb host contain no dashes. As a convenience ValidateUdid
//returns the udid with the dash removed so it can be used
//as a correct USB SerialNumber.
func ValidateUdid(udid string) (string, error) {
	udidLength := len(udid)
	if !(udidLength == 25 || udidLength == 40) {
		return udid, fmt.Errorf("Invalid length for udid:%s UDIDs must have 25 or 40 characters", udid)
	}
	if udidLength == 25 {
		if strings.Index(udid, "-") != 8 {
			return udid, fmt.Errorf("Invalid format for udid:%s 25 char UDIDs must contain a dash at position 8", udid)
		}
		removedDash := strings.Replace(udid, "-", "", 1)

		return removedDash + sixteenTimesZero, nil
	}
	return udid, nil
}

func (d *IosDevice) String() string {
	return fmt.Sprintf("'%s'  %s serial: %s, qt_mode:%t", d.ProductName, d.UsbInfo, d.SerialNumber, d.IsActivated())
}

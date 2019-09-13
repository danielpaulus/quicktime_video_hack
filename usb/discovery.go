package usb

import "errors"

type IosDevice struct {
}

func FindIosDevices() ([]IosDevice, error) {
	return nil, errors.New("bla")
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

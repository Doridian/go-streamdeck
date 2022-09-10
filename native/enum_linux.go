package native

import (
	"github.com/KarpelesLab/hid"
)

type linuxHIDEnumerator struct{}

func (e *linuxHIDEnumerator) Enumerate(vendorFilter uint16, productFilter uint16) ([]HIDDevice, error) {
	devArr := make([]HIDDevice, 0)

	hid.UsbWalk(func(dev hid.Device) {
		info := dev.Info()
		if vendorFilter > 0 && info.Vendor != vendorFilter {
			return
		}
		if productFilter > 0 && info.Product != productFilter {
			return
		}

		hidDev := &linuxHIDDevice{
			info: info,
			dev:  dev,
		}

		devArr = append(devArr, hidDev)
	})

	return devArr, nil
}

func NativeHIDEnumerator() HIDEnumerator {
	return &linuxHIDEnumerator{}
}

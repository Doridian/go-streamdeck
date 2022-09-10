package native

import "time"

type HIDDeviceHandle interface {
	GetFeatureReport(id byte) ([]byte, error)
	SetFeatureReport(payload []byte) error

	Read(timeout time.Duration) ([]byte, error)
	Write(packet []byte, timeout time.Duration) error

	Close() error
}

type HIDDevice interface {
	Vendor() uint16
	Product() uint16
	Open() (HIDDeviceHandle, error)
}

type HIDEnumerator interface {
	Enumerate(vendorFilter uint16, productFilter uint16) ([]HIDDevice, error)
}

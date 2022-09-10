// go:build darwin,!ios,cgo windows,cgo

package native

import (
	"time"

	"github.com/karalabe/hid"
)

type libHIDDeviceHandle struct {
	hdl *hid.Device
}

func (a *libHIDDeviceHandle) GetFeatureReport(id byte) ([]byte, error) {
	report := make([]byte, 256)
	report[0] = id
	_, err := a.hdl.GetFeatureReport(report)
	return report, err
}

func (a *libHIDDeviceHandle) SetFeatureReport(payload []byte) error {
	_, err := a.hdl.SendFeatureReport(payload)
	return err
}

func (a *libHIDDeviceHandle) Read(timeout time.Duration) ([]byte, error) {
	res := make([]byte, 256)
	n, err := a.hdl.Read(res)
	if err != nil {
		return nil, err
	}
	return res[:n], nil
}

func (a *libHIDDeviceHandle) Write(packet []byte, timeout time.Duration) error {
	_, err := a.hdl.Write(packet)
	return err
}

func (a *libHIDDeviceHandle) Close() error {
	return a.hdl.Close()
}

type libHIDDevice struct {
	info hid.DeviceInfo
}

func (d *libHIDDevice) Vendor() uint16 {
	return d.info.VendorID
}

func (d *libHIDDevice) Product() uint16 {
	return d.info.ProductID
}

func (d *libHIDDevice) Open() (HIDDeviceHandle, error) {
	hdl, err := d.info.Open()
	if err != nil {
		return nil, err
	}
	return &libHIDDeviceHandle{
		hdl: hdl,
	}, nil
}

type libHIDenumerator struct{}

func (e *libHIDenumerator) Enumerate(vendorFilter uint16, productFilter uint16) ([]HIDDevice, error) {
	nativeDevs := hid.Enumerate(vendorFilter, productFilter)
	devArr := make([]HIDDevice, 0, len(nativeDevs))

	for _, info := range nativeDevs {
		devArr = append(devArr, &libHIDDevice{
			info: info,
		})
	}

	return devArr, nil
}

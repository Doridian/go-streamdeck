package native

import (
	"time"

	"github.com/KarpelesLab/hid"
)

type linuxHIDDeviceHandle struct {
	hdl hid.Handle
}

func (h *linuxHIDDeviceHandle) GetFeatureReport(id byte) ([]byte, error) {
	return h.hdl.GetFeatureReport(int(id))
}

func (h *linuxHIDDeviceHandle) SetFeatureReport(payload []byte) error {
	// Yes, do NOT skip first byte of payload, even though it is the ID!
	// Exception: The ID is 0, in which case the device does not use numbered reports
	// So we have to account for that...
	// See: https://github.com/libusb/hidapi/blob/0f2cf886e5a91ff13f222a67d8931c0822244c9a/libusb/hid.c#L1433
	if payload[0] == 0 {
		return h.hdl.SetFeatureReport(0, payload[1:])
	}
	return h.hdl.SetFeatureReport(int(payload[0]), payload)
}

func (a *linuxHIDDeviceHandle) Read(timeout time.Duration) ([]byte, error) {
	return a.hdl.ReadInputPacket(timeout)
}

func (a *linuxHIDDeviceHandle) Write(packet []byte, timeout time.Duration) error {
	_, err := a.hdl.Write(packet, timeout)
	return err
}

func (h *linuxHIDDeviceHandle) Close() error {
	return h.hdl.Close()
}

type linuxHIDDevice struct {
	dev  hid.Device
	info hid.Info
}

func (d *linuxHIDDevice) Vendor() uint16 {
	return d.info.Vendor
}

func (d *linuxHIDDevice) Product() uint16 {
	return d.info.Product
}

func (d *linuxHIDDevice) Open() (HIDDeviceHandle, error) {
	hdl, err := d.dev.Open()
	if err != nil {
		return nil, err
	}
	return &linuxHIDDeviceHandle{
		hdl: hdl,
	}, nil
}

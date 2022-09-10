package streamdeck

import (
	"image"
	"time"

	"github.com/Doridian/go-streamdeck/native"
)

const usbTimeout = time.Second * 1

// Stream Deck Vendor & Product IDs.
//
//nolint:revive
const (
	VID_ELGATO              = 0x0fd9
	PID_STREAMDECK          = 0x0060
	PID_STREAMDECK_V2       = 0x006d
	PID_STREAMDECK_MK2      = 0x0080
	PID_STREAMDECK_MINI     = 0x0063
	PID_STREAMDECK_MINI_MK2 = 0x0090
	PID_STREAMDECK_XL       = 0x006c
)

// Firmware command IDs.
//
//nolint:revive
var (
	c_REV1_FIRMWARE   byte = 0x04
	c_REV1_RESET           = []byte{0x0b, 0x63}
	c_REV1_BRIGHTNESS      = []byte{0x05, 0x55, 0xaa, 0xd1, 0x01}

	c_REV2_FIRMWARE   byte = 0x05
	c_REV2_RESET           = []byte{0x03, 0x02}
	c_REV2_BRIGHTNESS      = []byte{0x03, 0x08}
)

// Device represents a single Stream Deck device.
type Device struct {
	Columns uint8
	Rows    uint8
	Keys    uint8
	Pixels  uint
	DPI     uint
	Padding uint

	featureReportSize   int
	firmwareOffset      int
	keyStateOffset      int
	translateKeyIndex   func(index, columns uint8) uint8
	imagePageSize       int
	imagePageHeaderSize int
	flipImage           func(image.Image) image.Image
	toImageFormat       func(image.Image) ([]byte, error)
	imagePageHeader     func(pageIndex int, keyIndex uint8, payloadLength int, lastPage bool) []byte

	getFirmwareCommand   byte
	resetCommand         []byte
	setBrightnessCommand []byte

	keyStateLegth int

	handle native.HIDDeviceHandle
	device native.HIDDevice
}

// Key holds the current status of a key on the device.
type Key struct {
	Index   uint8
	Pressed bool
}

// Devices returns all attached Stream Decks.
func Devices() ([]Device, error) {
	dd := []Device{}

	enum := native.NativeHIDEnumerator()
	devs, err := enum.Enumerate(VID_ELGATO, 0)
	if err != nil {
		return dd, err
	}

	for _, hidDev := range devs {
		pid := hidDev.Product()
		var dev Device

		switch pid {
		case PID_STREAMDECK:
			dev = Device{
				Columns:              5,
				Rows:                 3,
				Keys:                 15,
				Pixels:               72,
				DPI:                  124,
				Padding:              16,
				featureReportSize:    17,
				firmwareOffset:       5,
				keyStateOffset:       1,
				translateKeyIndex:    translateRightToLeft,
				imagePageSize:        7819,
				imagePageHeaderSize:  16,
				imagePageHeader:      rev1ImagePageHeader,
				flipImage:            flipHorizontally,
				toImageFormat:        toBMP,
				getFirmwareCommand:   c_REV1_FIRMWARE,
				resetCommand:         c_REV1_RESET,
				setBrightnessCommand: c_REV1_BRIGHTNESS,
			}
		case PID_STREAMDECK_MINI:
		case PID_STREAMDECK_MINI_MK2:
			dev = Device{
				Columns:              3,
				Rows:                 2,
				Keys:                 6,
				Pixels:               80,
				DPI:                  138,
				Padding:              16,
				featureReportSize:    17,
				firmwareOffset:       5,
				keyStateOffset:       1,
				translateKeyIndex:    identity,
				imagePageSize:        1024,
				imagePageHeaderSize:  16,
				imagePageHeader:      miniImagePageHeader,
				flipImage:            rotateCounterclockwise,
				toImageFormat:        toBMP,
				getFirmwareCommand:   c_REV1_FIRMWARE,
				resetCommand:         c_REV1_RESET,
				setBrightnessCommand: c_REV1_BRIGHTNESS,
			}
		case PID_STREAMDECK_V2:
		case PID_STREAMDECK_MK2:
			dev = Device{
				Columns:              5,
				Rows:                 3,
				Keys:                 15,
				Pixels:               72,
				DPI:                  124,
				Padding:              16,
				featureReportSize:    32,
				firmwareOffset:       6,
				keyStateOffset:       4,
				translateKeyIndex:    identity,
				imagePageSize:        1024,
				imagePageHeaderSize:  8,
				imagePageHeader:      rev2ImagePageHeader,
				flipImage:            flipHorizontallyAndVertically,
				toImageFormat:        toJPEG,
				getFirmwareCommand:   c_REV2_FIRMWARE,
				resetCommand:         c_REV2_RESET,
				setBrightnessCommand: c_REV2_BRIGHTNESS,
			}
		case PID_STREAMDECK_XL:
			dev = Device{
				Columns:              8,
				Rows:                 4,
				Keys:                 32,
				Pixels:               96,
				DPI:                  166,
				Padding:              16,
				featureReportSize:    32,
				firmwareOffset:       6,
				keyStateOffset:       4,
				translateKeyIndex:    identity,
				imagePageSize:        1024,
				imagePageHeaderSize:  8,
				imagePageHeader:      rev2ImagePageHeader,
				flipImage:            flipHorizontallyAndVertically,
				toImageFormat:        toJPEG,
				getFirmwareCommand:   c_REV2_FIRMWARE,
				resetCommand:         c_REV2_RESET,
				setBrightnessCommand: c_REV2_BRIGHTNESS,
			}
		}

		if dev.Columns > 0 {
			dev.keyStateLegth = int(dev.Columns) * int(dev.Rows)
			dev.device = hidDev
			dd = append(dd, dev)
		}
	}

	return dd, nil
}

// Open the device for input/output. This must be called before trying to
// communicate with the device.
func (d *Device) Open() error {
	var err error
	d.handle, err = d.device.Open()
	return err
}

// Close the connection with the device.
func (d *Device) Close() error {
	if d.handle != nil {
		return d.handle.Close()
	}
	return nil
}

// FirmwareVersion returns the firmware version of the device.
func (d *Device) FirmwareVersion() (string, error) {
	result, err := d.getFeatureReport(d.getFirmwareCommand)
	if err != nil {
		return "", err
	}
	return string(result[d.firmwareOffset:]), nil
}

// Resets the Stream Deck, clears all button images and shows the standby image.
func (d *Device) Reset() error {
	return d.setFeatureReport(d.resetCommand)
}

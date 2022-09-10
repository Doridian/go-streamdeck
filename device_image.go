package streamdeck

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
)

// SetImage sets the image of a button on the Stream Deck. The provided image
// needs to be in the correct resolution for the device. The index starts with
// 0 being the top-left button.
func (d *Device) SetImage(index uint8, img image.Image) error {
	imageData, err := d.ConvertImage(img)
	if err != nil {
		return err
	}

	return d.SetConvertedImage(index, imageData)
}

func (d *Device) SetConvertedImage(index uint8, imageData *ImageData) error {
	data := make([]byte, d.imagePageSize)
	translatedIndex := d.translateKeyIndex(index, d.Columns)

	var page int
	var lastPage bool
	for !lastPage {
		var payload []byte
		payload, lastPage = imageData.page(page)
		header := d.imagePageHeader(page, translatedIndex, len(payload), lastPage)

		copy(data, header)
		copy(data[len(header):], payload)

		err := d.handle.Write(data, usbTimeout)
		if err != nil {
			return fmt.Errorf("cannot write image page %d (%d bytes): %v",
				page, len(data), err)
		}

		page++
	}

	return nil
}

func (d *Device) ConvertImage(img image.Image) (*ImageData, error) {
	if img.Bounds().Dy() != int(d.Pixels) ||
		img.Bounds().Dx() != int(d.Pixels) {
		return nil, fmt.Errorf("supplied image has wrong dimensions, expected %[1]dx%[1]d pixels", d.Pixels)
	}

	imageBytes, err := d.toImageFormat(d.flipImage(img))
	if err != nil {
		return nil, fmt.Errorf("cannot convert image data: %v", err)
	}

	pageSize := d.imagePageSize - d.imagePageHeaderSize

	pageCount := len(imageBytes) / pageSize
	if len(imageBytes)%pageSize != 0 {
		pageCount++
	}

	return &ImageData{
		image:     imageBytes,
		pageSize:  pageSize,
		pageCount: pageSize,
	}, nil
}

// Clears the Stream Deck, setting a black image on all buttons.
func (d *Device) Clear() error {
	img := image.NewRGBA(image.Rect(0, 0, int(d.Pixels), int(d.Pixels)))
	draw.Draw(img, img.Bounds(), image.NewUniform(color.RGBA{0, 0, 0, 255}), image.Point{}, draw.Src)
	for i := uint8(0); i <= d.Columns*d.Rows; i++ {
		err := d.SetImage(i, img)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

package streamdeck

// ImageData allows to access raw image data in a byte array through pages of a
// given size.
type ImageData struct {
	image     []byte
	pageSize  int
	pageCount int
}

// page returns the page with the given index and an indication if this is the
// last page.
func (d ImageData) page(pageIndex int) ([]byte, bool) {
	offset := pageIndex * d.pageSize
	length := d.pageSize

	if offset+length > len(d.image) {
		length = len(d.image) - offset
	}

	if length <= 0 {
		return []byte{}, true
	}

	return d.image[offset : offset+length], pageIndex == d.pageCount-1
}

package streamdeck

// getFeatureReport from the device without worries about the correct payload
// size.
func (d *Device) getFeatureReport(id int) ([]byte, error) {
	b, err := d.handle.GetFeatureReport(id)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// sendFeatureReport to the device without worries about the correct payload
// size.
func (d *Device) sendFeatureReport(payload []byte) error {
	b := make([]byte, d.featureReportSize-1)
	copy(b, payload[1:])
	return d.handle.SetFeatureReport(int(payload[0]), payload)
}

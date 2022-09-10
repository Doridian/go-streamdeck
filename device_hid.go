package streamdeck

// getFeatureReport from the device without worries about the correct payload size.
func (d *Device) getFeatureReport(id byte) ([]byte, error) {
	return d.handle.GetFeatureReport(id)
}

// sendFeatureReport to the device without worries about the correct payload size.
func (d *Device) sendFeatureReport(payload []byte) error {
	b := make([]byte, d.featureReportSize-1)
	copy(b, payload[1:])
	return d.handle.SetFeatureReport(payload)
}

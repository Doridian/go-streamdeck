package streamdeck

// ReadKeys returns a channel, which it will use to emit key presses/releases.
func (d *Device) ReadKeys() (chan Key, error) {
	kch := make(chan Key)

	keyBufferLen := d.keyStateOffset + d.keyStateLegth
	oldKeyBuffer := make([]byte, keyBufferLen)

	go func() {
		for {
			keyBuffer, err := d.handle.Read(usbTimeout)
			if err != nil {
				if isErrorTimeout(err) {
					continue
				}
				close(kch)
				return
			}

			if len(keyBuffer) < keyBufferLen {
				continue
			}

			for i := d.keyStateOffset; i < keyBufferLen; i++ {
				if keyBuffer[i] != oldKeyBuffer[i] {
					keyIndex := uint8(i - d.keyStateOffset)
					kch <- Key{
						Index:   d.translateKeyIndex(keyIndex, d.Columns),
						Pressed: keyBuffer[i] == 1,
					}
				}
			}

			oldKeyBuffer = keyBuffer
		}
	}()

	return kch, nil
}

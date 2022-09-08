package streamdeck

import "time"

// ReadKeys returns a channel, which it will use to emit key presses/releases.
func (d *Device) ReadKeys() (chan Key, error) {
	kch := make(chan Key)

	keyBufferLen := d.keyStateOffset + d.keyStateLegth
	oldKeyBuffer := make([]byte, keyBufferLen)

	go func() {
		for {
			keyBuffer, err := d.handle.ReadInputPacket(usbTimeout)
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

			// don't trigger a key event if the device is asleep, but wake it
			if d.asleep {
				_ = d.Wake()
				// reset state so no spurious key events get triggered
				oldKeyBuffer = make([]byte, keyBufferLen)
				continue
			}

			d.sleepMutex.Lock()
			d.lastActionTime = time.Now()
			d.sleepMutex.Unlock()

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

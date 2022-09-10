package native

func NativeHIDEnumerator() HIDEnumerator {
	return &libHIDenumerator{}
}

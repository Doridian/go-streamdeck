package streamdeck

import "syscall"

func isErrorTimeout(err error) bool {
	syscallErr, ok := err.(syscall.Errno)
	if !ok {
		return false
	}
	return syscallErr.Timeout()
}

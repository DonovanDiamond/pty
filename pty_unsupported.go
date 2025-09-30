//go:build !linux && !darwin && !freebsd && !dragonfly && !netbsd && !openbsd && !solaris && !zos && !windows
// +build !linux,!darwin,!freebsd,!dragonfly,!netbsd,!openbsd,!solaris,!zos,!windows

package pty

import (
	"os"
)

func open() (pty, tty *os.File, err error) {
	return nil, nil, ErrUnsupported
}

//go:build windows
// +build windows

package pty

import (
	"syscall"
	"unsafe"
)

// Types from golang.org/x/sys/windows.

// Ref: https://pkg.go.dev/golang.org/x/sys/windows#Coord
type windowsCoord struct {
	X int16
	Y int16
}

// Ref: https://pkg.go.dev/golang.org/x/sys/windows#SmallRect
type windowsSmallRect struct {
	Left   int16
	Top    int16
	Right  int16
	Bottom int16
}

// Ref: https://pkg.go.dev/golang.org/x/sys/windows#ConsoleScreenBufferInfo
type windowsConsoleScreenBufferInfo struct {
	Size              windowsCoord
	CursorPosition    windowsCoord
	Attributes        uint16
	Window            windowsSmallRect
	MaximumWindowSize windowsCoord
}

func (c windowsCoord) Pack() uintptr {
	return uintptr((int32(c.Y) << 16) | int32(c.X))
}

// Setsize resizes t to ws.
func Setsize(t FdHolder, ws *Winsize) error {
	if err := resizePseudoConsole.Find(); err != nil {
		return err
	}

	//nolint:errcheck // Windows syscall: actual error status is in r0 HRESULT, not err.
	r0, _, _ := resizePseudoConsole.Call(
		t.Fd(),
		(windowsCoord{X: int16(ws.Cols), Y: int16(ws.Rows)}).Pack(),
	)
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}

		// S_OK: 0
		return syscall.Errno(r0)
	}

	return nil
}

// GetsizeFull returns the full terminal size description.
func GetsizeFull(t FdHolder) (size *Winsize, err error) {
	// GetConsoleScreenBufferInfo requires a console buffer handle, not a pseudo console handle.
	// For WindowsPty/WindowsTty, we need to get the actual console buffer handle.
	var consoleHandle uintptr

	switch v := t.(type) {
	case *WindowsPty:
		// Use the write pipe's file descriptor as it's connected to the console output
		consoleHandle = v.w.Fd()
	case *WindowsTty:
		// Use the write pipe's file descriptor as it's connected to the console output
		consoleHandle = v.w.Fd()
	default:
		// Fallback to Fd() for other types
		consoleHandle = t.Fd()
	}

	if err := getConsoleScreenBufferInfo.Find(); err != nil {
		return nil, err
	}

	var info windowsConsoleScreenBufferInfo
	//nolint:errcheck // Windows syscall: actual error status is in r0 HRESULT, not err.
	r0, _, _ := getConsoleScreenBufferInfo.Call(consoleHandle, uintptr(unsafe.Pointer(&info)))
	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}

		// S_OK: 0
		return nil, syscall.Errno(r0)
	}

	return &Winsize{
		Rows: uint16(info.Window.Bottom - info.Window.Top + 1),
		Cols: uint16(info.Window.Right - info.Window.Left + 1),
	}, nil
}

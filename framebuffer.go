// +build linux,cgo
package framebuffer

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

/*
#include <linux/fb.h>
*/
import "C"

// Device is a framebuffer device
type Device struct {
	FixScreenInfo FixScreenInfo
	Buffer        []byte
	file          *os.File
}

// Open framebuffer device
func Open(device string) (*Device, error) {
	file, err := os.OpenFile(device, os.O_RDWR, os.ModeDevice)
	if err != nil {
		return nil, fmt.Errorf("could not open device: %v", err)
	}
	fb := Device{
		file: file,
	}
	if err := fb.ioctl(C.FBIOGET_FSCREENINFO, unsafe.Pointer(&fb.FixScreenInfo)); err != nil {
		return nil, fmt.Errorf("ioctl FBIOGET_FSCREENINFO failed: %v", err)
	}
	mmap, err := syscall.Mmap(int(file.Fd()), 0, int(fb.FixScreenInfo.SmemLen), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return nil, fmt.Errorf("could not mmap: %v", err)
	}
	fb.Buffer = mmap
	return &fb, nil
}

// Close the framebuffer device
func (fb *Device) Close() error {
	if err := syscall.Munmap(fb.Buffer); err != nil {
		return err
	}
	return fb.file.Close()
}

// VarScreenInfo read the framebuffer settings
func (fb *Device) VarScreenInfo() (*VarScreenInfo, error) {
	v := VarScreenInfo{}
	if err := fb.ioctl(C.FBIOGET_VSCREENINFO, unsafe.Pointer(&v)); err != nil {
		return nil, fmt.Errorf("ioctl FBIOGET_FSCREENINFO failed: %v", err)
	}
	return &v, nil
}

// SetVarScreenInfo change framebuffer settings
func (fb *Device) SetVarScreenInfo(v *VarScreenInfo) error {
	return fb.ioctl(C.FBIOPUT_VSCREENINFO, unsafe.Pointer(&v))
}

// PanDisplay can be used to scroll inside the virtual resolution.
func (fb *Device) PanDisplay(v *VarScreenInfo) error {
	return fb.ioctl(C.FBIOPAN_DISPLAY, unsafe.Pointer(&v))
}

func (fb *Device) ioctl(cmd uintptr, arg unsafe.Pointer) error {
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fb.file.Fd(), cmd, uintptr(arg)); errno != 0 {
		return &os.SyscallError{"SYS_IOCTL", errno}
	}
	return nil
}

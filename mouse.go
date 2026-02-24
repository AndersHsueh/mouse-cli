package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

// uinput constants
const (
	UINPUT_MAX_NAME_SIZE = 80
	EV_REL               = 0x02
	EV_KEY               = 0x01
	EV_SYN               = 0x00

	// Relative mouse codes
	REL_X       = 0x00
	REL_Y       = 0x01
	REL_Z       = 0x02
	REL_WHEEL   = 0x08
	REL_HWHEEL  = 0x06

	// Mouse buttons
	BTN_LEFT   = 0x110
	BTN_RIGHT  = 0x111
	BTN_MIDDLE = 0x112
	BTN_SIDE   = 0x113
	BTN_EXTRA  = 0x114

	// Mouse button names
	MOUSE_LEFT   = "left"
	MOUSE_RIGHT  = "right"
	MOUSE_MIDDLE = "middle"
	MOUSE_SIDE   = "side"
	MOUSE_EXTRA  = "extra"
)

type uinputSetup struct {
	id      [16]byte
	name    [UINPUT_MAX_NAME_SIZE]byte
	ffEffectsMax uint32
}

type inputEvent struct {
	time  syscall.Timeval
	type_ uint16
	code  uint16
	value int32
}

// Ioctl constants
const (
	UI_DEV_SETUP = 0x40045539
	UI_DEV_CREATE = 0x4004553a
	UI_SET_EVBIT = 0x40045564
	UI_SET_RELBIT = 0x40045566
	UI_SET_KEYBIT = 0x40045565
)

// VirtualMouse represents a virtual mouse device
type VirtualMouse struct {
	file *os.File
}

// NewVirtualMouse creates a new virtual mouse
func NewVirtualMouse() (*VirtualMouse, error) {
	fd, err := os.OpenFile("/dev/uinput", os.O_WRONLY, 0644)
	if err != nil {
		fd, err = os.OpenFile("/dev/input/uinput", os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open uinput: %w (need root or input group)", err)
		}
	}

	vm := &VirtualMouse{file: fd}
	if err := vm.setupDevice(); err != nil {
		fd.Close()
		return nil, err
	}

	return vm, nil
}

func (vm *VirtualMouse) setupDevice() error {
	setup := uinputSetup{}
	copy(setup.name[:], []byte("virtual-mouse"))

	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		vm.file.Fd(),
		UI_DEV_SETUP,
		uintptr(unsafe.Pointer(&setup)),
	)
	if errno != 0 {
		return fmt.Errorf("ioctl UI_DEV_SETUP failed: %v", errno)
	}

	// Enable relative movement events
	if _, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		vm.file.Fd(),
		UI_SET_EVBIT,
		uintptr(EV_REL),
	); errno != 0 {
		return fmt.Errorf("ioctl UI_SET_EVBIT failed: %v", errno)
	}

	// Enable mouse button events
	if _, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		vm.file.Fd(),
		UI_SET_EVBIT,
		uintptr(EV_KEY),
	); errno != 0 {
		return fmt.Errorf("ioctl UI_SET_EVBIT (KEY) failed: %v", errno)
	}

	// Enable specific relative axes
	for _, code := range []int{REL_X, REL_Y, REL_WHEEL, REL_HWHEEL} {
		if _, _, errno := syscall.Syscall(
			syscall.SYS_IOCTL,
			vm.file.Fd(),
			UI_SET_RELBIT,
			uintptr(code),
		); errno != 0 {
			return fmt.Errorf("ioctl UI_SET_RELBIT failed: %v", errno)
		}
	}

	// Enable mouse buttons
	for _, code := range []int{BTN_LEFT, BTN_RIGHT, BTN_MIDDLE, BTN_SIDE, BTN_EXTRA} {
		if _, _, errno := syscall.Syscall(
			syscall.SYS_IOCTL,
			vm.file.Fd(),
			UI_SET_KEYBIT,
			uintptr(code),
		); errno != 0 {
			return fmt.Errorf("ioctl UI_SET_KEYBIT failed: %v", errno)
		}
	}

	// Create the device
	if _, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		vm.file.Fd(),
		UI_DEV_CREATE,
		0,
	); errno != 0 {
		return fmt.Errorf("ioctl UI_DEV_CREATE failed: %v", errno)
	}

	return nil
}

func (vm *VirtualMouse) sendEvent(code int, value int) error {
	event := inputEvent{
		time:  syscall.Timeval{},
		type_: EV_REL,
		code:  uint16(code),
		value: int32(value),
	}

	_, err := vm.file.Write((*[24]byte)(unsafe.Pointer(&event))[:])
	return err
}

func (vm *VirtualMouse) sendButtonEvent(code int, value int) error {
	event := inputEvent{
		time:  syscall.Timeval{},
		type_: EV_KEY,
		code:  uint16(code),
		value: int32(value),
	}

	_, err := vm.file.Write((*[24]byte)(unsafe.Pointer(&event))[:])
	return err
}

// Move moves the mouse by the specified delta
func (vm *VirtualMouse) Move(dx, dy int) error {
	if dx != 0 {
		if err := vm.sendEvent(REL_X, dx); err != nil {
			return err
		}
	}
	if dy != 0 {
		if err := vm.sendEvent(REL_Y, dy); err != nil {
			return err
		}
	}
	return nil
}

// MoveTo moves the mouse to absolute position (requires EV_ABS, not implemented for simplicity)
// For now, use Move for relative movement

// Scroll scrolls the mouse wheel
func (vm *VirtualMouse) Scroll(lines int) error {
	return vm.sendEvent(REL_WHEEL, lines)
}

// HScroll scrolls horizontally
func (vm *VirtualMouse) HScroll(lines int) error {
	return vm.sendEvent(REL_HWHEEL, lines)
}

// Click clicks a mouse button
func (vm *VirtualMouse) Click(button string) error {
	code, err := parseButton(button)
	if err != nil {
		return err
	}

	// Press
	if err := vm.sendButtonEvent(code, 1); err != nil {
		return err
	}
	// Release
	if err := vm.sendButtonEvent(code, 0); err != nil {
		return err
	}
	return nil
}

// Press holds down a mouse button
func (vm *VirtualMouse) Press(button string) error {
	code, err := parseButton(button)
	if err != nil {
		return err
	}
	return vm.sendButtonEvent(code, 1)
}

// Release releases a mouse button
func (vm *VirtualMouse) Release(button string) error {
	code, err := parseButton(button)
	if err != nil {
		return err
	}
	return vm.sendButtonEvent(code, 0)
}

// Close closes the virtual mouse device
func (vm *VirtualMouse) Close() error {
	if vm.file != nil {
		return vm.file.Close()
	}
	return nil
}

func parseButton(button string) (int, error) {
	switch strings.ToLower(button) {
	case "left", "l":
		return BTN_LEFT, nil
	case "right", "r":
		return BTN_RIGHT, nil
	case "middle", "m":
		return BTN_MIDDLE, nil
	case "side", "s":
		return BTN_SIDE, nil
	case "extra", "e":
		return BTN_EXTRA, nil
	default:
		return 0, fmt.Errorf("unknown button: %s (use left, right, middle)", button)
	}
}

func parseMovement(arg string) (dx, dy int, err error) {
	// Format: "x,y" or just a number (dx)
	parts := strings.Split(arg, ",")
	if len(parts) == 1 {
		// Single value = dx
		dx, err = strconv.Atoi(parts[0])
		return
	}
	if len(parts) == 2 {
		dx, err = strconv.Atoi(parts[0])
		if err != nil {
			return
		}
		dy, err = strconv.Atoi(parts[1])
		return
	}
	return 0, 0, fmt.Errorf("invalid movement format: %s (use '100' or '100,200')", arg)
}

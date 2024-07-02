package cli

import (
	"os"
	"syscall"
	"unsafe"

	"github.com/andrwui/hyprsources/constants"
)

// I didn't wanted to use libraries so here's something...
// I have no idea what this does, GPT wrote it. I tried to understand it so hard but i just can't right now.
func CheckInput() byte {
	// Open the terminal device file directly
	tty, err := os.OpenFile("/dev/tty", syscall.O_RDWR, 0)
	if err != nil {
		panic(err)
	}
	defer tty.Close()

	// Get the current terminal attributes
	var oldState syscall.Termios
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, tty.Fd(), syscall.TCGETS, uintptr(unsafe.Pointer(&oldState)), 0, 0, 0); err != 0 {
		panic(err)
	}

	// Create a new state based on the old one for raw mode
	newState := oldState
	newState.Iflag &^= syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK | syscall.ISTRIP | syscall.INLCR | syscall.IGNCR | syscall.ICRNL | syscall.IXON
	newState.Oflag &^= syscall.OPOST
	newState.Lflag &^= syscall.ECHO | syscall.ECHONL | syscall.ICANON | syscall.ISIG | syscall.IEXTEN
	newState.Cflag &^= syscall.CSIZE | syscall.PARENB
	newState.Cflag |= syscall.CS8
	newState.Cc[syscall.VMIN] = 1
	newState.Cc[syscall.VTIME] = 0

	// Set the terminal to raw mode
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, tty.Fd(), syscall.TCSETS, uintptr(unsafe.Pointer(&newState)), 0, 0, 0); err != 0 {
		panic(err)
	}

	// Ensure the terminal is restored to its original state on exit
	defer syscall.Syscall6(syscall.SYS_IOCTL, tty.Fd(), syscall.TCSETS, uintptr(unsafe.Pointer(&oldState)), 0, 0, 0)

	// Read from the terminal
	readBytes := make([]byte, 3)
	read, err := tty.Read(readBytes)
	if err != nil {
		panic(err)
	}

	// Check the input and return the appropriate byte
	if read == 3 {
		if _, ok := constants.Keys[readBytes[2]]; ok {
			return readBytes[2]
		}
	} else {
		return readBytes[0]
	}
	return 0
}

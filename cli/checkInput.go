package cli

import (
	"github.com/andrwui/hyprsources/constants"
	"github.com/pkg/term"
)

func CheckInput() byte {
	t, _ := term.Open("/dev/tty")

	err := term.RawMode(t)
	if err != nil {
		panic(err)
	}

	var read int
	readBytes := make([]byte, 3)
	read, err = t.Read(readBytes)
	if err != nil {
		panic(err)
	}

	t.Restore()
	t.Close()

	if read == 3 {
		if _, ok := constants.Keys[readBytes[2]]; ok {
			return readBytes[2]
		}
	} else {
		return readBytes[0]
	}
	return 0
}

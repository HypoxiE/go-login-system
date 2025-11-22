package utils

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func PressAnyKey(prints bool) {
	if prints {
		fmt.Println("Press any key...")
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	var b [1]byte
	os.Stdin.Read(b[:])
}

package utils

import (
	"fmt"
	"os"

	"github.com/HypoxiE/go-login-system/pkg/stdin"
	"github.com/HypoxiE/go-login-system/pkg/stdout"
	"golang.org/x/term"
)

func OldPressAnyKey(prints bool) {
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

func PressAnyKey(cin stdin.ConsoleInput, cout *stdout.ConsoleOutput) {
	if cout != nil {
		cout.LineOut("Press any key...")
		cout.Sync()
	}

	cin.Flush <- struct{}{}
	<-cin.LastSymbol

}

package core

import (
	"fmt"
	"os"

	"github.com/msteinert/pam"

	gpass "github.com/HypoxiE/go-login-system/pkg/get_password"
)

func StartPam(s pam.Style, msg string) (string, error) {
	switch s {
	case pam.PromptEchoOff:
		fmt.Print(msg)
		password, err := gpass.ReadPasswordWithStars("", os.Stdin, os.Stdout, int(os.Stdin.Fd()), true)
		if err != nil {
			fmt.Println("Error:", err)
			return "", nil
		}
		return password, nil

	case pam.PromptEchoOn:
		fmt.Print(msg)
		var input string
		fmt.Scanln(&input)
		return input, nil

	case pam.ErrorMsg:
		fmt.Fprintln(os.Stderr, msg)
		return "", nil

	case pam.TextInfo:
		fmt.Println(msg)
		return "", nil
	default:
		return "", nil
	}
}

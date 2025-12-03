package core

import (
	"github.com/msteinert/pam"

	gstr "github.com/HypoxiE/go-login-system/pkg/get_strings"
	"github.com/HypoxiE/go-login-system/pkg/stdin"
	"github.com/HypoxiE/go-login-system/pkg/stdout"
)

type PamInputs struct {
	COut *stdout.ConsoleOutput
	CIn  *stdin.ConsoleInput
}

func InitPI(cout *stdout.ConsoleOutput, cin *stdin.ConsoleInput) PamInputs {
	return PamInputs{
		COut: cout,
		CIn:  cin,
	}
}

func (pam_inp *PamInputs) StartPam(s pam.Style, msg string) (string, error) {
	cout, cin := pam_inp.COut, pam_inp.CIn
	switch s {
	case pam.PromptEchoOff:
		cout.TextOut(msg)
		password := gstr.ReadPasswordWithStars(cout, cin)
		return password, nil

	case pam.PromptEchoOn:
		cout.TextOut(msg)
		password := gstr.ReadString(cout, cin)
		return password, nil

	//case pam.ErrorMsg:
	//	fmt.Fprintln(os.Stderr, msg)
	//	return "", nil

	//case pam.TextInfo:
	//	fmt.Println(msg)
	//	return "", nil
	default:
		return "", nil
	}
}

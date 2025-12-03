package getstrings

import (
	"github.com/HypoxiE/go-login-system/pkg/stdin"
	"github.com/HypoxiE/go-login-system/pkg/stdout"
)

func ReadPasswordWithStars(cout *stdout.ConsoleOutput, cin *stdin.ConsoleInput) string {

	x := cout.CursorColumn
	for c := range cin.LastSymbol {
		if c == '\n' || c == '\r' {
			break
		} else if c == 0x7f || c == 0x08 {
			if cout.CursorColumn > x {
				cout.CursorColumn--
				cout.TextOut(" ")
				cout.CursorColumn--
				cout.ShowCursor()
			}
			continue
		}
		cout.TextOut("*")
		cout.ShowCursor()
	}
	pass := cin.GetForLine()
	cin.Flush <- struct{}{}
	cout.NewLine()
	cout.HideCursor()

	return pass
}

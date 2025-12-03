package getstrings

import (
	"github.com/HypoxiE/go-login-system/pkg/stdin"
	"github.com/HypoxiE/go-login-system/pkg/stdout"
)

func ReadString(cout *stdout.ConsoleOutput, cin *stdin.ConsoleInput, replace_symbol *rune) string {

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
		if replace_symbol == nil {
			cout.TextOut(string(c))
		} else {
			cout.TextOut(string(*replace_symbol))
		}
		cout.ShowCursor()
	}
	str := cin.GetForLine()
	cin.Flush <- struct{}{}
	cout.NewLine()
	cout.HideCursor()

	return str
}

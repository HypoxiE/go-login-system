package getstrings

import (
	"github.com/HypoxiE/go-login-system/pkg/stdin"
	"github.com/HypoxiE/go-login-system/pkg/stdout"
)

func ReadPasswordWithStars(cout *stdout.ConsoleOutput, cin *stdin.ConsoleInput) string {
	r_sym := '*'
	return ReadString(cout, cin, &r_sym)
}

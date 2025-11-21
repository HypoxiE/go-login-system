package getpassword

import (
	"fmt"
	"io"

	"golang.org/x/term"
)

// ReadPasswordWithStars("Password: ", os.Stdin, os.Stdout, int(os.Stdin.Fd()), true)
func ReadPasswordWithStars(prompt string, r io.Reader, w io.Writer, fd int, raw bool) (string, error) {
	var oldState *term.State
	var err error

	if raw {
		oldState, err = term.MakeRaw(fd)
		if err != nil {
			return "", err
		}
		defer term.Restore(fd, oldState)
	}

	fmt.Fprint(w, prompt)

	var pass []byte
	buf := make([]byte, 1)

	for {
		n, _ := r.Read(buf)
		if n == 0 {
			continue
		}
		c := buf[0]

		if c == '\n' || c == '\r' {
			fmt.Fprint(w, "\r\n")
			break
		}

		if c == 127 || c == 8 {
			if len(pass) > 0 {
				pass = pass[:len(pass)-1]
				fmt.Fprint(w, "\b \b")
			}
			continue
		}

		pass = append(pass, c)
		fmt.Fprint(w, "*")
	}

	return string(pass), nil
}

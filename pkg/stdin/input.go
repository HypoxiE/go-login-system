package stdin

import (
	"io"
	"unicode/utf8"

	"golang.org/x/term"
)

type ConsoleInput struct {
	Value      []rune
	LastSymbol chan rune
	Flush      chan struct{}

	StopOutput chan struct{}
}

func InitCIn() ConsoleInput {
	stop_sign := make(chan struct{})
	flush := make(chan struct{})
	last := make(chan rune, 1)

	return ConsoleInput{
		Value: []rune{},
		Flush: flush,

		LastSymbol: last,

		StopOutput: stop_sign,
	}
}

func (inp *ConsoleInput) MainLoop(r io.Reader, w io.Writer, fd int, raw bool) error {

	var oldState *term.State
	var err error

	if raw {
		oldState, err = term.MakeRaw(fd)
		if err != nil {
			return err
		}
		defer term.Restore(fd, oldState)
	}

	inputCh := make(chan rune)

	go func() {
		buf := make([]byte, 1)
		var utfbuf []byte
		for {
			n, err := r.Read(buf)
			if err != nil {
				close(inputCh)
				return
			}
			if n == 0 {
				continue
			}

			utfbuf = append(utfbuf, buf[0])

			r, size := utf8.DecodeRune(utfbuf)
			if r == utf8.RuneError && size == 1 {
				continue
			}
			inputCh <- r
			utfbuf = utfbuf[size:]
		}
	}()

	for {
		select {
		case c, ok := <-inputCh:
			if !ok {
				return nil
			}

			select {
			case <-inp.LastSymbol:
			default:
			}
			inp.LastSymbol <- c

			//if c == 0x03 {
			//	key_interrupt <- struct{}{}
			//}

			if c == 0x7f || c == 0x08 {
				if len(inp.Value) > 0 {
					inp.Value = inp.Value[:len(inp.Value)-1]
				}
				continue
			}

			if c == '\r' {
				inp.Value = append(inp.Value, '\n')
			}
			inp.Value = append(inp.Value, c)
		case <-inp.StopOutput:
			return nil
		case <-inp.Flush:
			inp.Value = []rune{}
			select {
			case <-inp.LastSymbol:
			default:
			}
		}
	}
}

func (inp ConsoleInput) GetLength() int {
	return len(inp.Value)
}

func (inp ConsoleInput) GetValue() string {
	return string(inp.Value)
}

func (inp ConsoleInput) GetForLine() string {
	val := ""
	for _, char := range inp.Value {
		if char == '\n' || char == '\r' {
			return val
		}
		val += string(char)
	}
	return val
}

func (inp *ConsoleInput) Clean() {
	inp.Flush <- struct{}{}
}

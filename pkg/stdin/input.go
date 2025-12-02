package stdin

import (
	"io"

	"golang.org/x/term"
)

type ConsoleInput struct {
	Value []byte
	Flush chan struct{}

	StopOutput chan struct{}
}

func InitCIn() ConsoleInput {
	stop_sign := make(chan struct{})
	flush := make(chan struct{})

	return ConsoleInput{
		Value: []byte{},
		Flush: flush,

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

	inputCh := make(chan byte)

	go func() {
		buf := make([]byte, 1)
		for {
			n, err := r.Read(buf)
			if err != nil {
				close(inputCh)
				return
			}
			if n > 0 {
				inputCh <- buf[0]
			}
		}
	}()

	for {
		select {
		case c, ok := <-inputCh:
			if !ok {
				return nil
			}
			if c == '\r' {
				inp.Value = append(inp.Value, '\n')
			}
			inp.Value = append(inp.Value, c)
		case <-inp.StopOutput:
			return nil
		case <-inp.Flush:
			inp.Value = []byte{}
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

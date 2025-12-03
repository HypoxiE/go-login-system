package stdout

import (
	"log"
	"sync"
	"time"

	"github.com/gdamore/tcell"
)

type ConsoleOutput struct {
	Screen tcell.Screen

	CursorMutex  sync.Mutex
	CursorColumn int
	CursorLine   int

	SyncMutex sync.Mutex
}

func (cout *ConsoleOutput) SetCursorPosition(x int, y int) {
	cout.CursorMutex.Lock()
	defer cout.CursorMutex.Unlock()

	cout.CursorColumn = x
	cout.CursorLine = y
}

func (cout *ConsoleOutput) SetCursorXPosition(x int) {
	cout.CursorMutex.Lock()
	defer cout.CursorMutex.Unlock()

	cout.CursorColumn = x
}

func (cout *ConsoleOutput) SetCursorYPosition(y int) {
	cout.CursorMutex.Lock()
	defer cout.CursorMutex.Unlock()

	cout.CursorLine = y
}

func (cout *ConsoleOutput) GetCursorPosition() (x int, y int) {
	cout.CursorMutex.Lock()
	defer cout.CursorMutex.Unlock()

	return cout.CursorColumn, cout.CursorLine
}

func InitCOut() ConsoleOutput {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	screen.Clear()
	screen.HideCursor()

	return ConsoleOutput{
		Screen:       screen,
		CursorColumn: 0,
		CursorLine:   0,
	}
}

func (cout *ConsoleOutput) NewLine() {
	cout.SetCursorPosition(0, cout.CursorLine+1)
}

func (cout *ConsoleOutput) LineOut(new_string string) {
	for _, r := range new_string {
		cout.Screen.SetContent(cout.CursorColumn, cout.CursorLine, r, nil, tcell.StyleDefault)
		cout.SetCursorXPosition(cout.CursorColumn + 1)
	}
	cout.NewLine()
}

func (cout *ConsoleOutput) FreeTextOut(x int, y int, new_string string, use_x_in_new_string bool) (x_end int, y_end int) {
	cursor_line := y
	cursor_col := x

	for _, r := range new_string {
		if r == '\n' {
			cursor_line += 1
			if use_x_in_new_string {
				cursor_col = x
			} else {
				cursor_col = 0
			}
			continue
		} else if r == '\r' {
			if use_x_in_new_string {
				cursor_col = x
			} else {
				cursor_col = 0
			}
			continue
		}
		cout.Screen.SetContent(cursor_col, cursor_line, r, nil, tcell.StyleDefault)
		cursor_col += 1
	}

	return cursor_col, cursor_line
}

func (cout *ConsoleOutput) TextOut(new_string string) {
	cout.SetCursorPosition(cout.FreeTextOut(cout.CursorColumn, cout.CursorLine, new_string, false))
}

func (cout *ConsoleOutput) TextOutLn(new_string string) {
	cout.TextOut(new_string)
	cout.NewLine()
}

func (cout *ConsoleOutput) TextOutSync(new_string string) {
	cout.CursorColumn, cout.CursorLine = cout.FreeTextOut(cout.CursorColumn, cout.CursorLine, new_string, false)
	cout.Sync()
}

func (cout *ConsoleOutput) SlowTextOut(x int, y int, new_string string, autosync bool, millisec_duration int) {
	for _, r := range new_string {
		x, y = cout.FreeTextOut(x, y, string(r), false)
		if autosync {
			cout.Sync()
		}
		time.Sleep(time.Duration(millisec_duration) * time.Millisecond)
	}
}

func (cout *ConsoleOutput) ShowCursor() {
	cout.Screen.ShowCursor(cout.CursorColumn, cout.CursorLine)
}
func (cout *ConsoleOutput) HideCursor() {
	cout.Screen.HideCursor()
}

func (cout *ConsoleOutput) Fini() {
	cout.Screen.Fini()
}

func (cout *ConsoleOutput) Sync() {
	cout.SyncMutex.Lock()
	defer cout.SyncMutex.Unlock()
	cout.Screen.Show()
}

func (cout *ConsoleOutput) SyncLoop(stop_signal chan struct{}, fps int) {
	runFlag := true

	for runFlag {
		select {
		case <-stop_signal:
			runFlag = false
		default:
			cout.Sync()
		}
		time.Sleep(time.Second / time.Duration(fps))
	}
}

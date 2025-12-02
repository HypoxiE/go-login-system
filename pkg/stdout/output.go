package stdout

import (
	"log"
	"strings"
	"time"

	"github.com/gdamore/tcell"
)

type ConsoleOutput struct {
	Screen tcell.Screen

	CursorColumn int
	CursorLine   int
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
	cout.CursorColumn = 0
	cout.CursorLine += 1
}

func (cout *ConsoleOutput) LineOut(new_string string) {
	for _, r := range new_string {
		cout.Screen.SetContent(cout.CursorColumn, cout.CursorLine, r, nil, tcell.StyleDefault)
		cout.CursorColumn += 1
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

	return cursor_line, cursor_col
}

func (cout *ConsoleOutput) TextOut(new_string string) {

	for _, r := range new_string {
		if r == '\n' {
			cout.CursorLine++
			cout.CursorColumn = 0
			continue
		} else if r == '\r' {
			cout.CursorColumn = 0
			continue
		}
		cout.Screen.SetContent(cout.CursorColumn, cout.CursorLine, r, nil, tcell.StyleDefault)
		cout.CursorColumn += 1
	}
}

func (cout *ConsoleOutput) SlowTextOut(new_string string) {

	x := cout.CursorColumn
	y := cout.CursorLine

	cout.CursorLine += strings.Count(new_string, "\n") + 1

	for _, r := range new_string {
		if r == '\n' {
			y++
			x = 0
			continue
		} else if r == '\r' {
			x = 0
			continue
		}
		cout.Screen.SetContent(x, y, r, nil, tcell.StyleDefault)
		x += 1
		cout.Sync()
		time.Sleep(2 * time.Millisecond)
	}
}

func (cout *ConsoleOutput) GetCursor() (x int, y int) {
	return cout.CursorColumn, cout.CursorLine
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
	cout.Screen.Show()
}

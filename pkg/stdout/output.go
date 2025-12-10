package stdout

import (
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell"
)

type ConsoleOutput struct {
	Screen tcell.Screen

	CursorMutex  sync.Mutex
	CursorColumn int
	CursorLine   int

	CurrentStyle tcell.Style

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

		CurrentStyle: tcell.StyleDefault,
	}
}

func (cout *ConsoleOutput) NewLine() {
	cout.SetCursorPosition(0, cout.CursorLine+1)
}

func (cout *ConsoleOutput) LineOut(new_string string) {
	for _, r := range new_string {
		cout.Screen.SetContent(cout.CursorColumn, cout.CursorLine, r, nil, cout.CurrentStyle)
		cout.SetCursorXPosition(cout.CursorColumn + 1)
	}
	cout.NewLine()
}

const ColorAnsiRunes = "0123456789;"

func (cout *ConsoleOutput) FreeTextOut(x int, y int, new_string string, use_x_in_new_string bool) (x_end int, y_end int) {
	cursor_line := y
	cursor_col := x

	style := tcell.StyleDefault

	char_i := 0

	rune_array := []rune(new_string)

	for char_i < len(rune_array) {
		if rune_array[char_i] == '\n' {
			cursor_line += 1
			if use_x_in_new_string {
				cursor_col = x
			} else {
				cursor_col = 0
			}
			char_i++
			continue
		} else if rune_array[char_i] == '\r' {
			if use_x_in_new_string {
				cursor_col = x
			} else {
				cursor_col = 0
			}
			char_i++
			continue
		} else if rune_array[char_i] == '\x1b' && rune_array[char_i+1] == '[' {
			char_i += 2
			color_ansi_code := ""
			for ; char_i < len(rune_array) && rune_array[char_i] != 'm'; char_i++ {

				// Проверка на часть незаконченного блока ansi кода
				pass_chapter := false
				for _, r := range ColorAnsiRunes {
					if r == rune_array[char_i] {
						pass_chapter = true
						break
					}
				}
				if !pass_chapter {
					break
				}

				color_ansi_code += string(rune_array[char_i])
			}
			codes := []int{}

			for attr := range strings.SplitSeq(color_ansi_code, ";") {
				c, _ := strconv.Atoi(attr)
				codes = append(codes, c)
			}
			cout.ApplyStyle(codes, &style)

			char_i++
			continue
		}
		cout.Screen.SetContent(cursor_col, cursor_line, rune_array[char_i], nil, style)
		cursor_col += 1
		char_i++
	}

	return cursor_col, cursor_line
}

var ansiColors [16]tcell.Color = [16]tcell.Color{
	tcell.ColorBlack,   // 0
	tcell.ColorMaroon,  // 1
	tcell.ColorGreen,   // 2
	tcell.ColorOlive,   // 3
	tcell.ColorNavy,    // 4
	tcell.ColorPurple,  // 5
	tcell.ColorTeal,    // 6
	tcell.ColorSilver,  // 7
	tcell.ColorGray,    // 8
	tcell.ColorRed,     // 9
	tcell.ColorLime,    // 10
	tcell.ColorYellow,  // 11
	tcell.ColorBlue,    // 12
	tcell.ColorFuchsia, // 13
	tcell.ColorAqua,    // 14
	tcell.ColorWhite,   // 15
}

func (cout *ConsoleOutput) ApplyStyle(styleAttrs []int, style *tcell.Style) {

	if style == nil {
		style = &cout.CurrentStyle
	}

	for _, c := range styleAttrs {
		switch c {
		case 0:
			*style = tcell.StyleDefault
		case 1:
			*style = style.Bold(true)
		case 4:
			*style = style.Underline(true)

		// basic fg 30-37
		case 30:
			*style = style.Foreground(ansiColors[0]) // black
		case 31:
			*style = style.Foreground(ansiColors[1]) // red
		case 32:
			*style = style.Foreground(ansiColors[2]) // green
		case 33:
			*style = style.Foreground(ansiColors[3]) // yellow
		case 34:
			*style = style.Foreground(ansiColors[4]) // blue
		case 35:
			*style = style.Foreground(ansiColors[5]) // magenta
		case 36:
			*style = style.Foreground(ansiColors[6]) // cyan
		case 37:
			*style = style.Foreground(ansiColors[7]) // white

		// basic bg 40-47
		case 40:
			*style = style.Background(ansiColors[0])
		case 41:
			*style = style.Background(ansiColors[1])
		case 42:
			*style = style.Background(ansiColors[2])
		case 43:
			*style = style.Background(ansiColors[3])
		case 44:
			*style = style.Background(ansiColors[4])
		case 45:
			*style = style.Background(ansiColors[5])
		case 46:
			*style = style.Background(ansiColors[6])
		case 47:
			*style = style.Background(ansiColors[7])

		// bright fg 90-97 -> map to tcell bright colors
		case 90:
			*style = style.Foreground(ansiColors[8])
		case 91:
			*style = style.Foreground(ansiColors[9])
		case 92:
			*style = style.Foreground(ansiColors[10])
		case 93:
			*style = style.Foreground(ansiColors[11])
		case 94:
			*style = style.Foreground(ansiColors[12])
		case 95:
			*style = style.Foreground(ansiColors[13])
		case 96:
			*style = style.Foreground(ansiColors[14])
		case 97:
			*style = style.Foreground(ansiColors[15])

		// bright bg 100-107
		case 100:
			*style = style.Background(ansiColors[8])
		case 101:
			*style = style.Background(ansiColors[9])
		case 102:
			*style = style.Background(ansiColors[10])
		case 103:
			*style = style.Background(ansiColors[11])
		case 104:
			*style = style.Background(ansiColors[12])
		case 105:
			*style = style.Background(ansiColors[13])
		case 106:
			*style = style.Background(ansiColors[14])
		case 107:
			*style = style.Background(ansiColors[15])

		}
	}
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

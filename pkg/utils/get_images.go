package utils

import (
	"fmt"
	"strings"
	"time"
)

func PaintAsciiGif(gif string, done chan struct{}) {
	lines := strings.Split(gif, "\n")[1:]

loop:
	for {
		count := 0

		for _, line := range lines {
			select {
			case <-done:
				break loop // выходим из внешнего бесконечного цикла
			default:
			}

			if strings.Contains(line, "\x1b[u") {
				i := strings.Index(line, "\x1b[u")

				left := line[:i]
				right := line[i+len("\x1b[u"):]

				fmt.Println(left)

				time.Sleep(100 * time.Millisecond)
				clearLines(count)
				count = 1

				fmt.Println(right)
			} else {
				fmt.Println(line)
			}

			count += 1
		}
		time.Sleep(100 * time.Millisecond)
		clearLines(count)
	}
}

func clearLines(n int) {
	for i := 0; i < n; i++ {
		fmt.Print("\033[A")
		fmt.Print("\033[K")
	}
}

package outputanimations

import (
	"time"

	"github.com/HypoxiE/go-login-system/pkg/stdout"
)

func GrayGifOutput(x, y int, cout *stdout.ConsoleOutput, frames []string, delay int, stop_gif_animation chan struct{}) {
	for i := uint(0); ; i++ {
		select {
		case <-stop_gif_animation:
			return
		default:
			cout.FreeTextOut(x, y, frames[i%uint(len(frames))], true)
		}
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
}

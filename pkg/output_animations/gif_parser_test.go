package outputanimations

import (
	_ "embed"
	"fmt"
	"testing"
)

//go:embed test_gif/w_p.txt
var wrongPasswordGif string

func TestGetRawGifInfo(t *testing.T) {
	frames, height := GetRawGifInfo(wrongPasswordGif)

	fmt.Println(frames, height)
}

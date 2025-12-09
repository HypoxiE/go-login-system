package outputanimations

import "strings"

// frames; frame_y_length
func GetRawGifInfo(gif string) ([]string, int) {
	var frames []string

	var i int
	// Находим начало
	for i < len(gif) {
		if gif[i] == 0x1B && gif[i+1] == 0x5B && gif[i+2] == 's' {
			i += 3
			break
		}
		i++
	}

	active_gif := 0
	frames = append(frames, "")

	for i < len(gif) {
		if gif[i] == 0x1B && gif[i+1] == 0x5B && gif[i+2] == 'u' {
			active_gif += 1
			frames = append(frames, "")
			i += 3
			continue
		} else if gif[i] == 0x1B && gif[i+1] == 0x5B && gif[i+2] == '?' &&
			gif[i+3] == '2' && gif[i+4] == '5' && gif[i+5] == 'h' {
			break
		}
		frames[active_gif] += string(gif[i])
		i++
	}

	return frames, strings.Count(frames[0], "\n")
}

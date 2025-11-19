package getpassword

import (
	"bytes"
	"strings"
	"testing"
)

func TestReadPasswordWithStars(t *testing.T) {
	var testEnter = [][2]string{
		{"a b\n", "a b"},
		{"qwer\n", "qwer"},
	}

	for _, st := range testEnter {
		input := bytes.NewBuffer([]byte(st[0]))
		output := &bytes.Buffer{}

		pass, err := ReadPasswordWithStars(
			"Password: ", input, output, 0, false,
		)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if pass != st[1] {
			t.Fatalf("want %q, got %q", st[1], pass)
		}

		if !strings.Contains(output.String(), "Password: ") {
			t.Fatalf("prompt missing")
		}
	}

}

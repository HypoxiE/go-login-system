package main

import (
	"fmt"
	"os"

	core "github.com/HypoxiE/go-login-system/pkg/core"
	"github.com/msteinert/pam"
)

func main() {
	t, err := pam.StartFunc("login", "", core.StartPam)

	if err != nil {
		fmt.Println("pam_start error:", err)
		os.Exit(1)
	}

	if err = t.Authenticate(0); err != nil {
		fmt.Println("auth failed:", err)
		os.Exit(1)
	}

	fmt.Println("Authenticated successfully")
}

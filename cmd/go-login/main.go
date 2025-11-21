package main

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"

	core "github.com/HypoxiE/go-login-system/pkg/core"
	"github.com/msteinert/pam"
)

func main() {
	fmt.Print("Login: ")
	var username string
	fmt.Scanln(&username)

	t, err := pam.StartFunc("login", username, core.StartPam)

	if err != nil {
		fmt.Println("pam_start error:", err)
		os.Exit(1)
	}

	if err = t.Authenticate(pam.Silent); err != nil {
		fmt.Println("auth failed:", err)
		os.Exit(1)
	}

	if err := t.AcctMgmt(pam.Silent); err != nil {
		fmt.Println("AcctMgmt error:", err)
		os.Exit(1)
	}

	if err := t.OpenSession(pam.Silent); err != nil {
		fmt.Println("pam_open_session error:", err)
		os.Exit(1)
	}

	usr, err := user.Lookup(username)
	if err != nil {
		fmt.Println("lookup error:", err)
		os.Exit(1)
	}
	uid, _ := strconv.Atoi(usr.Uid)
	gid, _ := strconv.Atoi(usr.Gid)
	syscall.Setgid(gid)
	syscall.Setuid(uid)

	os.Chdir(usr.HomeDir)
	os.Setenv("HOME", usr.HomeDir)
	os.Setenv("USER", usr.Username)
	os.Setenv("LOGNAME", usr.Username)
	os.Setenv("SHELL", "/bin/bash")

	shell := "/bin/bash"

	syscall.Exec(shell, []string{shell}, os.Environ())

	fmt.Println("Authenticated successfully")
}

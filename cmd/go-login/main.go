package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"

	_ "embed"

	core "github.com/HypoxiE/go-login-system/pkg/core"
	"github.com/HypoxiE/go-login-system/pkg/initialize"
	"github.com/HypoxiE/go-login-system/pkg/utils"
	"github.com/msteinert/pam"
)

var (
	//go:embed text_templates/welcome_screen.txt
	welcomeScreen string
	//go:embed text_templates/start_screen.txt
	startScreen string
)

func main() {
	fmt.Println(startScreen)

	fmt.Print("Login: ")
	var username string
	fmt.Scanln(&username)

	if username == "shutdown" {
		exec.Command("shutdown", "-h", "now").Run()
	}

	t, err := pam.StartFunc("login", username, core.StartPam)

	if err != nil {
		fmt.Println("pam_start error:", err)
		utils.PressAnyKey(false)
		os.Exit(1)
	}

	if err = t.Authenticate(pam.Silent); err != nil {
		if err.Error() == "Authentication failure" {
			fmt.Println("auth failed:", err)
		} else if err.Error() == "User not known to the underlying authentication module" {
			fmt.Println("auth failed:", err)
		} else {
			fmt.Println("auth failed:", err)
		}
		utils.PressAnyKey(false)
		os.Exit(1)
	}

	if err := t.AcctMgmt(pam.Silent); err != nil {
		fmt.Println("AcctMgmt error:", err)
		utils.PressAnyKey(false)
		os.Exit(1)
	}

	usr, err := user.Lookup(username)
	if err != nil {
		fmt.Println("lookup error:", err)
		utils.PressAnyKey(false)
		os.Exit(1)
	}

	groups, err := usr.GroupIds()
	if err != nil {
		fmt.Println("GroupIds error:", err)
		utils.PressAnyKey(false)
		os.Exit(1)
	}
	gids := make([]int, len(groups))
	for i, g := range groups {
		gi, _ := strconv.Atoi(g)
		gids[i] = gi
	}
	syscall.Setgroups(gids)

	if err := t.SetCred(pam.EstablishCred); err != nil {
		fmt.Println("set_cred error:", err)
		utils.PressAnyKey(false)
		os.Exit(1)
	}

	if err := t.OpenSession(pam.Silent); err != nil {
		fmt.Println("pam_open_session error:", err)
		utils.PressAnyKey(false)
		os.Exit(1)
	}

	env, err := t.GetEnvList()
	if err != nil {
		fmt.Println("getenv error:", err)
		utils.PressAnyKey(false)
		os.Exit(1)
	}

	os.Chdir(usr.HomeDir)

	initialize.InitEnv(env)

	uid, _ := strconv.Atoi(usr.Uid)
	gid, _ := strconv.Atoi(usr.Gid)
	syscall.Setgid(gid)
	syscall.Setuid(uid)

	fmt.Println(welcomeScreen)

	syscall.Exec(os.Getenv("SHELL"), []string{os.Getenv("SHELL")}, os.Environ())

	t.CloseSession(pam.Silent)
	t.SetCred(pam.DeleteCred)
}

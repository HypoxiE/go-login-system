package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"syscall"

	"github.com/HypoxiE/go-login-system/pkg/core"
	getstrings "github.com/HypoxiE/go-login-system/pkg/get_strings"
	"github.com/HypoxiE/go-login-system/pkg/initialize"
	"github.com/HypoxiE/go-login-system/pkg/stdin"
	"github.com/HypoxiE/go-login-system/pkg/stdout"
	"github.com/HypoxiE/go-login-system/pkg/utils"
	"github.com/msteinert/pam"
)

var (
	//go:embed text_templates/welcome_screen.txt
	welcomeScreen string
	//go:embed text_templates/start_screen.txt
	startScreen string

	//go:embed text_templates/wrong_password.txt
	wrongPasswordGif string
)

func run() int {
	cin := stdin.InitCIn()
	go cin.MainLoop(os.Stdin, os.Stdout, int(os.Stdin.Fd()), true)
	defer func() {
		cin.StopOutput <- struct{}{}
	}()
	cout := stdout.InitCOut()
	defer func() {
		cout.Fini()
	}()
	defer func() {
		if r := recover(); r != nil {
			cout.Fini()
			panic(r)
		}
	}()
	StopSyncLoop := make(chan struct{})
	go cout.SyncLoop(StopSyncLoop, 30)
	defer func() {
		StopSyncLoop <- struct{}{}
	}()
	{
		x, y := cout.GetCursorPosition()
		cout.SetCursorYPosition(cout.CursorLine + strings.Count(startScreen, "\n"))
		go cout.SlowTextOut(x, y, startScreen, false, 1)
		cout.NewLine()
	}

	cout.TextOut("Login: ")
	cout.ShowCursor()
	username := getstrings.ReadString(&cout, &cin)
	if username == "shutdown" {
		exec.Command("shutdown", "-h", "now").Run()
	}

	sp := core.InitPI(&cout, &cin)
	t, err := pam.StartFunc("login", username, sp.StartPam)
	if err != nil {
		cout.TextOutLn("pam_start error:" + err.Error())
		utils.PressAnyKey(cin, nil)
		return 1
	}

	if err = t.Authenticate(pam.Silent); err != nil {
		if err.Error() == "Authentication failure" {
			cout.TextOutLn("Wrong password, baka!")
			utils.PressAnyKey(cin, nil)
		} else if err.Error() == "User not known to the underlying authentication module" {
			cout.TextOutLn("auth failed:" + err.Error())
			utils.PressAnyKey(cin, nil)
		} else {
			cout.TextOutLn("auth failed:" + err.Error())
			utils.PressAnyKey(cin, nil)
		}
		return 1
	}

	if err := t.AcctMgmt(pam.Silent); err != nil {
		cout.TextOutLn("AcctMgmt error: " + err.Error())
		utils.PressAnyKey(cin, nil)
		return 1
	}

	usr, err := user.Lookup(username)
	if err != nil {
		cout.TextOutLn("lookup error: " + err.Error())
		utils.PressAnyKey(cin, nil)
		return 1
	}

	groups, err := usr.GroupIds()
	if err != nil {
		cout.TextOutLn("GroupIds error: " + err.Error())
		utils.PressAnyKey(cin, nil)
		return 1
	}
	gids := make([]int, len(groups))
	for i, g := range groups {
		gi, _ := strconv.Atoi(g)
		gids[i] = gi
	}
	syscall.Setgroups(gids)

	if err := t.SetCred(pam.EstablishCred); err != nil {
		cout.TextOutLn("set_cred error: " + err.Error())
		utils.PressAnyKey(cin, nil)
		return 1
	}

	if err := t.OpenSession(pam.Silent); err != nil {
		cout.TextOutLn("pam_open_session error: " + err.Error())
		utils.PressAnyKey(cin, nil)
		return 1
	}

	env, err := t.GetEnvList()
	if err != nil {
		cout.TextOutLn("getenv error: " + err.Error())
		utils.PressAnyKey(cin, nil)
		return 1
	}

	os.Chdir(usr.HomeDir)

	initialize.InitEnv(env)

	uid, _ := strconv.Atoi(usr.Uid)
	gid, _ := strconv.Atoi(usr.Gid)
	syscall.Setgid(gid)
	syscall.Setuid(uid)

	cin.StopOutput <- struct{}{}
	StopSyncLoop <- struct{}{}
	cout.Fini()

	fmt.Println(welcomeScreen)

	syscall.Exec(os.Getenv("SHELL"), []string{os.Getenv("SHELL")}, os.Environ())
	t.CloseSession(pam.Silent)
	t.SetCred(pam.DeleteCred)

	return 0
}

func main() {
	excode := run()
	os.Exit(excode)
}

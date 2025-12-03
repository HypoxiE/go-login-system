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
	gstrings "github.com/HypoxiE/go-login-system/pkg/get_strings"
	"github.com/HypoxiE/go-login-system/pkg/initialize"
	outanims "github.com/HypoxiE/go-login-system/pkg/output_animations"
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

	//go:embed text_templates/w_p.txt
	wrongPasswordGif string
)

func resetTTY() {
	exec.Command("stty", "sane").Run()
}

func run() int {

	resetTTY()

	cin := stdin.InitCIn()
	go cin.MainLoop(os.Stdin, os.Stdout, int(os.Stdin.Fd()), true)
	defersIsNeed := true
	defer func() {
		if defersIsNeed {
			cin.StopOutput <- struct{}{}
		}
	}()
	cout := stdout.InitCOut()
	defer func() {
		if defersIsNeed {
			cout.Fini()
		}
	}()
	StopSyncLoop := make(chan struct{})
	go cout.SyncLoop(StopSyncLoop, 30)
	defer func() {
		if defersIsNeed {
			StopSyncLoop <- struct{}{}
		}
	}()
	{
		x, y := cout.GetCursorPosition()
		cout.SetCursorYPosition(cout.CursorLine + strings.Count(startScreen, "\n"))
		go cout.SlowTextOut(x, y, startScreen, false, 1)
		cout.NewLine()
	}

	cout.TextOut("Login: ")
	cout.ShowCursor()
	username := gstrings.ReadString(&cout, &cin, nil)
	if username == "shutdown" {
		return 1000
	}

	sp := core.InitPI(&cout, &cin)
	t, err := pam.StartFunc("login", username, sp.StartPam)
	if err != nil {
		cout.TextOutLn("pam_start error: " + err.Error())
		utils.PressAnyKey(cin, nil)
		return 1
	}

	if err = t.Authenticate(pam.Silent); err != nil {
		if err.Error() == "Authentication failure" {
			cout.TextOutLn("Wrong password, baka!")
			stop_gif_animation := make(chan struct{})
			gifDeferIsNeed := true
			defer func() {
				if gifDeferIsNeed {
					close(stop_gif_animation)
				}
			}()
			{
				frames, y_len := outanims.GetRawGifInfo(wrongPasswordGif)
				x, y := cout.GetCursorPosition()
				go outanims.GrayGifOutput(x, y, &cout, frames, 100, stop_gif_animation)
				cout.SetCursorYPosition(cout.CursorLine + y_len)
				cout.NewLine()
			}
			utils.PressAnyKey(cin, nil)
			close(stop_gif_animation)
			gifDeferIsNeed = false

		} else if err.Error() == "User not known to the underlying authentication module" {
			cout.TextOutLn("auth failed: " + err.Error())
			utils.PressAnyKey(cin, nil)
		} else {
			cout.TextOutLn("auth failed: " + err.Error())
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

	close(cin.StopOutput)
	close(StopSyncLoop)
	cout.Fini()
	defersIsNeed = false

	fmt.Println(welcomeScreen)

	syscall.Exec(os.Getenv("SHELL"), []string{os.Getenv("SHELL")}, os.Environ())
	t.CloseSession(pam.Silent)
	t.SetCred(pam.DeleteCred)

	return 0
}

func main() {
	excode := run()
	if excode == 1000 {
		exec.Command("shutdown", "-h", "now").Run()
	}
	os.Exit(excode)
}

package main

import (
	_ "embed"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"sync"
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

	//go:embed text_templates/wrong_password.txt
	wrongPasswordGif string
)

func resetTTY() {
	exec.Command("stty", "sane").Run()
}

func wrongPassword(cout *stdout.ConsoleOutput, cin *stdin.ConsoleInput) {
	cout.TextOutLn("Wrong password, baka!")

	stop := make(chan struct{})
	frames, ylen := outanims.GetRawGifInfo(wrongPasswordGif)
	x, y := cout.GetCursorPosition()

	cout.SetCursorYPosition(cout.CursorLine + ylen)
	cout.NewLine()
	go outanims.GrayGifOutput(x, y, cout, frames, 100, stop)

	utils.PressAnyKey(*cin, nil)
	close(stop)
}

func run() int {

	resetTTY()
	os.Setenv("COLORTERM", "truecolor")

	log.SetOutput(io.Discard)
	f, err := os.OpenFile("/var/log/gologin.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)

	cin := stdin.InitCIn()
	go cin.MainLoop(os.Stdin, os.Stdout, int(os.Stdin.Fd()), true)
	cout := stdout.InitCOut()
	StopSyncLoop := make(chan struct{})
	go cout.SyncLoop(StopSyncLoop, 30)

	var finiOnce sync.Once
	var inputStopOnce sync.Once
	var outputSyncStopOnce sync.Once

	defer finiOnce.Do(cout.Fini)
	defer inputStopOnce.Do(cin.Stop)
	defer outputSyncStopOnce.Do(func() { StopSyncLoop <- struct{}{} })

	{
		x, y := cout.GetCursorPosition()
		cout.SetCursorYPosition(cout.CursorLine + strings.Count(startScreen, "\n"))
		go cout.SlowTextOut(x, y, startScreen, false, 1)
		cout.NewLine()
	}

	cout.TextOut("Login: ")
	username := gstrings.ReadString(&cout, &cin, nil)
	switch username {
	case "shutdown":
		return 1000
	case "wrong_password_test":
		wrongPassword(&cout, &cin)
		return 0
	case "colors_test":
		cout.TextOutLn("\x1b[30m Black foreground\x1b[0m        \x1b[40m Black background")
		cout.TextOutLn("\x1b[31m DarkRed foreground\x1b[0m       \x1b[41m DarkRed background")
		cout.TextOutLn("\x1b[32m DarkGreen foreground\x1b[0m     \x1b[42m DarkGreen background")
		cout.TextOutLn("\x1b[33m DarkYellow foreground\x1b[0m    \x1b[43m DarkYellow background")
		cout.TextOutLn("\x1b[34m DarkBlue foreground\x1b[0m      \x1b[44m DarkBlue background")
		cout.TextOutLn("\x1b[35m DarkMagenta foreground\x1b[0m   \x1b[45m DarkMagenta background")
		cout.TextOutLn("\x1b[36m DarkCyan foreground\x1b[0m      \x1b[46m DarkCyan background")
		cout.TextOutLn("\x1b[37m Gray/LightGray foreground\x1b[0m \x1b[47m Gray/LightGray background")

		cout.TextOutLn("\x1b[90m LightBlack (Gray) foreground\x1b[0m     \x1b[100m LightBlack background")
		cout.TextOutLn("\x1b[91m LightRed foreground\x1b[0m             \x1b[101m LightRed background")
		cout.TextOutLn("\x1b[92m LightGreen foreground\x1b[0m           \x1b[102m LightGreen background")
		cout.TextOutLn("\x1b[93m LightYellow foreground\x1b[0m          \x1b[103m LightYellow background")
		cout.TextOutLn("\x1b[94m LightBlue foreground\x1b[0m            \x1b[104m LightBlue background")
		cout.TextOutLn("\x1b[95m LightMagenta foreground\x1b[0m         \x1b[105m LightMagenta background")
		cout.TextOutLn("\x1b[96m LightCyan foreground\x1b[0m            \x1b[106m LightCyan background")
		cout.TextOutLn("\x1b[97m White foreground\x1b[0m                \x1b[107m White background")
		utils.PressAnyKey(cin, nil)
		return 0
	case "exit":
		return 0
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
			wrongPassword(&cout, &cin)
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

	finiOnce.Do(cout.Fini)
	inputStopOnce.Do(cin.Stop)
	outputSyncStopOnce.Do(func() { StopSyncLoop <- struct{}{} })

	fmt.Println(welcomeScreen)

	//cmd := exec.Command(os.Getenv("SHELL"))
	//cmd.Env = os.Environ()
	//cmd.Stdin = os.Stdin
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	//cmd.Run()

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

package main

import (
	_ "embed"
	"os"
	"os/exec"

	"github.com/HypoxiE/go-login-system/pkg/stdin"
	"github.com/HypoxiE/go-login-system/pkg/stdout"
	"github.com/HypoxiE/go-login-system/pkg/utils"
)

var (
	//go:embed text_templates/welcome_screen.txt
	welcomeScreen string
	//go:embed text_templates/start_screen.txt
	startScreen string

	//go:embed text_templates/wrong_password.txt
	wrongPasswordGif string
)

func main() {

	cin := stdin.InitCIn()
	cout := stdout.InitCOut()
	go cin.MainLoop(os.Stdin, os.Stdout, int(os.Stdin.Fd()), true)
	defer func() {
		cin.StopOutput <- struct{}{}
	}()
	defer func() {
		if r := recover(); r != nil {
			cin.StopOutput <- struct{}{}
			panic(r)
		}
	}()
	defer func() {
		cout.Fini()
	}()
	defer func() {
		if r := recover(); r != nil {
			cout.Fini()
			panic(r)
		}
	}()

	go cout.SlowTextOut(startScreen)
	cout.NewLine()
	cout.Sync()

	// input username
	cout.TextOut("Login: ")
	cout.ShowCursor()
	cout.Sync()

	x := cout.CursorColumn
	for c := range cin.LastSymbol {
		if c == '\n' || c == '\r' {
			break
		} else if c == 0x7f || c == 0x08 {
			if cout.CursorColumn > x {
				cout.CursorColumn--
				cout.TextOut(" ")
				cout.CursorColumn--
				cout.ShowCursor()
				cout.Sync()
			}
			continue
		}
		cout.TextOut(string(c))
		cout.ShowCursor()
		cout.Sync()
	}

	cout.HideCursor()
	cout.NewLine()
	cout.Sync()

	username := cin.GetForLine()

	if username == "shutdown" {
		exec.Command("shutdown", "-h", "now").Run()
	}

	utils.PressAnyKey(cin, &cout)

	//_, err := pam.StartFunc("login", username, core.StartPam)

	//if err != nil {
	//	fmt.Println("pam_start error:", err)
	//	utils.PressAnyKey(cin, nil)
	//	os.Exit(1)
	//}

	//if err = t.Authenticate(pam.Silent); err != nil {
	//	if err.Error() == "Authentication failure" {
	//		fmt.Println("Wrong password, baka!")

	//		done := make(chan struct{})
	//		go func() {
	//			reader := bufio.NewReader(os.Stdin)
	//			reader.ReadByte()
	//			close(done)
	//		}()
	//		utils.PaintAsciiGif(wrongPasswordGif, done)

	//	} else if err.Error() == "User not known to the underlying authentication module" {
	//		fmt.Println("auth failed:", err)
	//		utils.PressAnyKey(false)
	//	} else {
	//		fmt.Println("auth failed:", err)
	//		utils.PressAnyKey(false)
	//	}
	//	os.Exit(1)
	//}

	//if err := t.AcctMgmt(pam.Silent); err != nil {
	//	fmt.Println("AcctMgmt error:", err)
	//	utils.PressAnyKey(false)
	//	os.Exit(1)
	//}

	//usr, err := user.Lookup(username)
	//if err != nil {
	//	fmt.Println("lookup error:", err)
	//	utils.PressAnyKey(false)
	//	os.Exit(1)
	//}

	//groups, err := usr.GroupIds()
	//if err != nil {
	//	fmt.Println("GroupIds error:", err)
	//	utils.PressAnyKey(false)
	//	os.Exit(1)
	//}
	//gids := make([]int, len(groups))
	//for i, g := range groups {
	//	gi, _ := strconv.Atoi(g)
	//	gids[i] = gi
	//}
	//syscall.Setgroups(gids)

	//if err := t.SetCred(pam.EstablishCred); err != nil {
	//	fmt.Println("set_cred error:", err)
	//	utils.PressAnyKey(false)
	//	os.Exit(1)
	//}

	//if err := t.OpenSession(pam.Silent); err != nil {
	//	fmt.Println("pam_open_session error:", err)
	//	utils.PressAnyKey(false)
	//	os.Exit(1)
	//}

	//env, err := t.GetEnvList()
	//if err != nil {
	//	fmt.Println("getenv error:", err)
	//	utils.PressAnyKey(false)
	//	os.Exit(1)
	//}

	//os.Chdir(usr.HomeDir)

	//initialize.InitEnv(env)

	//uid, _ := strconv.Atoi(usr.Uid)
	//gid, _ := strconv.Atoi(usr.Gid)
	//syscall.Setgid(gid)
	//syscall.Setuid(uid)

	//if username == "hypoxie" {
	//	fmt.Println(welcomeScreen)
	//}

	//syscall.Exec(os.Getenv("SHELL"), []string{os.Getenv("SHELL")}, os.Environ())

	//t.CloseSession(pam.Silent)
	//t.SetCred(pam.DeleteCred)
}

package common

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

	"github.com/manifoldco/promptui"
)

func OpenURL(url string) error {
	var cmd string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	case "windows":
		cmd = "start"
	default:
		return fmt.Errorf("Unsupport open OS:%s, Please open the URL manually: %s", runtime.GOOS, url)
	}

	return exec.Command(cmd, url).Run()
}

func IsPromptQuit(err error) bool {
	switch err {
	case promptui.ErrInterrupt, io.EOF:
		return true
	case nil:
		return false
	default:
		return true
	}
}

func IsSelectQuit(err error) bool {
	switch err {
	case promptui.ErrEOF:
		return true
	case nil, promptui.ErrInterrupt:
		return false
	default:
		return true
	}
}

var clear map[string]func() //create a map for storing clear funcs

func Init() {
	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["darwin"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	clear, ok := clear[runtime.GOOS]
	if ok {
		clear()
	}
}

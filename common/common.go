package common

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/manifoldco/promptui"
)

func OpenURL(url string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", url).Run()
	case "windows":
		return exec.Command("cmd", "/c", "start", url).Run()
	}

	return nil
}

func DownloadMardownFile(url string) (string, error) {
	res, err := http.Get(url + ".md")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	file, err := os.Create("tmp.md")
	if err != nil {
		return "", err
	}

	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

func RemoveFile(file string) error {
	return os.Remove(file)
}

func DisplayMarkdown(file string) error {
	md, err := exec.Command("mdv", file).Output()
	if err != nil {
		return err
	}

	fmt.Println(string(md))
	return nil
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

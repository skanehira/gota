package common

import (
	"fmt"
	"io"
	"os/exec"
	"runtime"

	"github.com/manifoldco/promptui"
	"github.com/rivo/tview"
)

func OpenURL(url string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", url).Run()
	case "windows":
		return exec.Command("cmd", "/c", "start", url).Run()
	case "linux":
		return exec.Command("xdg-open", url).Run()
	}

	return nil
}

func ViewTerminal(body string) error {
	app := tview.NewApplication()
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)

	fmt.Fprintf(textView, "%s ", body)

	textView.SetBorder(true)
	if err := app.SetRoot(textView, true).Run(); err != nil {
		return err
	}

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

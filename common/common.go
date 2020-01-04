package common

import (
	"io"
	"os/exec"
	"runtime"

	"github.com/awesome-gocui/gocui"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
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
	g, err := gocui.NewGui(gocui.OutputNormal, false)
	if err != nil {
		return errors.Wrap(err, "error starting the interactive UI")
	}
	defer g.Close()

	ui, err := newUi(g)
	if err != nil {
		return err
	}

	ui.setContent(body)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
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

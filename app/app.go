package app

import (
	"errors"
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/skanehira/gota/common"
	"github.com/skanehira/gota/qiita"
)

var ErrNotFound = errors.New("not found")

type App struct {
	Prompt   *Prompt
	Selecter *Selecter
	Qiita    *qiita.Qiita
}

type Prompt struct {
	promptui.Prompt
}

type Selecter struct {
	promptui.Select
	Result qiita.Result
}

func New() *App {
	prompt := &Prompt{
		promptui.Prompt{
			Label: "Search",
			Templates: &promptui.PromptTemplates{
				Prompt:  "{{ . }} ",
				Valid:   "{{ . | cyan }} ",
				Invalid: "{{ . | red }} ",
				Success: `{{ "Searching..." | green }} `,
			},
		},
	}

	selecter := &Selecter{
		Select: promptui.Select{
			Label: "Result",
			Templates: &promptui.SelectTemplates{
				Label:    `{{ . }}`,
				Active:   `{{ .Title | red }}`,
				Inactive: ` {{ .Title | cyan }}`,
				Selected: `{{ .Title | yellow }} {{ .Url | green}}`,
				Details: `--------- Details ----------
{{ "URL:" | green }}	{{ .Url | green }}
{{ "Created:" | blue }}	{{ .CreatedAt | blue }}
{{ "Updated:" | yellow }}	{{ .UpdatedAt | yellow }}
{{ "User:" | red }}	{{ .User.Id | red }}
{{ "Tags:" | cyan }}	{{ range .Tags }}{{ .Name | cyan }} {{end}}

`,
			},
			Size: 20,
		},
	}

	selecter.Searcher = func(input string, index int) bool {
		item := selecter.Result.Items[index]
		name := strings.Replace(strings.ToLower(item.Title), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	return &App{
		Prompt:   prompt,
		Selecter: selecter,
		Qiita:    qiita.New(""),
	}
}

func (app *App) Run() error {

	for {

		// wait input search word
		word, err := app.prompt()
		if err != nil {
			return err
		}

		// new search conditions
		cond := &qiita.SearchCondition{
			Title:   word,
			Page:    1,
			PerPage: 100,
		}

		result, err := app.search(cond)

		if err != nil {
			if err == ErrNotFound {
				fmt.Println(err)
				continue
			}
			return err
		}

		if err := app.selecter(result); err != nil {
			return err
		}

	}
}

func (app *App) prompt() (string, error) {
	word, err := app.Prompt.Run()

	if common.IsPromptQuit(err) {
		return "", err
	}

	return word, nil
}

func (app *App) search(cond *qiita.SearchCondition) (qiita.Result, error) {
	result, err := app.Qiita.SearchItems(cond)

	if err != nil {
		fmt.Printf("search failed: %s\n", err)
		return result, err
	}

	if len(result.Items) < 1 {
		return result, ErrNotFound
	}
	return result, nil
}

func (app *App) selecter(result qiita.Result) error {
	app.Selecter.Items = result.Items
	app.Selecter.Result = result

	for {
		i, _, err := app.Selecter.Run()

		if common.IsSelectQuit(err) {
			return err
		}

		if err == promptui.ErrInterrupt {
			return nil
		}

		if err := app.confirm(result.Items[i].Url); err != nil {
			return err
		}
	}
}

func (app *App) confirm(url string) error {
	result := []struct {
		Selected string
	}{
		{"browser"},
		{"terminal"},
	}

	confirm := promptui.Select{
		Label: "Open browser or display terminal?",
		Templates: &promptui.SelectTemplates{
			Label:    `{{ . }}`,
			Active:   `{{ .Selected | red }}`,
			Inactive: ` {{ .Selected | cyan }}`,
			Selected: ` `,
		},
		Items: result,
		Size:  2,
	}

	i, _, err := confirm.Run()

	if common.IsSelectQuit(err) {
		return err
	}

	if err == promptui.ErrInterrupt {
		return nil
	}

	switch result[i].Selected {
	case "browser":
		if err := common.OpenURL(url); err != nil {
			fmt.Printf("open failed: %s", err)
			return err
		}
	case "terminal":
		file, err := common.DownloadMardownFile(url)
		if err != nil {
			fmt.Printf("open failed: %s", err)
			return err
		}

		return common.ViewMarkdown(file)
	}

	return nil
}

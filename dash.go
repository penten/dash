package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/penten/pocket"
	"html/template"
	"io"
	"os"
)

type layout struct {
	Body template.HTML
}

func displayPage(out io.Writer, body string) error {
	t, err := template.ParseFiles("templates/base.html")
	if err != nil {
		return err
	}

	return t.Execute(out, layout{Body: template.HTML(body)})
}

func nth(n int) func() bool {
	var i = 1
	return func() bool {
		if i == n {
			i = 1
			return true
		}
		i++
		return false
	}
}

var funcs = template.FuncMap{
	"fourth": nth(4),
}

func displayTemplate(file string, data interface{}) (string, error) {
	t, err := template.New("main").Funcs(funcs).ParseFiles("templates/" + file)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	err = t.ExecuteTemplate(&b, file, data)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func displayPocket() (string, error) {
	archive, erra := pocket.GetArticles(pocket_appkey, pocket_apptoken, map[string]string{"favorite": "0", "state": "archive", "count": "8", "sort": "newest"})
	favorite, errb := pocket.GetArticles(pocket_appkey, pocket_apptoken, map[string]string{"favorite": "1", "state": "archive", "count": "4", "sort": "newest"})

	if erra != nil || errb != nil {
		return "", errors.New("Unable to fetch articles from pocket")
	}

	a, erra := displayTemplate("archive.html", archive)
	b, errb := displayTemplate("favorites.html", favorite)

	if erra != nil || errb != nil {
		return "", errors.New("Unable to fetch display template")
	}

	return a + b, nil
}

func main() {
	f, err := os.Create("output.html")
	if err != nil {
		fmt.Printf("Error opening output file: " + err.Error())
		os.Exit(1)
	}

	body, err := displayPocket()
	if err != nil {
		fmt.Printf("Error displaying pocket: " + err.Error())
		os.Exit(1)
	}

	err = displayPage(f, body)
	if err != nil {
		fmt.Printf("Error displaying template: " + err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

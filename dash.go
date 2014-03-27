package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/penten/pocket"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"time"
)

type config struct {
	Pocket_appkey   string
	Pocket_apptoken string
}

var c config

type layout struct {
	Body  template.HTML
	Title string
}

func displayPage(out io.Writer, body string, title string) error {
	t, err := template.ParseFiles("templates/base.html")
	if err != nil {
		return err
	}

	return t.Execute(out, layout{Body: template.HTML(body), Title: title})
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

func currentTime() string {
	return time.Now().Format("Jan 2, 2006 at 3:04pm")
}

var funcs = template.FuncMap{
	"fourth":      nth(4),
	"currentTime": currentTime,
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
	archive, erra := pocket.GetArticles(c.Pocket_appkey, c.Pocket_apptoken, map[string]string{"favorite": "0", "state": "archive", "count": "8", "sort": "newest"})
	favorite, errb := pocket.GetArticles(c.Pocket_appkey, c.Pocket_apptoken, map[string]string{"favorite": "1", "state": "archive", "count": "4", "sort": "newest"})

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

func loadConfig() error {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &c)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading configuration " + err.Error())
		os.Exit(1)
	}

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

	err = displayPage(f, body, "Dashboard")
	if err != nil {
		fmt.Printf("Error displaying template: " + err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

type (
	ZenQuote struct {
		Quote  string `json:"q"`
		Author string `json:"a"`
		Html   string `json:"h"`
	}
)

const (
	readmeFile = "README.md"
)

var (
	replacer = regexp.MustCompile("(<!-- dailyquote:start -->\n)(.*\n)+(<!-- dailyquote:end -->\n)")
)

func getZenQuote() (ZenQuote, error) {
	if request, err := http.NewRequest(http.MethodGet, "https://zenquotes.io/api/random", nil); err == nil {
		request.Header.Set("Content-Type", "application/json; charset=utf-8")
		client := &http.Client{}
		if response, err := client.Do(request); err == nil {
			if response.StatusCode != http.StatusOK {
				return ZenQuote{}, errors.New("invalid status code for fetching response")
			}
			if body, err := io.ReadAll(response.Body); err == nil {
				var zenquote []ZenQuote
				err := json.Unmarshal(body, &zenquote)
				return zenquote[0], err
			} else {
				return ZenQuote{}, err
			}
		} else {
			return ZenQuote{}, err
		}
	} else {
		return ZenQuote{}, err
	}
}

func main() {
	var zenquote ZenQuote
	if zen, err := getZenQuote(); err == nil {
		zenquote = zen
	} else {
		panic(err)
	}

	if data, err := os.ReadFile(readmeFile); err == nil {
		quote := fmt.Sprintf("<p>%s</p>\n\n<p>%s</p>\n", zenquote.Quote, zenquote.Author)
		content := string(data)
		content = replacer.ReplaceAllString(content, fmt.Sprintf("${1}%s${3}", quote))
		if err := os.WriteFile(readmeFile, []byte(content), 0644); err != nil {
			panic(err)
		}
	} else {
		panic(err)
	}
}

package main

import (
	"log"
	"net/http"
	"net/url"

	proxy ".."
)

type searcher struct{}

func (searcher) Suggest(string) ([]string, error) {
	return []string{
		"sears",
		"search engines",
		"search engine",
		"search",
		"sears.com",
		"seattle times",
	}, nil
}

func (searcher) Search(string) (string, error) {
	return "https://google.com", nil
}

func main() {
	url, err := url.Parse("http://localhost:2020")
	if err != nil {
		panic(err)
	}
	searchProxy := proxy.New(proxy.Configuration{
		Searcher: searcher{},
		Host:     url,
		Metadata: proxy.Metadata{
			Favicon:     "https://assets-cdn.github.com/favicon.ico",
			Name:        "Example",
			Description: "example search",
		},
	})
	log.Fatal(http.ListenAndServe(":2020", searchProxy))
}

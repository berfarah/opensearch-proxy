package main

import (
	"log"
	"net/http"
	"net/url"
	"os"

	proxy ".."
	paper "../dropboxpaper"
)

func main() {
	url, err := url.Parse("http://localhost:2020")
	if err != nil {
		panic(err)
	}
	searchProxy := proxy.New(proxy.Configuration{
		Searcher: paper.New(os.Getenv("TOKEN")),
		Host:     url,
		Metadata: proxy.Metadata{
			Favicon:     "https://paper.dropbox.com/favicon.ico?v3",
			Name:        "Paper",
			Description: "Search Paper",
		},
	})
	log.Fatal(http.ListenAndServe(":2020", searchProxy))
}

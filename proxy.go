package proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

// Searcher ...
type Searcher interface {
	// Search returns a redirect URL.
	Search(string) (string, error)
	// Suggest must return search suggestions.
	Suggest(string) ([]string, error)
}

// Configuration is the configuration for the search engine.
type Configuration struct {
	Host *url.URL
	Searcher
	Metadata
	// Some field for the logic - func or interface
}

// Metadata is the data around an OpenSearch search engine.
type Metadata struct {
	// Favicon path to be used as the search engine image.
	Favicon string
	// Name of the search engine.
	Name string
	// Description for the search engine.
	Description string
}

// Proxy is the server that acts as a proxy for OpenSearch.
type Proxy struct {
	*http.ServeMux
	host       *url.URL
	searcher   Searcher
	root       []byte
	definition []byte
}

// New creates a new proxy.
func New(c Configuration) *Proxy {
	definition, err := generateDefinition(c)
	if err != nil {
		panic(fmt.Sprintf("Invalid definition template: %v", err))
	}
	root, err := generateRoot(c)
	if err != nil {
		panic(fmt.Sprintf("Invalid root template: %v", err))
	}

	return &Proxy{
		searcher:   c.Searcher,
		host:       c.Host,
		root:       root,
		definition: definition,
	}
}

func logHTTPError(r *http.Request, msg string, err error) {
	log.Printf("ERROR: %s %s: %s (%v)\n", r.Method, r.URL.Path, msg, err)
}

func (p *Proxy) handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	if _, err := w.Write(p.root); err != nil {
		logHTTPError(r, "Couldn't serve root", err)
	}
}

func (p *Proxy) handleDefinition(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/opensearchdescription+xml")
	if _, err := w.Write(p.definition); err != nil {
		logHTTPError(r, "Couldn't serve opensearch definition", err)
	}
}

func (p *Proxy) handleSearch(w http.ResponseWriter, r *http.Request) {
	queryString := r.URL.Query().Get("q")
	url, err := p.searcher.Search(queryString)
	if err != nil {
		logHTTPError(r, "Couldn't fetch search", err)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func (p *Proxy) handleSuggest(w http.ResponseWriter, r *http.Request) {
	queryString := r.URL.Query().Get("q")
	suggestions, err := p.searcher.Suggest(queryString)
	if err != nil {
		logHTTPError(r, "Couldn't fetch suggestions", err)
		return
	}

	var urls []string
	for _, suggestion := range suggestions {
		urls = append(urls, fmt.Sprintf("%s/search?q=%s", p.host, url.QueryEscape(suggestion)))
	}

	response, err := json.Marshal([]interface{}{queryString, suggestions, make([]string, len(suggestions)), urls})
	if err != nil {
		logHTTPError(r, "Couldn't convert suggestions to JSON", err)
		return
	}

	if _, err := w.Write(response); err != nil {
		logHTTPError(r, "Couldn't serve json", err)
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()
	router.HandleFunc("/", p.handleRoot).Methods("GET")
	router.HandleFunc("/search.xml", p.handleDefinition).Methods("GET")
	router.HandleFunc("/suggest", p.handleSuggest).Methods("GET")
	router.HandleFunc("/search", p.handleSearch).Methods("GET")
	router.ServeHTTP(w, r)
}

var _ http.Handler = &Proxy{}

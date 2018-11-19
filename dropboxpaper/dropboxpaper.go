package dropboxpaper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// doc is a dropbox paper document.
type doc struct {
	ID             string    `json:"doc_id"`
	Title          string    `json:"title"`
	Snippets       []string  `json:"snippets"`
	CreatedDate    time.Time `json:"created_date"`
	LastEditedDate time.Time `json:"last_edited_date"`
	CreatorName    string    `json:"creator_name"`
	LastEditorName string    `json:"last_editor_name"`
}

// Searcher is a Dropbox Paper searcher.
type Searcher struct {
	apiToken string
	cache    map[string]doc
}

// New creates a new Searcher.
func New(APIToken string) *Searcher {
	return &Searcher{
		apiToken: APIToken,
		cache:    make(map[string]doc),
	}
}

func (s *Searcher) query(query string) ([]doc, error) {
	type queryPayload struct {
		Query string `json:"query"`
		Limit int    `json:"limit"`
	}

	payload := queryPayload{query, 10}
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("json marshal: %s", err.Error())
	}

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/2/paper/docs/search", bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("new request: %s", err.Error())
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiToken))
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %s", err.Error())
	}

	var response struct {
		Documents []doc `json:"docs"`
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode error: %s", err.Error())
	}
	defer res.Body.Close()

	return response.Documents, nil
}

func (s *Searcher) Search(query string) (string, error) {
	doc, ok := s.cache[query]
	if !ok {
		return "", fmt.Errorf("No doc found for %s", query)
	}
	return fmt.Sprintf("https://paper.dropbox.com/doc/%s", doc.ID), nil
}

// Suggest returns search results from Dropbox.
func (s *Searcher) Suggest(query string) ([]string, error) {
	docs, err := s.query(query)
	if err != nil {
		return nil, err
	}
	var results []string
	for _, doc := range docs {
		s.cache[doc.Title] = doc
		results = append(results, doc.Title)
	}

	return results, nil
}

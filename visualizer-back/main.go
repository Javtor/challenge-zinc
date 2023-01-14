package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

const (
	searchEndpoint = "http://localhost:4080/api/email/_search"
	authUsername   = "admin"
	authPassword   = "Complexpass#123"
)

var httpClient = &http.Client{}

type SearchRequest struct {
	SearchType string      `json:"search_type"`
	Query      interface{} `json:"query"`
	From       int         `json:"from"`
	MaxResults int         `json:"max_results"`
}

type SearchResponse struct {
	Emails []Email `json:"emails"`
	Total  int     `json:"total"`
	From   int     `json:"from"`
	Size   int     `json:"size"`
}

type Email struct {
	Subject string   `json:"subject"`
	From    string   `json:"from"`
	To      []string `json:"to"`
	Body    string   `json:"body"`
}

type ZincSearchResponse struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source struct {
				Header map[string][]string `json:"header"`
				Body   string              `json:"body"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func createSearchRequest(query string, page, perPage int) (*http.Request, error) {
	reqBody := SearchRequest{
		SearchType: "querystring",
		Query: map[string]interface{}{
			"term": query,
		},
		From:       (page - 1) * perPage,
		MaxResults: perPage,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %w", err)
	}

	req, err := http.NewRequest("POST", searchEndpoint, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.SetBasicAuth(authUsername, authPassword)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func decodeSearchResponse(resp *http.Response) (ZincSearchResponse, error) {
	var zincResponse ZincSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&zincResponse); err != nil {
		return ZincSearchResponse{}, fmt.Errorf("error decoding search results: %w", err)
	}
	return zincResponse, nil
}

func mapZincResponseToSearchResponse(zincResponse ZincSearchResponse, page, perPage int) (SearchResponse, error) {
	var searchResults SearchResponse
	searchResults.Total = zincResponse.Hits.Total.Value
	searchResults.From = (page - 1) * perPage
	searchResults.Size = perPage
	searchResults.Emails = make([]Email, len(zincResponse.Hits.Hits))
	for i, hit := range zincResponse.Hits.Hits {
		searchResults.Emails[i].Subject = hit.Source.Header["Subject"][0]
		searchResults.Emails[i].From = hit.Source.Header["From"][0]
		searchResults.Emails[i].To = hit.Source.Header["To"]
		searchResults.Emails[i].Body = hit.Source.Body
	}
	return searchResults, nil
}

func performSearch(query string, page, perPage int) (SearchResponse, error) {
	req, err := createSearchRequest(query, page, perPage)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("error creating search request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("error making search request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return SearchResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	zincResponse, err := decodeSearchResponse(resp)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("error decoding search results: %w", err)
	}

	searchResults, err := mapZincResponseToSearchResponse(zincResponse, page, perPage)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("error mapping zinc response to search response: %w", err)
	}
	return searchResults, nil
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	perPage, err := strconv.Atoi(r.URL.Query().Get("per_page"))
	if err != nil || perPage < 1 {
		perPage = 10
	}

	results, err := performSearch(query, page, perPage)
	if err != nil {
		render.JSON(w, r, err.Error())
		return
	}
	render.JSON(w, r, &results)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/search", searchHandler)

	http.ListenAndServe(":3000", r)
}

package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Search struct is used to search emails that match the given query
type Search struct {
	// searchEndpoint is the endpoint where the search request is sent
	searchEndpoint string
	// authUsername is the username used for basic auth
	authUsername string
	// authPassword is the password used for basic auth
	authPassword string
}

// NewSearch returns a new instance of Search struct with the default values
func NewSearch() *Search {
	return &Search{
		searchEndpoint: "http://zinc:4080/api/email/_search",
		authUsername:   "admin",
		authPassword:   "Complexpass#123",
	}
}

// ZincSearchRequest represents the body of a search request to Zincsearch
type ZincSearchRequest struct {
	SearchType string            `json:"search_type"`
	Query      map[string]string `json:"query"`
	From       int               `json:"from"`
	MaxResults int               `json:"max_results"`
}

// ZincSearchResponse represents the response from a search request to Zincsearch
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

// EmailSearchResponse represents the response that will be returned
type EmailSearchResponse struct {
	Emails []Email `json:"emails"`
	Total  int     `json:"total"`
	From   int     `json:"from"`
	Size   int     `json:"size"`
}

// Email represents an email
type Email struct {
	Subject string   `json:"subject"`
	From    string   `json:"from"`
	To      []string `json:"to"`
	Body    string   `json:"body"`
}

// createSearchRequest creates an HTTP request for a search with the given query and page/perPage parameters
func (s *Search) createSearchRequest(query string, page, perPage int) (*http.Request, error) {
	reqBody := ZincSearchRequest{
		SearchType: "querystring",
		Query: map[string]string{
			"term": query,
		},
		From:       (page - 1) * perPage,
		MaxResults: perPage,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %w", err)
	}

	req, err := http.NewRequest("POST", s.searchEndpoint, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.SetBasicAuth(s.authUsername, s.authPassword)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// mapZincResponseToEmailSearchResponse maps the ZincSearchResponse to an EmailSearchResponse struct
func (s *Search) mapZincResponseToEmailSearchResponse(zincResponse ZincSearchResponse, page, perPage int) (EmailSearchResponse, error) {
	var searchResults EmailSearchResponse
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

// PerformSearch performs a search for emails that match the given query, and returns a struct containing the search results.
// It takes in 3 parameters:
// - query: a string representing the search query
// - page: an int representing the page number of the results
// - perPage: an int representing the number of results per page
func (s *Search) PerformSearch(query string, page, perPage int) (EmailSearchResponse, error) {
	log.Printf("Performing search for query: %s", query)
	req, err := s.createSearchRequest(query, page, perPage)
	if err != nil {
		return EmailSearchResponse{}, err
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Error performing search: %v", err)
		return EmailSearchResponse{}, fmt.Errorf("error performing search: %w", err)
	}
	log.Printf("Search completed successfully")
	defer resp.Body.Close()

	var zincResponse ZincSearchResponse
	err = json.NewDecoder(resp.Body).Decode(&zincResponse)
	if err != nil {
		return EmailSearchResponse{}, fmt.Errorf("error decoding response: %w", err)
	}

	return s.mapZincResponseToEmailSearchResponse(zincResponse, page, perPage)
}

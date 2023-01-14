package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Search struct {
	searchEndpoint string
	authUsername   string
	authPassword   string
}

func NewSearch() *Search {
	return &Search{
		searchEndpoint: "http://localhost:4080/api/email/_search",
		authUsername:   "admin",
		authPassword:   "Complexpass#123",
	}
}

type ZincSearchRequest struct {
	SearchType string            `json:"search_type"`
	Query      map[string]string `json:"query"`
	From       int               `json:"from"`
	MaxResults int               `json:"max_results"`
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

type EmailSearchResponse struct {
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

func (s *Search) PerformSearch(query string, page, perPage int) (EmailSearchResponse, error) {
	req, err := s.createSearchRequest(query, page, perPage)
	if err != nil {
		return EmailSearchResponse{}, err
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return EmailSearchResponse{}, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return EmailSearchResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var zincResp ZincSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&zincResp); err != nil {
		return EmailSearchResponse{}, fmt.Errorf("error decoding response: %w", err)
	}

	return s.mapZincResponseToEmailSearchResponse(zincResp, page, perPage)
}

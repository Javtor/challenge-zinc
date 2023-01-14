package search

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateSearchRequest(t *testing.T) {
	// Create a new search client
	s := NewSearch()

	// Test creating a search request with valid parameters
	req, err := s.createSearchRequest("test query", 1, 10)
	if err != nil {
		t.Errorf("Error creating search request: %v", err)
	}
	if req.Method != "POST" {
		t.Errorf("Expected request method to be 'POST', got %s", req.Method)
	}
	if req.URL.String() != s.searchEndpoint {
		t.Errorf("Expected request URL to be %s, got %s", s.searchEndpoint, req.URL)
	}
	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type header to be 'application/json', got %s", req.Header.Get("Content-Type"))
	}
}

func TestMapZincResponseToEmailSearchResponse(t *testing.T) {
	// Create a new search client
	s := NewSearch()

	// Test mapping a valid ZincSearchResponse
	zincResp := ZincSearchResponse{}
	zincResp.Hits.Total.Value = 1
	zincResp.Hits.Hits = []struct {
		Source struct {
			Header map[string][]string `json:"header"`
			Body   string              `json:"body"`
		} `json:"_source"`
	}{
		{
			Source: struct {
				Header map[string][]string `json:"header"`
				Body   string              `json:"body"`
			}{
				Header: map[string][]string{
					"Subject": {"Test Subject"},
					"From":    {"test@example.com"},
					"To":      {"recipient@example.com"},
				},
				Body: "Test body",
			},
		},
	}
	resp, err := s.mapZincResponseToEmailSearchResponse(zincResp, 1, 10)
	if err != nil {
		t.Errorf("Error mapping ZincSearchResponse: %v", err)
	}
	if len(resp.Emails) != 1 {
		t.Errorf("Expected 1 email, got %d", len(resp.Emails))
	}
	if resp.Emails[0].Subject != "Test Subject" {
		t.Errorf("Expected email subject to be 'Test Subject', got %s", resp.Emails[0].Subject)
	}
	if resp.Emails[0].From != "test@example.com" {
		t.Errorf("Expected email from to be 'test@example.com', got %s", resp.Emails[0].From)
	}
	if resp.Emails[0].To[0] != "recipient@example.com" {
		t.Errorf("Expected email to to be 'recipient@example.com', got %s", resp.Emails[0].To[0])
	}
	if resp.Emails[0].Body != "Test body" {
		t.Errorf("Expected email body to be 'Test body', got %s", resp.Emails[0].Body)
	}

	// Test mapping a ZincSearchResponse with no hits
	zincResp.Hits.Total.Value = 0
	zincResp.Hits.Hits = nil
	resp, err = s.mapZincResponseToEmailSearchResponse(zincResp, 1, 10)
	if err != nil {
		t.Errorf("Error mapping ZincSearchResponse: %v", err)
	}
	if len(resp.Emails) != 0 {
		t.Errorf("Expected 0 emails, got %d", len(resp.Emails))
	}
	if resp.Total != 0 {
		t.Errorf("Expected total to be 0, got %d", resp.Total)
	}
}

func TestPerformSearch(t *testing.T) {
	// Create a test server that will respond to our search request
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ZincSearchRequest
		json.NewDecoder(r.Body).Decode(&req)

		resp := ZincSearchResponse{}
		resp.Hits.Total.Value = 1
		resp.Hits.Hits = []struct {
			Source struct {
				Header map[string][]string `json:"header"`
				Body   string              `json:"body"`
			} `json:"_source"`
		}{
			{
				Source: struct {
					Header map[string][]string `json:"header"`
					Body   string              `json:"body"`
				}{
					Header: map[string][]string{
						"Subject": {"Test Subject"},
						"From":    {"test@example.com"},
						"To":      {"recipient@example.com"},
					},
					Body: "Test body",
				},
			},
		}

		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	// Create a new search client with the test server's URL
	s := &Search{searchEndpoint: ts.URL}

	// Perform a search and check the response
	resp, err := s.PerformSearch("test query", 1, 10)
	if err != nil {
		t.Errorf("Error performing search: %v", err)
	}
	if len(resp.Emails) != 1 {
		t.Errorf("Expected 1 email, got %d", len(resp.Emails))
	}
	if resp.Emails[0].Subject != "Test Subject" {
		t.Errorf("Expected email subject to be 'Test Subject', got %s", resp.Emails[0].Subject)
	}
	if resp.Emails[0].Body != "Test body" {
		t.Errorf("Expected email body to be 'Test body', got %s", resp.Emails[0].Body)
	}
}

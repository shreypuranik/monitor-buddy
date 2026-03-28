package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCrawlURLs_Up(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	config := URLsConfig{
		Websites: []WebsiteConfig{
			{Name: "Test Site", URL: srv.URL, SiteID: 1},
		},
	}

	results := crawlURLs(config, srv.Client())

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].StatusCode != 200 {
		t.Errorf("expected status 200, got %d", results[0].StatusCode)
	}
	if !results[0].Up {
		t.Error("expected Up to be true")
	}
}

func TestCrawlURLs_Down(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	config := URLsConfig{
		Websites: []WebsiteConfig{
			{Name: "Broken Site", URL: srv.URL, SiteID: 2},
		},
	}

	results := crawlURLs(config, srv.Client())

	if results[0].StatusCode != 500 {
		t.Errorf("expected status 500, got %d", results[0].StatusCode)
	}
	if results[0].Up {
		t.Error("expected Up to be false")
	}
}

func TestCrawlURLs_Unreachable(t *testing.T) {
	// Use an address that immediately refuses connections.
	config := URLsConfig{
		Websites: []WebsiteConfig{
			{Name: "Dead Site", URL: "http://127.0.0.1:1", SiteID: 3},
		},
	}

	results := crawlURLs(config, &http.Client{})

	if results[0].StatusCode != 0 {
		t.Errorf("expected status 0, got %d", results[0].StatusCode)
	}
	if results[0].Up {
		t.Error("expected Up to be false")
	}
}

func TestCrawlURLs_MultiplePreservesOrder(t *testing.T) {
	codes := []int{200, 404, 503}
	servers := make([]*httptest.Server, len(codes))
	for i, code := range codes {
		code := code
		servers[i] = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(code)
		}))
		defer servers[i].Close()
	}

	config := URLsConfig{
		Websites: []WebsiteConfig{
			{Name: "A", URL: servers[0].URL, SiteID: 1},
			{Name: "B", URL: servers[1].URL, SiteID: 2},
			{Name: "C", URL: servers[2].URL, SiteID: 3},
		},
	}

	results := crawlURLs(config, servers[0].Client())

	for i, want := range codes {
		if results[i].StatusCode != want {
			t.Errorf("result[%d]: expected status %d, got %d", i, want, results[i].StatusCode)
		}
	}
}

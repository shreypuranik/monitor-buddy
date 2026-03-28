package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// buildTestConfig returns a URLsConfig pointing at a local test server.
func buildTestConfig(url string) URLsConfig {
	return URLsConfig{
		Websites: []WebsiteConfig{
			{Name: "Test Site", URL: url, SiteID: 1},
		},
	}
}

func TestAPIStatusHandler(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer upstream.Close()

	config := buildTestConfig(upstream.URL)
	client := upstream.Client()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		statuses := crawlURLs(config, client)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(statuses)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}

	var statuses []SiteStatus
	if err := json.Unmarshal(rec.Body.Bytes(), &statuses); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if len(statuses) != 1 {
		t.Fatalf("expected 1 status, got %d", len(statuses))
	}
	if statuses[0].StatusCode != 200 {
		t.Errorf("expected status_code 200, got %d", statuses[0].StatusCode)
	}
	if !statuses[0].Up {
		t.Error("expected up to be true")
	}
}

func TestUIHandler(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer upstream.Close()

	config := buildTestConfig(upstream.URL)
	client := upstream.Client()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		statuses := crawlURLs(config, client)
		w.Header().Set("Content-Type", "text/html")
		tmplErr := uiTmpl.Execute(w, statuses)
		if tmplErr != nil {
			t.Errorf("template execution error: %v", tmplErr)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Monitor Buddy") {
		t.Error("expected page title 'Monitor Buddy' in response body")
	}
	if !strings.Contains(body, "Test Site") {
		t.Error("expected site name 'Test Site' in response body")
	}
	if !strings.Contains(body, "card up") {
		t.Error("expected green 'up' card in response body")
	}
}

func TestLoadConfigParsesYAML(t *testing.T) {
	config, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig() error: %v", err)
	}
	if len(config.Websites) == 0 {
		t.Error("expected at least one website in urls.yaml")
	}
	for _, site := range config.Websites {
		if site.Name == "" {
			t.Error("website missing name")
		}
		if site.URL == "" {
			t.Error("website missing url")
		}
		if site.SiteID == 0 {
			t.Error("website missing site_id")
		}
	}
}

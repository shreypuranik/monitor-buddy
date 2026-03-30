package main

import (
	"net/http"
	"sync"
	"time"
)

// HTTPClient is satisfied by *http.Client, making crawlURLs testable.
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type WebsiteConfig struct {
	Name   string `yaml:"name"`
	URL    string `yaml:"url"`
	SiteID int    `yaml:"site_id"`
}

type URLsConfig struct {
	Websites []WebsiteConfig `yaml:"websites"`
}

type SiteStatus struct {
	Name           string `json:"name"`
	URL            string `json:"url"`
	SiteID         int    `json:"site_id"`
	StatusCode     int    `json:"status_code"`
	Up             bool   `json:"up"`
	ResponseTimeMs int64  `json:"response_time_ms"`
}

func crawlURLs(config URLsConfig, client HTTPClient) []SiteStatus {
	results := make([]SiteStatus, len(config.Websites))
	var wg sync.WaitGroup

	for i, site := range config.Websites {
		wg.Add(1)
		go func(idx int, s WebsiteConfig) {
			defer wg.Done()
			status := SiteStatus{
				Name:   s.Name,
				URL:    s.URL,
				SiteID: s.SiteID,
			}
			start := time.Now()
			resp, err := client.Get(s.URL)
			status.ResponseTimeMs = time.Since(start).Milliseconds()
			if err != nil {
				status.StatusCode = 0
				status.Up = false
			} else {
				resp.Body.Close()
				status.StatusCode = resp.StatusCode
				status.Up = resp.StatusCode >= 200 && resp.StatusCode < 400
			}
			results[idx] = status
		}(i, site)
	}

	wg.Wait()
	return results
}

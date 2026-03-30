package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

var uiTmpl = template.Must(template.New("ui").Parse(uiTemplate))

func loadConfig() (URLsConfig, error) {
	data, err := os.ReadFile("urls.yaml")
	if err != nil {
		return URLsConfig{}, err
	}
	var config URLsConfig
	err = yaml.Unmarshal(data, &config)
	return config, err
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	client := &http.Client{}

	http.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		statuses := crawlURLs(config, client)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(statuses)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		statuses := crawlURLs(config, client)
		tmpl := template.Must(template.New("ui").Parse(uiTemplate))
		w.Header().Set("Content-Type", "text/html")
		tmpl.Execute(w, statuses)
	})

	fmt.Println("Monitor Buddy running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

const uiTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Monitor Buddy</title>
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
      background: #0f1117;
      color: #e2e8f0;
      min-height: 100vh;
      padding: 2rem;
    }
    h1 {
      font-size: 1.8rem;
      font-weight: 700;
      margin-bottom: 0.25rem;
      color: #f8fafc;
    }
    .subtitle {
      color: #64748b;
      margin-bottom: 2rem;
      font-size: 0.9rem;
    }
    .grid {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 1rem;
    }
    .card {
      border-radius: 12px;
      padding: 1.25rem 1.5rem;
      border: 1px solid;
      transition: transform 0.1s;
    }
    .card:hover { transform: translateY(-2px); }
    .card.up {
      background: #052e16;
      border-color: #16a34a;
    }
    .card.down {
      background: #2d0a0a;
      border-color: #dc2626;
    }
    .card-header {
      display: flex;
      align-items: center;
      gap: 0.6rem;
      margin-bottom: 0.75rem;
    }
    .dot {
      width: 12px;
      height: 12px;
      border-radius: 50%;
      flex-shrink: 0;
    }
    .dot.up { background: #22c55e; box-shadow: 0 0 6px #22c55e; }
    .dot.down { background: #ef4444; box-shadow: 0 0 6px #ef4444; }
    .site-name {
      font-weight: 600;
      font-size: 1rem;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    .site-url {
      font-size: 0.78rem;
      color: #94a3b8;
      margin-bottom: 0.75rem;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    .badge {
      display: inline-block;
      padding: 0.2rem 0.6rem;
      border-radius: 999px;
      font-size: 0.8rem;
      font-weight: 600;
    }
    .badge.up { background: #16a34a33; color: #4ade80; border: 1px solid #16a34a; }
    .badge.down { background: #dc262633; color: #f87171; border: 1px solid #dc2626; }
  </style>
</head>
<body>
  <h1>Monitor Buddy</h1>
  <p class="subtitle">Live status for all monitored sites</p>
  <div class="grid">
    {{range .}}
    <div class="card {{if .Up}}up{{else}}down{{end}}">
      <div class="card-header">
        <div class="dot {{if .Up}}up{{else}}down{{end}}"></div>
        <span class="site-name">{{.Name}}</span>
      </div>
      <div class="site-url">{{.URL}}</div>
      {{if eq .StatusCode 0}}
        <span class="badge down">Unreachable</span>
      {{else}}
        <span class="badge {{if .Up}}up{{else}}down{{end}}">HTTP {{.StatusCode}}</span>
        <span class="badge {{if .Up}}up{{else}}down{{end}}">{{.ResponseTimeMs}}ms</span>
      {{end}}
    </div>
    {{end}}
  </div>
</body>
</html>`

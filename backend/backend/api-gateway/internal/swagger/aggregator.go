package swagger

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type ServiceSpec struct {
	Name    string // "messenger", "media"
	URL     string // "http://messenger-service:8081/swagger/doc.json"
	Prefix  string // "/api/v1/messenger"
}

type Aggregator struct {
	services []ServiceSpec
	cache    map[string]interface{}
	mu       sync.RWMutex
	ttl      time.Duration
	lastFetch time.Time
}

func NewAggregator(services []ServiceSpec) *Aggregator {
	return &Aggregator{
		services: services,
		cache:    make(map[string]interface{}),
		ttl:      5 * time.Minute,
	}
}

func (a *Aggregator) RegisterRoutes(e *echo.Echo) {
	e.GET("/swagger/doc.json", a.GetMergedSpec)
	e.GET("/swagger/*", a.ServeUI)
}

func (a *Aggregator) GetMergedSpec(c echo.Context) error {
	a.mu.RLock()
	if time.Since(a.lastFetch) < a.ttl && a.cache != nil {
		cached := a.cache
		a.mu.RUnlock()
		return c.JSON(http.StatusOK, cached)
	}
	a.mu.RUnlock()

	merged, err := a.merge()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	a.mu.Lock()
	a.cache = merged
	a.lastFetch = time.Now()
	a.mu.Unlock()

	return c.JSON(http.StatusOK, merged)
}
func (a *Aggregator) merge() (map[string]interface{}, error) {
	merged := map[string]interface{}{
		"swagger": "2.0",
		"info": map[string]interface{}{
			"title":       "Suda Messenger API",
			"description": "Unified API — Messenger, Media, Transactions",
			"version":     "1.0",
		},
		"host":     "localhost:8080",
		"basePath": "/",
		"schemes":  []string{"http"},
		"paths":       map[string]interface{}{},
		"definitions": map[string]interface{}{},
		"tags":        []interface{}{},
		"securityDefinitions": map[string]interface{}{
			"BearerAuth": map[string]interface{}{
				"type": "apiKey",
				"in":   "header",
				"name": "Authorization",
			},
		},
	}

	mergedPaths := merged["paths"].(map[string]interface{})
	mergedDefs := merged["definitions"].(map[string]interface{})
	mergedTags := []interface{}{}

	for _, svc := range a.services {
		spec, err := a.fetchSpec(svc.URL)
		if err != nil {
			log.Warn().Str("service", svc.Name).Err(err).Msg("Failed to fetch swagger spec")
			continue
		}

		// Merge paths with prefix
		if paths, ok := spec["paths"].(map[string]interface{}); ok {
			for path, methods := range paths {
				fullPath := svc.Prefix + path
				mergedPaths[fullPath] = methods
			}
		}

		// Merge definitions (Swagger 2.0)
		if definitions, ok := spec["definitions"].(map[string]interface{}); ok {
			for name, def := range definitions {
				mergedDefs[name] = def
			}
		}

		// Merge components.schemas (OpenAPI 3.0 fallback)
		if components, ok := spec["components"].(map[string]interface{}); ok {
			if schemas, ok := components["schemas"].(map[string]interface{}); ok {
				for name, schema := range schemas {
					mergedDefs[name] = schema
				}
			}
		}

		// Merge tags
		if tags, ok := spec["tags"].([]interface{}); ok {
			mergedTags = append(mergedTags, tags...)
		}
	}

	merged["tags"] = mergedTags
	return merged, nil
}

func (a *Aggregator) fetchSpec(url string) (map[string]interface{}, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d from %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	var spec map[string]interface{}
	if err := json.Unmarshal(body, &spec); err != nil {
		return nil, fmt.Errorf("parse json: %w", err)
	}

	return spec, nil
}

// ServeUI отдаёт Swagger UI
func (a *Aggregator) ServeUI(c echo.Context) error {
	return c.HTML(http.StatusOK, swaggerUIHTML)
}

const swaggerUIHTML = `<!DOCTYPE html>
<html>
<head>
    <title>Suda API Documentation</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
    <style>
        body { margin: 0; background: #fafafa; }
        .swagger-ui .topbar { display: none; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
        SwaggerUIBundle({
            url: "/swagger/doc.json",
            dom_id: '#swagger-ui',
            deepLinking: true,
            presets: [
                SwaggerUIBundle.presets.apis,
                SwaggerUIBundle.SwaggerUIStandalonePreset
            ],
            layout: "BaseLayout"
        });
    </script>
</body>
</html>`
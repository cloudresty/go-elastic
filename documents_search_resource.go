package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v9/esapi"
)

// SearchResource provides search operations across indices
type SearchResource struct {
	client *Client
}

// extractIndicesFromOptions extracts indices from search options, defaults to "_all"
func extractIndicesFromOptions(options []SearchOption) []string {
	// Create a temporary map to collect indices
	temp := make(map[string]any)
	for _, option := range options {
		option(temp)
	}

	if indices, exists := temp["indices"]; exists {
		switch v := indices.(type) {
		case string:
			return []string{v}
		case []string:
			return v
		case []any:
			result := make([]string, len(v))
			for i, idx := range v {
				result[i] = fmt.Sprint(idx)
			}
			return result
		}
	}

	// Default to all indices
	return []string{"_all"}
}

// Scroll returns a SearchScroll resource for scroll operations
func (sr *SearchResource) Scroll(options ...SearchOption) *SearchScroll {
	return &SearchScroll{
		client: sr.client,
	}
}

// Search performs a search across the specified indices
func (sr *SearchResource) Search(ctx context.Context, query map[string]any, options ...SearchOption) (*SearchResponse, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	// Build search body using existing BuildSearchQuery function
	searchBody := BuildSearchQuery(query, options...)

	bodyBytes, err := json.Marshal(searchBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search query: %w", err)
	}

	// Extract indices from options, default to "_all"
	indices := extractIndicesFromOptions(options)

	req := esapi.SearchRequest{
		Index: indices,
		Body:  bytes.NewReader(bodyBytes),
	}

	res, err := req.Do(ctx, sr.client.client)
	if err != nil {
		sr.client.config.Logger.Error("Search failed - indices: %s, error: %s", strings.Join(indices, ","), err.Error())
		return nil, fmt.Errorf("search request failed: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			sr.client.config.Logger.Warn("Failed to close response body - error: %s", err.Error())
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		sr.client.config.Logger.Error("Search failed - indices: %s, status: %s, response: %s", strings.Join(indices, ","), res.Status(), string(bodyBytes))
		return nil, fmt.Errorf("search failed: %s - %s", res.Status(), string(bodyBytes))
	}

	var searchResponse SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&searchResponse); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	sr.client.config.Logger.Debug("Search completed successfully - indices: %s, hits: %d, total: %d, took: %d", strings.Join(indices, ","), len(searchResponse.Hits.Hits), int(searchResponse.Hits.Total.Value), searchResponse.Took)

	return &searchResponse, nil
}

// Count returns the number of documents matching the query
func (sr *SearchResource) Count(ctx context.Context, query map[string]any, options ...SearchOption) (int64, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	var bodyBytes []byte
	var err error

	if query != nil {
		countBody := map[string]any{"query": query}
		bodyBytes, err = json.Marshal(countBody)
		if err != nil {
			return 0, fmt.Errorf("failed to marshal count query: %w", err)
		}
	}

	// Extract indices from options, default to "_all"
	indices := extractIndicesFromOptions(options)

	req := esapi.CountRequest{
		Index: indices,
	}

	if bodyBytes != nil {
		req.Body = bytes.NewReader(bodyBytes)
	}

	res, err := req.Do(ctx, sr.client.client)
	if err != nil {
		sr.client.config.Logger.Error("Count failed - indices: %s, error: %s", strings.Join(indices, ","), err.Error())
		return 0, fmt.Errorf("count request failed: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			sr.client.config.Logger.Warn("Failed to close response body - error: %s", err.Error())
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		sr.client.config.Logger.Error("Count failed - indices: %s, status: %s, response: %s", strings.Join(indices, ","), res.Status(), string(bodyBytes))
		return 0, fmt.Errorf("count failed: %s - %s", res.Status(), string(bodyBytes))
	}

	var countResponse struct {
		Count int64 `json:"count"`
	}

	if err := json.NewDecoder(res.Body).Decode(&countResponse); err != nil {
		return 0, fmt.Errorf("failed to decode count response: %w", err)
	}

	sr.client.config.Logger.Debug("Count completed successfully - indices: %s, count: %d", strings.Join(indices, ","), int(countResponse.Count))

	return countResponse.Count, nil
}

// startScrollSearch initiates a scroll search and returns the initial response
func (sr *SearchResource) startScrollSearch(ctx context.Context, query map[string]any, scrollTime time.Duration, options ...SearchOption) (*SearchResponse, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	// Build search body using existing BuildSearchQuery function
	searchBody := BuildSearchQuery(query, options...)

	// Set default scroll size if not specified
	if _, hasSize := searchBody["size"]; !hasSize {
		searchBody["size"] = 1000
	}

	bodyBytes, err := json.Marshal(searchBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search query: %w", err)
	}

	// Extract indices from options, default to "_all"
	indices := extractIndicesFromOptions(options)

	req := esapi.SearchRequest{
		Index:  indices,
		Body:   bytes.NewReader(bodyBytes),
		Scroll: scrollTime,
	}

	res, err := req.Do(ctx, sr.client.client)
	if err != nil {
		sr.client.config.Logger.Error("Scroll search failed - indices: %s, error: %s", strings.Join(indices, ","), err.Error())
		return nil, fmt.Errorf("scroll search request failed: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			sr.client.config.Logger.Warn("Failed to close response body - error: %s", err.Error())
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		sr.client.config.Logger.Error("Scroll search failed - indices: %s, status: %s, response: %s", strings.Join(indices, ","), res.Status(), string(bodyBytes))
		return nil, fmt.Errorf("scroll search failed: %s - %s", res.Status(), string(bodyBytes))
	}

	var searchResponse SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&searchResponse); err != nil {
		return nil, fmt.Errorf("failed to decode scroll search response: %w", err)
	}

	sr.client.config.Logger.Debug("Scroll search started successfully - indices: %s, scroll_id: %s, initial_hits: %d, total: %d, took: %d", strings.Join(indices, ","), searchResponse.ScrollID, len(searchResponse.Hits.Hits), int(searchResponse.Hits.Total.Value), searchResponse.Took)

	return &searchResponse, nil
}

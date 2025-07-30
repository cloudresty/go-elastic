package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/elastic/go-elasticsearch/v9/esapi"
)

// SearchScroll provides scroll-related operations for search
type SearchScroll struct {
	client *Client
}

// Start starts a scroll search for processing large result sets
func (ss *SearchScroll) Start(ctx context.Context, query map[string]any, scrollTime time.Duration, options ...SearchOption) (*SearchResponse, error) {
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

	res, err := req.Do(ctx, ss.client.client)
	if err != nil {
		return nil, fmt.Errorf("scroll search request failed: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			ss.client.config.Logger.Warn("Failed to close response body - error: %s", err.Error())
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("scroll search failed: %s - %s", res.Status(), string(bodyBytes))
	}

	var searchResponse SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&searchResponse); err != nil {
		return nil, fmt.Errorf("failed to decode scroll search response: %w", err)
	}

	return &searchResponse, nil
}

// Continue continues a scroll search using the scroll ID
func (ss *SearchScroll) Continue(ctx context.Context, scrollID string, scrollTime time.Duration) (*SearchResponse, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	req := esapi.ScrollRequest{
		ScrollID: scrollID,
		Scroll:   scrollTime,
	}

	res, err := req.Do(ctx, ss.client.client)
	if err != nil {
		ss.client.config.Logger.Error("Scroll continue failed - scroll_id: %s, error: %s", scrollID, err.Error())
		return nil, fmt.Errorf("scroll continue request failed: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			ss.client.config.Logger.Warn("Failed to close response body - error: %s", err.Error())
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		ss.client.config.Logger.Error("Scroll continue failed - scroll_id: %s, status: %s, response: %s", scrollID, res.Status(), string(bodyBytes))
		return nil, fmt.Errorf("scroll continue failed: %s - %s", res.Status(), string(bodyBytes))
	}

	var searchResponse SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&searchResponse); err != nil {
		return nil, fmt.Errorf("failed to decode scroll continue response: %w", err)
	}

	ss.client.config.Logger.Debug("Scroll continue completed successfully - scroll_id: %s, hits: %d, took: %d", scrollID, len(searchResponse.Hits.Hits), searchResponse.Took)

	return &searchResponse, nil
}

// Clear clears a specific scroll context
func (ss *SearchScroll) Clear(ctx context.Context, scrollID string) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	req := esapi.ClearScrollRequest{
		ScrollID: []string{scrollID},
	}

	res, err := req.Do(ctx, ss.client.client)
	if err != nil {
		return fmt.Errorf("clear scroll request failed: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			ss.client.config.Logger.Warn("Failed to close response body - error: %s", err.Error())
		}
	}()

	if res.IsError() {
		ss.client.config.Logger.Warn("Clear scroll failed - scroll_id: %s, status: %s", scrollID, res.Status())
		return fmt.Errorf("clear scroll failed: %s", res.Status())
	}

	ss.client.config.Logger.Debug("Scroll cleared successfully - scroll_id: %s", scrollID)

	return nil
}

// ClearAll clears all scroll contexts
func (ss *SearchScroll) ClearAll(ctx context.Context) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	req := esapi.ClearScrollRequest{
		ScrollID: []string{"_all"},
	}

	res, err := req.Do(ctx, ss.client.client)
	if err != nil {
		return fmt.Errorf("clear all scrolls request failed: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			ss.client.config.Logger.Warn("Failed to close response body - error: %s", err.Error())
		}
	}()

	if res.IsError() {
		return fmt.Errorf("clear all scrolls failed: %s", res.Status())
	}

	return nil
}

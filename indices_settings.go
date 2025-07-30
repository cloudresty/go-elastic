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

// IndexSettings provides settings-related operations for an index
type IndexSettings struct {
	client    *Client
	indexName string
}

// Get retrieves the index settings
func (is *IndexSettings) Get(ctx context.Context) (map[string]any, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	req := esapi.IndicesGetSettingsRequest{
		Index: []string{is.indexName},
	}

	res, err := req.Do(ctx, is.client.client)
	if err != nil {
		return nil, fmt.Errorf("failed to get index settings: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			is.client.config.Logger.Warn("Failed to close response body - error: %s",
				err.Error())
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("failed to get settings for index '%s': %s - %s", is.indexName, res.Status(), string(bodyBytes))
	}

	var result map[string]any
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode settings response: %w", err)
	}

	// Extract settings for this specific index
	if indexData, exists := result[is.indexName]; exists {
		if indexMap, ok := indexData.(map[string]any); ok {
			if settings, exists := indexMap["settings"]; exists {
				if settingsMap, ok := settings.(map[string]any); ok {
					return settingsMap, nil
				}
			}
		}
	}

	return result, nil
}

// Update updates the index settings
func (is *IndexSettings) Update(ctx context.Context, settings map[string]any) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	bodyBytes, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	req := esapi.IndicesPutSettingsRequest{
		Index: []string{is.indexName},
		Body:  bytes.NewReader(bodyBytes),
	}

	res, err := req.Do(ctx, is.client.client)
	if err != nil {
		return fmt.Errorf("failed to update index settings: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			is.client.config.Logger.Warn("Failed to close response body - error: %s",
				err.Error())
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return fmt.Errorf("failed to update settings for index '%s': %s - %s", is.indexName, res.Status(), string(bodyBytes))
	}

	return nil
}

// Refresh refreshes the index settings (re-reads from cluster state)
func (is *IndexSettings) Refresh(ctx context.Context) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	req := esapi.IndicesRefreshRequest{
		Index: []string{is.indexName},
	}

	res, err := req.Do(ctx, is.client.client)
	if err != nil {
		return fmt.Errorf("failed to refresh index settings: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			is.client.config.Logger.Warn("Failed to close response body - error: %s",
				err.Error())
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return fmt.Errorf("failed to refresh settings for index '%s': %s - %s", is.indexName, res.Status(), string(bodyBytes))
	}

	return nil
}

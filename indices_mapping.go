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

// IndexMapping provides mapping-related operations for an index
type IndexMapping struct {
	client    *Client
	indexName string
}

// Get retrieves the index mapping
func (im *IndexMapping) Get(ctx context.Context) (map[string]any, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	req := esapi.IndicesGetMappingRequest{
		Index: []string{im.indexName},
	}

	res, err := req.Do(ctx, im.client.client)
	if err != nil {
		return nil, fmt.Errorf("failed to get index mapping: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			im.client.config.Logger.Warn("Failed to close response body - error: %s",
				err.Error())
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("failed to get mapping for index '%s': %s - %s", im.indexName, res.Status(), string(bodyBytes))
	}

	var result map[string]any
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode mapping response: %w", err)
	}

	// Extract mapping for this specific index
	if indexData, exists := result[im.indexName]; exists {
		if indexMap, ok := indexData.(map[string]any); ok {
			if mappings, exists := indexMap["mappings"]; exists {
				if mappingsMap, ok := mappings.(map[string]any); ok {
					return mappingsMap, nil
				}
			}
		}
	}

	return result, nil
}

// Update updates the index mapping
func (im *IndexMapping) Update(ctx context.Context, mapping map[string]any) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	bodyBytes, err := json.Marshal(mapping)
	if err != nil {
		return fmt.Errorf("failed to marshal mapping: %w", err)
	}

	req := esapi.IndicesPutMappingRequest{
		Index: []string{im.indexName},
		Body:  bytes.NewReader(bodyBytes),
	}

	res, err := req.Do(ctx, im.client.client)
	if err != nil {
		return fmt.Errorf("failed to update index mapping: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			im.client.config.Logger.Warn("Failed to close response body - error: %s",
				err.Error())
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return fmt.Errorf("failed to update mapping for index '%s': %s - %s", im.indexName, res.Status(), string(bodyBytes))
	}

	return nil
}

// Create creates the index mapping (only works if index doesn't exist)
func (im *IndexMapping) Create(ctx context.Context, mapping map[string]any) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	// Check if index exists first
	exists, err := im.client.Indices().Exists(ctx, im.indexName)
	if err != nil {
		return fmt.Errorf("failed to check if index exists: %w", err)
	}

	if exists {
		return fmt.Errorf("cannot create mapping for existing index '%s', use Update() instead", im.indexName)
	}

	// Create index with mapping
	return im.client.Indices().Create(ctx, im.indexName, mapping)
}

// GetField retrieves the mapping for a specific field
func (im *IndexMapping) GetField(ctx context.Context, fieldName string) (map[string]any, error) {
	mapping, err := im.Get(ctx)
	if err != nil {
		return nil, err
	}

	// Navigate through the mapping structure to find the field
	if properties, exists := mapping["properties"]; exists {
		if propsMap, ok := properties.(map[string]any); ok {
			if fieldMapping, exists := propsMap[fieldName]; exists {
				if fieldMap, ok := fieldMapping.(map[string]any); ok {
					return fieldMap, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("field '%s' not found in mapping", fieldName)
}

// AddField adds a new field to the mapping
func (im *IndexMapping) AddField(ctx context.Context, fieldName string, fieldMapping map[string]any) error {
	updateMapping := map[string]any{
		"properties": map[string]any{
			fieldName: fieldMapping,
		},
	}

	return im.Update(ctx, updateMapping)
}

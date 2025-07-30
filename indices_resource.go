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

// IndexResource provides index management operations
type IndexResource struct {
	client *Client
	name   string
}

// Name returns the index name
func (ir *IndexResource) Name() string {
	return ir.name
}

// Create creates the index with optional mapping
func (ir *IndexResource) Create(ctx context.Context, mapping map[string]any) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	// Check if index already exists
	exists, err := ir.Exists(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if index exists: %w", err)
	}
	if exists {
		return fmt.Errorf("index '%s' already exists", ir.name)
	}

	var body io.Reader
	if mapping != nil {
		bodyBytes, err := json.Marshal(mapping)
		if err != nil {
			return fmt.Errorf("failed to marshal mapping: %w", err)
		}
		body = bytes.NewReader(bodyBytes)
	}

	req := esapi.IndicesCreateRequest{
		Index: ir.name,
		Body:  body,
	}

	res, err := req.Do(ctx, ir.client.client)
	if err != nil {
		ir.client.config.Logger.Error("Failed to create index - index: %s, error: %s", ir.name, err.Error())
		return fmt.Errorf("failed to create index: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			ir.client.config.Logger.Warn("Failed to close response body - error: %s", err.Error())
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		ir.client.config.Logger.Error("Failed to create index - index: %s, status: %s, response: %s", ir.name, res.Status(), string(bodyBytes))
		return fmt.Errorf("failed to create index '%s': %s - %s", ir.name, res.Status(), string(bodyBytes))
	}

	ir.client.config.Logger.Info("Index created successfully - index: %s", ir.name)

	return nil
}

// Delete deletes the index
func (ir *IndexResource) Delete(ctx context.Context) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	req := esapi.IndicesDeleteRequest{
		Index: []string{ir.name},
	}

	res, err := req.Do(ctx, ir.client.client)
	if err != nil {
		ir.client.config.Logger.Error("Failed to delete index - index: %s, error: %s", ir.name, err.Error())
		return fmt.Errorf("failed to delete index: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			ir.client.config.Logger.Warn("Failed to close response body - error: %s", err.Error())
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		ir.client.config.Logger.Error("Failed to delete index - index: %s, status: %s, response: %s", ir.name, res.Status(), string(bodyBytes))
		return fmt.Errorf("failed to delete index '%s': %s - %s", ir.name, res.Status(), string(bodyBytes))
	}

	ir.client.config.Logger.Info("Index deleted successfully - index: %s", ir.name)

	return nil
}

// Exists checks if the index exists
func (ir *IndexResource) Exists(ctx context.Context) (bool, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	req := esapi.IndicesExistsRequest{
		Index: []string{ir.name},
	}

	res, err := req.Do(ctx, ir.client.client)
	if err != nil {
		return false, fmt.Errorf("failed to check index existence: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			ir.client.config.Logger.Warn("Failed to close response body - error: %s", err.Error())
		}
	}()

	return res.StatusCode == 200, nil
}

// Settings returns an IndexSettings resource for this index
func (ir *IndexResource) Settings() *IndexSettings {
	return &IndexSettings{
		client:    ir.client,
		indexName: ir.name,
	}
}

// Mapping returns an IndexMapping resource for this index
func (ir *IndexResource) Mapping() *IndexMapping {
	return &IndexMapping{
		client:    ir.client,
		indexName: ir.name,
	}
}

// Document returns a Document resource for this index
func (ir *IndexResource) Document() *Document {
	return &Document{
		client: ir.client,
		index:  ir.name,
	}
}

// Search performs a search on this index
func (ir *IndexResource) Search(ctx context.Context, query map[string]any, options ...SearchOption) (*SearchResponse, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	// Use the existing search functionality from the Index type
	idx := &Index{
		client: ir.client,
		name:   ir.name,
	}

	return idx.Search(ctx, query, options...)
}

// Count returns the document count for this index
func (ir *IndexResource) Count(ctx context.Context, query map[string]any) (int64, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	// Use the existing count functionality from the Index type
	idx := &Index{
		client: ir.client,
		name:   ir.name,
	}

	return idx.Count(ctx, query)
}

// Index settings helpers

// DefaultIndexSettings returns commonly used index settings
func DefaultIndexSettings() map[string]any {
	return map[string]any{
		"number_of_shards":   1,
		"number_of_replicas": 1,
		"refresh_interval":   "1s",
	}
}

// Resource-level operations for single index efficiency

// Close closes this index (makes it unavailable but preserves data)
func (ir *IndexResource) Close(ctx context.Context) error {
	return ir.client.Indices().Close(ctx, ir.name)
}

// Open opens this previously closed index
func (ir *IndexResource) Open(ctx context.Context) error {
	return ir.client.Indices().Open(ctx, ir.name)
}

// Refresh forces a refresh of this index
func (ir *IndexResource) Refresh(ctx context.Context) error {
	return ir.client.Indices().Refresh(ctx, ir.name)
}

// Flush forces a flush of this index to disk
func (ir *IndexResource) Flush(ctx context.Context) error {
	return ir.client.Indices().Flush(ctx, ir.name)
}

// Stats returns statistics for this index
func (ir *IndexResource) Stats(ctx context.Context) (map[string]any, error) {
	return ir.client.Indices().Stats(ctx, ir.name)
}

// Clone creates a copy of this index
func (ir *IndexResource) Clone(ctx context.Context, targetIndex string) error {
	return ir.client.Indices().Clone(ctx, ir.name, targetIndex)
}

// Reindex copies documents from this index to a target index
func (ir *IndexResource) Reindex(ctx context.Context, targetIndex string, options ...map[string]any) error {
	return ir.client.Indices().Reindex(ctx, ir.name, targetIndex, options...)
}

// Shrink reduces the number of shards in this index
func (ir *IndexResource) Shrink(ctx context.Context, targetIndex string, targetShards int) error {
	return ir.client.Indices().Shrink(ctx, ir.name, targetIndex, targetShards)
}

// Analyze tests how text is analyzed in this index
func (ir *IndexResource) Analyze(ctx context.Context, text, analyzer string) (map[string]any, error) {
	return ir.client.Indices().Analyze(ctx, ir.name, text, analyzer)
}

// Aliases returns all aliases pointing to this index
func (ir *IndexResource) Aliases(ctx context.Context) (map[string]any, error) {
	allAliases, err := ir.client.Indices().Aliases(ctx)
	if err != nil {
		return nil, err
	}

	// Filter to only aliases for this index
	result := make(map[string]any)
	if indexData, exists := allAliases[ir.name]; exists {
		result[ir.name] = indexData
	}

	return result, nil
}

// AddAlias adds an alias pointing to this index
func (ir *IndexResource) AddAlias(ctx context.Context, aliasName string) error {
	return ir.client.Indices().Alias(ctx, aliasName, ir.name)
}

// RemoveAlias removes an alias from this index
func (ir *IndexResource) RemoveAlias(ctx context.Context, aliasName string) error {
	return ir.client.Indices().RemoveAlias(ctx, aliasName, ir.name)
}

// Rollover creates a new index when conditions are met (assuming this index is an alias)
func (ir *IndexResource) Rollover(ctx context.Context, options ...map[string]any) (map[string]any, error) {
	return ir.client.Indices().Rollover(ctx, ir.name, options...)
}

package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v9/esapi"
)

// Index wraps an Elasticsearch index with enhanced functionality
// This provides a convenient API for Elasticsearch index operations
type Index struct {
	client *Client
	name   string
}

// Name returns the index name
func (idx *Index) Name() string {
	return idx.name
}

// Mapping returns an IndexMapping resource for this index
// This provides access to mapping operations using the resource-oriented pattern
func (idx *Index) Mapping() *IndexMapping {
	return &IndexMapping{
		client:    idx.client,
		indexName: idx.name,
	}
}

// IndexMany indexes multiple documents
func (idx *Index) IndexMany(ctx context.Context, documents []map[string]any) (*BulkResponse, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second) //nolint:ineffassign,staticcheck
		defer cancel()
	}

	if len(documents) == 0 {
		return nil, fmt.Errorf("no documents provided")
	}

	// Build bulk operations
	operations := make([]map[string]any, 0, len(documents)*2)
	for _, doc := range documents {
		// Enhance document
		enhanced := idx.client.enhanceDocument(doc)

		// Get or generate ID
		var docID string
		if id, exists := enhanced["_id"]; exists {
			if idStr, ok := id.(string); ok {
				docID = idStr
			}
		}

		// Add index operation
		indexOp := map[string]any{
			"index": map[string]any{
				"_index": idx.name,
			},
		}
		if docID != "" {
			indexOp["index"].(map[string]any)["_id"] = docID
		}

		operations = append(operations, indexOp, enhanced)
	}

	// Use the new BulkResource API
	bulkResource := &BulkResource{
		client: idx.client,
		index:  idx.name,
	}

	response, err := bulkResource.ExecuteRaw(ctx, operations)
	if err != nil {
		idx.client.config.Logger.Error("Failed to index documents - error: %s, index: %s, count: %d", err.Error(), idx.name, len(documents))
		return nil, err
	}

	idx.client.config.Logger.Debug("Documents indexed successfully - index: %s, count: %d", idx.name, len(documents))

	return response, nil
}

// Search performs a search query
func (idx *Index) Search(ctx context.Context, query map[string]any, options ...SearchOption) (*SearchResponse, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second) //nolint:ineffassign,staticcheck
		defer cancel()
	}

	searchResource := &SearchResource{
		client: idx.client,
	}
	response, err := searchResource.Search(ctx, query, append(options, WithIndices(idx.name))...)
	if err != nil {
		idx.client.config.Logger.Error("Failed to search documents - error: %s, index: %s", err.Error(), idx.name)
		return nil, err
	}

	idx.client.config.Logger.Debug("Search completed successfully - index: %s, hits: %d", idx.name, response.Hits.Total.Value)

	return response, nil
}

// Count counts documents matching a query
func (idx *Index) Count(ctx context.Context, query map[string]any) (int64, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second) //nolint:ineffassign,staticcheck
		defer cancel()
	}

	// Use the _count API
	countQuery := map[string]any{
		"query": query,
	}

	// Convert query to JSON
	queryBytes, err := json.Marshal(countQuery)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal count query: %w", err)
	}

	// Make the count request using the underlying client
	req := esapi.CountRequest{
		Index: []string{idx.name},
		Body:  bytes.NewReader(queryBytes),
	}

	response, err := req.Do(ctx, idx.client.GetClient())
	if err != nil {
		idx.client.config.Logger.Error("Failed to count documents - error: %s, index: %s", err.Error(), idx.name)
		return 0, fmt.Errorf("failed to count documents: %w", err)
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			idx.client.config.Logger.Warn("Failed to close response body - error: %s", err.Error())
		}
	}()

	if response.IsError() {
		return 0, fmt.Errorf("count request failed: %s", response.String())
	}

	var countResponse struct {
		Count int64 `json:"count"`
	}

	if err := json.NewDecoder(response.Body).Decode(&countResponse); err != nil {
		return 0, fmt.Errorf("failed to decode count response: %w", err)
	}

	idx.client.config.Logger.Debug("Documents counted successfully - index: %s, count: %d", idx.name, int(countResponse.Count))

	return countResponse.Count, nil
}

// Delete deletes the index
func (idx *Index) Delete(ctx context.Context) error {
	return idx.client.Indices().Delete(ctx, idx.name)
}

// Exists checks if the index exists
func (idx *Index) Exists(ctx context.Context) (bool, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second) //nolint:ineffassign,staticcheck
		defer cancel()
	}

	return idx.client.Indices().Exists(ctx, idx.name)
}

package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/cloudresty/emit"
	"github.com/elastic/go-elasticsearch/v9/esapi"
)

// IndicesService methods

// Create creates a new index with optional mapping
func (s *IndicesService) Create(ctx context.Context, indexName string, mapping map[string]any) error {
	indexResource := &IndexResource{
		client: s.client,
		name:   indexName,
	}
	return indexResource.Create(ctx, mapping)
}

// Delete deletes an index
func (s *IndicesService) Delete(ctx context.Context, indexName string) error {
	indexResource := &IndexResource{
		client: s.client,
		name:   indexName,
	}
	return indexResource.Delete(ctx)
}

// Exists checks if an index exists
func (s *IndicesService) Exists(ctx context.Context, indexName string) (bool, error) {
	indexResource := &IndexResource{
		client: s.client,
		name:   indexName,
	}
	return indexResource.Exists(ctx)
}

// Get returns an IndexResource for direct access to index operations
func (s *IndicesService) Get(indexName string) *IndexResource {
	return &IndexResource{
		client: s.client,
		name:   indexName,
	}
}

// List returns detailed information about all indices
func (s *IndicesService) List(ctx context.Context) ([]IndexInfo, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	req := esapi.CatIndicesRequest{
		Format: "json",
		H:      []string{"index", "status", "health", "pri", "rep", "docs.count", "store.size"},
	}

	res, err := req.Do(ctx, s.client.client)
	if err != nil {
		emit.Error.StructuredFields("Failed to list indices",
			emit.ZString("error", err.Error()))
		return nil, fmt.Errorf("failed to list indices: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			emit.Warn.StructuredFields("Failed to close response body",
				emit.ZString("error", err.Error()))
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		emit.Error.StructuredFields("Failed to list indices",
			emit.ZString("status", res.Status()),
			emit.ZString("response", string(bodyBytes)))
		return nil, fmt.Errorf("failed to list indices: %s - %s", res.Status(), string(bodyBytes))
	}

	var indices []IndexInfo
	if err := json.NewDecoder(res.Body).Decode(&indices); err != nil {
		return nil, fmt.Errorf("failed to decode indices response: %w", err)
	}

	emit.Debug.StructuredFields("Indices listed successfully",
		emit.ZInt("count", len(indices)))

	return indices, nil
}

// Close closes an index (makes it unavailable for read/write but preserves data)
func (s *IndicesService) Close(ctx context.Context, indexName string) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	req := esapi.IndicesCloseRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(ctx, s.client.client)
	if err != nil {
		return fmt.Errorf("failed to close index: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			emit.Warn.StructuredFields("Failed to close response body",
				emit.ZString("error", err.Error()))
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return fmt.Errorf("failed to close index '%s': %s - %s", indexName, res.Status(), string(bodyBytes))
	}

	return nil
}

// Open opens a previously closed index
func (s *IndicesService) Open(ctx context.Context, indexName string) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	req := esapi.IndicesOpenRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(ctx, s.client.client)
	if err != nil {
		return fmt.Errorf("failed to open index: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			emit.Warn.StructuredFields("Failed to close response body",
				emit.ZString("error", err.Error()))
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return fmt.Errorf("failed to open index '%s': %s - %s", indexName, res.Status(), string(bodyBytes))
	}

	return nil
}

// Refresh forces a refresh of specified indices (or all if none specified)
func (s *IndicesService) Refresh(ctx context.Context, indexNames ...string) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	req := esapi.IndicesRefreshRequest{
		Index: indexNames, // Empty slice means all indices
	}

	res, err := req.Do(ctx, s.client.client)
	if err != nil {
		return fmt.Errorf("failed to refresh indices: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			emit.Warn.StructuredFields("Failed to close response body",
				emit.ZString("error", err.Error()))
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return fmt.Errorf("failed to refresh indices: %s - %s", res.Status(), string(bodyBytes))
	}

	return nil
}

// Stats returns statistics for specified indices (or all if none specified)
func (s *IndicesService) Stats(ctx context.Context, indexNames ...string) (map[string]any, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	req := esapi.IndicesStatsRequest{
		Index: indexNames, // Empty slice means all indices
	}

	res, err := req.Do(ctx, s.client.client)
	if err != nil {
		return nil, fmt.Errorf("failed to get indices stats: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			emit.Warn.StructuredFields("Failed to close response body",
				emit.ZString("error", err.Error()))
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("failed to get indices stats: %s - %s", res.Status(), string(bodyBytes))
	}

	var result map[string]any
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode stats response: %w", err)
	}

	return result, nil
}

// Clone creates a copy of an existing index
func (s *IndicesService) Clone(ctx context.Context, sourceIndex, targetIndex string) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	req := esapi.IndicesCloneRequest{
		Index:  sourceIndex,
		Target: targetIndex,
	}

	res, err := req.Do(ctx, s.client.client)
	if err != nil {
		return fmt.Errorf("failed to clone index: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			emit.Warn.StructuredFields("Failed to close response body",
				emit.ZString("error", err.Error()))
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return fmt.Errorf("failed to clone index '%s' to '%s': %s - %s", sourceIndex, targetIndex, res.Status(), string(bodyBytes))
	}

	return nil
}

// Reindex copies documents from a source index to a target index
func (s *IndicesService) Reindex(ctx context.Context, sourceIndex, targetIndex string, options ...map[string]any) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Minute) // Longer timeout for reindex
		defer cancel()
	}

	// Build reindex body
	reindexBody := map[string]any{
		"source": map[string]any{
			"index": sourceIndex,
		},
		"dest": map[string]any{
			"index": targetIndex,
		},
	}

	// Apply options (e.g., query filters, size, etc.)
	if len(options) > 0 {
		for key, value := range options[0] {
			if key == "query" {
				if source, ok := reindexBody["source"].(map[string]any); ok {
					source["query"] = value
				}
			} else {
				reindexBody[key] = value
			}
		}
	}

	bodyBytes, err := json.Marshal(reindexBody)
	if err != nil {
		return fmt.Errorf("failed to marshal reindex body: %w", err)
	}

	req := esapi.ReindexRequest{
		Body: bytes.NewReader(bodyBytes),
	}

	res, err := req.Do(ctx, s.client.client)
	if err != nil {
		return fmt.Errorf("failed to reindex: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			emit.Warn.StructuredFields("Failed to close response body",
				emit.ZString("error", err.Error()))
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return fmt.Errorf("failed to reindex from '%s' to '%s': %s - %s", sourceIndex, targetIndex, res.Status(), string(bodyBytes))
	}

	return nil
}

// Aliases returns all index aliases
func (s *IndicesService) Aliases(ctx context.Context) (map[string]any, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	req := esapi.IndicesGetAliasRequest{}

	res, err := req.Do(ctx, s.client.client)
	if err != nil {
		return nil, fmt.Errorf("failed to get aliases: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			emit.Warn.StructuredFields("Failed to close response body",
				emit.ZString("error", err.Error()))
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("failed to get aliases: %s - %s", res.Status(), string(bodyBytes))
	}

	var result map[string]any
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode aliases response: %w", err)
	}

	return result, nil
}

// Alias creates or updates an alias pointing to one or more indices
func (s *IndicesService) Alias(ctx context.Context, aliasName string, indexNames ...string) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	if len(indexNames) == 0 {
		return fmt.Errorf("at least one index name must be provided")
	}

	// Build alias actions
	actions := make([]map[string]any, 0, len(indexNames))
	for _, indexName := range indexNames {
		actions = append(actions, map[string]any{
			"add": map[string]any{
				"index": indexName,
				"alias": aliasName,
			},
		})
	}

	aliasBody := map[string]any{
		"actions": actions,
	}

	bodyBytes, err := json.Marshal(aliasBody)
	if err != nil {
		return fmt.Errorf("failed to marshal alias body: %w", err)
	}

	req := esapi.IndicesUpdateAliasesRequest{
		Body: bytes.NewReader(bodyBytes),
	}

	res, err := req.Do(ctx, s.client.client)
	if err != nil {
		return fmt.Errorf("failed to update aliases: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			emit.Warn.StructuredFields("Failed to close response body",
				emit.ZString("error", err.Error()))
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return fmt.Errorf("failed to create alias '%s': %s - %s", aliasName, res.Status(), string(bodyBytes))
	}

	return nil
}

// RemoveAlias removes an alias from one or more indices
func (s *IndicesService) RemoveAlias(ctx context.Context, aliasName string, indexNames ...string) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	if len(indexNames) == 0 {
		return fmt.Errorf("at least one index name must be provided")
	}

	// Build alias actions
	actions := make([]map[string]any, 0, len(indexNames))
	for _, indexName := range indexNames {
		actions = append(actions, map[string]any{
			"remove": map[string]any{
				"index": indexName,
				"alias": aliasName,
			},
		})
	}

	aliasBody := map[string]any{
		"actions": actions,
	}

	bodyBytes, err := json.Marshal(aliasBody)
	if err != nil {
		return fmt.Errorf("failed to marshal alias body: %w", err)
	}

	req := esapi.IndicesUpdateAliasesRequest{
		Body: bytes.NewReader(bodyBytes),
	}

	res, err := req.Do(ctx, s.client.client)
	if err != nil {
		return fmt.Errorf("failed to update aliases: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			emit.Warn.StructuredFields("Failed to close response body",
				emit.ZString("error", err.Error()))
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return fmt.Errorf("failed to remove alias '%s': %s - %s", aliasName, res.Status(), string(bodyBytes))
	}

	return nil
}

// Analyze tests how text is analyzed in a specific index
func (s *IndicesService) Analyze(ctx context.Context, indexName, text, analyzer string) (map[string]any, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	analyzeBody := map[string]any{
		"text":     text,
		"analyzer": analyzer,
	}

	bodyBytes, err := json.Marshal(analyzeBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal analyze body: %w", err)
	}

	req := esapi.IndicesAnalyzeRequest{
		Index: indexName,
		Body:  bytes.NewReader(bodyBytes),
	}

	res, err := req.Do(ctx, s.client.client)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze text: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			emit.Warn.StructuredFields("Failed to close response body",
				emit.ZString("error", err.Error()))
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("failed to analyze text in index '%s': %s - %s", indexName, res.Status(), string(bodyBytes))
	}

	var result map[string]any
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode analyze response: %w", err)
	}

	return result, nil
}

// Shrink reduces the number of shards in an index
func (s *IndicesService) Shrink(ctx context.Context, sourceIndex, targetIndex string, targetShards int) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Minute) // Longer timeout for shrink
		defer cancel()
	}

	shrinkBody := map[string]any{
		"settings": map[string]any{
			"index.number_of_shards": targetShards,
		},
	}

	bodyBytes, err := json.Marshal(shrinkBody)
	if err != nil {
		return fmt.Errorf("failed to marshal shrink body: %w", err)
	}

	req := esapi.IndicesShrinkRequest{
		Index:  sourceIndex,
		Target: targetIndex,
		Body:   bytes.NewReader(bodyBytes),
	}

	res, err := req.Do(ctx, s.client.client)
	if err != nil {
		return fmt.Errorf("failed to shrink index: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			emit.Warn.StructuredFields("Failed to close response body",
				emit.ZString("error", err.Error()))
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return fmt.Errorf("failed to shrink index '%s' to '%s': %s - %s", sourceIndex, targetIndex, res.Status(), string(bodyBytes))
	}

	return nil
}

// Flush forces a flush of specified indices (or all if none specified)
func (s *IndicesService) Flush(ctx context.Context, indexNames ...string) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 2*time.Minute) // Longer timeout for flush
		defer cancel()
	}

	req := esapi.IndicesFlushRequest{
		Index: indexNames, // Empty slice means all indices
	}

	res, err := req.Do(ctx, s.client.client)
	if err != nil {
		return fmt.Errorf("failed to flush indices: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			emit.Warn.StructuredFields("Failed to close response body",
				emit.ZString("error", err.Error()))
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return fmt.Errorf("failed to flush indices: %s - %s", res.Status(), string(bodyBytes))
	}

	return nil
}

// Rollover creates a new index when conditions are met and updates alias
func (s *IndicesService) Rollover(ctx context.Context, aliasName string, options ...map[string]any) (map[string]any, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	rolloverBody := map[string]any{}

	// Apply options (conditions like max_age, max_size, etc.)
	if len(options) > 0 {
		for key, value := range options[0] {
			rolloverBody[key] = value
		}
	}

	var body io.Reader
	if len(rolloverBody) > 0 {
		bodyBytes, err := json.Marshal(rolloverBody)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal rollover body: %w", err)
		}
		body = bytes.NewReader(bodyBytes)
	}

	req := esapi.IndicesRolloverRequest{
		Alias: aliasName,
		Body:  body,
	}

	res, err := req.Do(ctx, s.client.client)
	if err != nil {
		return nil, fmt.Errorf("failed to rollover index: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			emit.Warn.StructuredFields("Failed to close response body",
				emit.ZString("error", err.Error()))
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("failed to rollover alias '%s': %s - %s", aliasName, res.Status(), string(bodyBytes))
	}

	var result map[string]any
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode rollover response: %w", err)
	}

	return result, nil
}

// Template Management Methods

// CreateTemplate creates an index template
func (s *IndicesService) CreateTemplate(ctx context.Context, name string, template map[string]any) error {
	clusterResource := &ClusterResource{
		client: s.client,
	}
	return clusterResource.CreateTemplate(ctx, name, template)
}

// GetTemplate retrieves an index template
func (s *IndicesService) GetTemplate(ctx context.Context, name string) (map[string]any, error) {
	clusterResource := &ClusterResource{
		client: s.client,
	}
	return clusterResource.GetTemplate(ctx, name)
}

// DeleteTemplate deletes an index template
func (s *IndicesService) DeleteTemplate(ctx context.Context, name string) error {
	clusterResource := &ClusterResource{
		client: s.client,
	}
	return clusterResource.DeleteTemplate(ctx, name)
}

// ListTemplates lists all index templates
func (s *IndicesService) ListTemplates(ctx context.Context) (map[string]any, error) {
	clusterResource := &ClusterResource{
		client: s.client,
	}
	return clusterResource.ListTemplates(ctx)
}

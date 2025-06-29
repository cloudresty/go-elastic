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

// ClusterResource provides cluster-level operations
type ClusterResource struct {
	client *Client
}

// Health returns the cluster health
func (cr *ClusterResource) Health(ctx context.Context) (*ClusterHealth, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	req := esapi.ClusterHealthRequest{}

	res, err := req.Do(ctx, cr.client.client)
	if err != nil {
		emit.Error.StructuredFields("Failed to get cluster health",
			emit.ZString("error", err.Error()))
		return nil, fmt.Errorf("failed to get cluster health: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			emit.Warn.StructuredFields("Failed to close response body",
				emit.ZString("error", err.Error()))
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		emit.Error.StructuredFields("Failed to get cluster health",
			emit.ZString("status", res.Status()),
			emit.ZString("response", string(bodyBytes)))
		return nil, fmt.Errorf("cluster health request failed: %s - %s", res.Status(), string(bodyBytes))
	}

	var health ClusterHealth
	if err := json.NewDecoder(res.Body).Decode(&health); err != nil {
		return nil, fmt.Errorf("failed to decode cluster health response: %w", err)
	}

	emit.Debug.StructuredFields("Cluster health retrieved successfully",
		emit.ZString("status", health.Status),
		emit.ZInt("active_primary_shards", health.ActivePrimaryShards),
		emit.ZInt("active_shards", health.ActiveShards))

	return &health, nil
}

// Stats returns cluster statistics
func (cr *ClusterResource) Stats(ctx context.Context) (*ClusterStats, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	req := esapi.ClusterStatsRequest{}

	res, err := req.Do(ctx, cr.client.client)
	if err != nil {
		emit.Error.StructuredFields("Failed to get cluster stats",
			emit.ZString("error", err.Error()))
		return nil, fmt.Errorf("failed to get cluster stats: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			emit.Warn.StructuredFields("Failed to close response body",
				emit.ZString("error", err.Error()))
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		emit.Error.StructuredFields("Failed to get cluster stats",
			emit.ZString("status", res.Status()),
			emit.ZString("response", string(bodyBytes)))
		return nil, fmt.Errorf("cluster stats request failed: %s - %s", res.Status(), string(bodyBytes))
	}

	var stats ClusterStats
	if err := json.NewDecoder(res.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("failed to decode cluster stats response: %w", err)
	}

	emit.Debug.StructuredFields("Cluster stats retrieved successfully",
		emit.ZString("cluster_name", stats.ClusterName),
		emit.ZString("status", stats.Status))

	return &stats, nil
}

// CreateTemplate creates an index template
func (cr *ClusterResource) CreateTemplate(ctx context.Context, name string, template map[string]any) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	bodyBytes, err := json.Marshal(template)
	if err != nil {
		return fmt.Errorf("failed to marshal template: %w", err)
	}

	req := esapi.IndicesPutIndexTemplateRequest{
		Name: name,
		Body: io.NopCloser(bytes.NewReader(bodyBytes)),
	}

	res, err := req.Do(ctx, cr.client.client)
	if err != nil {
		emit.Error.StructuredFields("Failed to create index template",
			emit.ZString("template", name),
			emit.ZString("error", err.Error()))
		return fmt.Errorf("failed to create index template: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			emit.Warn.StructuredFields("Failed to close response body",
				emit.ZString("error", err.Error()))
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		emit.Error.StructuredFields("Failed to create index template",
			emit.ZString("template", name),
			emit.ZString("status", res.Status()),
			emit.ZString("response", string(bodyBytes)))
		return fmt.Errorf("failed to create template '%s': %s - %s", name, res.Status(), string(bodyBytes))
	}

	emit.Info.StructuredFields("Index template created successfully",
		emit.ZString("template", name))

	return nil
}

// GetTemplate retrieves an index template
func (cr *ClusterResource) GetTemplate(ctx context.Context, name string) (map[string]any, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	req := esapi.IndicesGetIndexTemplateRequest{
		Name: name,
	}

	res, err := req.Do(ctx, cr.client.client)
	if err != nil {
		return nil, fmt.Errorf("failed to get index template: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			emit.Warn.StructuredFields("Failed to close response body",
				emit.ZString("error", err.Error()))
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("failed to get template '%s': %s - %s", name, res.Status(), string(bodyBytes))
	}

	var result map[string]any
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode template response: %w", err)
	}

	return result, nil
}

// DeleteTemplate deletes an index template
func (cr *ClusterResource) DeleteTemplate(ctx context.Context, name string) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	req := esapi.IndicesDeleteIndexTemplateRequest{
		Name: name,
	}

	res, err := req.Do(ctx, cr.client.client)
	if err != nil {
		emit.Error.StructuredFields("Failed to delete index template",
			emit.ZString("template", name),
			emit.ZString("error", err.Error()))
		return fmt.Errorf("failed to delete index template: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			emit.Warn.StructuredFields("Failed to close response body",
				emit.ZString("error", err.Error()))
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		emit.Error.StructuredFields("Failed to delete index template",
			emit.ZString("template", name),
			emit.ZString("status", res.Status()),
			emit.ZString("response", string(bodyBytes)))
		return fmt.Errorf("failed to delete template '%s': %s - %s", name, res.Status(), string(bodyBytes))
	}

	emit.Info.StructuredFields("Index template deleted successfully",
		emit.ZString("template", name))

	return nil
}

// ListTemplates lists all index templates
func (cr *ClusterResource) ListTemplates(ctx context.Context) (map[string]any, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	req := esapi.IndicesGetIndexTemplateRequest{}

	res, err := req.Do(ctx, cr.client.client)
	if err != nil {
		return nil, fmt.Errorf("failed to list index templates: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			emit.Warn.StructuredFields("Failed to close response body",
				emit.ZString("error", err.Error()))
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("failed to list templates: %s - %s", res.Status(), string(bodyBytes))
	}

	var result map[string]any
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode templates response: %w", err)
	}

	return result, nil
}

// Settings returns cluster settings (persistent, transient, and default)
func (cr *ClusterResource) Settings(ctx context.Context) (map[string]any, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	req := esapi.ClusterGetSettingsRequest{
		IncludeDefaults: func() *bool { b := true; return &b }(),
	}

	res, err := req.Do(ctx, cr.client.client)
	if err != nil {
		emit.Error.StructuredFields("Failed to get cluster settings",
			emit.ZString("error", err.Error()))
		return nil, fmt.Errorf("failed to get cluster settings: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			emit.Warn.StructuredFields("Failed to close response body",
				emit.ZString("error", err.Error()))
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		emit.Error.StructuredFields("Failed to get cluster settings",
			emit.ZString("status", res.Status()),
			emit.ZString("response", string(bodyBytes)))
		return nil, fmt.Errorf("cluster settings request failed: %s - %s", res.Status(), string(bodyBytes))
	}

	var result map[string]any
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode cluster settings response: %w", err)
	}

	emit.Debug.StructuredFields("Cluster settings retrieved successfully")

	return result, nil
}

// AllocationExplain explains why a shard is unassigned or can't be moved
func (cr *ClusterResource) AllocationExplain(ctx context.Context, body map[string]any) (map[string]any, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	var req esapi.ClusterAllocationExplainRequest

	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal allocation explain body: %w", err)
		}
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	res, err := req.Do(ctx, cr.client.client)
	if err != nil {
		emit.Error.StructuredFields("Failed to get allocation explanation",
			emit.ZString("error", err.Error()))
		return nil, fmt.Errorf("failed to get allocation explanation: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			emit.Warn.StructuredFields("Failed to close response body",
				emit.ZString("error", err.Error()))
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		emit.Error.StructuredFields("Failed to get allocation explanation",
			emit.ZString("status", res.Status()),
			emit.ZString("response", string(bodyBytes)))
		return nil, fmt.Errorf("allocation explain request failed: %s - %s", res.Status(), string(bodyBytes))
	}

	var result map[string]any
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode allocation explain response: %w", err)
	}

	emit.Debug.StructuredFields("Allocation explanation retrieved successfully")

	return result, nil
}

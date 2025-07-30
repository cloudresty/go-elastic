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
		cr.client.config.Logger.Error("Failed to get cluster health - error: %s", err.Error())
		return nil, fmt.Errorf("failed to get cluster health: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			cr.client.config.Logger.Warn("Failed to close response body - error: %s", err.Error())
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		cr.client.config.Logger.Error("Failed to get cluster health - status: %s, response: %s", res.Status(), string(bodyBytes))
		return nil, fmt.Errorf("cluster health request failed: %s - %s", res.Status(), string(bodyBytes))
	}

	var health ClusterHealth
	if err := json.NewDecoder(res.Body).Decode(&health); err != nil {
		return nil, fmt.Errorf("failed to decode cluster health response: %w", err)
	}

	cr.client.config.Logger.Debug("Cluster health retrieved successfully - status: %s, active_primary_shards: %d, active_shards: %d", health.Status, health.ActivePrimaryShards, health.ActiveShards)

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
		cr.client.config.Logger.Error("Failed to get cluster stats - error: %s", err.Error())
		return nil, fmt.Errorf("failed to get cluster stats: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			cr.client.config.Logger.Warn("Failed to close response body - error: %s", err.Error())
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		cr.client.config.Logger.Error("Failed to get cluster stats - status: %s, response: %s", res.Status(), string(bodyBytes))
		return nil, fmt.Errorf("cluster stats request failed: %s - %s", res.Status(), string(bodyBytes))
	}

	var stats ClusterStats
	if err := json.NewDecoder(res.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("failed to decode cluster stats response: %w", err)
	}

	cr.client.config.Logger.Debug("Cluster stats retrieved successfully - cluster_name: %s, status: %s", stats.ClusterName, stats.Status)

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
		cr.client.config.Logger.Error("Failed to create index template", map[string]interface{}{
			"template": name,
			"error":    err.Error(),
		})
		return fmt.Errorf("failed to create index template: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			cr.client.config.Logger.Warn("Failed to close response body", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		cr.client.config.Logger.Error("Failed to create index template", map[string]interface{}{
			"template": name,
			"status":   res.Status(),
			"response": string(bodyBytes),
		})
		return fmt.Errorf("failed to create template '%s': %s - %s", name, res.Status(), string(bodyBytes))
	}

	cr.client.config.Logger.Info("Index template created successfully", map[string]interface{}{
		"template": name,
	})

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
			cr.client.config.Logger.Warn("Failed to close response body", map[string]interface{}{
				"error": err.Error(),
			})
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
		cr.client.config.Logger.Error("Failed to delete index template", map[string]interface{}{
			"template": name,
			"error":    err.Error(),
		})
		return fmt.Errorf("failed to delete index template: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			cr.client.config.Logger.Warn("Failed to close response body", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		cr.client.config.Logger.Error("Failed to delete index template", map[string]interface{}{
			"template": name,
			"status":   res.Status(),
			"response": string(bodyBytes),
		})
		return fmt.Errorf("failed to delete template '%s': %s - %s", name, res.Status(), string(bodyBytes))
	}

	cr.client.config.Logger.Info("Index template deleted successfully", map[string]interface{}{
		"template": name,
	})

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
			cr.client.config.Logger.Warn("Failed to close response body", map[string]interface{}{
				"error": err.Error(),
			})
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
		cr.client.config.Logger.Error("Failed to get cluster settings", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to get cluster settings: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			cr.client.config.Logger.Warn("Failed to close response body", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		cr.client.config.Logger.Error("Failed to get cluster settings", map[string]interface{}{
			"status":   res.Status(),
			"response": string(bodyBytes),
		})
		return nil, fmt.Errorf("cluster settings request failed: %s - %s", res.Status(), string(bodyBytes))
	}

	var result map[string]any
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode cluster settings response: %w", err)
	}

	cr.client.config.Logger.Debug("Cluster settings retrieved successfully", nil)

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
		cr.client.config.Logger.Error("Failed to get allocation explanation", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to get allocation explanation: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			cr.client.config.Logger.Warn("Failed to close response body", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		cr.client.config.Logger.Error("Failed to get allocation explanation", map[string]interface{}{
			"status":   res.Status(),
			"response": string(bodyBytes),
		})
		return nil, fmt.Errorf("allocation explain request failed: %s - %s", res.Status(), string(bodyBytes))
	}

	var result map[string]any
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode allocation explain response: %w", err)
	}

	cr.client.config.Logger.Debug("Allocation explanation retrieved successfully", nil)

	return result, nil
}

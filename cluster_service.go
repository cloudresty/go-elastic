package elastic

import (
	"context"
)

// ClusterService methods

// Health returns cluster health information
func (s *ClusterService) Health(ctx context.Context) (*ClusterHealth, error) {
	clusterResource := &ClusterResource{
		client: s.client,
	}
	return clusterResource.Health(ctx)
}

// Stats returns cluster statistics
func (s *ClusterService) Stats(ctx context.Context) (*ClusterStats, error) {
	clusterResource := &ClusterResource{
		client: s.client,
	}
	return clusterResource.Stats(ctx)
}

// Settings returns cluster settings (persistent, transient, and default)
func (s *ClusterService) Settings(ctx context.Context) (map[string]any, error) {
	clusterResource := &ClusterResource{
		client: s.client,
	}
	return clusterResource.Settings(ctx)
}

// AllocationExplain explains why a shard is unassigned or can't be moved
func (s *ClusterService) AllocationExplain(ctx context.Context, options ...map[string]any) (map[string]any, error) {
	clusterResource := &ClusterResource{
		client: s.client,
	}

	var body map[string]any
	if len(options) > 0 {
		body = options[0]
	}

	return clusterResource.AllocationExplain(ctx, body)
}

// Package elastic provides a modern, production-ready Go package for Elasticsearch operations
// with environment-first configuration, ULID document IDs, auto-reconnection, and comprehensive production features.
//
// This package provides a clean, intuitive API for Elasticsearch operations while maintaining high performance
// and production-ready features.
//
// Key Features:
//   - Environment-first configuration using cloudresty/go-env
//   - Elasticsearch's native random ID generation for optimal shard distribution (default)
//   - Optional ULID support for time-ordered IDs (with performance warnings)
//   - Auto-reconnection with intelligent retry and exponential backoff
//   - Zero-allocation logging with cloudresty/emit
//   - Production-ready features (graceful shutdown, health checks, metrics)
//   - Simple, intuitive function names following Go best practices
//   - Comprehensive error handling and logging
//   - Built-in connection pooling and compression
//   - Index management utilities
//   - Search and aggregation helpers
//   - Bulk operations for high-throughput scenarios
//
// Environment Variables:
//   - ELASTICSEARCH_HOSTS: Elasticsearch server hosts (default: localhost, supports single host or comma-separated multiple hosts)
//   - ELASTICSEARCH_PORT: Elasticsearch server port (default: 9200)
//   - ELASTICSEARCH_USERNAME: Authentication username
//   - ELASTICSEARCH_PASSWORD: Authentication password
//   - ELASTICSEARCH_API_KEY: API key for authentication
//   - ELASTICSEARCH_SERVICE_TOKEN: Service token for authentication
//   - ELASTICSEARCH_CLOUD_ID: Elastic Cloud ID
//   - ELASTICSEARCH_INDEX_PREFIX: Prefix for all index names
//   - ELASTICSEARCH_ID_MODE: ID generation mode (elastic=default, ulid=time-ordered, custom=user-provided)
//   - ELASTICSEARCH_TLS_ENABLED: Enable TLS (default: false)
//   - ELASTICSEARCH_TLS_INSECURE: Allow insecure TLS (default: false)
//   - ELASTICSEARCH_COMPRESSION_ENABLED: Enable compression (default: true)
//   - ELASTICSEARCH_RETRY_ON_STATUS: Retry on these HTTP status codes
//   - ELASTICSEARCH_MAX_RETRIES: Maximum number of retries (default: 3)
//   - ELASTICSEARCH_CONNECTION_NAME: Connection identifier for logging
//   - ELASTICSEARCH_APP_NAME: Application name for connection metadata
//   - ELASTICSEARCH_LOG_LEVEL: Logging level (default: info)
//
// Basic Usage:
//
//	client, err := elastic.NewClient()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
//
//	// Index a document with auto-generated ULID ID
//	doc := map[string]any{
//	    "title": "My Document",
//	    "content": "This is the content",
//	}
//	result, err := client.Index("my-index", doc)
//
//	// Search documents
//	query := map[string]any{
//	    "query": map[string]any{
//	        "match": map[string]any{
//	            "title": "My Document",
//	        },
//	    },
//	}
//	results, err := client.Search("my-index", query)
//
// For more examples and detailed documentation, see the docs/ directory.
package elastic

// Utility functions

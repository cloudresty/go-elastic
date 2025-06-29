# Production Features

[Home](../README.md) &nbsp;/&nbsp; [Docs](README.md) &nbsp;/&nbsp; Production Features

&nbsp;

This document covers all production-ready features designed for
high-availability, fault-tolerant Elasticsearch deployments.

&nbsp;

## Auto-Reconnection

Intelligent reconnection with exponential backoff for network resilience.

### Basic Auto-Reconnection

```go
package main

import (
    "log"
    "time"

    "github.com/cloudresty/go-elastic"
)

func main() {
    config := &elastic.Config{
        Hosts:                []string{"elasticsearch.example.com:9200"},
        ReconnectEnabled:     true,
        ReconnectDelay:       5 * time.Second,
        MaxReconnectDelay:    1 * time.Minute,
        ReconnectBackoff:     2.0,
        MaxReconnectAttempts: 10,
    }

    client, err := elastic.NewClient(elastic.WithConfig(config))
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
}
```

### Environment Configuration

```bash
export ELASTICSEARCH_HOSTS=elasticsearch.example.com:9200
export ELASTICSEARCH_RECONNECT_ENABLED=true
export ELASTICSEARCH_RECONNECT_DELAY=5s
export ELASTICSEARCH_MAX_RECONNECT_DELAY=1m
export ELASTICSEARCH_RECONNECT_BACKOFF=2.0
export ELASTICSEARCH_MAX_RECONNECT_ATTEMPTS=10
```

```go
package main

import (
    "log"

    "github.com/cloudresty/go-elastic"
)

func main() {
    // Environment-first approach (recommended)
    client, err := elastic.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
}
```

### Reconnection Features

• **Exponential Backoff**: Intelligent delay progression to avoid overwhelming the server
• **Maximum Delay Cap**: Prevents excessively long wait times
• **Attempt Limiting**: Configurable maximum reconnection attempts
• **Connection State Tracking**: Monitor reconnection status and count
• **Automatic Recovery**: Seamless operation resumption after reconnection

&nbsp;

## Health Checks

Comprehensive health monitoring for proactive issue detection.

### Basic Health Checks

```go
package main

import (
    "context"
    "log"

    "github.com/cloudresty/go-elastic"
)

func main() {
    // Environment-first approach (recommended)
    client, err := elastic.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    ctx := context.Background()

    // Simple connection test
    if err := client.Ping(ctx); err != nil {
        log.Printf("Elasticsearch unreachable: %s", err)
    } else {
        log.Println("Elasticsearch connection is active")
    }

    // Detailed cluster health check
    health, err := client.Cluster().Health(ctx)
    if err != nil {
        log.Printf("Failed to get cluster health: %s", err)
    } else if health.Status != "green" {
        log.Printf("Cluster not healthy: %s", health.Status)
    }
}
```

### Automated Health Monitoring

```bash
# Environment-first approach
export ELASTICSEARCH_HEALTH_CHECK_ENABLED=true
export ELASTICSEARCH_HEALTH_CHECK_INTERVAL=30s
export ELASTICSEARCH_HOSTS=localhost:9200
```

```go
package main

import (
    "log"
    "time"

    "github.com/cloudresty/go-elastic"
)

func main() {
    // Environment configuration is loaded automatically
    client, err := elastic.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
    // Health checks run automatically in the background

    // Or programmatic configuration
    config := &elastic.Config{
        HealthCheckEnabled:  true,
        HealthCheckInterval: 30 * time.Second,
        Hosts:              []string{"localhost:9200"},
    }

    client2, err := elastic.NewClient(elastic.WithConfig(config))
    if err != nil {
        log.Fatal(err)
    }
    defer client2.Close()
}
```

### Health Check Environment Variables

```bash
export ELASTICSEARCH_HEALTH_CHECK_ENABLED=true
export ELASTICSEARCH_HEALTH_CHECK_INTERVAL=30s
```

### Health Check Features

• **Automated Monitoring**: Background health checks at configurable intervals
• **Connection Validation**: Ping operations to verify connectivity
• **Error Detection**: Early detection of connection issues
• **Status Reporting**: Detailed health status with error information
• **Reconnection Triggering**: Automatic reconnection on health failures

&nbsp;

## Timeout Configuration

Comprehensive timeout controls for production reliability.

### Environment Variables

```bash
# Timeout environment variables
export ELASTICSEARCH_CONNECT_TIMEOUT=30s
export ELASTICSEARCH_REQUEST_TIMEOUT=60s
export ELASTICSEARCH_HOSTS=localhost:9200
```

```go
package main

import (
    "log"

    "github.com/cloudresty/go-elastic"
)

func main() {
    // Environment-first approach (recommended)
    client, err := elastic.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
}
```

### Programmatic Configuration

```go
package main

import (
    "log"
    "time"

    "github.com/cloudresty/go-elastic"
)

func main() {
    config := &elastic.Config{
        ConnectTimeout: 30 * time.Second,
        RequestTimeout: 60 * time.Second,
        Hosts:         []string{"localhost:9200"},
    }

    client, err := elastic.NewClient(elastic.WithConfig(config))
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
}
```

### Timeout Types

| Type | Description | Default | Production Recommendation |
|------|-------------|---------|---------------------------|
| Connect | Initial connection establishment | 10s | 30s for production |
| Request | Individual operation timeout | 30s | Based on operation complexity |

### Production Timeout Recommendations

```bash
# High-traffic production environment
export ELASTICSEARCH_HOSTS=prod-es-1.example.com:9200,prod-es-2.example.com:9200
export ELASTICSEARCH_CONNECT_TIMEOUT=45s
export ELASTICSEARCH_REQUEST_TIMEOUT=120s

# Low-latency environment
export ELASTICSEARCH_HOSTS=fast-es.example.com:9200
export ELASTICSEARCH_CONNECT_TIMEOUT=15s
export ELASTICSEARCH_REQUEST_TIMEOUT=30s
```

### Timeout Best Practices

• Set realistic timeouts based on your network conditions
• Consider operation complexity when setting request timeouts
• Monitor timeout errors to identify infrastructure issues
• Use different timeouts per environment (dev vs. staging vs. production)
• Account for cluster failover time in connection timeouts

&nbsp;

## Graceful Shutdown

Production-ready graceful shutdown with coordinated resource cleanup.

### Basic Graceful Shutdown

```go
package main

import (
    "log"
    "time"

    "github.com/cloudresty/go-elastic"
)

func main() {
    // Environment-first approach
    client, err := elastic.NewClient()
    if err != nil {
        log.Fatal(err)
    }

    // Set up signal handling
    shutdownManager := elastic.NewShutdownManager(&elastic.ShutdownConfig{
        Timeout: 30 * time.Second,
    })

    shutdownManager.SetupSignalHandler()
    shutdownManager.Register(client)

    // Application logic here...

    // Wait for shutdown signal
    shutdownManager.Wait()
}
```

### Advanced Coordinated Shutdown

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/cloudresty/go-elastic"
)

func main() {
    // Multiple clients and resources - environment configuration
    // ELASTICSEARCH_HOSTS=primary-es:9200
    // SERVICE_B_ELASTICSEARCH_HOSTS=service-b-es:9200
    clientA, err := elastic.NewClient()
    if err != nil {
        log.Fatal(err)
    }

    clientB, err := elastic.NewClient(elastic.FromEnvWithPrefix("SERVICE_B_"))
    if err != nil {
        log.Fatal(err)
    }

    shutdownManager := elastic.NewShutdownManager(&elastic.ShutdownConfig{
        Timeout:          30 * time.Second,
        GracePeriod:      5 * time.Second,
        ForceKillTimeout: 10 * time.Second,
    })

    // Register all resources
    shutdownManager.Register(clientA, clientB)
    shutdownManager.SetupSignalHandler()

    // Start background workers
    shutdownManager.Go(func(ctx context.Context) {
        backgroundWorker(ctx)
    })
    shutdownManager.Go(func(ctx context.Context) {
        healthChecker(ctx)
    })

    // Wait for shutdown
    shutdownManager.Wait()
}

func backgroundWorker(ctx context.Context) {
    // Background work implementation
}

func healthChecker(ctx context.Context) {
    // Health check implementation
}
```

### Shutdown Features

• **Signal Handling**: Automatic SIGINT/SIGTERM signal processing
• **In-Flight Tracking**: Waits for pending operations to complete
• **Timeout Protection**: Prevents indefinite waiting during shutdown
• **Component Coordination**: Unified shutdown across multiple clients
• **Zero Data Loss**: Ensures operation completion before exit

&nbsp;

## Performance Characteristics

Optimized for high-throughput, low-latency operations.

### Connection Pooling

```go
package main

import (
    "log"
    "time"

    "github.com/cloudresty/go-elastic"
)

func main() {
    config := &elastic.Config{
        MaxIdleConns:        200,             // Maximum idle connections
        MaxIdleConnsPerHost: 10,              // Maximum idle per host
        IdleConnTimeout:     5 * time.Minute, // Idle connection timeout
        MaxConnLifetime:     30 * time.Minute, // Maximum connection age
        Hosts:              []string{"localhost:9200"},
    }

    client, err := elastic.NewClient(elastic.WithConfig(config))
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
}
```

### Compression Configuration

```go
package main

import (
    "log"

    "github.com/cloudresty/go-elastic"
)

func main() {
    config := &elastic.Config{
        CompressionEnabled: true,
        Hosts:             []string{"localhost:9200"},
    }

    client, err := elastic.NewClient(elastic.WithConfig(config))
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Environment configuration (recommended)
    // export ELASTICSEARCH_COMPRESSION_ENABLED=true
    // export ELASTICSEARCH_HOSTS=localhost:9200
}
```

### Retry Configuration

```go
package main

import (
    "log"

    "github.com/cloudresty/go-elastic"
)

func main() {
    config := &elastic.Config{
        MaxRetries:    3,
        RetryOnStatus: []int{502, 503, 504, 429}, // Retry on these HTTP status codes
        Hosts:        []string{"localhost:9200"},
    }

    client, err := elastic.NewClient(elastic.WithConfig(config))
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Environment configuration (recommended)
    // export ELASTICSEARCH_MAX_RETRIES=3
    // export ELASTICSEARCH_RETRY_ON_STATUS=502,503,504,429
    // export ELASTICSEARCH_HOSTS=localhost:9200
}
```

### Performance Benchmarks

| Operation | Throughput | Latency | Notes |
|-----------|------------|---------|-------|
| Index | 30K ops/sec | <10ms | With Elasticsearch native IDs |
| Get | 80K ops/sec | <3ms | Simple document retrieval |
| Search | 20K ops/sec | <20ms | Simple queries with indexes |
| Bulk | 100K docs/sec | <50ms | Batched operations |

**Note:** Benchmarks performed on Elasticsearch 8.x, 16 CPU cores, 32GB RAM

&nbsp;

## ID Generation Strategies

Optimized document ID generation for performance and scalability.

### Elasticsearch Native IDs (Default - Recommended)

```go
package main

import (
    "log"

    "github.com/cloudresty/go-elastic"
)

func main() {
    config := &elastic.Config{
        IDMode: elastic.IDModeElastic, // Default
        Hosts:  []string{"localhost:9200"},
    }

    client, err := elastic.NewClient(elastic.WithConfig(config))
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Environment configuration (recommended)
    // export ELASTICSEARCH_ID_MODE=elastic
    // export ELASTICSEARCH_HOSTS=localhost:9200
}
```

**Advantages:**
• **Optimal shard distribution** - Random IDs ensure even load across shards
• **Best write performance** - No hotspotting issues
• **Zero configuration** - Works out of the box

### ULID IDs (Use with Caution)

```go
package main

import (
    "log"

    "github.com/cloudresty/go-elastic"
)

func main() {
    config := &elastic.Config{
        IDMode: elastic.IDModeULID,
        Hosts:  []string{"localhost:9200"},
    }

    client, err := elastic.NewClient(elastic.WithConfig(config))
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Environment configuration (recommended)
    // export ELASTICSEARCH_ID_MODE=ulid
    // export ELASTICSEARCH_HOSTS=localhost:9200
}
```

**⚠️ Warning:** Can cause shard hotspotting in multi-shard indices due to time-ordering. Only use when you need sortable IDs and understand the performance implications.

### Custom IDs

```go
package main

import (
    "log"

    "github.com/cloudresty/go-elastic"
)

func main() {
    config := &elastic.Config{
        IDMode: elastic.IDModeCustom,
        Hosts:  []string{"localhost:9200"},
    }

    client, err := elastic.NewClient(elastic.WithConfig(config))
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Provide your own _id in documents
    doc := map[string]any{
        "_id": "custom-id-123",
        "data": "your data",
    }
    _ = doc // Use doc in your application
}
```

For detailed information, see [ID Generation Documentation](id-generation.md).

&nbsp;

## Production Checklist

Essential items for production deployment:

### Completed Features

• **Implement proper logging** - emit library integrated with structured, high-performance logging
• **Configure appropriate timeouts** - connection and request timeouts
• **Implement graceful shutdown** - shutdown manager, in-flight tracking, signal handling
• **Environment-first configuration** - zero-config setup with ELASTICSEARCH_* environment variables
• **Auto-reconnection** - intelligent retry with configurable backoff
• **Health checks** - automated monitoring and status reporting
• **Multiple ID strategies** - Elasticsearch native (default), ULID, and custom IDs
• **Connection pooling** - Configurable connection pool settings for performance
• **Compression support** - HTTP compression for reduced bandwidth
• **Retry logic** - Configurable retry on specific HTTP status codes

### Additional Recommended Items

• Set up monitoring and metrics
• Test failover scenarios
• Monitor memory usage
• Set up alerting for connection failures
• Configure load balancing for Elasticsearch cluster
• Implement backup and disaster recovery
• Performance testing under load
• Security audit and hardening
• Index lifecycle management
• Shard allocation strategies

&nbsp;

## Multi-Instance Support

Support for multiple Elasticsearch connections in the same application:

```go
package main

import (
    "log"

    "github.com/cloudresty/go-elastic"
)

func main() {
    // Default instance
    defaultClient, err := elastic.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    defer defaultClient.Close()

    // Payments service instance
    paymentsClient, err := elastic.NewClient(elastic.FromEnvWithPrefix("PAYMENTS_"))
    if err != nil {
        log.Fatal(err)
    }
    defer paymentsClient.Close()

    // Orders service instance
    ordersClient, err := elastic.NewClient(elastic.FromEnvWithPrefix("ORDERS_"))
    if err != nil {
        log.Fatal(err)
    }
    defer ordersClient.Close()

    // Each uses separate environment variables:
    // ELASTICSEARCH_HOSTS vs PAYMENTS_ELASTICSEARCH_HOSTS vs ORDERS_ELASTICSEARCH_HOSTS
}
```

This allows each service to connect to different Elasticsearch clusters with completely separate configurations.

&nbsp;

---

&nbsp;

An open source project brought to you by the [Cloudresty](https://cloudresty.com/) team.

[Website](https://cloudresty.com/) &nbsp;|&nbsp; [LinkedIn](https://www.linkedin.com/company/cloudresty) &nbsp;|&nbsp; [BlueSky](https://bsky.app/profile/cloudresty.com) &nbsp;|&nbsp; [GitHub](https://github.com/cloudresty)

&nbsp;

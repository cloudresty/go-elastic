# Getting Started with go-elastic

[Home](../README.md) &nbsp;/&nbsp; [Docs](README.md) &nbsp;/&nbsp; Getting Started

&nbsp;

This guide will help you get up and running with the go-elastic package quickly using the modern, resource-oriented API.

&nbsp;

## Prerequisites

- Go 1.21 or later
- Elasticsearch 8.x running locally or remotely
- Basic understanding of Go and Elasticsearch concepts

🔝 [back to top](#getting-started-with-go-elastic)

&nbsp;

## Quick Installation

```bash
go mod init your-project
go get github.com/cloudresty/go-elastic
```

🔝 [back to top](#getting-started-with-go-elastic)

&nbsp;

## Basic Setup

&nbsp;

### 1. Environment Variables (Recommended)

Create a `.env` file in your project root:

```bash
ELASTICSEARCH_HOSTS=localhost:9200
ELASTICSEARCH_ID_MODE=elastic
ELASTICSEARCH_CONNECTION_NAME=my-app
```

🔝 [back to top](#getting-started-with-go-elastic)

&nbsp;

### 2. Basic Client Creation

```go
package main

import (
    "context"
    "log"

    "github.com/cloudresty/go-elastic"
)

func main() {
    // Create client from environment variables
    client, err := elastic.NewClient(elastic.FromEnv())
    if err != nil {
        log.Fatal("Failed to create client:", err)
    }
    defer client.Close()

    // Test the connection
    ctx := context.Background()
    if err := client.Ping(ctx); err != nil {
        log.Fatal("Failed to ping Elasticsearch:", err)
    }

    log.Println("✓ Connected to Elasticsearch successfully!")
}
```

🔝 [back to top](#getting-started-with-go-elastic)

&nbsp;

## Your First Document

&nbsp;

### Index a Document

```go
// Define your document structure
type User struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Age      int    `json:"age"`
    Active   bool   `json:"active"`
}

// Create and index a document
user := User{
    Name:   "John Doe",
    Email:  "john@example.com",
    Age:    30,
    Active: true,
}

// Use the Documents service for document operations
result, err := client.Documents().Create(ctx, "users", user)
if err != nil {
    log.Fatal("Failed to index document:", err)
}

log.Printf("Document indexed with ID: %s", result.ID)
```

🔝 [back to top](#getting-started-with-go-elastic)

&nbsp;

### Search Documents

```go
import "github.com/cloudresty/go-elastic/query"

// Create a type-safe search using the fluent query builder
searchQuery := query.New().
    Must(query.Match("name", "John")).
    Must(query.Term("active", true)).
    Build()

// Use typed search for better Go experience
typedDocs := elastic.For[User](client.Documents())
result, err := typedDocs.Search(ctx, searchQuery, elastic.WithIndices("users"))
if err != nil {
    log.Fatal("Search failed:", err)
}

log.Printf("Found %d documents", result.Total())

// Process typed search results - no manual unmarshaling needed!
users := result.Documents()
for i, user := range users {
    log.Printf("User %d: %s (%s) - Age: %d",
        i+1, user.Name, user.Email, user.Age)
}
```

🔝 [back to top](#getting-started-with-go-elastic)

&nbsp;

## Configuration Options

&nbsp;

### Using Environment Variables

The package supports extensive configuration through environment variables:

```bash
# Connection (hosts with ports required)
ELASTICSEARCH_HOSTS=localhost:9200,localhost:9201
ELASTICSEARCH_USERNAME=your-username
ELASTICSEARCH_PASSWORD=your-password

# Performance
ELASTICSEARCH_MAX_RETRIES=5
ELASTICSEARCH_COMPRESSION_ENABLED=true
ELASTICSEARCH_CONNECT_TIMEOUT=30s

# ID Generation
ELASTICSEARCH_ID_MODE=elastic  # or 'ulid' or 'custom'

# Health Checks
ELASTICSEARCH_HEALTH_CHECK_ENABLED=true
ELASTICSEARCH_HEALTH_CHECK_INTERVAL=30s
```

🔝 [back to top](#getting-started-with-go-elastic)

&nbsp;

### Using Configuration Struct

```go
import "time"

config := &elastic.Config{
    Hosts:              []string{"elasticsearch.example.com:9200"},
    Username:           "myuser",
    Password:           "mypassword",
    ConnectionName:     "my-production-app",
    IDMode:             elastic.IDModeElastic,
    MaxRetries:         5,
    CompressionEnabled: true,
    ConnectTimeout:     30 * time.Second,
}

client, err := elastic.NewClient(elastic.WithConfig(config))
```

🔝 [back to top](#getting-started-with-go-elastic)

&nbsp;

## Modern Search Experience

&nbsp;

### Fluent Query Builder

The package includes a powerful, type-safe query builder in a dedicated `query` sub-package:

```go
import "github.com/cloudresty/go-elastic/query"

// Build complex queries with a fluent API
searchQuery := query.New().
    Must(query.Match("name", "John")).
    Filter(
        query.Range("age").Gte(18).Lte(65).Build(),
        query.Term("active", true),
    ).
    Build()

// You can also build single queries directly
termQuery := query.Term("status", "published")
matchQuery := query.Match("title", "elasticsearch")
rangeQuery := query.Range("price").Gte(10).Lte(100).Build()
```

🔝 [back to top](#getting-started-with-go-elastic)

&nbsp;

### Typed Search Results

Get strongly-typed results without manual unmarshaling:

```go
// Define your document type
type Product struct {
    Name        string  `json:"name"`
    Category    string  `json:"category"`
    Price       float64 `json:"price"`
    InStock     bool    `json:"in_stock"`
}

// Perform typed search
typedDocs := elastic.For[Product](client.Documents())
result, err := typedDocs.Search(ctx, searchQuery, elastic.WithIndices("products"))

// Get typed results - no casting needed!
products := result.Documents()
for _, product := range products {
    fmt.Printf("Product: %s - $%.2f\n", product.Name, product.Price)
}
```

🔝 [back to top](#getting-started-with-go-elastic)

&nbsp;

## ID Generation Strategies

&nbsp;

### Elasticsearch Native (Recommended)

```go
// Default behavior - optimal shard distribution
client, _ := elastic.NewClient()

doc := map[string]any{"title": "My Document"}
result, _ := client.Documents().Create(context.Background(), "my-index", doc)
// result.ID will be something like: "7J9FxIwBDq1XwjOTVxRk"
```

🔝 [back to top](#getting-started-with-go-elastic)

&nbsp;

### ULID (Time-ordered)

```bash
export ELASTICSEARCH_ID_MODE=ulid
```

```go
client, _ := elastic.NewClient()

doc := map[string]any{"title": "My Document"}
result, _ := client.Documents().Create(context.Background(), "my-index", doc)
// result.ID will be something like: "01ARZ3NDEKTSV4RRFFQ69G5FAV"
```

🔝 [back to top](#getting-started-with-go-elastic)

&nbsp;

### Custom IDs

```bash
export ELASTICSEARCH_ID_MODE=custom
```

```go
client, _ := elastic.NewClient()

doc := map[string]any{"title": "My Document"}
result, _ := client.Documents().CreateWithID(context.Background(), "my-index", "custom-id-123", doc)
// result.ID will be: "custom-id-123"
```

🔝 [back to top](#getting-started-with-go-elastic)

&nbsp;

## Production Features

&nbsp;

### Health Checks

```go
import "time"

// Health checks run automatically when enabled
config := &elastic.Config{
    HealthCheckEnabled:  true,
    HealthCheckInterval: 30 * time.Second,
}

client, _ := elastic.NewClient(elastic.WithConfig(config))

// Manual health check
if err := client.Ping(ctx); err == nil {
    log.Println("✓ Elasticsearch is reachable")
}

// For detailed cluster health information
health, err := client.Cluster().Health(ctx)
if err == nil && health.Status == "green" {
    log.Println("✓ Cluster is healthy")
}
```

🔝 [back to top](#getting-started-with-go-elastic)

&nbsp;

### Auto-Reconnection

```go
import "time"

// Auto-reconnection is enabled by default
config := &elastic.Config{
    ReconnectEnabled:     true,
    ReconnectDelay:       5 * time.Second,
    MaxReconnectDelay:    1 * time.Minute,
    ReconnectBackoff:     2.0,
    MaxReconnectAttempts: 10,
}

client, _ := elastic.NewClient(elastic.WithConfig(config))
// Client will automatically reconnect on connection failures
```

🔝 [back to top](#getting-started-with-go-elastic)

&nbsp;

### Graceful Shutdown

```go
import (
    "os"
    "os/signal"
    "syscall"
)

// Set up graceful shutdown
shutdown := elastic.NewShutdownManager()
defer shutdown.WaitForShutdown()

// Register your client
shutdown.RegisterClient(client)

// Handle OS signals
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

go func() {
    <-sigChan
    log.Println("Shutting down gracefully...")
    shutdown.Shutdown()
}()
```

🔝 [back to top](#getting-started-with-go-elastic)

&nbsp;

## Common Patterns

&nbsp;

### Error Handling

```go
result, err := client.Documents().Create(ctx, "users", user)
if err != nil {
    switch {
    case elastic.IsConnectionError(err):
        log.Println("Connection error - check Elasticsearch availability")
    case elastic.IsTimeoutError(err):
        log.Println("Timeout error - consider increasing timeouts")
    case elastic.IsValidationError(err):
        log.Println("Validation error - check document structure")
    default:
        log.Printf("Unexpected error: %v", err)
    }
    return
}
```

🔝 [back to top](#getting-started-with-go-elastic)

&nbsp;

### Bulk Operations

```go
// Use the modern Bulk API for efficient batch operations
documents := []User{
    {Name: "Alice", Email: "alice@example.com", Age: 25},
    {Name: "Bob", Email: "bob@example.com", Age: 30},
    {Name: "Charlie", Email: "charlie@example.com", Age: 35},
}

// Create a bulk operation
bulkOp := client.Documents().Bulk("users")
for _, doc := range documents {
    bulkOp.Index(doc)
}

// Execute the bulk operation
result, err := bulkOp.Do(ctx)
if err != nil {
    log.Printf("Bulk operation failed: %v", err)
    return
}

log.Printf("Bulk operation completed: %d items processed", len(result.Items))
for _, item := range result.Items {
    if item.Index.Error != nil {
        log.Printf("Error indexing document: %s", item.Index.Error.Reason)
    } else {
        log.Printf("Indexed document with ID: %s", item.Index.ID)
    }
}
```

🔝 [back to top](#getting-started-with-go-elastic)

&nbsp;

### Advanced Search with Query Builder

```go
import "github.com/cloudresty/go-elastic/query"

// Use the fluent query builder for complex searches
searchQuery := query.New().
    Filter(
        query.Range("age").Gte(18).Lte(65).Build(),
        query.Term("active", true),
    ).
    Build()

// Execute typed search with aggregations and search options
typedDocs := elastic.For[User](client.Documents())
result, err := typedDocs.Search(ctx, searchQuery,
    elastic.WithIndices("users"),
    elastic.WithSize(50),
    elastic.WithSort(map[string]any{"age": "asc"}),
    elastic.WithAggregations(map[string]any{
        "age_groups": map[string]any{
            "histogram": map[string]any{
                "field":    "age",
                "interval": 10,
            },
        },
    }),
)

if err != nil {
    log.Fatal("Search failed:", err)
}

// Process typed results
users := result.Documents()
log.Printf("Found %d users", len(users))

// Access aggregations
if ageGroups, ok := result.Aggregations()["age_groups"]; ok {
    log.Printf("Age group aggregation: %+v", ageGroups)
}
```

🔝 [back to top](#getting-started-with-go-elastic)

&nbsp;

### Search Iteration

For large result sets, use the search iterator:

```go
import "github.com/cloudresty/go-elastic/query"

// Create a search iterator for paginating through results
searchQuery := query.MatchAll().Build()
typedDocs := elastic.For[User](client.Documents())

iterator := typedDocs.SearchIterator(ctx, searchQuery, elastic.WithIndices("users"))
defer iterator.Close()

// Iterate through all results
for iterator.Next() {
    result := iterator.Result()
    users := result.Documents()

    log.Printf("Processing batch of %d users", len(users))
    for _, user := range users {
        // Process each user
        fmt.Printf("User: %s\n", user.Name)
    }
}

if err := iterator.Err(); err != nil {
    log.Printf("Iterator error: %v", err)
}
```

🔝 [back to top](#getting-started-with-go-elastic)

&nbsp;

## Resource-Oriented API

&nbsp;

### Document Operations

```go
// Get the Documents service
docs := client.Documents()

// Create a document
createResp, err := docs.Create(ctx, "users", user)

// Get a document
doc, err := docs.Get(ctx, "users", createResp.ID)

// Update a document
updateResp, err := docs.Update(ctx, "users", createResp.ID, map[string]any{
    "age": 31,
})

// Delete a document
deleteResp, err := docs.Delete(ctx, "users", createResp.ID)

// Multi-get multiple documents
ids := []string{createResp.ID, "other-id"}
multiResp, err := docs.MultiGet(ctx, "users", ids)
```

🔝 [back to top](#getting-started-with-go-elastic)

&nbsp;

### Index Management

```go
// Get the Indices service
indices := client.Indices()

// Create an index with mapping
mapping := map[string]any{
    "mappings": map[string]any{
        "properties": map[string]any{
            "name":  map[string]any{"type": "text"},
            "age":   map[string]any{"type": "integer"},
            "email": map[string]any{"type": "keyword"},
        },
    },
}
err := indices.Create(ctx, "users", mapping)

// Check if index exists
exists, err := indices.Exists(ctx, "users")

// Get index information
info, err := indices.Get(ctx, "users")

// Delete an index
err = indices.Delete(ctx, "users")

// List all indices
list, err := indices.List(ctx)
```

🔝 [back to top](#getting-started-with-go-elastic)

&nbsp;

## Next Steps

1. **Read the [API Reference](api-reference.md)** - Complete API documentation
2. **Explore Examples** - Check the `examples/` directory for comprehensive demos
3. **Review [Production Features](production-features.md)** - Production deployment guidance
4. **Configure [Environment Variables](environment-variables.md)** - All supported variables
5. **Learn [Environment Configuration](environment-configuration.md)** - Setup patterns and examples

🔝 [back to top](#getting-started-with-go-elastic)

&nbsp;

## Common Issues

&nbsp;

### Connection Refused

```bash
# Check if Elasticsearch is running
curl http://localhost:9200

# Or use Docker
docker run -d -p 9200:9200 -e "discovery.type=single-node" \
  docker.elastic.co/elasticsearch/elasticsearch:8.11.0
```

🔝 [back to top](#getting-started-with-go-elastic)

&nbsp;

### Authentication Errors

```bash
# Set credentials
export ELASTICSEARCH_USERNAME=elastic
export ELASTICSEARCH_PASSWORD=changeme

# Or use API key
export ELASTICSEARCH_API_KEY=your-api-key
```

🔝 [back to top](#getting-started-with-go-elastic)

&nbsp;

### Index Not Found

```go
// Create index with mapping before indexing documents
mapping := map[string]any{
    "mappings": map[string]any{
        "properties": map[string]any{
            "name": map[string]any{"type": "text"},
            "age":  map[string]any{"type": "integer"},
        },
    },
}

err := client.Indices().Create(ctx, "users", mapping)
```

🔝 [back to top](#getting-started-with-go-elastic)

&nbsp;

## Getting Help

- **Documentation**: Check the [docs/](README.md) directory
- **Examples**: Run examples with `make run-<example-name>`
- **Issues**: Report issues on [GitHub](https://github.com/cloudresty/go-elastic/issues)
- **Community**: Join discussions on [GitHub Discussions](https://github.com/cloudresty/go-elastic/discussions)

&nbsp;

---

&nbsp;

An open source project brought to you by the [Cloudresty](https://cloudresty.com/) team.

[Website](https://cloudresty.com/) &nbsp;|&nbsp; [LinkedIn](https://www.linkedin.com/company/cloudresty) &nbsp;|&nbsp; [BlueSky](https://bsky.app/profile/cloudresty.com) &nbsp;|&nbsp; [GitHub](https://github.com/cloudresty)

&nbsp;

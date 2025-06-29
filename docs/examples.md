# Examples

[Home](../README.md) &nbsp;/&nbsp; [Docs](README.md) &nbsp;/&nbsp; Examples

&nbsp;

This document provides comprehensive examples demonstrating the modern resource-oriented API
of the go-elastic package with fluent query builder and typed results.

&nbsp;

## Available Examples

The package includes several examples in the `examples/` directory:

- `basic-client/` - Basic Elasticsearch client setup using environment variables
- `env-config/` - Environment variable configuration examples
- `search-experience-demo/` - Best-in-class search experience with query builder
- `unified-search-demo/` - Unified search API demonstration
- `resource-api-demo/` - Resource-oriented API examples
- `production-features/` - Production features demonstration with health checks
- `id-demo/` - ID generation strategies demonstration

[🔝 back to top](#examples)

&nbsp;

## Quick Start Examples

&nbsp;

### Basic Client Setup

```go
package main

import (
    "context"
    "os"

    "github.com/cloudresty/emit"
    elastic "github.com/cloudresty/go-elastic"
)

func main() {
    // Create client from environment variables
    client, err := elastic.NewClient(elastic.FromEnv())
    if err != nil {
        emit.Error.StructuredFields("Failed to create client",
            emit.ZString("error", err.Error()))
        os.Exit(1)
    }
    defer client.Close()

    // Check cluster health
    health, err := client.Cluster().Health(context.Background())
    if err != nil {
        emit.Error.StructuredFields("Failed to check health",
            emit.ZString("error", err.Error()))
        os.Exit(1)
    }

    emit.Info.StructuredFields("Elasticsearch client connected successfully",
        emit.ZString("cluster_name", health.ClusterName),
        emit.ZString("status", health.Status))
}
```

[🔝 back to top](#examples)

&nbsp;

### Basic CRUD Operations

```go
package main

import (
    "context"
    "os"

    "github.com/cloudresty/emit"
    elastic "github.com/cloudresty/go-elastic"
)

type User struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Active   bool   `json:"active"`
    Age      int    `json:"age"`
}

func main() {
    client, err := elastic.NewClient(elastic.FromEnv())
    if err != nil {
        emit.Error.StructuredFields("Failed to create client",
            emit.ZString("error", err.Error()))
        os.Exit(1)
    }
    defer client.Close()

    ctx := context.Background()
    docs := client.Documents()

    // Create - ID is automatically generated based on ID mode
    user := User{
        Name:   "John Doe",
        Email:  "john@example.com",
        Active: true,
        Age:    30,
    }

    // Index document (Create)
    result, err := docs.Create(ctx, "users", user)
    if err != nil {
        emit.Error.StructuredFields("Failed to index user",
            emit.ZString("error", err.Error()))
        os.Exit(1)
    }

    emit.Info.StructuredFields("Indexed user successfully",
        emit.ZString("document_id", result.ID),
        emit.ZString("result", result.Result))

    // Read - get the document by ID
    var retrievedUser User
    getResult, err := docs.Get(ctx, "users", result.ID, &retrievedUser)
    if err != nil {
        emit.Error.StructuredFields("Failed to get user",
            emit.ZString("error", err.Error()))
        os.Exit(1)
    }

    emit.Info.StructuredFields("Retrieved user successfully",
        emit.ZString("name", retrievedUser.Name),
        emit.ZBool("found", getResult.Found))

    // Update - modify the user
    updatedUser := User{
        Name:   "John Smith",
        Email:  "john.smith@example.com",
        Active: true,
        Age:    31,
    }

    updateResult, err := docs.Update(ctx, "users", result.ID, updatedUser)
    if err != nil {
        emit.Error.StructuredFields("Failed to update user",
            emit.ZString("error", err.Error()))
        os.Exit(1)
    }

    emit.Info.StructuredFields("Updated user successfully",
        emit.ZString("result", updateResult.Result))

    // Delete - remove the document
    deleteResult, err := docs.Delete(ctx, "users", result.ID)
    if err != nil {
        emit.Error.StructuredFields("Failed to delete user",
            emit.ZString("error", err.Error()))
        os.Exit(1)
    }

    emit.Info.StructuredFields("Deleted user successfully",
        emit.ZString("result", deleteResult.Result))
}
```

[🔝 back to top](#examples)

&nbsp;

## Environment Configuration

&nbsp;

### Running Examples

To run the examples:

```bash
# Set up environment (use ELASTICSEARCH_HOSTS with ports)
export ELASTICSEARCH_HOSTS=localhost:9200
export ELASTICSEARCH_INDEX_PREFIX=myapp
export ELASTICSEARCH_CONNECTION_NAME=example-client

# Run basic client example
go run examples/basic-client/main.go

# Run search experience demo
go run examples/search-experience-demo/main.go

# Run unified search demo
go run examples/unified-search-demo/main.go

# Run production features example
go run examples/production-features/main.go

# Run environment config example
go run examples/env-config/main.go
```

[🔝 back to top](#examples)

&nbsp;

### Custom Connection Names

```go
package main

import (
    "context"
    "os"

    "github.com/cloudresty/emit"
    elastic "github.com/cloudresty/go-elastic"
)

func main() {
    client, err := elastic.NewClient(
        elastic.FromEnv(),
        elastic.WithConnectionName("user-service-v1.2.3"),
    )
    if err != nil {
        emit.Error.StructuredFields("Failed to create client",
            emit.ZString("error", err.Error()))
        os.Exit(1)
    }
    defer client.Close()

    emit.Info.StructuredFields("Connected with custom name",
        emit.ZString("connection_name", client.Name()))
}
```

[🔝 back to top](#examples)

&nbsp;

## Production Configuration

&nbsp;

### With Environment Variables

```bash
# Production Elasticsearch setup
export ELASTICSEARCH_HOSTS=prod-cluster.elasticsearch.net:9200
export ELASTICSEARCH_USERNAME=prod-user
export ELASTICSEARCH_PASSWORD=prod-password
export ELASTICSEARCH_TLS_ENABLED=true
export ELASTICSEARCH_CONNECTION_NAME=production-service
export ELASTICSEARCH_INDEX_PREFIX=production
```

```go
package main

import (
    "context"
    "os"

    "github.com/cloudresty/emit"
    elastic "github.com/cloudresty/go-elastic"
)

func main() {
    client, err := elastic.NewClient(elastic.FromEnv())
    if err != nil {
        emit.Error.StructuredFields("Failed to create production client",
            emit.ZString("error", err.Error()))
        os.Exit(1)
    }
    defer client.Close()

    emit.Info.Msg("Production client configured successfully")
}
```

[🔝 back to top](#examples)

&nbsp;

## Search Operations - Modern Query Builder API

&nbsp;

### Basic Search with Query Builder

```go
package main

import (
    "context"
    "os"

    "github.com/cloudresty/emit"
    elastic "github.com/cloudresty/go-elastic"
    "github.com/cloudresty/go-elastic/query"
)

type Article struct {
    Title       string `json:"title"`
    Content     string `json:"content"`
    Category    string `json:"category"`
    PublishDate string `json:"publish_date"`
    Views       int    `json:"views"`
}

func main() {
    client, err := elastic.NewClient(elastic.FromEnv())
    if err != nil {
        emit.Error.StructuredFields("Failed to create client",
            emit.ZString("error", err.Error()))
        os.Exit(1)
    }
    defer client.Close()

    ctx := context.Background()

    // Simple match query using the fluent query builder
    searchQuery := query.New().
        Must(query.Match("title", "elasticsearch"))

    // Execute search with typed results
    typedDocs := elastic.For[Article](client.Documents())
    results, err := typedDocs.Search(ctx, searchQuery,
        elastic.WithIndices("articles"),
        elastic.WithSize(10),
        elastic.WithSort(map[string]any{"views": "desc"}),
    )
    if err != nil {
        emit.Error.StructuredFields("Search failed",
            emit.ZString("error", err.Error()))
        return
    }

    emit.Info.StructuredFields("Search completed",
        emit.ZInt64("total_hits", results.TotalHits()),
        emit.ZBool("has_hits", results.HasHits()))

    // Process typed results
    if results.HasHits() {
        articles := results.Documents()
        for _, article := range articles {
            emit.Info.StructuredFields("Found article",
                emit.ZString("title", article.Title),
                emit.ZString("category", article.Category))
        }
    }
}
```

[🔝 back to top](#examples)

&nbsp;

### Complex Query with Filters and Aggregations

```go
package main

import (
    "context"
    "os"

    "github.com/cloudresty/emit"
    elastic "github.com/cloudresty/go-elastic"
    "github.com/cloudresty/go-elastic/query"
)

type Product struct {
    Name        string  `json:"name"`
    Category    string  `json:"category"`
    Price       float64 `json:"price"`
    Rating      float64 `json:"rating"`
    InStock     bool    `json:"in_stock"`
    Description string  `json:"description"`
}

func main() {
    client, err := elastic.NewClient(elastic.FromEnv())
    if err != nil {
        emit.Error.StructuredFields("Failed to create client",
            emit.ZString("error", err.Error()))
        os.Exit(1)
    }
    defer client.Close()

    ctx := context.Background()

    // Complex query with bool, filters, and should clauses
    complexQuery := query.New().
        Must(
            query.MultiMatch("laptop computer", "name", "description"),
            query.Range("rating").Gte(4.0).Build(),
        ).
        Filter(
            query.Term("in_stock", true),
            query.Range("price").Gte(500.0).Lte(2000.0).Build(),
        ).
        Should(
            query.Term("category", "electronics"),
            query.Term("category", "computers"),
        ).
        MinimumShouldMatch(1)

    // Execute search with aggregations
    typedDocs := elastic.For[Product](client.Documents())
    results, err := typedDocs.Search(ctx, complexQuery,
        elastic.WithIndices("products"),
        elastic.WithSize(20),
        elastic.WithSort(map[string]any{"rating": "desc", "price": "asc"}),
        elastic.WithAggregation("categories", elastic.NewTermsAggregation("category.keyword", 5)),
        elastic.WithAggregation("avg_price", elastic.NewAvgAggregation("price")),
        elastic.WithAggregation("price_ranges", elastic.NewRangeAggregation("price").
            AddRange("budget", 0, 1000).
            AddRange("mid_range", 1000, 2000).
            AddRange("premium", 2000, 5000)),
    )
    if err != nil {
        emit.Error.StructuredFields("Complex search failed",
            emit.ZString("error", err.Error()))
        return
    }

    emit.Info.StructuredFields("Complex search completed",
        emit.ZInt64("total_hits", results.TotalHits()),
        emit.ZFloat64("max_score", *results.MaxScore()))

    // Process aggregations
    if agg, ok := results.Aggregations()["avg_price"]; ok {
        if avgPrice, exists := agg["value"]; exists {
            emit.Info.StructuredFields("Average price",
                emit.ZFloat64("avg_price", avgPrice.(float64)))
        }
    }

    // Process typed documents
    if results.HasHits() {
        first, hasFirst := results.First()
        if hasFirst {
            emit.Info.StructuredFields("Top product",
                emit.ZString("name", first.Name),
                emit.ZFloat64("price", first.Price),
                emit.ZFloat64("rating", first.Rating))
        }
    }
}
```

[🔝 back to top](#examples)

&nbsp;

### Bulk Operations

```go
package main

import (
    "context"
    "os"

    "github.com/cloudresty/emit"
    elastic "github.com/cloudresty/go-elastic"
)

type Document struct {
    ID      string `json:"id"`
    Title   string `json:"title"`
    Content string `json:"content"`
    Status  string `json:"status"`
}

func main() {
    client, err := elastic.NewClient(elastic.FromEnv())
    if err != nil {
        emit.Error.StructuredFields("Failed to create client",
            emit.ZString("error", err.Error()))
        os.Exit(1)
    }
    defer client.Close()

    ctx := context.Background()
    docs := client.Documents()

    // Prepare documents for bulk operation
    documents := []Document{
        {ID: "doc1", Title: "First Document", Content: "Content of first document", Status: "published"},
        {ID: "doc2", Title: "Second Document", Content: "Content of second document", Status: "draft"},
        {ID: "doc3", Title: "Third Document", Content: "Content of third document", Status: "published"},
    }

    // Create bulk builder
    bulk := docs.Bulk()

    // Add operations to bulk
    for _, doc := range documents {
        bulk.Index("documents", doc.ID, doc)
    }

    // Execute bulk operation
    result, err := bulk.Execute(ctx)
    if err != nil {
        emit.Error.StructuredFields("Bulk operation failed",
            emit.ZString("error", err.Error()))
        os.Exit(1)
    }

    emit.Info.StructuredFields("Bulk operation completed",
        emit.ZInt("took_ms", result.Took),
        emit.ZBool("errors", result.Errors),
        emit.ZInt("items_count", len(result.Items)))

    // Process bulk results
    for i, item := range result.Items {
        for operation, details := range item {
            emit.Info.StructuredFields("Bulk item result",
                emit.ZString("operation", operation),
                emit.ZString("id", details.ID),
                emit.ZString("result", details.Result),
                emit.ZInt("status", details.Status))
        }
    }
}
```

[🔝 back to top](#examples)

&nbsp;

### Multi-Document Operations

```go
package main

import (
    "context"
    "os"

    "github.com/cloudresty/emit"
    elastic "github.com/cloudresty/go-elastic"
)

type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

func main() {
    client, err := elastic.NewClient(elastic.FromEnv())
    if err != nil {
        emit.Error.StructuredFields("Failed to create client",
            emit.ZString("error", err.Error()))
        os.Exit(1)
    }
    defer client.Close()

    ctx := context.Background()
    docs := client.Documents()

    // Multi-get operation
    var users []User
    result, err := docs.MultiGet(ctx, "users", []string{"user1", "user2", "user3"}, &users)
    if err != nil {
        emit.Error.StructuredFields("Multi-get failed",
            emit.ZString("error", err.Error()))
        return
    }

    emit.Info.StructuredFields("Multi-get completed",
        emit.ZInt("found_count", len(users)))

    for i, user := range users {
        if i < len(result.Docs) && result.Docs[i].Found {
            emit.Info.StructuredFields("Retrieved user",
                emit.ZString("name", user.Name),
                emit.ZString("email", user.Email))
        }
    }
}
```

[🔝 back to top](#examples)

&nbsp;

### Index Management

```go
package main

import (
    "context"
    "os"

    "github.com/cloudresty/emit"
    elastic "github.com/cloudresty/go-elastic"
)

func main() {
    client, err := elastic.NewClient(elastic.FromEnv())
    if err != nil {
        emit.Error.StructuredFields("Failed to create client",
            emit.ZString("error", err.Error()))
        os.Exit(1)
    }
    defer client.Close()

    ctx := context.Background()
    indices := client.Indices()

    // Create index with mapping
    mapping := map[string]any{
        "properties": map[string]any{
            "title": map[string]any{
                "type": "text",
                "analyzer": "standard",
            },
            "category": map[string]any{
                "type": "keyword",
            },
            "publish_date": map[string]any{
                "type": "date",
            },
            "views": map[string]any{
                "type": "integer",
            },
        },
    }

    // Create index
    created, err := indices.Create(ctx, "articles", mapping)
    if err != nil {
        emit.Error.StructuredFields("Failed to create index",
            emit.ZString("error", err.Error()))
        return
    }

    emit.Info.StructuredFields("Index creation result",
        emit.ZBool("acknowledged", created.Acknowledged))

    // Check if index exists
    exists, err := indices.Exists(ctx, "articles")
    if err != nil {
        emit.Error.StructuredFields("Failed to check index existence",
            emit.ZString("error", err.Error()))
        return
    }

    emit.Info.StructuredFields("Index exists check",
        emit.ZBool("exists", exists))

    // Get index settings
    settings, err := indices.GetSettings(ctx, "articles")
    if err != nil {
        emit.Error.StructuredFields("Failed to get index settings",
            emit.ZString("error", err.Error()))
        return
    }

    emit.Info.StructuredFields("Retrieved index settings",
        emit.ZInt("settings_count", len(settings)))

    // List indices
    indexList, err := indices.List(ctx)
    if err != nil {
        emit.Error.StructuredFields("Failed to list indices",
            emit.ZString("error", err.Error()))
        return
    }

    emit.Info.StructuredFields("Index list",
        emit.ZInt("total_indices", len(indexList)))

    for _, index := range indexList {
        emit.Info.StructuredFields("Index info",
            emit.ZString("name", index.Name),
            emit.ZString("status", index.Status),
            emit.ZString("health", index.Health))
    }
}
```

[🔝 back to top](#examples)

&nbsp;

## Advanced Examples

&nbsp;

### Search with Pagination

```go
package main

import (
    "context"
    "os"

    "github.com/cloudresty/emit"
    elastic "github.com/cloudresty/go-elastic"
    "github.com/cloudresty/go-elastic/query"
)

type BlogPost struct {
    Title       string `json:"title"`
    Author      string `json:"author"`
    PublishDate string `json:"publish_date"`
    Tags        []string `json:"tags"`
}

func main() {
    client, err := elastic.NewClient(elastic.FromEnv())
    if err != nil {
        emit.Error.StructuredFields("Failed to create client",
            emit.ZString("error", err.Error()))
        os.Exit(1)
    }
    defer client.Close()

    ctx := context.Background()

    // Search with pagination
    searchQuery := query.New().
        Must(query.Match("title", "golang"))

    typedDocs := elastic.For[BlogPost](client.Documents())

    pageSize := 5
    page := 0

    for {
        results, err := typedDocs.Search(ctx, searchQuery,
            elastic.WithIndices("blog_posts"),
            elastic.WithSize(pageSize),
            elastic.WithFrom(page*pageSize),
            elastic.WithSort(map[string]any{"publish_date": "desc"}),
        )
        if err != nil {
            emit.Error.StructuredFields("Search failed",
                emit.ZString("error", err.Error()))
            break
        }

        if !results.HasHits() {
            emit.Info.Msg("No more results")
            break
        }

        emit.Info.StructuredFields("Page results",
            emit.ZInt("page", page+1),
            emit.ZInt("hits_in_page", len(results.Documents())),
            emit.ZInt64("total_hits", results.TotalHits()))

        // Process results
        for _, post := range results.Documents() {
            emit.Info.StructuredFields("Blog post",
                emit.ZString("title", post.Title),
                emit.ZString("author", post.Author))
        }

        page++

        // Stop if we've reached the last page
        if len(results.Documents()) < pageSize {
            break
        }
    }
}
```

[🔝 back to top](#examples)

&nbsp;

### Error Handling and Retry Logic

```go
package main

import (
    "context"
    "os"
    "time"

    "github.com/cloudresty/emit"
    elastic "github.com/cloudresty/go-elastic"
    "github.com/cloudresty/go-elastic/query"
)

type Document struct {
    ID      string `json:"id"`
    Content string `json:"content"`
    Status  string `json:"status"`
}

func main() {
    client, err := elastic.NewClient(elastic.FromEnv())
    if err != nil {
        emit.Error.StructuredFields("Failed to create client",
            emit.ZString("error", err.Error()))
        os.Exit(1)
    }
    defer client.Close()

    ctx := context.Background()
    docs := client.Documents()

    // Create document with retry logic
    doc := Document{
        ID:      "important-doc",
        Content: "This is important content",
        Status:  "published",
    }

    maxRetries := 3
    retryDelay := time.Second

    for attempt := 0; attempt < maxRetries; attempt++ {
        result, err := docs.Create(ctx, "documents", doc)
        if err != nil {
            emit.Error.StructuredFields("Create attempt failed",
                emit.ZInt("attempt", attempt+1),
                emit.ZString("error", err.Error()))

            if attempt < maxRetries-1 {
                emit.Info.StructuredFields("Retrying",
                    emit.ZDuration("delay", retryDelay))
                time.Sleep(retryDelay)
                retryDelay *= 2 // Exponential backoff
                continue
            }

            emit.Error.Msg("Max retries exceeded")
            os.Exit(1)
        }

        emit.Info.StructuredFields("Document created successfully",
            emit.ZString("id", result.ID),
            emit.ZInt("attempt", attempt+1))
        break
    }

    // Search with error handling
    searchQuery := query.New().
        Must(query.Term("status", "published"))

    typedDocs := elastic.For[Document](client.Documents())
    results, err := typedDocs.Search(ctx, searchQuery,
        elastic.WithIndices("documents"),
        elastic.WithSize(10),
    )

    if err != nil {
        emit.Error.StructuredFields("Search failed",
            emit.ZString("error", err.Error()))
        return
    }

    if !results.HasHits() {
        emit.Warn.Msg("No documents found")
        return
    }

    emit.Info.StructuredFields("Search successful",
        emit.ZInt64("total_hits", results.TotalHits()))
}
```

[🔝 back to top](#examples)

&nbsp;

## Summary

These examples demonstrate the modern, resource-oriented API of go-elastic with:

- **Environment-first configuration** using `elastic.FromEnv()`
- **Resource-oriented services** like `client.Documents()` and `client.Indices()`
- **Fluent query builder** with `query.New()` and composable query methods
- **Type-safe search results** using `elastic.For[T]()` for strongly-typed documents
- **Functional options** for search configuration (`WithIndices`, `WithSize`, etc.)
- **Comprehensive error handling** with structured logging
- **Production-ready patterns** including bulk operations, pagination, and retry logic

For more detailed examples, see the `examples/` directory in the repository.

[🔝 back to top](#examples)

&nbsp;

---

&nbsp;

An open source project brought to you by the [Cloudresty](https://cloudresty.com/) team.

[Website](https://cloudresty.com/) &nbsp;|&nbsp; [LinkedIn](https://www.linkedin.com/company/cloudresty) &nbsp;|&nbsp; [BlueSky](https://bsky.app/profile/cloudresty.com) &nbsp;|&nbsp; [GitHub](https://github.com/cloudresty)

&nbsp;

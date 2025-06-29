# Go Elastic

[Home](README.md) &nbsp;/

&nbsp;

A modern, production-ready Go package for Elasticsearch operations with environment-first configuration, ULID IDs, and comprehensive production features.

&nbsp;

[![Go Reference](https://pkg.go.dev/badge/github.com/cloudresty/go-elastic.svg)](https://pkg.go.dev/github.com/cloudresty/go-elastic)
[![Go Tests](https://github.com/cloudresty/go-elastic/actions/workflows/ci.yaml/badge.svg)](https://github.com/cloudresty/go-elastic/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/cloudresty/go-elastic)](https://goreportcard.com/report/github.com/cloudresty/go-elastic)
[![GitHub Tag](https://img.shields.io/github/v/tag/cloudresty/go-elastic?label=Version)](https://github.com/cloudresty/go-elastic/tags)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

&nbsp;

## Table of Contents

- [Key Features](#key-features)
- [Quick Start](#quick-start)
  - [Installation](#installation)
  - [Basic Usage](#basic-usage)
  - [Environment Configuration](#environment-configuration)
- [Documentation](#documentation)
- [Why This Package?](#why-this-package)
- [Production Usage](#production-usage)
- [Requirements](#requirements)
- [Contributing](#contributing)
- [Security](#security)
- [License](#license)

&nbsp;

## Key Features

- **Best-in-Class Search Experience**: Three-pillar approach with fluent query builder, composable search API, and rich typed results
- **Environment-First**: Configure via environment variables for cloud-native deployments
- **Multiple ID Strategies**: Elasticsearch native (default), ULID, and custom ID generation
- **Auto-Reconnection**: Intelligent retry with configurable exponential backoff
- **Production-Ready**: Graceful shutdown, timeouts, health checks, bulk operations
- **Resource-Oriented API**: Clean, idiomatic Go patterns with DocumentsService, IndicesService, ClusterService
- **High Performance**: Optimized connection pooling, compression, and efficient operations
- **Fully Tested**: Comprehensive test coverage with CI/CD pipeline

🔝 [back to top](#go-elastic)

&nbsp;

## Quick Start

&nbsp;

### Installation

```bash
go get github.com/cloudresty/go-elastic
```

🔝 [back to top](#go-elastic)

&nbsp;

### Basic Usage

```go
package main

import (
    "context"
    "github.com/cloudresty/go-elastic"
    "github.com/cloudresty/go-elastic/query"
)

type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

func main() {
    // Client - uses ELASTICSEARCH_* environment variables automatically
    client, err := elastic.NewClient()
    if err != nil {
        panic(err)
    }
    defer client.Close()

    ctx := context.Background()

    // Index a document with auto-generated ID (Elasticsearch native by default)
    user := User{
        Name:  "John Doe",
        Email: "john@example.com",
        Age:   30,
    }
    result, err := client.Documents().Create(ctx, "users", user)
    if err != nil {
        panic(err)
    }

    // BEST-IN-CLASS SEARCH EXPERIENCE

    // 1. Fluent Query Builder - Type-safe, readable queries
    searchQuery := query.New().
        Must(
            query.Match("name", "John"),
            query.Range("age").Gte(18).Build(),
        ).
        Filter(query.Term("active", true))

    // 2. Composable Search API - Rich options, clean syntax
    typedDocs := elastic.For[User](client.Documents())
    results, err := typedDocs.Search(
        ctx,
        searchQuery,
        elastic.WithIndices("users"),
        elastic.WithSize(10),
        elastic.WithSort(map[string]any{"age": "desc"}),
        elastic.WithAggregation("avg_age", elastic.NewAvgAggregation("age")),
    )

    // 3. Rich, Typed Results - Effortless data extraction
    if results.HasHits() {
        users := results.Documents()     // []User - typed slice
        firstUser, _ := results.First()  // User - typed document

        results.Each(func(hit elastic.TypedHit[User]) {
            println(hit.Source.Name, hit.Source.Email)
        })

        adults := results.Filter(func(u User) bool {
            return u.Age >= 18
        })
    }
}
```

🔝 [back to top](#go-elastic)

&nbsp;

### Environment Configuration

Set environment variables for your deployment:

```bash
export ELASTICSEARCH_HOSTS=localhost:9200
export ELASTICSEARCH_CONNECTION_NAME=my-service
export ELASTICSEARCH_TLS_ENABLED=false
export ELASTICSEARCH_CONNECT_TIMEOUT=10s
export ELASTICSEARCH_REQUEST_TIMEOUT=30s
```

🔝 [back to top](#go-elastic)

&nbsp;

## Documentation

| Document | Description |
|----------|-------------|
| [API Reference](docs/api-reference.md) | Complete function reference and usage patterns |
| [Getting Started](docs/getting-started.md) | Step-by-step guide to get up and running quickly |
| [Environment Configuration](docs/environment-configuration.md) | Environment variables and deployment configurations |
| [Environment Variables](docs/environment-variables.md) | Complete reference of all environment variables |
| [Production Features](docs/production-features.md) | Auto-reconnection, graceful shutdown, health checks, bulk operations |
| [ID Generation](docs/id-generation.md) | Document ID generation strategies and performance considerations |
| [Examples](docs/examples.md) | Comprehensive examples and usage patterns |

🔝 [back to top](#go-elastic)

&nbsp;

## Why This Package?

This package is designed for modern cloud-native applications that require robust, high-performance Elasticsearch operations. It leverages the power of Elasticsearch while providing a developer-friendly API that integrates seamlessly with environment-based configurations.

🔝 [back to top](#go-elastic)

&nbsp;

### Environment-First Design

Perfect for modern cloud deployments with Docker, Kubernetes, and CI/CD pipelines. No more hardcoded connection strings.

🔝 [back to top](#go-elastic)

&nbsp;

### Smart ID Generation

Multiple ID generation strategies: Elasticsearch native (default, recommended), ULID for sortable IDs, and custom IDs for specific use cases. Optimized for performance and shard distribution.

🔝 [back to top](#go-elastic)

&nbsp;

### Production-Ready

Built-in support for high availability, graceful shutdown, automatic reconnection, and comprehensive timeout controls.

🔝 [back to top](#go-elastic)

&nbsp;

### Performance Optimized

Optimized connection pooling, HTTP compression, configurable retry logic, and efficient ID generation for high-throughput scenarios.

🔝 [back to top](#go-elastic)

&nbsp;

## Production Usage

```go
// Use custom environment prefix for multi-service deployments
client, err := elastic.NewClient(elastic.FromEnvWithPrefix("SEARCH_"))
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// Health checks and monitoring
if err := client.Ping(context.Background()); err != nil {
    log.Printf("Elasticsearch connection failed: %v", err)
} else {
    log.Println("Elasticsearch connection is healthy")
}

// Graceful shutdown with signal handling
shutdownManager := elastic.NewShutdownManager(&elastic.ShutdownConfig{
    Timeout: 30 * time.Second,
})
shutdownManager.SetupSignalHandler()
shutdownManager.Register(client)
shutdownManager.Wait() // Blocks until SIGINT/SIGTERM
```

🔝 [back to top](#go-elastic)

&nbsp;

## Search Experience Philosophy

go-elastic delivers a **best-in-class search experience** built on three foundational pillars:

🔝 [back to top](#go-elastic)

&nbsp;

### Pillar 1: Fluent Query Builder

Build complex queries with a type-safe, chainable API that reads like natural language:

```go
import "github.com/cloudresty/go-elastic/query"

// Simple queries
userQuery := query.Term("status", "active")
searchQuery := query.Match("title", "golang tutorial")

// Complex bool queries
complexQuery := query.New().
    Must(
        query.MultiMatch("programming guide", "title", "description"),
        query.Range("rating").Gte(4.0).Build(),
    ).
    Filter(
        query.Term("published", true),
        query.Range("price").Lte(50.0).Build(),
    ).
    Should(
        query.Term("category", "programming"),
        query.Term("category", "technology"),
    ).
    MinimumShouldMatch(1)
```

🔝 [back to top](#go-elastic)

&nbsp;

### Pillar 2: Composable Search API

A single, powerful search method with functional options for ultimate flexibility:

```go
// The ONLY way to search: clean, readable, type-safe
typedDocs := elastic.For[Product](client.Documents())
results, err := typedDocs.Search(
    ctx,
    queryBuilder,
    elastic.WithIndices("products"),
    elastic.WithSize(20),
    elastic.WithSort(map[string]any{"rating": "desc"}),
    elastic.WithAggregation("categories", elastic.NewTermsAggregation("category")),
)
```

🔝 [back to top](#go-elastic)

&nbsp;

### Pillar 3: Rich, Typed Results

Smart, structured responses with built-in helpers for effortless data extraction:

```go
// Get typed results using the fluent API
typedDocs := elastic.For[Product](client.Documents())
results, err := typedDocs.Search(ctx, queryBuilder, options...)

// Rich result operations
products := results.Documents()        // []Product - clean slice
total := results.TotalHits()          // int - total count
first, hasFirst := results.First()    // Product, bool - safe access

// Functional operations
expensive := results.Filter(func(p Product) bool {
    return p.Price > 100.0
})

names := results.Map(func(p Product) Product {
    p.Name = strings.ToUpper(p.Name)
    return p
})

// Iterate with metadata
results.Each(func(hit elastic.TypedHit[Product]) {
    fmt.Printf("Product: %s (Score: %.2f)\n",
        hit.Source.Name, *hit.Score)
})
```

**The Result**: Search operations that are not just less verbose, but genuinely enjoyable to write and maintain.

&nbsp;

🔝 [back to top](#go-elastic)

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create a feature branch
3. Add tests for your changes
4. Ensure all tests pass
5. Submit a pull request

🔝 [back to top](#go-elastic)

&nbsp;

## Security

If you discover a security vulnerability, please report it via email to [security@cloudresty.com](mailto:security@cloudresty.com).

🔝 [back to top](#go-elastic)

&nbsp;

## Requirements

- Go 1.24+ (recommended)
- Elasticsearch 8.0+ (recommended)

🔝 [back to top](#go-elastic)

&nbsp;

## License

This project is licensed under the MIT License - see the [LICENSE.txt](LICENSE.txt) file for details.

🔝 [back to top](#go-elastic)

&nbsp;

---

&nbsp;

An open source project brought to you by the [Cloudresty](https://cloudresty.com) team.

[Website](https://cloudresty.com) &nbsp;|&nbsp; [LinkedIn](https://www.linkedin.com/company/cloudresty) &nbsp;|&nbsp; [BlueSky](https://bsky.app/profile/cloudresty.com) &nbsp;|&nbsp; [GitHub](https://github.com/cloudresty)

&nbsp;

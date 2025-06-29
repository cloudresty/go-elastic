# API Reference

[Home](../README.md) &nbsp;/&nbsp; [Docs](README.md) &nbsp;/&nbsp; API Reference

&nbsp;

This document provides the definitive API reference for the `go-elastic` package. It is designed to be your primary resource for mastering the entire library, from client creation and environment configuration to advanced search, bulk, and index management operations.

We've crafted this API to be powerful, consistent, and idiomatically Go. You'll discover how our resource-oriented services, fluent query builder, and type-safe search results can help you write cleaner, more reliable, and more productive Elasticsearch applications.

&nbsp;

## Core Functions

&nbsp;

### Client Creation

| Function | Description |
|----------|-------------|
| `NewClient(options...)` | Creates a client with functional options (use `FromEnv()` to load from environment variables) |

🔝 [back to top](#api-reference)

&nbsp;

#### Client Options

| Option | Description |
|--------|-------------|
| `WithConfig(config *Config)` | Sets a custom configuration for the client |
| `WithHosts(hosts ...string)` | Sets custom hosts (overrides environment) |
| `WithPort(port int)` | Sets a custom port (overrides environment) |
| `WithCredentials(username, password string)` | Sets username and password for authentication (overrides environment) |
| `WithAPIKey(apiKey string)` | Sets API key for authentication (overrides environment) |
| `WithCloudID(cloudID string)` | Sets Elastic Cloud ID (overrides environment) |
| `WithTLS(enabled bool)` | Enables or disables TLS (overrides environment) |
| `WithConnectionName(name string)` | Sets a connection name for logging and identification |

🔝 [back to top](#api-reference)

&nbsp;

### Connection Operations

| Function | Description |
|----------|-------------|
| `client.Name() string` | Get the configured connection name for logging and identification |
| `client.Ping(ctx context.Context) error` | Test connection with context and update internal state |
| `client.Stats() ConnectionStats` | Get connection statistics (reconnect count, last reconnect time, etc.) |
| `client.Close() error` | Close the client and stop background routines |

🔝 [back to top](#api-reference)

&nbsp;

## Environment Configuration

| Function | Description |
|----------|-------------|
| `FromEnv() ClientOption` | Load configuration from `ELASTICSEARCH_*` environment variables (functional option) |
| `FromEnvWithPrefix(prefix string) ClientOption` | Load configuration with custom prefix (e.g., `LOGS_ELASTICSEARCH_*`) (functional option) |

🔝 [back to top](#api-reference)

&nbsp;

## Cluster Operations

| Function | Description |
|----------|-------------|
| `cluster.Health(ctx context.Context) (*ClusterHealth, error)` | Get comprehensive cluster health information |
| `cluster.Stats(ctx context.Context) (*ClusterStats, error)` | Get cluster statistics |
| `cluster.Settings(ctx context.Context) (*ClusterSettings, error)` | Get cluster settings (persistent, transient, and default) |
| `cluster.AllocationExplain(ctx context.Context, options ...AllocationExplainOption) (*AllocationExplain, error)` | Explain why a shard is unassigned or can't be moved |

🔝 [back to top](#api-reference)

&nbsp;

## Document Operations

&nbsp;

### Basic Operations

| Function | Description |
|----------|-------------|
| `documents.Create(ctx context.Context, indexName string, document any) (*IndexResponse, error)` | Create a new document with auto-generated ID |
| `documents.CreateWithID(ctx context.Context, indexName, documentID string, document any) (*IndexResponse, error)` | Create a document with specific ID (fails if exists) |
| `documents.Index(ctx context.Context, indexName, documentID string, document any) (*IndexResponse, error)` | Create or replace a document with specific ID |
| `documents.Get(ctx context.Context, indexName, documentID string) (map[string]any, error)` | Get a document by ID |
| `documents.Update(ctx context.Context, indexName, documentID string, document any) (*IndexResponse, error)` | Partially update a document |
| `documents.Delete(ctx context.Context, indexName, documentID string) (*DeleteResponse, error)` | Delete a document by ID |
| `documents.Exists(ctx context.Context, indexName, documentID string) (bool, error)` | Check if a document exists (more efficient than `Get`) |
| `documents.MultiGet(ctx context.Context, indexName string, documentIDs []string) ([]map[string]any, error)` | Retrieve multiple documents by IDs |
| `documents.UpdateByQuery(ctx context.Context, indexName string, query, script map[string]any) (map[string]any, error)` | Update all documents matching a query |
| `documents.DeleteByQuery(ctx context.Context, indexName string, query map[string]any) (map[string]any, error)` | Delete all documents matching a query |

🔝 [back to top](#api-reference)

&nbsp;

### Search Operations

| Function | Description |
|----------|-------------|
| `For[T any](service *DocumentsService) *TypedDocuments[T]` | Create a typed search interface for fluent method-style calls |
| `typedDocs.Search(ctx context.Context, queryBuilder *query.Builder, options ...SearchOption) (*SearchResult[T], error)` | **THE** search method - typed, builder-required, rich results |
| `typedDocs.Scroll(ctx context.Context, queryBuilder *query.Builder, scrollTime time.Duration, options ...SearchOption) (*TypedSearchIterator[T], error)` | Create a typed search iterator using a query builder |
| `service.Count(ctx context.Context, queryBuilder *query.Builder, options ...SearchOption) (int64, error)` | Count documents using a query builder |

🔝 [back to top](#api-reference)

&nbsp;

### The Unified Search Experience

`go-elastic` provides a single, unambiguous way to search that eliminates choice paralysis and guides users toward the best practices:

&nbsp;

#### One Search Method to Rule Them All

```go
// The ONLY search method you need - fluent, typed, and intuitive
typedDocs := elastic.For[Product](client.Documents())
results, err := typedDocs.Search(
    ctx,
    queryBuilder,  // ← Query builder REQUIRED (no more map[string]any!)
    options...,    // ← Rich functional options
)
// ↑ Returns SearchResult[Product] - rich, typed results
```

🔝 [back to top](#api-reference)

&nbsp;

**Why This Design?**

- **Zero Choice Paralysis**: One obvious, powerful way to search
- **Forces Best Practices**: Builder required (no error-prone maps), types required (no manual assertions)
- **Maximum Productivity**: IDE autocomplete, type safety, rich result methods
- **Fluent and Intuitive**: Method-style API that reads naturally

🔝 [back to top](#api-reference)

&nbsp;

#### Pillar 1: Fluent Query Builder (package `query`)

| Function | Description |
|----------|-------------|
| `query.New()` | Create a new `bool` query builder |
| `query.Term(field, value)` | Create a `term` query builder |
| `query.Terms(field, values...)` | Create a `terms` query builder |
| `query.Match(field, text)` | Create a `match` query builder |
| `query.MatchPhrase(field, text)` | Create a `match_phrase` query builder |
| `query.MultiMatch(text, fields...)` | Create a `multi_match` query builder |
| `query.Range(field)` | Create a `range` query builder with fluent methods |
| `query.Exists(field)` | Create an `exists` query builder |
| `query.MatchAll()` | Create a `match_all` query builder |
| `query.MatchNone()` | Create a `match_none` query builder |

🔝 [back to top](#api-reference)

&nbsp;

**Range Query Builder Methods:**

| Method | Description |
|--------|-------------|
| `rangeBuilder.Gte(value)` | Set greater than or equal to value |
| `rangeBuilder.Gt(value)` | Set greater than value |
| `rangeBuilder.Lte(value)` | Set less than or equal to value |
| `rangeBuilder.Lt(value)` | Set less than value |
| `rangeBuilder.Format(format)` | Set date format for date fields |
| `rangeBuilder.TimeZone(tz)` | Set timezone for date fields |
| `rangeBuilder.Build()` | Convert to query builder |

🔝 [back to top](#api-reference)

&nbsp;

**Bool Query Builder Methods:**

| Method | Description |
|--------|-------------|
| `builder.Must(queries...)` | Add queries to `must` clause |
| `builder.Filter(queries...)` | Add queries to `filter` clause |
| `builder.Should(queries...)` | Add queries to `should` clause |
| `builder.MustNot(queries...)` | Add queries to `must_not` clause |
| `builder.MinimumShouldMatch(count)` | Set minimum should match count |
| `builder.Build()` | Get the query as `map[string]any` |

🔝 [back to top](#api-reference)

&nbsp;

#### Pillar 2: Composable Search Call

The single `Search[T]()` function supports all search scenarios through functional options.

🔝 [back to top](#api-reference)

&nbsp;

#### Pillar 3: Rich, Typed Result Set

| Type | Description |
|------|-------------|
| `SearchResult[T]` | Generic search result with typed documents |
| `TypedHit[T]` | Individual search hit with typed source |
| `TypedSearchIterator[T]` | Typed scroll iterator for large result sets |
| `TypedDocuments[T]` | Typed search interface for method-style calls |

🔝 [back to top](#api-reference)

&nbsp;

**SearchResult[T] Methods:**

| Method | Description |
|--------|-------------|
| `result.Documents()` | Get slice of typed documents |
| `result.DocumentIDs()` | Get slice of document IDs |
| `result.DocumentsWithIDs()` | Get slice of `DocumentWithID[T]` |
| `result.TotalHits()` | Get total number of hits |
| `result.HasHits()` | Check if there are any hits |
| `result.MaxScore()` | Get maximum relevance score |
| `result.First()` | Get first document (if available) |
| `result.Last()` | Get last document (if available) |
| `result.Each(fn)` | Iterate over all hits |
| `result.Map(fn)` | Transform all documents |
| `result.Filter(fn)` | Filter documents by predicate |

🔝 [back to top](#api-reference)

&nbsp;

**TypedDocuments[T] Methods:**

| Method | Description |
|--------|-------------|
| `typedDocs.Search(ctx, queryBuilder, options...)` | Typed search with method-style API |
| `typedDocs.Scroll(ctx, queryBuilder, scrollTime, options...)` | Typed scroll with method-style API |

🔝 [back to top](#api-reference)

&nbsp;

#### Search Options

| Option | Description |
|--------|-------------|
| `WithIndices(indices ...string) SearchOption` | Search specific indices (supports single or multiple indices) |
| `WithSize(size int) SearchOption` | Set the number of hits to return |
| `WithFrom(from int) SearchOption` | Set the starting offset for pagination |
| `WithSort(sorts ...map[string]any) SearchOption` | Add sorting to the search (can be called multiple times) |
| `WithAggregations(aggs map[string]any) SearchOption` | Add aggregations to the search |
| `WithSource(includes ...string) SearchOption` | Include specific fields in results (can be called multiple times) |
| `WithTimeout(timeout time.Duration) SearchOption` | Set search timeout |

🔝 [back to top](#api-reference)

&nbsp;

#### Search Iterator Methods

| Method | Description |
|--------|-------------|
| `iterator.Next(ctx context.Context) bool` | Advance to next document (returns true if available) |
| `iterator.Scan(dest any) error` | Unmarshal current document into destination |
| `iterator.Current() map[string]any` | Get current document as `map[string]any` |
| `iterator.CurrentHit() *TypedHit[T]` | Get current Hit with metadata |
| `iterator.Err() error` | Get any error that occurred during iteration |
| `iterator.TotalHits() int64` | Get total number of hits found |
| `iterator.ProcessedHits() int64` | Get number of hits processed so far |
| `iterator.Close(ctx context.Context) error` | Clean up scroll context (called automatically) |

🔝 [back to top](#api-reference)

&nbsp;

### Bulk Operations

| Function | Description |
|----------|-------------|
| `documents.Bulk(indexName string) *BulkIndexer` | Create a `BulkIndexer` for chaining bulk operations |

🔝 [back to top](#api-reference)

&nbsp;

#### BulkIndexer Methods

| Method | Description |
|--------|-------------|
| `bulkIndexer.Create(document any) *BulkIndexer` | Add a create operation with auto-generated ID |
| `bulkIndexer.CreateWithID(id string, document any) *BulkIndexer` | Add a create operation with specific ID |
| `bulkIndexer.Index(id string, document any) *BulkIndexer` | Add an index operation (create or replace) |
| `bulkIndexer.Update(id string, document any) *BulkIndexer` | Add an update operation |
| `bulkIndexer.UpdateWithScript(id string, script map[string]any) *BulkIndexer` | Add an update operation with script |
| `bulkIndexer.Delete(id string) *BulkIndexer` | Add a delete operation |
| `bulkIndexer.Do(ctx context.Context) (*BulkResponse, error)` | Execute all accumulated operations |

🔝 [back to top](#api-reference)

&nbsp;

## Index Management

All methods are part of the `IndicesService` and are accessed via `client.Indices()`.

&nbsp;

### Core Index Operations

| Method | Description |
|--------|-------------|
| `indices.Create(ctx context.Context, indexName string, mapping map[string]any) error` | Create an index with optional mapping |
| `indices.Delete(ctx context.Context, indexName string) error` | Delete one or more indices |
| `indices.Exists(ctx context.Context, indexName string) (bool, error)` | Check if an index exists |
| `indices.Get(indexName string) *IndexResource` | Get detailed information about one or more indices |
| `indices.List(ctx context.Context) ([]IndexInfo, error)` | Get detailed information about all indices |
| `indices.Close(ctx context.Context, indexName string) error` | Close one or more indices |
| `indices.Open(ctx context.Context, indexName string) error` | Open previously closed indices |

🔝 [back to top](#api-reference)

&nbsp;

### Lifecycle and Maintenance

| Method | Description |
|--------|-------------|
| `indices.Refresh(ctx, indexNames...)` | Force refresh of indices (or all if none specified) |
| `indices.Flush(ctx, indexNames...)` | Force flush to disk (or all if none specified) |
| `indices.Stats(ctx, indexNames...)` | Get statistics for indices (or all if none specified) |
| `indices.Clone(ctx, sourceIndex, targetIndex)` | Create a copy of an existing index |
| `indices.Reindex(ctx, sourceIndex, targetIndex, options...)` | Copy documents between indices with optional filtering |
| `indices.Rollover(ctx, aliasName, options...)` | Create a new index for a data stream or alias |
| `indices.Shrink(ctx, sourceIndex, targetIndex, shards)` | Reduce the number of primary shards |

🔝 [back to top](#api-reference)

&nbsp;

### Alias Management

| Method | Description |
|--------|-------------|
| `indices.GetAliases(ctx)` | Get all aliases in the cluster |
| `indices.AddAlias(ctx, indexNames, aliasName)` | Add an alias to one or more indices |
| `indices.RemoveAlias(ctx, indexNames, aliasName)` | Remove an alias from one or more indices |

🔝 [back to top](#api-reference)

&nbsp;

### Mapping and Settings

| Method | Description |
|--------|-------------|
| `indices.GetMapping(ctx, indexName)` | Get the mapping for an index |
| `indices.UpdateMapping(ctx, indexName, mapping)` | Update the mapping for an index |
| `indices.GetSettings(ctx, indexName)` | Get the settings for an index |
| `indices.UpdateSettings(ctx, indexName, settings)` | Update the settings for an index |
| `indices.Analyze(ctx, indexName, text, analyzer)` | Test how text is analyzed with a specific analyzer |

🔝 [back to top](#api-reference)

&nbsp;

### Template Management

| Method | Description |
|--------|-------------|
| `indices.CreateTemplate(ctx, name, template)` | Create an index template |
| `indices.GetTemplate(ctx, name)` | Retrieve an index template |
| `indices.DeleteTemplate(ctx, name)` | Delete an index template |
| `indices.ListTemplates(ctx)` | List all index templates |

> **Note**: For index-scoped document operations (search, CRUD), use the existing `DocumentsService` with `WithIndices(indexName)` option. For example: `client.Documents().Search(ctx, query, elastic.WithIndices("my-index"))`.

🔝 [back to top](#api-reference)

&nbsp;

---

&nbsp;

An open source project brought to you by the [Cloudresty](https://cloudresty.com/) team.

[Website](https://cloudresty.com/) &nbsp;|&nbsp; [LinkedIn](https://www.linkedin.com/company/cloudresty) &nbsp;|&nbsp; [BlueSky](https://bsky.app/profile/cloudresty.com) &nbsp;|&nbsp; [GitHub](https://github.com/cloudresty)

&nbsp;

# Document ID Generation Strategies

[Home](../README.md) &nbsp;/&nbsp; [Docs](README.md) &nbsp;/&nbsp; ID Generation Strategies

&nbsp;

This document explains the different document ID generation strategies available in go-elastic and when to use each one.

&nbsp;

## Available ID Modes

&nbsp;

### 1. Elasticsearch Native IDs (Default - Recommended)

Mode: `elastic`
Environment Variable: `ELASTICSEARCH_ID_MODE=elastic`
Environment-first approach (Recommended):

```bash
export ELASTICSEARCH_ID_MODE=elastic
```

```go
// Uses environment configuration
client, err := elastic.NewClient(elastic.FromEnv())
```

🔝 [back to top](#document-id-generation-strategies)

&nbsp;

Programmatic configuration:

```go
config := &elastic.Config{
    Hosts:  []string{"localhost:9200"},
    IDMode: elastic.IDModeElastic,
}
client, err := elastic.NewClient(elastic.WithConfig(config))
```

🔝 [back to top](#document-id-generation-strategies)

&nbsp;

How it works:

- Elasticsearch generates random, unique IDs automatically
- No `_id` field is set in the document - Elasticsearch handles it

&nbsp;

Advantages:

- **Optimal shard distribution** - Random IDs ensure even load across shards
- **Best write performance** - No hotspotting issues
- **Battle-tested** - Elasticsearch's default behavior
- **No additional dependencies** - Uses native ES functionality

&nbsp;

Disadvantages:

- IDs are not time-ordered or sortable
- No semantic meaning in the ID

&nbsp;

When to use:

- **Multi-shard indices** (most production use cases)
- **High write throughput** requirements
- **When you don't need sortable IDs**
- **Default choice for most applications**

🔝 [back to top](#document-id-generation-strategies)

&nbsp;

### 2. ULID (Time-Ordered IDs)

Mode: `ulid`
Environment Variable: `ELASTICSEARCH_ID_MODE=ulid`
Environment-first approach (Recommended):

```bash
export ELASTICSEARCH_ID_MODE=ulid
```

```go
// Uses environment configuration
client, err := elastic.NewClient(elastic.FromEnv())
```

🔝 [back to top](#document-id-generation-strategies)

&nbsp;

Programmatic configuration:

```go
config := &elastic.Config{
    Hosts:  []string{"localhost:9200"},
    IDMode: elastic.IDModeULID,
}
client, err := elastic.NewClient(elastic.WithConfig(config))
```

🔝 [back to top](#document-id-generation-strategies)

&nbsp;

How it works:

- Generates Universally Unique Lexicographically Sortable Identifiers
- Time-based prefix ensures chronological ordering

&nbsp;

Advantages:

- **Time-ordered** - IDs sort chronologically
- **URL-safe** - Can be used in URLs
- **Efficient range queries** - When querying by ID ranges

&nbsp;

Disadvantages:

- **Potential shard hotspotting** - Similar prefixes route to same shards
- **Uneven write load** - Newer documents may overload specific shards
- **Performance impact** - Can reduce write throughput in multi-shard setups

&nbsp;

When to use:

- **Single-shard indices**
- **Time-based indices** (daily/monthly rotation)
- **When you need sortable IDs** for application logic
- **When you control shard routing** manually
- **Low to medium write volume** scenarios

&nbsp;

Performance Warning:

In multi-shard indices, ULID can cause write hotspotting because documents with similar timestamps (created around the same time) will have similar ID prefixes and may route to the same shard.

🔝 [back to top](#document-id-generation-strategies)

&nbsp;

### 3. Custom IDs

Mode: `custom`
Environment Variable: `ELASTICSEARCH_ID_MODE=custom`
Environment-first approach (Recommended):

```bash
export ELASTICSEARCH_ID_MODE=custom
```

```go
// Uses environment configuration
client, err := elastic.NewClient(elastic.FromEnv())

// You must provide _id in your documents
doc := map[string]any{
    "_id": "user-12345",
    "name": "John Doe",
}
```

🔝 [back to top](#document-id-generation-strategies)

&nbsp;

Programmatic configuration:

```go
config := &elastic.Config{
    Hosts:  []string{"localhost:9200"},
    IDMode: elastic.IDModeCustom,
}
client, err := elastic.NewClient(elastic.WithConfig(config))

// You must provide _id in your documents
doc := map[string]any{
    "_id": "user-12345",
    "name": "John Doe",
}
```

🔝 [back to top](#document-id-generation-strategies)

&nbsp;

How it works:

- You provide the `_id` field in your documents
- No automatic ID generation

&nbsp;

Advantages:

- **Full control** over ID format and content
- **Semantic IDs** - Can embed meaning in the ID
- **Predictable** - You know exactly what the ID will be

&nbsp;

Disadvantages:

- **Manual management** - You must ensure uniqueness
- **Potential hotspotting** - Depending on your ID scheme
- **Additional complexity** - Error handling for ID conflicts

&nbsp;

When to use:

- **When you have natural unique identifiers** (user IDs, order numbers)
- **When you need predictable IDs** for external integrations
- **When documents have meaningful business keys**

&nbsp;

## Performance Considerations

&nbsp;

### Shard Distribution

Elasticsearch uses the document ID to determine which shard to route the document to. The routing algorithm is essentially:

```text
shard = hash(document_id) % number_of_shards
```

&nbsp;

Random IDs (elastic mode):

- Even distribution across all shards
- Consistent write performance

&nbsp;

Time-ordered IDs (ULID mode):

- Documents created at similar times have similar prefixes
- Similar prefixes → similar hash values → same shard
- Can cause 80/20 or 90/10 write distribution instead of even distribution

🔝 [back to top](#document-id-generation-strategies)

&nbsp;

### Write Throughput Impact

In a 5-shard index with high write volume:

| ID Mode | Shard Distribution | Write Performance |
|---------|-------------------|-------------------|
| `elastic` | ~20% per shard | Optimal |
| `ulid` | 60-90% on 1-2 shards | Reduced |
| `custom` | Depends on ID scheme | Varies |

🔝 [back to top](#document-id-generation-strategies)

&nbsp;

## Best Practices

&nbsp;

### 1. Default Choice: Elasticsearch Native

Unless you have specific requirements for sortable IDs, use the default `elastic` mode:

```bash
# No environment variable needed - this is the default
# export ELASTICSEARCH_ID_MODE=elastic
```

🔝 [back to top](#document-id-generation-strategies)

&nbsp;

### 2. ULID for Time-Based Workloads

Use ULID when you have time-based access patterns and can manage the sharding implications:

```bash
export ELASTICSEARCH_ID_MODE=ulid
```

Consider using time-based indices (daily/monthly) to mitigate hotspotting:

- `logs-2024.01.15`
- `logs-2024.01.16`
- etc.

🔝 [back to top](#document-id-generation-strategies)

&nbsp;

### 3. Custom IDs for Business Keys

Use custom IDs when you have natural business identifiers:

```bash
export ELASTICSEARCH_ID_MODE=custom
```

```go
// Good: Distributed business keys
userDoc := map[string]any{
    "_id": "user-uuid-" + uuid.New().String(),
    "email": "user@example.com",
}

// Avoid: Sequential keys that can cause hotspotting
orderDoc := map[string]any{
    "_id": fmt.Sprintf("order-%d", sequentialOrderNumber),
    "amount": 100.00,
}
```

🔝 [back to top](#document-id-generation-strategies)

&nbsp;

### 4. Monitoring Shard Distribution

Monitor your shard usage to detect hotspotting:

```bash
curl -X GET "localhost:9200/_cat/shards/your-index?v&h=index,shard,docs,store&s=shard"
```

Look for significant imbalances in document counts across shards.

🔝 [back to top](#document-id-generation-strategies)

&nbsp;

## Practical Usage Examples

&nbsp;

### Using ID Modes with Document Operations

```go
package main

import (
    "context"
    "fmt"
    "os"

    elastic "github.com/cloudresty/go-elastic"
)

type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

func main() {
    // Environment-first approach - ID mode set via ELASTICSEARCH_ID_MODE
    client, err := elastic.NewClient(elastic.FromEnv())
    if err != nil {
        panic(err)
    }
    defer client.Close()

    ctx := context.Background()
    docs := client.Documents()

    user := User{
        Name:  "John Doe",
        Email: "john@example.com",
        Age:   30,
    }

    // Create document - ID generation depends on ELASTICSEARCH_ID_MODE:
    // - elastic: Random Elasticsearch-generated ID
    // - ulid: Time-ordered ULID
    // - custom: Must include "_id" field in user struct or document
    result, err := docs.Create(ctx, "users", user)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Created document with ID: %s\n", result.ID)
}
```

🔝 [back to top](#document-id-generation-strategies)

&nbsp;

### Custom ID Mode Example

```go
// For custom ID mode, include _id in your document
doc := map[string]any{
    "_id":   "user-12345",
    "name":  "John Doe",
    "email": "john@example.com",
    "age":   30,
}

result, err := docs.Create(ctx, "users", doc)
// result.ID will be "user-12345"
```

🔝 [back to top](#document-id-generation-strategies)

&nbsp;

## Migration Between ID Modes

If you need to change ID modes for an existing index:

1. **Create a new index** with the desired ID mode
2. **Reindex data** using Elasticsearch's reindex API
3. **Update application** to use the new index
4. **Delete old index** after verification

```json
POST _reindex
{
  "source": {
    "index": "old-index"
  },
  "dest": {
    "index": "new-index"
  }
}
```

🔝 [back to top](#document-id-generation-strategies)

&nbsp;

## Conclusion

- **For most applications:** Use `elastic` mode (default)
- **For time-series data with time-based indices:** Consider `ulid` mode
- **For applications with natural business keys:** Use `custom` mode
- **Always monitor shard distribution** in production

The default `elastic` mode is chosen specifically to provide the best performance characteristics for the majority of Elasticsearch use cases.

🔝 [back to top](#document-id-generation-strategies)

&nbsp;

---

&nbsp;

An open source project brought to you by the [Cloudresty](https://cloudresty.com/) team.

[Website](https://cloudresty.com/) &nbsp;|&nbsp; [LinkedIn](https://www.linkedin.com/company/cloudresty) &nbsp;|&nbsp; [BlueSky](https://bsky.app/profile/cloudresty.com) &nbsp;|&nbsp; [GitHub](https://github.com/cloudresty)

&nbsp;

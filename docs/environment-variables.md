# Environment Variables

[Home](../README.md) &nbsp;/&nbsp; [Docs](README.md) &nbsp;/&nbsp; Environment Variables

&nbsp;

This document provides a comprehensive reference for all environment variables supported by the go-elastic package. For practical examples and usage scenarios, see the [Environment Configuration Guide](environment-configuration.md).

&nbsp;

## Connection Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `ELASTICSEARCH_HOSTS` | localhost:9200 | Elasticsearch hosts with ports (single: "host:port" or multiple: "host1:port1,host2:port2") |
| `ELASTICSEARCH_USERNAME` | "" | Username for authentication |
| `ELASTICSEARCH_PASSWORD` | "" | Password for authentication |
| `ELASTICSEARCH_API_KEY` | "" | API key for authentication |
| `ELASTICSEARCH_SERVICE_TOKEN` | "" | Service token for authentication |
| `ELASTICSEARCH_CLOUD_ID` | "" | Elastic Cloud ID |
| `ELASTICSEARCH_CONNECTION_NAME` | "" | Connection identifier for logging and monitoring |
| `ELASTICSEARCH_APP_NAME` | go-elastic-app | Application name for connection metadata |

[🔝 back to top](#environment-variables)

&nbsp;

## Security Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `ELASTICSEARCH_TLS_ENABLED` | false | Enable TLS/HTTPS connections |
| `ELASTICSEARCH_TLS_INSECURE` | false | Allow insecure TLS connections (skip certificate verification) |

[🔝 back to top](#environment-variables)

&nbsp;

## Document Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `ELASTICSEARCH_ID_MODE` | elastic | ID generation strategy: `elastic`, `ulid`, or `custom` |

[🔝 back to top](#environment-variables)

&nbsp;

## Connection Pool Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `ELASTICSEARCH_MAX_RETRIES` | 3 | Maximum retry attempts for failed requests |
| `ELASTICSEARCH_RETRY_ON_STATUS` | "" | HTTP status codes to retry on (comma-separated, e.g., "502,503,504") |
| `ELASTICSEARCH_MAX_IDLE_CONNS` | 100 | Maximum idle connections in the pool |
| `ELASTICSEARCH_MAX_IDLE_CONNS_PER_HOST` | 10 | Maximum idle connections per host |
| `ELASTICSEARCH_IDLE_CONN_TIMEOUT` | 90s | Idle connection timeout |
| `ELASTICSEARCH_MAX_CONN_LIFETIME` | 0s | Maximum connection lifetime (0 = no limit) |

[🔝 back to top](#environment-variables)

&nbsp;

## Timeout Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `ELASTICSEARCH_CONNECT_TIMEOUT` | 10s | Initial connection timeout |
| `ELASTICSEARCH_REQUEST_TIMEOUT` | 30s | Request operation timeout |

[🔝 back to top](#environment-variables)

&nbsp;

## Reconnection Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `ELASTICSEARCH_RECONNECT_ENABLED` | true | Enable automatic reconnection on connection loss |
| `ELASTICSEARCH_RECONNECT_DELAY` | 5s | Initial delay before reconnection attempts |
| `ELASTICSEARCH_MAX_RECONNECT_DELAY` | 1m | Maximum delay between reconnection attempts |
| `ELASTICSEARCH_RECONNECT_BACKOFF` | 2.0 | Backoff multiplier for reconnection delays |
| `ELASTICSEARCH_MAX_RECONNECT_ATTEMPTS` | 10 | Maximum number of reconnection attempts |

[🔝 back to top](#environment-variables)

&nbsp;

## Health Check Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `ELASTICSEARCH_HEALTH_CHECK_ENABLED` | true | Enable periodic health checks |
| `ELASTICSEARCH_HEALTH_CHECK_INTERVAL` | 30s | Interval between health checks |

[🔝 back to top](#environment-variables)

&nbsp;

## Performance Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `ELASTICSEARCH_COMPRESSION_ENABLED` | true | Enable request/response compression |
| `ELASTICSEARCH_DISCOVER_NODES_ON_START` | false | Enable node discovery on client startup |

[🔝 back to top](#environment-variables)

&nbsp;

## Logging Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `ELASTICSEARCH_LOG_LEVEL` | info | Logging level: `debug`, `info`, `warn`, `error` |
| `ELASTICSEARCH_LOG_FORMAT` | json | Log output format: `json` or `text` |

[🔝 back to top](#environment-variables)

&nbsp;

## Custom Prefixes

All environment variables can be prefixed with a custom string when using `FromEnvWithPrefix()`:

```go
// Use MYAPP_ELASTICSEARCH_* instead of ELASTICSEARCH_*
client, err := elastic.NewClient(elastic.FromEnvWithPrefix("MYAPP_"))
```

For example:

- `MYAPP_ELASTICSEARCH_HOSTS` instead of `ELASTICSEARCH_HOSTS`
- `ORDERS_ELASTICSEARCH_USERNAME` instead of `ELASTICSEARCH_USERNAME`
- `PAYMENTS_ELASTICSEARCH_API_KEY` instead of `ELASTICSEARCH_API_KEY`

[🔝 back to top](#environment-variables)

&nbsp;

---

&nbsp;

An open source project brought to you by the [Cloudresty](https://cloudresty.com/) team.

[Website](https://cloudresty.com/) &nbsp;|&nbsp; [LinkedIn](https://www.linkedin.com/company/cloudresty) &nbsp;|&nbsp; [BlueSky](https://bsky.app/profile/cloudresty.com) &nbsp;|&nbsp; [GitHub](https://github.com/cloudresty)

&nbsp;

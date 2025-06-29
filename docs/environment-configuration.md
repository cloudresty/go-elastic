# Environment Configuration

[Home](../README.md) &nbsp;/&nbsp; [Docs](README.md) &nbsp;/&nbsp; Environment Configuration

&nbsp;

This guide provides practical examples and usage scenarios for configuring your Elasticsearch client using environment variables. For a complete reference of all supported variables, see the [Environment Variables Reference](environment-variables.md).

&nbsp;

## Quick Start

```go
// Load with default ELASTICSEARCH_ prefix
client, err := elastic.NewClient(elastic.FromEnv())
if err != nil {
    log.Fatal(err)
}

// Load and customize with functional options
client, err := elastic.NewClient(
    elastic.FromEnv(),
    elastic.WithConnectionName("my-custom-service"),
)
```

**Required environment variable:**

```bash
export ELASTICSEARCH_HOSTS=localhost:9200
```

[🔝 back to top](#environment-configuration)

&nbsp;

## Custom Prefix

```go
// Use custom prefix (e.g., MYAPP_ELASTICSEARCH_HOSTS instead of ELASTICSEARCH_HOSTS)
client, err := elastic.NewClient(elastic.FromEnvWithPrefix("MYAPP_"))
```

**Custom prefix example:**

```bash
export MYAPP_ELASTICSEARCH_HOSTS=elasticsearch.example.com:9200
export MYAPP_ELASTICSEARCH_USERNAME=myuser
export MYAPP_ELASTICSEARCH_PASSWORD=mypass
```

[🔝 back to top](#environment-configuration)

&nbsp;

## Connection Examples

&nbsp;

### Single Host

```bash
# Simple single host
export ELASTICSEARCH_HOSTS=localhost:9200

# Single host with authentication
export ELASTICSEARCH_HOSTS=elasticsearch.example.com:9200
export ELASTICSEARCH_USERNAME=myuser
export ELASTICSEARCH_PASSWORD=secretpass

# Single host with API key
export ELASTICSEARCH_HOSTS=elasticsearch.example.com:9200
export ELASTICSEARCH_API_KEY=VnVhQ2ZHY0JDZGJrUW0tZTVhT3g6dWkybHA2MWlhR20ta2NhZUNGMGg3Zw==

# Single host with TLS
export ELASTICSEARCH_HOSTS=elasticsearch.example.com:9200
export ELASTICSEARCH_TLS_ENABLED=true
export ELASTICSEARCH_USERNAME=prod_user
export ELASTICSEARCH_PASSWORD=prod_password
```

[🔝 back to top](#environment-configuration)

&nbsp;

### Multiple Hosts (Cluster)

```bash
# Multiple hosts on same port
export ELASTICSEARCH_HOSTS=es1.cluster.com:9200,es2.cluster.com:9200,es3.cluster.com:9200

# Multiple hosts on different ports
export ELASTICSEARCH_HOSTS=es1.example.com:9200,es2.example.com:9201,es3.example.com:9202

# Multiple hosts with authentication
export ELASTICSEARCH_HOSTS=node1.cluster.com:9200,node2.cluster.com:9200,node3.cluster.com:9200
export ELASTICSEARCH_USERNAME=cluster_user
export ELASTICSEARCH_PASSWORD=cluster_password
export ELASTICSEARCH_TLS_ENABLED=true
```

[🔝 back to top](#environment-configuration)

&nbsp;

### Elastic Cloud

```bash
# Using Elastic Cloud ID
export ELASTICSEARCH_CLOUD_ID=my-deployment:dXMtY2VudHJhbDEuZ2NwLmNsb3VkLmVzLmlvJGFiY2RlZjEyMzQ1Njc4OTA=
export ELASTICSEARCH_USERNAME=elastic
export ELASTICSEARCH_PASSWORD=changeme

# Using service token
export ELASTICSEARCH_CLOUD_ID=my-deployment:dXMtY2VudHJhbDEuZ2NwLmNsb3VkLmVzLmlvJGFiY2RlZjEyMzQ1Njc4OTA=
export ELASTICSEARCH_SERVICE_TOKEN=AAEAAWVsYXN0aWMvZmxlZXQtc2VydmVyL3Rva2VuLTE6TnNhLXRydEdUeWl2Yjh3VzFJWXpfQQ
```

[🔝 back to top](#environment-configuration)

&nbsp;

## Application Configuration Examples

### Development Environment

```bash
# Minimal development setup
export ELASTICSEARCH_HOSTS=localhost:9200
export ELASTICSEARCH_ID_MODE=elastic
export ELASTICSEARCH_LOG_LEVEL=debug
export ELASTICSEARCH_CONNECTION_NAME=myapp-dev
export ELASTICSEARCH_APP_NAME=myapp-development
```

[🔝 back to top](#environment-configuration)

&nbsp;

### Production Environment

```bash
# Production cluster configuration
export ELASTICSEARCH_HOSTS=es1.prod.com:9200,es2.prod.com:9200,es3.prod.com:9200
export ELASTICSEARCH_USERNAME=prod_user
export ELASTICSEARCH_PASSWORD=secure_production_password
export ELASTICSEARCH_TLS_ENABLED=true
export ELASTICSEARCH_ID_MODE=elastic
export ELASTICSEARCH_CONNECTION_NAME=myapp-production
export ELASTICSEARCH_APP_NAME=myapp-production

# Performance and reliability settings
export ELASTICSEARCH_MAX_RETRIES=5
export ELASTICSEARCH_CONNECT_TIMEOUT=30s
export ELASTICSEARCH_REQUEST_TIMEOUT=60s
export ELASTICSEARCH_COMPRESSION_ENABLED=true
export ELASTICSEARCH_DISCOVER_NODES_ON_START=true

# Health monitoring
export ELASTICSEARCH_HEALTH_CHECK_ENABLED=true
export ELASTICSEARCH_HEALTH_CHECK_INTERVAL=30s

# Structured logging
export ELASTICSEARCH_LOG_LEVEL=info
export ELASTICSEARCH_LOG_FORMAT=json
```

[🔝 back to top](#environment-configuration)

&nbsp;

### High-Performance Environment

```bash
# Cluster configuration for high throughput
export ELASTICSEARCH_HOSTS=es1.perf.com:9200,es2.perf.com:9200,es3.perf.com:9200
export ELASTICSEARCH_USERNAME=perf_user
export ELASTICSEARCH_PASSWORD=perf_password

# Optimized connection pool
export ELASTICSEARCH_MAX_IDLE_CONNS=200
export ELASTICSEARCH_MAX_IDLE_CONNS_PER_HOST=50
export ELASTICSEARCH_IDLE_CONN_TIMEOUT=120s
export ELASTICSEARCH_MAX_CONN_LIFETIME=30m

# Retry configuration for resilience
export ELASTICSEARCH_MAX_RETRIES=3
export ELASTICSEARCH_RETRY_ON_STATUS=502,503,504
export ELASTICSEARCH_RECONNECT_ENABLED=true
export ELASTICSEARCH_MAX_RECONNECT_ATTEMPTS=15

# Performance optimizations
export ELASTICSEARCH_COMPRESSION_ENABLED=true
export ELASTICSEARCH_DISCOVER_NODES_ON_START=true
```

[🔝 back to top](#environment-configuration)

&nbsp;

## Multi-Service Architecture

[🔝 back to top](#environment-configuration)

&nbsp;

### Service A (Payments)

```bash
# Payments service configuration
export PAYMENTS_ELASTICSEARCH_HOSTS=elasticsearch-payments.internal:9200
export PAYMENTS_ELASTICSEARCH_USERNAME=payments_user
export PAYMENTS_ELASTICSEARCH_PASSWORD=payments_password
export PAYMENTS_ELASTICSEARCH_ID_MODE=ulid
export PAYMENTS_ELASTICSEARCH_CONNECTION_NAME=payments-service
export PAYMENTS_ELASTICSEARCH_APP_NAME=payments-microservice
export PAYMENTS_ELASTICSEARCH_LOG_LEVEL=info
```

[🔝 back to top](#environment-configuration)

&nbsp;

### Service B (Orders)

```bash
# Orders service configuration
export ORDERS_ELASTICSEARCH_HOSTS=elasticsearch-orders.internal:9200
export ORDERS_ELASTICSEARCH_USERNAME=orders_user
export ORDERS_ELASTICSEARCH_PASSWORD=orders_password
export ORDERS_ELASTICSEARCH_ID_MODE=elastic
export ORDERS_ELASTICSEARCH_CONNECTION_NAME=orders-service
export ORDERS_ELASTICSEARCH_APP_NAME=orders-microservice
export ORDERS_ELASTICSEARCH_LOG_LEVEL=info
```

[🔝 back to top](#environment-configuration)

&nbsp;

### Application Code

```go
// Initialize both services with different prefixes
paymentsClient, err := elastic.NewClient(elastic.FromEnvWithPrefix("PAYMENTS_"))
if err != nil {
    log.Fatal("Failed to create payments client:", err)
}

ordersClient, err := elastic.NewClient(elastic.FromEnvWithPrefix("ORDERS_"))
if err != nil {
    log.Fatal("Failed to create orders client:", err)
}
```

[🔝 back to top](#environment-configuration)

&nbsp;

## Container Deployments

&nbsp;

### Environment File (.env)

```bash
# Production environment file
ELASTICSEARCH_HOSTS=elasticsearch.production.com:9200
ELASTICSEARCH_USERNAME=myuser
ELASTICSEARCH_PASSWORD=mypassword
ELASTICSEARCH_TLS_ENABLED=true
ELASTICSEARCH_ID_MODE=elastic
ELASTICSEARCH_CONNECTION_NAME=my-production-service
ELASTICSEARCH_APP_NAME=my-production-app
ELASTICSEARCH_MAX_RETRIES=5
ELASTICSEARCH_CONNECT_TIMEOUT=15s
ELASTICSEARCH_COMPRESSION_ENABLED=true
ELASTICSEARCH_DISCOVER_NODES_ON_START=true
ELASTICSEARCH_HEALTH_CHECK_ENABLED=true
ELASTICSEARCH_LOG_LEVEL=info
ELASTICSEARCH_LOG_FORMAT=json
```

[🔝 back to top](#environment-configuration)

&nbsp;

### Docker Compose

```yaml
version: '3.8'
services:
  my-app:
    image: my-app:latest
    environment:
      ELASTICSEARCH_HOSTS: elasticsearch:9200
      ELASTICSEARCH_ID_MODE: elastic
      ELASTICSEARCH_USERNAME: app_user
      ELASTICSEARCH_PASSWORD: secure_password
      ELASTICSEARCH_CONNECTION_NAME: my-app-instance
      ELASTICSEARCH_APP_NAME: my-app
      ELASTICSEARCH_LOG_LEVEL: debug
    depends_on:
      - elasticsearch

  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.11.0
    environment:
      discovery.type: single-node
      xpack.security.enabled: false
      ES_JAVA_OPTS: "-Xms512m -Xmx512m"
```

[🔝 back to top](#environment-configuration)

&nbsp;

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
spec:
  template:
    spec:
      containers:
      - name: my-app
        image: my-app:latest
        env:
        - name: ELASTICSEARCH_HOSTS
          value: "elasticsearch-service:9200"
        - name: ELASTICSEARCH_ID_MODE
          value: "elastic"
        - name: ELASTICSEARCH_USERNAME
          valueFrom:
            secretKeyRef:
              name: elasticsearch-secret
              key: username
        - name: ELASTICSEARCH_PASSWORD
          valueFrom:
            secretKeyRef:
              name: elasticsearch-secret
              key: password
        - name: ELASTICSEARCH_CONNECTION_NAME
          value: "my-app-pod"
        - name: ELASTICSEARCH_APP_NAME
          value: "my-app"
        - name: ELASTICSEARCH_TLS_ENABLED
          value: "true"
        - name: ELASTICSEARCH_LOG_LEVEL
          value: "info"
```

[🔝 back to top](#environment-configuration)

&nbsp;

## Best Practices

&nbsp;

### Security

- **Never hardcode credentials** - Always use environment variables or secrets
- **Use HTTPS in production** - Set `ELASTICSEARCH_TLS_ENABLED=true`
- **Rotate passwords regularly** - Update environment variables during deployments
- **Use API keys when possible** - More secure than username/password authentication

```bash
# Secure production setup
export ELASTICSEARCH_HOSTS=elasticsearch.prod.com:9200
export ELASTICSEARCH_API_KEY=VnVhQ2ZHY0JDZGJrUW0tZTVhT3g6dWkybHA2MWlhR20ta2NhZUNGMGg3Zw==
export ELASTICSEARCH_TLS_ENABLED=true
export ELASTICSEARCH_CONNECTION_NAME=secure-app
```

[🔝 back to top](#environment-configuration)

&nbsp;

### Performance

- **Size connection pools appropriately** - Adjust `ELASTICSEARCH_MAX_IDLE_CONNS` based on your workload
- **Configure timeouts** - Set realistic timeouts for your network conditions
- **Enable compression** - Use `ELASTICSEARCH_COMPRESSION_ENABLED=true` for network-bound workloads
- **Use node discovery carefully** - Enable `ELASTICSEARCH_DISCOVER_NODES_ON_START=true` only in stable environments

```bash
# High-performance setup
export ELASTICSEARCH_HOSTS=es1.fast.com:9200,es2.fast.com:9200,es3.fast.com:9200
export ELASTICSEARCH_MAX_IDLE_CONNS=200
export ELASTICSEARCH_MAX_IDLE_CONNS_PER_HOST=50
export ELASTICSEARCH_CONNECT_TIMEOUT=30s
export ELASTICSEARCH_REQUEST_TIMEOUT=60s
export ELASTICSEARCH_COMPRESSION_ENABLED=true
export ELASTICSEARCH_DISCOVER_NODES_ON_START=true
```

[🔝 back to top](#environment-configuration)

&nbsp;

### Monitoring

- **Use descriptive connection names** - Set `ELASTICSEARCH_CONNECTION_NAME` for better monitoring
- **Enable health checks** - Use `ELASTICSEARCH_HEALTH_CHECK_ENABLED=true`
- **Configure structured logging** - Use `ELASTICSEARCH_LOG_FORMAT=json` in production
- **Set appropriate log levels** - Use `debug` for development, `info` for production

```bash
# Monitoring-friendly setup
export ELASTICSEARCH_HOSTS=elasticsearch.monitored.com:9200
export ELASTICSEARCH_CONNECTION_NAME=monitored-app-prod
export ELASTICSEARCH_APP_NAME=monitored-application
export ELASTICSEARCH_HEALTH_CHECK_ENABLED=true
export ELASTICSEARCH_HEALTH_CHECK_INTERVAL=30s
export ELASTICSEARCH_LOG_LEVEL=info
export ELASTICSEARCH_LOG_FORMAT=json
```

[🔝 back to top](#environment-configuration)

&nbsp;

---

&nbsp;

An open source project brought to you by the [Cloudresty](https://cloudresty.com/) team.

[Website](https://cloudresty.com/) &nbsp;|&nbsp; [LinkedIn](https://www.linkedin.com/company/cloudresty) &nbsp;|&nbsp; [BlueSky](https://bsky.app/profile/cloudresty.com) &nbsp;|&nbsp; [GitHub](https://github.com/cloudresty)

&nbsp;

# Production Features Demo

This example demonstrates production-ready features of the go-elastic library using the `emit` logging library.

## Features Demonstrated

- **Custom Logger Integration**: Shows how to implement the `elastic.Logger` interface using the `emit` logging library
- **Multi-Client Configuration**: Configure multiple Elasticsearch clients with different environment prefixes
- **Graceful Shutdown**: Proper application shutdown with configurable timeouts
- **Background Health Checking**: Continuous monitoring of client health
- **Production Configuration**: Environment-based configuration using prefixes

## Setup

1. Copy the environment configuration:
   ```bash
   cp .env.example .env
   ```

2. Configure your Elasticsearch instances in `.env`

3. Build and run:
   ```bash
   go build .
   ./production-features
   ```

## Key Implementation Details

### EmitLogger Implementation

The example shows how to create a custom logger that implements the `elastic.Logger` interface:

```go
type EmitLogger struct{}

func (e *EmitLogger) Info(msg string, fields ...any) {
    if len(fields) > 0 {
        emit.Info.KeyValue(msg, fields...)
    } else {
        emit.Info.Msg(msg)
    }
}
// ... similar for Warn, Error, Debug
```

### Client Configuration

```go
paymentsClient, err := elastic.NewClient(
    elastic.FromEnvWithPrefix("PAYMENTS_"),
    elastic.WithConnectionName("payments-cluster"),
    elastic.WithLogger(logger),
)
```

### Shutdown Management

```go
shutdownManager := elastic.NewShutdownManager(&elastic.ShutdownConfig{
    Timeout:          30 * time.Second,
    GracePeriod:      5 * time.Second,
    ForceKillTimeout: 10 * time.Second,
}, logger)
```

## Dependencies

This example uses its own `go.mod` to independently manage the `emit` logging library dependency while using the local go-elastic library via a replace directive.

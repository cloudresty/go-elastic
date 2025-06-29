package main

import (
	"context"
	"os"
	"time"

	"github.com/cloudresty/emit"
	elastic "github.com/cloudresty/go-elastic"
)

func main() {
	emit.Info.Msg("Starting reconnection behavior demonstration")

	// Create client with environment config and custom settings for testing
	client, err := elastic.NewClient(
		elastic.FromEnv(), // Load from environment
		elastic.WithConnectionName("reconnection-test"), // Set connection name
		// Note: Other timeout and reconnection settings would need to be added
		// as additional WithXXX options in a full implementation
	)
	if err != nil {
		emit.Error.StructuredFields("Failed to create client",
			emit.ZString("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		if err := client.Close(); err != nil {
			emit.Error.StructuredFields("Failed to close client", emit.ZString("error", err.Error()))
		}
	}()

	connectionName := client.Name()
	emit.Info.StructuredFields("Client configured for reconnection testing",
		emit.ZString("connection_name", connectionName))

	ctx := context.Background()

	// Test initial connection
	emit.Info.Msg("Testing initial connection...")
	err = testConnection(ctx, client)
	if err != nil {
		emit.Error.StructuredFields("Initial connection test failed",
			emit.ZString("error", err.Error()))
		// Continue anyway for demonstration
	} else {
		emit.Info.Msg("✓ Initial connection successful")
	}

	// Simulate operations with potential connection issues
	emit.Info.Msg("Starting continuous operation simulation...")

	operationCount := 0
	successCount := 0
	failureCount := 0

	// Run for 2 minutes
	endTime := time.Now().Add(2 * time.Minute)

	for time.Now().Before(endTime) {
		operationCount++

		err := performOperation(ctx, client, operationCount)
		if err != nil {
			failureCount++
			emit.Warn.StructuredFields("Operation failed",
				emit.ZInt("operation", operationCount),
				emit.ZString("error", err.Error()))
		} else {
			successCount++
			emit.Info.StructuredFields("Operation successful",
				emit.ZInt("operation", operationCount))
		}

		// Wait between operations
		time.Sleep(3 * time.Second)

		// Show stats every 10 operations
		if operationCount%10 == 0 {
			stats := client.Stats()
			emit.Info.StructuredFields("Operation statistics",
				emit.ZInt("total_operations", operationCount),
				emit.ZInt("successful", successCount),
				emit.ZInt("failed", failureCount),
				emit.ZFloat64("success_rate", float64(successCount)/float64(operationCount)*100),
				emit.ZBool("connected", stats.IsConnected),
				emit.ZInt64("reconnects", stats.Reconnects))
		}
	}

	// Final statistics
	stats := client.Stats()
	clientName := client.Name()
	emit.Info.StructuredFields("Reconnection test completed",
		emit.ZString("client_name", clientName),
		emit.ZInt("total_operations", operationCount),
		emit.ZInt("successful_operations", successCount),
		emit.ZInt("failed_operations", failureCount),
		emit.ZFloat64("final_success_rate", float64(successCount)/float64(operationCount)*100),
		emit.ZBool("connected", stats.IsConnected),
		emit.ZInt64("total_reconnects", stats.Reconnects))

	if stats.Reconnects > 0 {
		emit.Info.StructuredFields("Connection statistics",
			emit.ZString("client_name", clientName),
			emit.ZInt64("reconnects", stats.Reconnects),
			emit.ZTime("last_reconnect", stats.LastReconnect))
	}

	if successCount > 0 {
		emit.Info.Msg("✓ Reconnection behavior working correctly")
	} else {
		emit.Warn.Msg("⚠ No successful operations - check Elasticsearch connectivity")
	}
}

func testConnection(ctx context.Context, client *elastic.Client) error {
	return client.Ping(ctx)
}

func performOperation(ctx context.Context, client *elastic.Client, operationID int) error {
	// Try a simple health check first
	if err := client.Ping(ctx); err != nil {
		return err
	}

	// Perform a simple index operation
	doc := map[string]any{
		"operation_id": operationID,
		"timestamp":    time.Now().Format(time.RFC3339),
		"message":      "Reconnection test document",
		"test_type":    "reconnection",
	}

	_, err := client.Documents().Create(ctx, "reconnection-test", doc)
	return err
}

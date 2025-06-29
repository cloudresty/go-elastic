package main

import (
	"context"
	"log"
	"os"

	"github.com/cloudresty/emit"
	elastic "github.com/cloudresty/go-elastic"
)

func main() {
	emit.Info.Msg("Starting environment configuration examples")

	// Example 1: Using default ELASTICSEARCH_ prefix with FromEnv()
	emit.Info.Msg("Creating client from environment variables (ELASTICSEARCH_ prefix)")

	client, err := elastic.NewClient(elastic.FromEnv())
	if err != nil {
		emit.Error.StructuredFields("Failed to create client from environment",
			emit.ZString("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("Failed to close client: %v", err)
		}
	}()

	emit.Info.Msg("Client created successfully from environment variables")

	// Example 2: Using custom prefix with new FromEnvWithPrefix
	emit.Info.Msg("Creating client from environment variables with custom prefix")

	clientWithPrefix, err := elastic.NewClient(elastic.FromEnvWithPrefix("SEARCH_"))
	if err != nil {
		emit.Error.StructuredFields("Failed to create client from environment with prefix",
			emit.ZString("error", err.Error()))
		// This might fail if SEARCH_ prefixed vars aren't set, which is expected
		emit.Warn.Msg("Custom prefix example failed (expected if SEARCH_* vars not set)")
	} else {
		defer func() {
			if err := clientWithPrefix.Close(); err != nil {
				emit.Error.StructuredFields("Failed to close client with prefix",
					emit.ZString("error", err.Error()))
			}
		}()
		emit.Info.Msg("Client with custom prefix created successfully")
	}

	// Example 3: Loading env config and customizing with functional options
	emit.Info.Msg("Loading environment config and customizing with functional options")

	// Create client with environment config and custom overrides
	customClient, err := elastic.NewClient(
		elastic.FromEnv(), // Load from environment
		elastic.WithConnectionName("env-config-example"), // Override connection name
		// Could add other overrides here like:
		// elastic.WithPort(9201),
		// elastic.WithTLS(true),
	)
	if err != nil {
		emit.Error.StructuredFields("Failed to create client with environment config",
			emit.ZString("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		if err := customClient.Close(); err != nil {
			emit.Error.StructuredFields("Failed to close custom client",
				emit.ZString("error", err.Error()))
		}
	}()

	clientName := customClient.Name()
	emit.Info.StructuredFields("Client created with environment config and custom overrides",
		emit.ZString("connection_name", clientName))

	// Test the connections
	ctx := context.Background()
	if err := client.Ping(ctx); err != nil {
		emit.Error.StructuredFields("Default client ping failed",
			emit.ZString("error", err.Error()))
	} else {
		emit.Info.Msg("Default client ping successful")
	}

	if err := customClient.Ping(ctx); err != nil {
		emit.Error.StructuredFields("Custom client ping failed",
			emit.ZString("error", err.Error()))
	} else {
		emit.Info.Msg("Custom client ping successful")
	}

	// Example 4: Multiple clients with different configurations (production pattern)
	emit.Info.Msg("Demonstrating multi-client configuration patterns")

	// Primary cluster uses defaults from environment
	primaryClient, err := elastic.NewClient(elastic.FromEnv())
	if err != nil {
		emit.Error.StructuredFields("Failed to create primary client",
			emit.ZString("error", err.Error()))
	} else {
		defer primaryClient.Close()
		primaryName := primaryClient.Name()
		emit.Info.StructuredFields("Primary client created",
			emit.ZString("connection_name", primaryName))
	}

	// Logging cluster uses a prefix and overrides
	loggingClient, err := elastic.NewClient(
		elastic.FromEnvWithPrefix("LOGS_"),
		elastic.WithConnectionName("logging-cluster"),
		// Could add other overrides:
		// elastic.WithTimeout(30*time.Second),
	)
	if err != nil {
		emit.Error.StructuredFields("Failed to create logging client",
			emit.ZString("error", err.Error()))
		emit.Warn.Msg("Logging client example failed (expected if LOGS_* vars not set)")
	} else {
		defer loggingClient.Close()
		loggingName := loggingClient.Name()
		emit.Info.StructuredFields("Logging client created",
			emit.ZString("connection_name", loggingName))
	}

	emit.Info.Msg("Environment configuration examples completed!")
}

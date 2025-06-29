package main

import (
	"context"
	"os"
	"time"

	"github.com/cloudresty/emit"
	elastic "github.com/cloudresty/go-elastic"
)

func main() {
	emit.Info.Msg("Starting ID generation strategies demonstration")

	ctx := context.Background()

	// Example 1: Elastic Native IDs (Default - Recommended)
	emit.Info.Msg("=== Elasticsearch Native ID Generation (Default) ===")
	err := demonstrateElasticIDs(ctx)
	if err != nil {
		emit.Error.StructuredFields("Elastic IDs demo failed",
			emit.ZString("error", err.Error()))
	}

	time.Sleep(1 * time.Second)

	// Example 2: ULID IDs (Time-ordered with hotspotting warning)
	emit.Info.Msg("=== ULID ID Generation (Time-ordered) ===")
	err = demonstrateULIDIDs(ctx)
	if err != nil {
		emit.Error.StructuredFields("ULID IDs demo failed",
			emit.ZString("error", err.Error()))
	}

	time.Sleep(1 * time.Second)

	// Example 3: Custom IDs (User-provided)
	emit.Info.Msg("=== Custom ID Generation (User-provided) ===")
	err = demonstrateCustomIDs(ctx)
	if err != nil {
		emit.Error.StructuredFields("Custom IDs demo failed",
			emit.ZString("error", err.Error()))
	}

	time.Sleep(1 * time.Second)

	// Example 4: Environment-based configuration
	emit.Info.Msg("=== Environment-based ID Configuration ===")
	err = demonstrateEnvironmentIDs(ctx)
	if err != nil {
		emit.Error.StructuredFields("Environment IDs demo failed",
			emit.ZString("error", err.Error()))
	}

	emit.Info.Msg("ID generation demonstration completed!")
}

func demonstrateElasticIDs(ctx context.Context) error {
	emit.Info.Msg("Using Elasticsearch native ID generation (optimal shard distribution)")

	// Default config uses Elastic native IDs
	client, err := elastic.NewClient()
	if err != nil {
		return err
	}
	defer func() {
		if err := client.Close(); err != nil {
			emit.Error.StructuredFields("Failed to close client", emit.ZString("error", err.Error()))
		}
	}()

	emit.Info.Msg("✓ Client configured for Elasticsearch native IDs")

	// Index multiple documents to show different IDs
	for i := 0; i < 5; i++ {
		doc := map[string]any{
			"title":       "Document with Elastic ID",
			"description": "This document uses Elasticsearch's native ID generation",
			"sequence":    i + 1,
			"timestamp":   time.Now().Format(time.RFC3339),
			"id_type":     "elastic_native",
		}

		result, err := client.Documents().Create(ctx, "elastic-ids-demo", doc)
		if err != nil {
			emit.Error.StructuredFields("Failed to index document",
				emit.ZInt("sequence", i+1),
				emit.ZString("error", err.Error()))
			continue
		}

		emit.Info.StructuredFields("Document indexed with Elastic ID",
			emit.ZInt("sequence", i+1),
			emit.ZString("document_id", result.ID),
			emit.ZString("result", result.Result),
			emit.ZString("characteristics", "Random, optimal distribution"))
	}

	return nil
}

func demonstrateULIDIDs(ctx context.Context) error {
	emit.Warn.Msg("ULID IDs provide time-ordering but may cause hotspotting in high-write scenarios")

	config := &elastic.Config{
		Hosts:  []string{"localhost:9200"},
		IDMode: elastic.IDModeULID,
	}

	client, err := elastic.NewClient(elastic.WithConfig(config))
	if err != nil {
		return err
	}
	defer func() {
		if err := client.Close(); err != nil {
			emit.Error.StructuredFields("Failed to close ULID client", emit.ZString("error", err.Error()))
		}
	}()

	emit.Info.Msg("✓ Client configured for ULID generation")

	// Index multiple documents to show ULID progression
	for i := 0; i < 5; i++ {
		doc := map[string]any{
			"title":       "Document with ULID",
			"description": "This document uses ULID generation for time-ordering",
			"sequence":    i + 1,
			"timestamp":   time.Now().Format(time.RFC3339),
			"id_type":     "ulid",
		}

		result, err := client.Documents().Create(ctx, "ulid-demo", doc)
		if err != nil {
			emit.Error.StructuredFields("Failed to index document",
				emit.ZInt("sequence", i+1),
				emit.ZString("error", err.Error()))
			continue
		}

		emit.Info.StructuredFields("Document indexed with ULID",
			emit.ZInt("sequence", i+1),
			emit.ZString("document_id", result.ID),
			emit.ZString("result", result.Result),
			emit.ZString("characteristics", "Time-ordered, lexicographically sortable"))

		// Small delay to show time progression in ULIDs
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

func demonstrateCustomIDs(ctx context.Context) error {
	emit.Info.Msg("Using custom user-provided IDs")

	config := &elastic.Config{
		Hosts:  []string{"localhost:9200"},
		IDMode: elastic.IDModeCustom,
	}

	client, err := elastic.NewClient(elastic.WithConfig(config))
	if err != nil {
		return err
	}
	defer func() {
		if err := client.Close(); err != nil {
			emit.Error.StructuredFields("Failed to close custom ID client", emit.ZString("error", err.Error()))
		}
	}()

	emit.Info.Msg("✓ Client configured for custom ID mode")

	// Define custom IDs for different use cases
	customIDs := []struct {
		id          string
		description string
		useCase     string
	}{
		{"user-12345", "User profile document", "User management"},
		{"order-2023-12-001", "E-commerce order", "Order tracking"},
		{"session-abc123def456", "User session data", "Session management"},
		{"product-sku-LAPTOP001", "Product catalog entry", "Inventory management"},
		{"log-2023-12-01-001", "Application log entry", "Logging system"},
	}

	for i, item := range customIDs {
		doc := map[string]any{
			"title":       "Document with Custom ID",
			"description": item.description,
			"use_case":    item.useCase,
			"sequence":    i + 1,
			"timestamp":   time.Now().Format(time.RFC3339),
			"id_type":     "custom",
		}

		result, err := client.Documents().CreateWithID(ctx, "custom-ids-demo", item.id, doc)
		if err != nil {
			emit.Error.StructuredFields("Failed to index document",
				emit.ZInt("sequence", i+1),
				emit.ZString("custom_id", item.id),
				emit.ZString("error", err.Error()))
			continue
		}

		emit.Info.StructuredFields("Document indexed with custom ID",
			emit.ZInt("sequence", i+1),
			emit.ZString("document_id", result.ID),
			emit.ZString("result", result.Result),
			emit.ZString("use_case", item.useCase),
			emit.ZString("characteristics", "User-controlled, semantic meaning"))
	}

	return nil
}

func demonstrateEnvironmentIDs(ctx context.Context) error {
	emit.Info.Msg("Demonstrating environment-based ID configuration")

	// Get current environment setting
	idMode := os.Getenv("ELASTICSEARCH_ID_MODE")
	if idMode == "" {
		idMode = "elastic" // default
	}

	emit.Info.StructuredFields("Environment configuration",
		emit.ZString("ELASTICSEARCH_ID_MODE", idMode),
		emit.ZString("note", "Client will use this mode automatically"))

	// Create client using environment variables
	client, err := elastic.NewClient()
	if err != nil {
		return err
	}
	defer func() {
		if err := client.Close(); err != nil {
			emit.Error.StructuredFields("Failed to close environment client", emit.ZString("error", err.Error()))
		}
	}()

	emit.Info.StructuredFields("Client created from environment",
		emit.ZString("id_mode", idMode))

	// Index a document using environment-determined ID mode
	doc := map[string]any{
		"title":       "Environment-configured Document",
		"description": "This document's ID generation depends on environment variables",
		"timestamp":   time.Now().Format(time.RFC3339),
		"id_type":     "environment_" + idMode,
		"env_mode":    idMode,
	}

	result, err := client.Documents().Create(ctx, "env-config-demo", doc)
	if err != nil {
		return err
	}

	emit.Info.StructuredFields("Document indexed with environment ID mode",
		emit.ZString("document_id", result.ID),
		emit.ZString("result", result.Result),
		emit.ZString("configured_mode", idMode),
		emit.ZString("benefits", "Easy deployment configuration"))

	// Show how to override for different environments
	emit.Info.Msg("Example environment configurations:")
	emit.Info.Msg("  Production:   export ELASTICSEARCH_ID_MODE=elastic")
	emit.Info.Msg("  Development:  export ELASTICSEARCH_ID_MODE=ulid")
	emit.Info.Msg("  Integration:  export ELASTICSEARCH_ID_MODE=custom")

	return nil
}

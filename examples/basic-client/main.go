package main

import (
	"context"
	"log"
	"os"

	"github.com/cloudresty/emit"
	elastic "github.com/cloudresty/go-elastic"
)

func main() {
	emit.Info.Msg("Starting basic Elasticsearch client example")

	// Create client from environment variables
	client, err := elastic.NewClient()
	if err != nil {
		emit.Error.StructuredFields("Failed to create client",
			emit.ZString("error", err.Error()),
			emit.ZString("hint", "Set ELASTICSEARCH_* environment variables or defaults will be used"))
		os.Exit(1)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("Failed to close client: %v", err)
		}
	}()

	emit.Info.Msg("Elasticsearch client connected successfully")

	// Test connection
	ctx := context.Background()
	err = client.Ping(ctx)
	if err != nil {
		emit.Error.StructuredFields("Failed to ping Elasticsearch",
			emit.ZString("error", err.Error()))
		os.Exit(1)
	}

	emit.Info.Msg("Elasticsearch ping successful")

	// Create an index
	indexName := "test-index"
	err = client.Indices().Create(ctx, indexName, nil)
	if err != nil {
		emit.Warn.StructuredFields("Failed to create index (may already exist)",
			emit.ZString("index", indexName),
			emit.ZString("error", err.Error()))
	} else {
		emit.Info.StructuredFields("Index created successfully",
			emit.ZString("index", indexName))
	}

	// Index a document using the new resource-oriented API
	doc := map[string]any{
		"title":   "Basic Example Document",
		"content": "This is a test document for the basic client example",
		"author":  "Go Elastic",
		"tags":    []string{"example", "test", "go-elastic"},
	}

	result, err := client.Documents().Create(ctx, indexName, doc)
	if err != nil {
		emit.Error.StructuredFields("Failed to index document",
			emit.ZString("error", err.Error()))
		os.Exit(1)
	}

	emit.Info.StructuredFields("Document indexed successfully",
		emit.ZString("index", result.Index),
		emit.ZString("document_id", result.ID),
		emit.ZString("result", result.Result))

	// Search for the document using the new resource-oriented API
	query := elastic.MatchQuery("title", "Basic Example")

	searchResult, err := client.Search(indexName).Search(ctx, query, elastic.WithSize(5))
	if err != nil {
		emit.Error.StructuredFields("Failed to search documents",
			emit.ZString("error", err.Error()))
		os.Exit(1)
	}

	emit.Info.StructuredFields("Search completed successfully",
		emit.ZString("index", indexName),
		emit.ZInt("hits", searchResult.Hits.Total.Value),
		emit.ZInt("took", searchResult.Took))

	if len(searchResult.Hits.Hits) > 0 {
		for i, hit := range searchResult.Hits.Hits {
			emit.Info.StructuredFields("Search result",
				emit.ZInt("hit", i+1),
				emit.ZString("document_id", hit.ID),
				emit.ZFloat64("score", hit.Score))
		}
	}

	emit.Info.Msg("Basic client example completed successfully!")
}

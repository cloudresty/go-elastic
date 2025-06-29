package main

import (
	"context"
	"time"

	"github.com/cloudresty/emit"
	"github.com/cloudresty/go-elastic"
)

func main() {
	emit.Info.Msg("Starting go-elastic utilities demonstration")

	demoConnectionUtilities()
	demoQueryBuilders()
	demoAggregationBuilders()
	demoErrorHandling()

	emit.Info.Msg("Utilities demonstration completed successfully!")
}

func demoConnectionUtilities() {
	emit.Info.Msg("\n1. Connection Utilities Demo")
	emit.Info.Msg("----------------------------")

	// Production-ready connection with proper configuration
	client, err := elastic.NewClient()
	if err != nil {
		emit.Warn.StructuredFields("Connection failed (expected if no Elasticsearch running)",
			emit.ZString("error", err.Error()))
	} else {
		defer func() {
			if err := client.Close(); err != nil {
				emit.Warn.StructuredFields("Error closing client", emit.ZString("error", err.Error()))
			}
		}()
		emit.Info.Msg("Production-ready connection established successfully")

		// Test connectivity using the client's ping method
		err = client.Ping(context.Background())
		if err != nil {
			emit.Warn.StructuredFields("Ping failed (expected if no Elasticsearch running)",
				emit.ZString("error", err.Error()))
		} else {
			emit.Info.Msg("Ping successful")
		}
	}
}

func demoQueryBuilders() {
	emit.Info.Msg("\n2. Query Builders Demo")
	emit.Info.Msg("----------------------")

	// Create proper JSON documents for Elasticsearch
	doc := map[string]any{
		"title":   "Sample Document",
		"tags":    []string{"search", "elasticsearch", "go"},
		"created": time.Now(),
		"status":  "published",
	}
	emit.Info.StructuredFields("Created document using JSON structure",
		emit.ZString("title", doc["title"].(string)))

	// Basic queries - these create proper Elasticsearch JSON query structures
	emit.Info.Msg("Creating Elasticsearch queries:")

	_ = elastic.MatchQuery("title", "elasticsearch")
	_ = elastic.TermQuery("status", "published")
	_ = elastic.RangeQuery("date", map[string]any{"gte": "2023-01-01"}, nil)
	_ = elastic.ExistsQuery("author")

	// New enhanced query builders
	_ = elastic.MatchAllQuery()
	_ = elastic.MatchPhraseQuery("title", "sample document")
	_ = elastic.MultiMatchQuery("elasticsearch", "title", "content", "tags")
	_ = elastic.TermsQuery("status", "published", "draft")

	// Complex bool query
	boolQuery := elastic.BoolQuery()
	boolQuery = elastic.WithMust(boolQuery, elastic.MatchQuery("title", "elasticsearch"))
	boolQuery = elastic.WithFilter(boolQuery, elastic.ExistsQuery("author"))

	// Build complete search query
	searchQuery := elastic.BuildSearchQuery(
		boolQuery,
		elastic.WithSize(10),
		elastic.WithFrom(0),
		elastic.WithSort([]map[string]any{elastic.SortDesc("created")}...),
	)

	// Convenience search builders
	_ = elastic.SimpleSearch("elasticsearch tutorial", []string{"title", "content"}, 20)
	_ = elastic.PaginatedSearch(elastic.MatchAllQuery(), 2, 10, "created", false)

	emit.Info.StructuredFields("Built complex search query", emit.ZInt("size", 10))
	_ = searchQuery
}

func demoAggregationBuilders() {
	emit.Info.Msg("\n3. Aggregation Builders Demo")
	emit.Info.Msg("-----------------------------")

	// Aggregations for analytics
	emit.Info.Msg("Creating Elasticsearch aggregations:")

	_ = elastic.TermsAggregation("category", 10)
	_ = elastic.DateHistogramAggregation("created", "day")
	_ = elastic.StatsAggregation("views")
	_ = elastic.AvgAggregation("score")
	_ = elastic.CardinalityAggregation("user_id")
	_ = elastic.FiltersAggregation(map[string]map[string]any{
		"published": elastic.TermQuery("status", "published"),
		"recent":    elastic.RangeQuery("created", map[string]any{"gte": "now-7d"}, nil),
	})

	emit.Info.Msg("Aggregations created for analytics queries")
}

func demoErrorHandling() {
	emit.Info.Msg("\n4. Error Handling Demo")
	emit.Info.Msg("----------------------")

	// Test error handling patterns using existing functions
	emit.Info.Msg("Testing error handling patterns...")

	// These are the error checking utilities available
	_ = elastic.IsNotFoundError
	_ = elastic.IsConflictError
	_ = elastic.IsTimeoutError
	_ = elastic.IsConnectionError
	_ = elastic.IsIndexNotFoundError
	_ = elastic.IsDocumentExistsError
	_ = elastic.IsMappingError
	_ = elastic.IsNetworkError

	emit.Info.Msg("Error handling patterns demonstrated")
}

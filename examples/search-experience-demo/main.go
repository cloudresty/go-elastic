package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cloudresty/go-elastic"
	"github.com/cloudresty/go-elastic/query"
)

// User represents a user document
type User struct {
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Age      int       `json:"age"`
	City     string    `json:"city"`
	JoinDate time.Time `json:"join_date"`
	Active   bool      `json:"active"`
}

// Product represents a product document
type Product struct {
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	InStock     bool    `json:"in_stock"`
	Rating      float64 `json:"rating"`
}

func main() {
	// Connect to Elasticsearch
	client, err := elastic.NewClient(elastic.FromEnv())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Demo the three pillars of our search experience
	fmt.Println("🔍 Demonstrating Best-in-Class Search Experience")
	fmt.Println("================================================")

	// PILLAR 1: Fluent Query Builder
	fmt.Println("\n📝 PILLAR 1: Fluent, Type-Safe Query Builder")
	fmt.Println("--------------------------------------------")

	demoQueryBuilder()

	// PILLAR 2: Composable Search Call
	fmt.Println("\n🔧 PILLAR 2: Composable Search API")
	fmt.Println("----------------------------------")

	demoComposableSearch(client, ctx)

	// PILLAR 3: Rich, Typed Result Set
	fmt.Println("\n🎯 PILLAR 3: Rich, Typed Results")
	fmt.Println("-------------------------------")

	demoTypedResults(client, ctx)

	// BONUS: All Three Pillars Together
	fmt.Println("\n✨ COMPLETE EXPERIENCE: All Pillars Together")
	fmt.Println("===========================================")

	demoCompleteExperience(client, ctx)
}

func demoQueryBuilder() {
	// Simple queries
	fmt.Println("Simple queries:")

	termQuery := query.Term("status", "active")
	fmt.Printf("Term query: %s\n", termQuery.String())

	matchQuery := query.Match("title", "elasticsearch guide")
	fmt.Printf("Match query: %s\n", matchQuery.String())

	// Range queries with fluent builder
	priceRangeQuery := query.Range("price").Gte(10.0).Lte(100.0).Build()
	fmt.Printf("Range query: %s\n", priceRangeQuery.String())

	dateRangeQuery := query.Range("created_at").
		Gte("2023-01-01").
		Format("yyyy-MM-dd").
		TimeZone("UTC").
		Build()
	fmt.Printf("Date range query: %s\n", dateRangeQuery.String())

	// Complex bool queries
	fmt.Println("\nComplex bool queries:")

	complexQuery := query.New().
		Must(
			query.Match("description", "high quality"),
			query.Range("rating").Gte(4.0).Build(),
		).
		Filter(
			query.Term("in_stock", true),
			query.Range("price").Lte(500.0).Build(),
		).
		Should(
			query.Term("category", "electronics"),
			query.Term("category", "books"),
		).
		MinimumShouldMatch(1)

	fmt.Printf("Complex query: %s\n", complexQuery.String())

	// Multi-field search
	multiMatch := query.MultiMatch("programming golang", "title", "description", "tags")
	fmt.Printf("Multi-match query: %s\n", multiMatch.String())
}

func demoComposableSearch(client *elastic.Client, ctx context.Context) {
	docs := client.Documents()

	// The ONLY way to search - fluent typed API!
	fmt.Println("Unified typed search API:")

	builderQuery := query.New().
		Must(query.Match("name", "John")).
		Filter(query.Range("age").Gte(18).Build())

	// Create typed documents interface and search
	typedDocs := elastic.For[User](docs)
	response, err := typedDocs.Search(ctx, builderQuery,
		elastic.WithIndices("users"),
		elastic.WithSize(10),
		elastic.WithSort(map[string]any{"age": "desc"}),
	)

	if err != nil {
		fmt.Printf("Search error: %v\n", err)
	} else {
		fmt.Printf("Found %d adult users named John\n", response.TotalHits())
		if response.HasHits() {
			users := response.Documents()
			fmt.Printf("Users: %v\n", len(users))
		}
	}
}

func demoTypedResults(client *elastic.Client, ctx context.Context) {
	docs := client.Documents()

	// Unified typed search - the ONLY way!
	fmt.Println("Unified typed search:")

	// Build query using the fluent builder
	searchQuery := query.New().
		Must(query.Match("category", "electronics"))

	typedDocs := elastic.For[Product](docs)
	typedResult, err := typedDocs.Search(ctx, searchQuery, elastic.WithIndices("products"))
	if err != nil {
		fmt.Printf("Search error: %v\n", err)
		return
	}

	fmt.Printf("Typed result has %d products\n", typedResult.TotalHits())

	// Rich result methods
	if typedResult.HasHits() {
		fmt.Printf("Max score: %.2f\n", *typedResult.MaxScore())

		first, hasFirst := typedResult.First()
		if hasFirst {
			fmt.Printf("First product: %s ($%.2f)\n", first.Name, first.Price)
		}

		// All documents as slice
		products := typedResult.Documents()
		fmt.Printf("All products: %d items\n", len(products))

		// Documents with IDs
		withIDs := typedResult.DocumentsWithIDs()
		if len(withIDs) > 0 {
			fmt.Printf("First product with ID: %s - %s\n", withIDs[0].ID, withIDs[0].Document.Name)
		}

		// Functional operations
		highRated := typedResult.Filter(func(p Product) bool {
			return p.Rating >= 4.0
		})
		fmt.Printf("High-rated products: %d\n", len(highRated))

		names := typedResult.Map(func(p Product) Product {
			p.Name = "Product: " + p.Name
			return p
		})
		fmt.Printf("Mapped products: %d\n", len(names))
	}
}

func demoCompleteExperience(client *elastic.Client, ctx context.Context) {
	docs := client.Documents()

	// Build a sophisticated query using the fluent builder
	searchQuery := query.New().
		Must(
			query.MultiMatch("high quality laptop", "name", "description"),
			query.Range("rating").Gte(4.0).Build(),
		).
		Filter(
			query.Term("in_stock", true),
			query.Range("price").Gte(500.0).Lte(2000.0).Build(),
		).
		Should(
			query.Term("category", "electronics"),
			query.Match("description", "gaming"),
		).
		MinimumShouldMatch(1)

	fmt.Printf("Complex search query:\n%s\n", searchQuery.String())

	// Execute the search with rich options and get typed results
	typedDocs := elastic.For[Product](docs)
	result, err := typedDocs.Search(
		ctx,
		searchQuery,
		elastic.WithIndices("products"),
		elastic.WithSize(20),
		elastic.WithSort(map[string]any{"rating": "desc"}),
		elastic.WithSort(map[string]any{"price": "asc"}),
		elastic.WithSource("name", "price", "rating", "category"),
		elastic.WithAggregation("by_category", elastic.NewTermsAggregation("category.keyword").Size(10)),
		elastic.WithAggregation("price_ranges", elastic.NewRangeAggregation("price").
			AddRange("budget", nil, &[]float64{500}[0]).
			AddRange("mid_range", &[]float64{500}[0], &[]float64{1500}[0]).
			AddRange("premium", &[]float64{1500}[0], nil)),
		elastic.WithAggregation("avg_rating", elastic.NewAvgAggregation("rating")),
	)

	if err != nil {
		fmt.Printf("Search failed: %v\n", err)
		return
	}

	// Rich result processing
	fmt.Printf("\n🎉 Search Results Summary:\n")
	fmt.Printf("Total products found: %d\n", result.TotalHits())
	fmt.Printf("Results returned: %d\n", len(result.Documents()))
	fmt.Printf("Search took: %d ms\n", result.Took)

	if result.HasHits() {
		fmt.Printf("Max relevance score: %.3f\n", *result.MaxScore())

		// Show top results
		fmt.Printf("\n📊 Top Products:\n")
		result.Each(func(hit elastic.TypedHit[Product]) {
			fmt.Printf("- %s: $%.2f (⭐ %.1f) [Score: %.3f]\n",
				hit.Source.Name,
				hit.Source.Price,
				hit.Source.Rating,
				*hit.Score,
			)
		})

		// Statistical analysis
		products := result.Documents()
		var totalPrice float64
		var totalRating float64

		for _, product := range products {
			totalPrice += product.Price
			totalRating += product.Rating
		}

		fmt.Printf("\n📈 Statistics:\n")
		fmt.Printf("Average price: $%.2f\n", totalPrice/float64(len(products)))
		fmt.Printf("Average rating: %.2f\n", totalRating/float64(len(products)))

		// Premium products
		premium := result.Filter(func(p Product) bool {
			return p.Price > 1000.0 && p.Rating >= 4.5
		})
		fmt.Printf("Premium products (>$1000, ⭐4.5+): %d\n", len(premium))

		// Show aggregation results
		if result.Aggregations != nil {
			fmt.Printf("\n📊 Aggregation Results:\n")

			// Category breakdown
			if categoryAgg, ok := result.Aggregations["by_category"].(map[string]any); ok {
				if buckets, ok := categoryAgg["buckets"].([]any); ok {
					fmt.Printf("Categories:\n")
					for _, bucket := range buckets {
						if b, ok := bucket.(map[string]any); ok {
							fmt.Printf("  - %s: %v products\n", b["key"], b["doc_count"])
						}
					}
				}
			}

			// Average rating
			if avgRatingAgg, ok := result.Aggregations["avg_rating"].(map[string]any); ok {
				if avgValue, ok := avgRatingAgg["value"].(float64); ok {
					fmt.Printf("Average rating: ⭐ %.2f\n", avgValue)
				}
			}
		}
	}

	fmt.Printf("\n✅ Complete search experience demonstrated!\n")
	fmt.Printf("   🔧 Fluent query builder for complex queries\n")
	fmt.Printf("   🔍 Composable search API with rich options\n")
	fmt.Printf("   🎯 Typed results with functional operations\n")
	fmt.Printf("   📊 Aggregations for analytics and insights\n")
}

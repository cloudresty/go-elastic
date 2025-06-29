package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudresty/go-elastic"
	"github.com/cloudresty/go-elastic/query"
)

// User represents a user document
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	City string `json:"city"`
}

func main() {
	// Connect to Elasticsearch
	client, err := elastic.NewClient(elastic.FromEnv())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	fmt.Println("🎯 Demonstrating the NEW Unified Search API")
	fmt.Println("==========================================")

	// Build a query using the fluent builder
	searchQuery := query.New().
		Must(query.Match("name", "John")).
		Filter(query.Range("age").Gte(18).Build())

	// THE ONLY way to search - clean, fluent, and typed!
	typedDocs := elastic.For[User](client.Documents())
	results, err := typedDocs.Search(
		ctx,
		searchQuery,
		elastic.WithIndices("users"),
		elastic.WithSize(10),
		elastic.WithSort(map[string]any{"age": "desc"}),
	)

	if err != nil {
		fmt.Printf("Search error: %v\n", err)
		return
	}

	fmt.Printf("✅ Found %d adult users named John\n", results.TotalHits())

	if results.HasHits() {
		users := results.Documents()
		fmt.Printf("📊 Retrieved %d users\n", len(users))

		if first, hasFirst := results.First(); hasFirst {
			fmt.Printf("👤 First user: %s, %d years old from %s\n",
				first.Name, first.Age, first.City)
		}
	}

	fmt.Println("\n🚀 API Benefits:")
	fmt.Println("   ✅ Zero choice paralysis - ONE way to search")
	fmt.Println("   ✅ Fluent and intuitive method-style API")
	fmt.Println("   ✅ Type safety with Go generics")
	fmt.Println("   ✅ Forces best practices (builder required)")
	fmt.Println("   ✅ Rich, typed results with helper methods")
}

package main

import (
	"context"
	"log"

	"github.com/cloudresty/go-elastic"
	"github.com/cloudresty/go-elastic/query"
)

func main() {
	// Create client with default settings
	client, err := elastic.NewClient()
	if err != nil {
		log.Fatalf("Error creating client: %s", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("Error closing client: %s", err)
		}
	}()

	ctx := context.Background()

	// --- Index Operations ---
	// Access the Indices service and call its Create method
	log.Println("Creating index...")
	err = client.Indices().Create(ctx, "my-new-index", nil)
	if err != nil {
		log.Printf("Error creating index (may already exist): %s", err)
	}

	// Check if index exists
	exists, err := client.Indices().Exists(ctx, "my-new-index")
	if err != nil {
		log.Fatalf("Error checking index existence: %s", err)
	}
	log.Printf("Index exists: %v", exists)

	// --- Document Operations ---
	// Access the Documents service and call its Create method
	log.Println("Creating document...")
	myDoc := map[string]any{
		"title":   "A New Document",
		"author":  "John Doe",
		"content": "This is the content of the document",
	}

	response, err := client.Documents().Create(ctx, "my-new-index", myDoc)
	if err != nil {
		log.Fatalf("Error creating document: %s", err)
	}
	log.Printf("Document created with ID: %s", response.ID)

	// Get the document we just created
	log.Println("Retrieving document...")
	doc, err := client.Documents().Get(ctx, "my-new-index", response.ID)
	if err != nil {
		log.Fatalf("Error getting document: %s", err)
	}
	log.Printf("Retrieved document: %+v", doc)

	// --- Search Operations ---
	// Search for documents
	log.Println("Searching documents...")
	searchQuery := query.Match("title", "Document")

	// Type to represent our document
	type MyDocument struct {
		Title string `json:"title"`
	}

	typedDocs := elastic.For[MyDocument](client.Documents())
	searchResponse, err := typedDocs.Search(ctx, searchQuery, elastic.WithIndices("my-new-index"))
	if err != nil {
		log.Fatalf("Error searching documents: %s", err)
	}
	log.Printf("Search found %d documents", searchResponse.TotalHits())

	// --- Cluster Operations ---
	// Get cluster health
	log.Println("Getting cluster health...")
	health, err := client.Cluster().Health(ctx)
	if err != nil {
		log.Fatalf("Error getting cluster health: %s", err)
	}
	log.Printf("Cluster status: %s", health.Status)

	// --- Bulk Operations ---
	// Perform bulk operations using the fluent builder
	log.Println("Performing bulk operations...")
	bulkResponse, err := client.Documents().Bulk("my-new-index").
		CreateWithID("doc1", map[string]any{
			"title":   "Bulk Document 1",
			"content": "Content 1",
		}).
		CreateWithID("doc2", map[string]any{
			"title":   "Bulk Document 2",
			"content": "Content 2",
		}).
		Do(ctx)
	if err != nil {
		log.Fatalf("Error executing bulk operations: %s", err)
	}
	log.Printf("Bulk operations completed, has errors: %v", bulkResponse.Errors)

	log.Println("Resource-oriented API demo completed successfully!")
}

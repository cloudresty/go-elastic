package elastic

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudresty/go-elastic/query"
)

// DocumentsService search methods

// TypedDocuments provides a typed interface to document operations for a specific type T
// This enables fluent method-style API calls for typed operations
type TypedDocuments[T any] struct {
	service *DocumentsService
}

// For returns a typed documents interface for method-style calls with a specific type
// Usage: typedDocs := elastic.For[User](client.Documents())
//
//	result, err := typedDocs.Search(ctx, queryBuilder, options...)
func For[T any](service *DocumentsService) *TypedDocuments[T] {
	return &TypedDocuments[T]{service: service}
}

// Search performs a typed search using a query builder and returns rich, typed results
// This is THE unified search method that requires the query builder
func (t *TypedDocuments[T]) Search(ctx context.Context, queryBuilder *query.Builder, options ...SearchOption) (*SearchResult[T], error) {
	searchResource := &SearchResource{
		client: t.service.client,
	}

	// Execute the search with the builder's query
	response, err := searchResource.Search(ctx, queryBuilder.Build(), options...)
	if err != nil {
		return nil, err
	}

	// Convert to typed result
	return ConvertSearchResponse[T](response)
}

// Scroll creates a new typed search iterator for paginated results using the scroll API
func (t *TypedDocuments[T]) Scroll(ctx context.Context, queryBuilder *query.Builder, scrollTime time.Duration, options ...SearchOption) (*TypedSearchIterator[T], error) {
	searchResource := &SearchResource{
		client: t.service.client,
	}

	// Start the initial scroll search
	initialResponse, err := searchResource.startScrollSearch(ctx, queryBuilder.Build(), scrollTime, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to start scroll search: %w", err)
	}

	// Create and return the typed iterator
	iterator := &TypedSearchIterator[T]{
		client:       t.service.client,
		scrollID:     initialResponse.ScrollID,
		scrollTime:   scrollTime,
		currentIndex: -1, // Start before first element
		totalHits:    int64(initialResponse.Hits.Total.Value),
	}

	// Convert initial hits to typed hits
	typedResult, err := ConvertSearchResponse[T](initialResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to convert initial scroll response: %w", err)
	}
	iterator.currentHits = typedResult.Hits.Hits

	return iterator, nil
}

// Count returns the count of documents matching a query builder
func (s *DocumentsService) Count(ctx context.Context, queryBuilder *query.Builder, options ...SearchOption) (int64, error) {
	searchResource := &SearchResource{
		client: s.client,
	}
	return searchResource.Count(ctx, queryBuilder.Build(), options...)
}

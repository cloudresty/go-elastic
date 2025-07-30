package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// SearchResult represents a rich, typed search result with generic document support
type SearchResult[T any] struct {
	Took         int            `json:"took"`
	TimedOut     bool           `json:"timed_out"`
	ScrollID     string         `json:"_scroll_id,omitempty"`
	Shards       SearchShards   `json:"_shards"`
	Hits         TypedHits[T]   `json:"hits"`
	Aggregations map[string]any `json:"aggregations,omitempty"`
	Suggest      map[string]any `json:"suggest,omitempty"`
}

// TypedHits represents the hits section with typed documents
type TypedHits[T any] struct {
	Total    SearchTotal   `json:"total"`
	MaxScore *float64      `json:"max_score"`
	Hits     []TypedHit[T] `json:"hits"`
}

// TypedHit represents a single search hit with typed source
type TypedHit[T any] struct {
	Index       string              `json:"_index"`
	Type        string              `json:"_type,omitempty"`
	ID          string              `json:"_id"`
	Score       *float64            `json:"_score"`
	Source      T                   `json:"_source"`
	Sort        []any               `json:"sort,omitempty"`
	Fields      map[string]any      `json:"fields,omitempty"`
	Highlight   map[string][]string `json:"highlight,omitempty"`
	InnerHits   map[string]any      `json:"inner_hits,omitempty"`
	Explanation map[string]any      `json:"_explanation,omitempty"`
}

// SearchShards represents shard information from a search response
type SearchShards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

// SearchTotal represents the total hits information
type SearchTotal struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

// Documents returns a slice of the typed documents from the search result
func (sr *SearchResult[T]) Documents() []T {
	docs := make([]T, len(sr.Hits.Hits))
	for i, hit := range sr.Hits.Hits {
		docs[i] = hit.Source
	}
	return docs
}

// DocumentIDs returns a slice of document IDs from the search result
func (sr *SearchResult[T]) DocumentIDs() []string {
	ids := make([]string, len(sr.Hits.Hits))
	for i, hit := range sr.Hits.Hits {
		ids[i] = hit.ID
	}
	return ids
}

// DocumentsWithIDs returns a slice of DocumentWithID containing both document and ID
func (sr *SearchResult[T]) DocumentsWithIDs() []DocumentWithID[T] {
	docs := make([]DocumentWithID[T], len(sr.Hits.Hits))
	for i, hit := range sr.Hits.Hits {
		docs[i] = DocumentWithID[T]{
			ID:       hit.ID,
			Document: hit.Source,
		}
	}
	return docs
}

// DocumentWithID combines a document with its Elasticsearch ID
type DocumentWithID[T any] struct {
	ID       string `json:"id"`
	Document T      `json:"document"`
}

// TotalHits returns the total number of hits
func (sr *SearchResult[T]) TotalHits() int {
	return sr.Hits.Total.Value
}

// HasHits returns true if there are any hits
func (sr *SearchResult[T]) HasHits() bool {
	return len(sr.Hits.Hits) > 0
}

// MaxScore returns the maximum score from the search results
func (sr *SearchResult[T]) MaxScore() *float64 {
	return sr.Hits.MaxScore
}

// Each calls the provided function for each hit in the search result
func (sr *SearchResult[T]) Each(fn func(hit TypedHit[T])) {
	for _, hit := range sr.Hits.Hits {
		fn(hit)
	}
}

// Map transforms each document using the provided function
func (sr *SearchResult[T]) Map(fn func(T) T) []T {
	mapped := make([]T, len(sr.Hits.Hits))
	for i, hit := range sr.Hits.Hits {
		mapped[i] = fn(hit.Source)
	}
	return mapped
}

// Filter returns documents that match the provided predicate
func (sr *SearchResult[T]) Filter(fn func(T) bool) []T {
	var filtered []T
	for _, hit := range sr.Hits.Hits {
		if fn(hit.Source) {
			filtered = append(filtered, hit.Source)
		}
	}
	return filtered
}

// First returns the first document if available
func (sr *SearchResult[T]) First() (T, bool) {
	var zero T
	if len(sr.Hits.Hits) == 0 {
		return zero, false
	}
	return sr.Hits.Hits[0].Source, true
}

// Last returns the last document if available
func (sr *SearchResult[T]) Last() (T, bool) {
	var zero T
	if len(sr.Hits.Hits) == 0 {
		return zero, false
	}
	return sr.Hits.Hits[len(sr.Hits.Hits)-1].Source, true
}

// ConvertSearchResponse converts a generic SearchResponse to a typed SearchResult[T]
func ConvertSearchResponse[T any](response *SearchResponse) (*SearchResult[T], error) {
	typedResult := &SearchResult[T]{
		Took:     response.Took,
		TimedOut: response.TimedOut,
		ScrollID: response.ScrollID,
		Shards: SearchShards{
			Total:      response.Shards.Total,
			Successful: response.Shards.Successful,
			Skipped:    response.Shards.Skipped,
			Failed:     response.Shards.Failed,
		},
		Hits: TypedHits[T]{
			Total: SearchTotal{
				Value:    response.Hits.Total.Value,
				Relation: response.Hits.Total.Relation,
			},
			MaxScore: &response.Hits.MaxScore,
			Hits:     make([]TypedHit[T], len(response.Hits.Hits)),
		},
		Aggregations: response.Aggregations,
	}

	// Convert hits to typed hits
	for i, hit := range response.Hits.Hits {
		var doc T
		if hit.Source != nil {
			// Parse the source into the typed document
			sourceBytes, err := json.Marshal(hit.Source)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal hit source: %w", err)
			}

			if err := json.Unmarshal(sourceBytes, &doc); err != nil {
				return nil, fmt.Errorf("failed to unmarshal hit source to type %T: %w", doc, err)
			}
		}

		typedResult.Hits.Hits[i] = TypedHit[T]{
			Index:  hit.Index,
			Type:   hit.Type,
			ID:     hit.ID,
			Score:  &hit.Score,
			Source: doc,
		}
	}

	return typedResult, nil
}

// TypedSearchIterator provides a typed iterator pattern for scrolling through large result sets
type TypedSearchIterator[T any] struct {
	client        *Client
	scrollID      string
	scrollTime    time.Duration
	currentHits   []TypedHit[T]
	currentIndex  int
	done          bool
	err           error
	totalHits     int64
	processedHits int64
}

// Next advances the iterator to the next document
// Returns true if there is a next document, false when iteration is complete
func (tsi *TypedSearchIterator[T]) Next(ctx context.Context) bool {
	if tsi.err != nil || tsi.done {
		return false
	}

	// If we have more hits in the current batch, advance to next
	if tsi.currentIndex < len(tsi.currentHits)-1 {
		tsi.currentIndex++
		tsi.processedHits++
		return true
	}

	// If no scroll ID, we're done
	if tsi.scrollID == "" {
		tsi.done = true
		return false
	}

	// Need to fetch next batch
	if err := tsi.fetchNextBatch(ctx); err != nil {
		tsi.err = err
		return false
	}

	// Check if we got new hits
	if len(tsi.currentHits) == 0 {
		tsi.done = true
		return false
	}

	// Reset to first hit of new batch
	tsi.currentIndex = 0
	tsi.processedHits++
	return true
}

// Scan unmarshals the current document into the destination
func (tsi *TypedSearchIterator[T]) Scan(dest *T) error {
	if tsi.currentIndex < 0 || tsi.currentIndex >= len(tsi.currentHits) {
		return fmt.Errorf("no current document - call Next() first")
	}

	*dest = tsi.currentHits[tsi.currentIndex].Source
	return nil
}

// Current returns the current document
func (tsi *TypedSearchIterator[T]) Current() T {
	if tsi.currentIndex < 0 || tsi.currentIndex >= len(tsi.currentHits) {
		var zero T
		return zero
	}
	return tsi.currentHits[tsi.currentIndex].Source
}

// CurrentHit returns the current hit with metadata
func (tsi *TypedSearchIterator[T]) CurrentHit() TypedHit[T] {
	if tsi.currentIndex < 0 || tsi.currentIndex >= len(tsi.currentHits) {
		return TypedHit[T]{}
	}
	return tsi.currentHits[tsi.currentIndex]
}

// Err returns any error that occurred during iteration
func (tsi *TypedSearchIterator[T]) Err() error {
	return tsi.err
}

// TotalHits returns the total number of hits found by the search
func (tsi *TypedSearchIterator[T]) TotalHits() int64 {
	return tsi.totalHits
}

// ProcessedHits returns the number of hits processed so far
func (tsi *TypedSearchIterator[T]) ProcessedHits() int64 {
	return tsi.processedHits
}

// Close cleans up the scroll context
func (tsi *TypedSearchIterator[T]) Close(ctx context.Context) error {
	if tsi.scrollID == "" {
		return nil
	}

	searchScroll := &SearchScroll{
		client: tsi.client,
	}

	if err := searchScroll.Clear(ctx, tsi.scrollID); err != nil {
		tsi.client.config.Logger.Warn("Failed to clear scroll context - scroll_id: %s, error: %s", tsi.scrollID, err.Error())
		return err
	}

	tsi.scrollID = ""
	return nil
}

// fetchNextBatch retrieves the next batch of results using the scroll API
func (tsi *TypedSearchIterator[T]) fetchNextBatch(ctx context.Context) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	// Use the SearchScroll to get next batch
	searchScroll := &SearchScroll{
		client: tsi.client,
	}

	response, err := searchScroll.Continue(ctx, tsi.scrollID, tsi.scrollTime)
	if err != nil {
		return fmt.Errorf("failed to continue scroll: %w", err)
	}

	// Update scroll ID for next iteration
	tsi.scrollID = response.ScrollID

	// Convert response to typed hits
	typedResult, err := ConvertSearchResponse[T](response)
	if err != nil {
		return fmt.Errorf("failed to convert scroll response: %w", err)
	}

	// Update current batch
	tsi.currentHits = typedResult.Hits.Hits
	tsi.currentIndex = -1 // Will be incremented to 0 by Next()

	tsi.client.config.Logger.Debug("Fetched next typed scroll batch - scroll_id: %s, batch_size: %d, processed_total: %d", tsi.scrollID, len(tsi.currentHits), tsi.processedHits)

	return nil
}

package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cloudresty/emit"
)

// SearchIterator provides an iterator pattern for scrolling through large result sets
type SearchIterator struct {
	client        *Client
	scrollID      string
	scrollTime    time.Duration
	currentHits   []Hit
	currentIndex  int
	done          bool
	err           error
	totalHits     int64
	processedHits int64
}

// Next advances the iterator to the next document
// Returns true if there is a next document, false when iteration is complete
func (si *SearchIterator) Next(ctx context.Context) bool {
	if si.err != nil || si.done {
		return false
	}

	// If we have more hits in the current batch, advance to next
	if si.currentIndex < len(si.currentHits)-1 {
		si.currentIndex++
		si.processedHits++
		return true
	}

	// If no scroll ID, we're done
	if si.scrollID == "" {
		si.done = true
		return false
	}

	// Need to fetch next batch
	if err := si.fetchNextBatch(ctx); err != nil {
		si.err = err
		return false
	}

	// Check if we got any hits in the new batch
	if len(si.currentHits) == 0 {
		si.done = true
		si.clearScroll(ctx) // Clean up scroll context
		return false
	}

	// Reset to first hit in new batch
	si.currentIndex = 0
	si.processedHits++
	return true
}

// Scan unmarshals the current document into the provided destination
func (si *SearchIterator) Scan(dest any) error {
	if si.err != nil {
		return si.err
	}

	if si.currentIndex >= len(si.currentHits) {
		return fmt.Errorf("no current document to scan")
	}

	currentHit := si.currentHits[si.currentIndex]

	// Marshal the source back to JSON, then unmarshal to dest
	sourceBytes, err := json.Marshal(currentHit.Source)
	if err != nil {
		return fmt.Errorf("failed to marshal document source: %w", err)
	}

	if err := json.Unmarshal(sourceBytes, dest); err != nil {
		return fmt.Errorf("failed to unmarshal document into destination: %w", err)
	}

	return nil
}

// Current returns the raw current document as a map
func (si *SearchIterator) Current() map[string]any {
	if si.currentIndex >= len(si.currentHits) {
		return nil
	}
	return si.currentHits[si.currentIndex].Source
}

// CurrentHit returns the full current Hit with metadata
func (si *SearchIterator) CurrentHit() *Hit {
	if si.currentIndex >= len(si.currentHits) {
		return nil
	}
	return &si.currentHits[si.currentIndex]
}

// Err returns any error that occurred during iteration
func (si *SearchIterator) Err() error {
	return si.err
}

// TotalHits returns the total number of hits found by the search
func (si *SearchIterator) TotalHits() int64 {
	return si.totalHits
}

// ProcessedHits returns the number of hits processed so far
func (si *SearchIterator) ProcessedHits() int64 {
	return si.processedHits
}

// Close cleans up the scroll context (called automatically when iteration completes)
func (si *SearchIterator) Close(ctx context.Context) error {
	if si.scrollID != "" {
		return si.clearScroll(ctx)
	}
	return nil
}

// fetchNextBatch retrieves the next batch of results using the scroll API
func (si *SearchIterator) fetchNextBatch(ctx context.Context) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	// Use the SearchScroll to get next batch
	searchScroll := &SearchScroll{
		client: si.client,
	}

	response, err := searchScroll.Continue(ctx, si.scrollID, si.scrollTime)
	if err != nil {
		return fmt.Errorf("failed to continue scroll: %w", err)
	}

	// Update scroll ID for next iteration
	si.scrollID = response.ScrollID

	// Update current hits
	si.currentHits = response.Hits.Hits

	emit.Debug.StructuredFields("Fetched next scroll batch",
		emit.ZString("scroll_id", si.scrollID),
		emit.ZInt("batch_size", len(si.currentHits)),
		emit.ZInt64("processed_total", si.processedHits))

	return nil
}

// clearScroll cleans up the scroll context on Elasticsearch
func (si *SearchIterator) clearScroll(ctx context.Context) error {
	if si.scrollID == "" {
		return nil
	}

	searchScroll := &SearchScroll{
		client: si.client,
	}

	if err := searchScroll.Clear(ctx, si.scrollID); err != nil {
		emit.Warn.StructuredFields("Failed to clear scroll context",
			emit.ZString("scroll_id", si.scrollID),
			emit.ZString("error", err.Error()))
		return err
	}

	si.scrollID = ""
	return nil
}

package elastic

import (
	"time"

	"github.com/cloudresty/ulid"
)

// ULID utility functions

// GenerateULID generates a new ULID string
// This is useful when you want to generate ULIDs outside of document operations
func GenerateULID() string {
	return generateULID()
}

func generateULID() string {
	id, _ := ulid.New()
	return id
}

// GenerateULIDFromTime generates a ULID with a specific timestamp
// This is useful for testing or when you need deterministic time-based IDs
func GenerateULIDFromTime(t time.Time) string {
	id, _ := ulid.NewTime(uint64(t.UnixMilli()))
	return id
}

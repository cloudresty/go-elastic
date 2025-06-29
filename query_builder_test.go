package elastic

import (
	"encoding/json"
	"testing"

	"github.com/cloudresty/go-elastic/query"
)

func TestQueryBuilder(t *testing.T) {
	// Test basic bool query builder functionality
	q := query.New().
		Must(query.Match("name", "John")).
		Must(query.Term("active", true)).
		Build()

	// Verify the query structure
	if q == nil {
		t.Fatal("Query should not be nil")
	}

	// Convert to JSON to verify structure
	jsonBytes, err := json.Marshal(q)
	if err != nil {
		t.Fatalf("Failed to marshal query: %v", err)
	}

	jsonStr := string(jsonBytes)
	t.Logf("Generated query: %s", jsonStr)

	// Basic verification that it contains expected elements
	if jsonStr == "" {
		t.Fatal("Query JSON should not be empty")
	}
}

func TestTermQuery(t *testing.T) {
	q := query.Term("status", "active")
	result := q.Build()

	if result == nil {
		t.Fatal("Term query result should not be nil")
	}

	// Check if it has the term structure
	term, ok := result["term"]
	if !ok {
		t.Fatal("Query should have 'term' field")
	}

	termMap, ok := term.(map[string]any)
	if !ok {
		t.Fatal("Term should be a map")
	}

	if termMap["status"] != "active" {
		t.Fatalf("Expected status=active, got %v", termMap["status"])
	}
}

func TestMatchQuery(t *testing.T) {
	q := query.Match("title", "elasticsearch")
	result := q.Build()

	if result == nil {
		t.Fatal("Match query result should not be nil")
	}

	// Check if it has the match structure
	match, ok := result["match"]
	if !ok {
		t.Fatal("Query should have 'match' field")
	}

	matchMap, ok := match.(map[string]any)
	if !ok {
		t.Fatal("Match should be a map")
	}

	if matchMap["title"] != "elasticsearch" {
		t.Fatalf("Expected title=elasticsearch, got %v", matchMap["title"])
	}
}

func TestRangeQuery(t *testing.T) {
	q := query.Range("age").Gte(18).Lte(65).Build()
	result := q.Build()

	if result == nil {
		t.Fatal("Range query result should not be nil")
	}

	// Check if it has the range structure
	rangeQuery, ok := result["range"]
	if !ok {
		t.Fatal("Query should have 'range' field")
	}

	rangeMap, ok := rangeQuery.(map[string]any)
	if !ok {
		t.Fatal("Range should be a map")
	}

	ageRange, ok := rangeMap["age"].(map[string]any)
	if !ok {
		t.Fatal("Age range should be a map")
	}

	if ageRange["gte"] != 18 {
		t.Fatalf("Expected gte=18, got %v", ageRange["gte"])
	}

	if ageRange["lte"] != 65 {
		t.Fatalf("Expected lte=65, got %v", ageRange["lte"])
	}
}

func TestPanicSafeBoolMethods(t *testing.T) {
	// Test that calling Must() on a Term query panics with a helpful message
	defer func() {
		if r := recover(); r != nil {
			panicMsg := r.(string)
			expected := "query: cannot call Must() on a non-bool query builder (e.g., a Term, Match, or Range query)"
			if panicMsg != expected {
				t.Errorf("Expected panic message %q, got %q", expected, panicMsg)
			}
		} else {
			t.Error("Expected Must() on Term query to panic, but it didn't")
		}
	}()

	// This should panic because Term() creates a non-bool query
	query.Term("status", "active").Must(query.Match("name", "test"))
}

func TestPanicSafeFilterMethod(t *testing.T) {
	// Test that calling Filter() on a Match query panics with a helpful message
	defer func() {
		if r := recover(); r != nil {
			panicMsg := r.(string)
			expected := "query: cannot call Filter() on a non-bool query builder (e.g., a Term, Match, or Range query)"
			if panicMsg != expected {
				t.Errorf("Expected panic message %q, got %q", expected, panicMsg)
			}
		} else {
			t.Error("Expected Filter() on Match query to panic, but it didn't")
		}
	}()

	// This should panic because Match() creates a non-bool query
	query.Match("title", "test").Filter(query.Term("active", true))
}

func TestPanicSafeShouldMethod(t *testing.T) {
	// Test that calling Should() on a Range query panics with a helpful message
	defer func() {
		if r := recover(); r != nil {
			panicMsg := r.(string)
			expected := "query: cannot call Should() on a non-bool query builder (e.g., a Term, Match, or Range query)"
			if panicMsg != expected {
				t.Errorf("Expected panic message %q, got %q", expected, panicMsg)
			}
		} else {
			t.Error("Expected Should() on Range query to panic, but it didn't")
		}
	}()

	// This should panic because Range().Build() creates a non-bool query
	query.Range("age").Gte(18).Build().Should(query.Term("active", true))
}

func TestPanicSafeMinimumShouldMatch(t *testing.T) {
	// Test that calling MinimumShouldMatch() on a Term query panics with a helpful message
	defer func() {
		if r := recover(); r != nil {
			panicMsg := r.(string)
			expected := "query: cannot call MinimumShouldMatch() on a non-bool query builder (e.g., a Term, Match, or Range query)"
			if panicMsg != expected {
				t.Errorf("Expected panic message %q, got %q", expected, panicMsg)
			}
		} else {
			t.Error("Expected MinimumShouldMatch() on Term query to panic, but it didn't")
		}
	}()

	// This should panic because Term() creates a non-bool query
	query.Term("status", "active").MinimumShouldMatch(1)
}

// Package query provides a fluent, type-safe way to build Elasticsearch queries
package query

import (
	"encoding/json"
)

// Builder represents a query builder that constructs Elasticsearch queries
type Builder struct {
	query map[string]any
}

// New creates a new query builder with a bool query
func New() *Builder {
	return &Builder{
		query: map[string]any{
			"bool": map[string]any{
				"must":     []any{},
				"must_not": []any{},
				"should":   []any{},
				"filter":   []any{},
			},
		},
	}
}

// Must adds one or more queries to the must clause
func (b *Builder) Must(queries ...*Builder) *Builder {
	// Use the safe "comma-ok" type assertion
	boolQuery, ok := b.query["bool"].(map[string]any)
	if !ok {
		// If the builder is not a bool query, panic with a helpful message
		panic("query: cannot call Must() on a non-bool query builder (e.g., a Term, Match, or Range query)")
	}

	must, _ := boolQuery["must"].([]any) // We can be sure this exists from New()

	for _, q := range queries {
		must = append(must, q.Build())
	}

	boolQuery["must"] = must
	return b
}

// Filter adds one or more queries to the filter clause
func (b *Builder) Filter(queries ...*Builder) *Builder {
	// Use the safe "comma-ok" type assertion
	boolQuery, ok := b.query["bool"].(map[string]any)
	if !ok {
		// If the builder is not a bool query, panic with a helpful message
		panic("query: cannot call Filter() on a non-bool query builder (e.g., a Term, Match, or Range query)")
	}

	filter, _ := boolQuery["filter"].([]any) // We can be sure this exists from New()

	for _, q := range queries {
		filter = append(filter, q.Build())
	}

	boolQuery["filter"] = filter
	return b
}

// Should adds one or more queries to the should clause
func (b *Builder) Should(queries ...*Builder) *Builder {
	// Use the safe "comma-ok" type assertion
	boolQuery, ok := b.query["bool"].(map[string]any)
	if !ok {
		// If the builder is not a bool query, panic with a helpful message
		panic("query: cannot call Should() on a non-bool query builder (e.g., a Term, Match, or Range query)")
	}

	should, _ := boolQuery["should"].([]any) // We can be sure this exists from New()

	for _, q := range queries {
		should = append(should, q.Build())
	}

	boolQuery["should"] = should
	return b
}

// MustNot adds one or more queries to the must_not clause
func (b *Builder) MustNot(queries ...*Builder) *Builder {
	// Use the safe "comma-ok" type assertion
	boolQuery, ok := b.query["bool"].(map[string]any)
	if !ok {
		// If the builder is not a bool query, panic with a helpful message
		panic("query: cannot call MustNot() on a non-bool query builder (e.g., a Term, Match, or Range query)")
	}

	mustNot, _ := boolQuery["must_not"].([]any) // We can be sure this exists from New()

	for _, q := range queries {
		mustNot = append(mustNot, q.Build())
	}

	boolQuery["must_not"] = mustNot
	return b
}

// MinimumShouldMatch sets the minimum number of should clauses that must match
func (b *Builder) MinimumShouldMatch(count int) *Builder {
	// Use the safe "comma-ok" type assertion
	boolQuery, ok := b.query["bool"].(map[string]any)
	if !ok {
		// If the builder is not a bool query, panic with a helpful message
		panic("query: cannot call MinimumShouldMatch() on a non-bool query builder (e.g., a Term, Match, or Range query)")
	}

	boolQuery["minimum_should_match"] = count
	return b
}

// Build returns the internal query map
func (b *Builder) Build() map[string]any {
	return b.query
}

// MarshalJSON implements json.Marshaler
func (b *Builder) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.query)
}

// String returns a JSON representation of the query
func (b *Builder) String() string {
	bytes, _ := json.MarshalIndent(b.query, "", "  ")
	return string(bytes)
}

// Term creates a term query builder
func Term(field string, value any) *Builder {
	return &Builder{
		query: map[string]any{
			"term": map[string]any{
				field: value,
			},
		},
	}
}

// Terms creates a terms query builder
func Terms(field string, values ...any) *Builder {
	return &Builder{
		query: map[string]any{
			"terms": map[string]any{
				field: values,
			},
		},
	}
}

// Match creates a match query builder
func Match(field string, text string) *Builder {
	return &Builder{
		query: map[string]any{
			"match": map[string]any{
				field: text,
			},
		},
	}
}

// MatchPhrase creates a match_phrase query builder
func MatchPhrase(field string, text string) *Builder {
	return &Builder{
		query: map[string]any{
			"match_phrase": map[string]any{
				field: text,
			},
		},
	}
}

// MultiMatch creates a multi_match query builder
func MultiMatch(text string, fields ...string) *Builder {
	return &Builder{
		query: map[string]any{
			"multi_match": map[string]any{
				"query":  text,
				"fields": fields,
			},
		},
	}
}

// MatchAll creates a match_all query builder
func MatchAll() *Builder {
	return &Builder{
		query: map[string]any{
			"match_all": map[string]any{},
		},
	}
}

// MatchNone creates a match_none query builder
func MatchNone() *Builder {
	return &Builder{
		query: map[string]any{
			"match_none": map[string]any{},
		},
	}
}

// Exists creates an exists query builder
func Exists(field string) *Builder {
	return &Builder{
		query: map[string]any{
			"exists": map[string]any{
				"field": field,
			},
		},
	}
}

// IDs creates an ids query builder
func IDs(ids ...string) *Builder {
	return &Builder{
		query: map[string]any{
			"ids": map[string]any{
				"values": ids,
			},
		},
	}
}

// Prefix creates a prefix query builder
func Prefix(field string, prefix string) *Builder {
	return &Builder{
		query: map[string]any{
			"prefix": map[string]any{
				field: prefix,
			},
		},
	}
}

// Wildcard creates a wildcard query builder
func Wildcard(field string, pattern string) *Builder {
	return &Builder{
		query: map[string]any{
			"wildcard": map[string]any{
				field: pattern,
			},
		},
	}
}

// Regexp creates a regexp query builder
func Regexp(field string, pattern string) *Builder {
	return &Builder{
		query: map[string]any{
			"regexp": map[string]any{
				field: pattern,
			},
		},
	}
}

// Fuzzy creates a fuzzy query builder
func Fuzzy(field string, value string) *Builder {
	return &Builder{
		query: map[string]any{
			"fuzzy": map[string]any{
				field: value,
			},
		},
	}
}

// RangeBuilder provides a fluent interface for building range queries
type RangeBuilder struct {
	field string
	query map[string]any
}

// Range creates a new range query builder for the specified field
func Range(field string) *RangeBuilder {
	return &RangeBuilder{
		field: field,
		query: map[string]any{},
	}
}

// Gte sets the greater than or equal to value
func (r *RangeBuilder) Gte(value any) *RangeBuilder {
	r.query["gte"] = value
	return r
}

// Gt sets the greater than value
func (r *RangeBuilder) Gt(value any) *RangeBuilder {
	r.query["gt"] = value
	return r
}

// Lte sets the less than or equal to value
func (r *RangeBuilder) Lte(value any) *RangeBuilder {
	r.query["lte"] = value
	return r
}

// Lt sets the less than value
func (r *RangeBuilder) Lt(value any) *RangeBuilder {
	r.query["lt"] = value
	return r
}

// Format sets the date format for date range queries
func (r *RangeBuilder) Format(format string) *RangeBuilder {
	r.query["format"] = format
	return r
}

// TimeZone sets the timezone for date range queries
func (r *RangeBuilder) TimeZone(tz string) *RangeBuilder {
	r.query["time_zone"] = tz
	return r
}

// Build converts the range builder to a query builder
func (r *RangeBuilder) Build() *Builder {
	return &Builder{
		query: map[string]any{
			"range": map[string]any{
				r.field: r.query,
			},
		},
	}
}

// Helper functions for Bool query clauses

// Must creates a must clause
func Must(queries ...*Builder) *Builder {
	builder := New()
	return builder.Must(queries...)
}

// Filter creates a filter clause
func Filter(queries ...*Builder) *Builder {
	builder := New()
	return builder.Filter(queries...)
}

// Should creates a should clause
func Should(queries ...*Builder) *Builder {
	builder := New()
	return builder.Should(queries...)
}

// MustNot creates a must_not clause
func MustNot(queries ...*Builder) *Builder {
	builder := New()
	return builder.MustNot(queries...)
}

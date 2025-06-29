package elastic

import (
	"strings"
)

// Common search helpers

// MatchQuery creates a match query
func MatchQuery(field string, value any) map[string]any {
	return map[string]any{
		"match": map[string]any{
			field: value,
		},
	}
}

// TermQuery creates a term query
func TermQuery(field string, value any) map[string]any {
	return map[string]any{
		"term": map[string]any{
			field: value,
		},
	}
}

// RangeQuery creates a range query
func RangeQuery(field string, gte, lte any) map[string]any {
	rangeQuery := map[string]any{}
	if gte != nil {
		rangeQuery["gte"] = gte
	}
	if lte != nil {
		rangeQuery["lte"] = lte
	}
	return map[string]any{
		"range": map[string]any{
			field: rangeQuery,
		},
	}
}

// BoolQuery creates a bool query
func BoolQuery() map[string]any {
	return map[string]any{
		"bool": map[string]any{
			"must":     []any{},
			"must_not": []any{},
			"should":   []any{},
			"filter":   []any{},
		},
	}
}

// WithMust adds must clauses to a bool query
func WithMust(boolQuery map[string]any, queries ...map[string]any) map[string]any {
	if boolMap, ok := boolQuery["bool"].(map[string]any); ok {
		if must, ok := boolMap["must"].([]any); ok {
			for _, query := range queries {
				must = append(must, query)
			}
			boolMap["must"] = must
		}
	}
	return boolQuery
}

// WithFilter adds filter clauses to a bool query
func WithFilter(boolQuery map[string]any, queries ...map[string]any) map[string]any {
	if boolMap, ok := boolQuery["bool"].(map[string]any); ok {
		if filter, ok := boolMap["filter"].([]any); ok {
			for _, query := range queries {
				filter = append(filter, query)
			}
			boolMap["filter"] = filter
		}
	}
	return boolQuery
}

// Sort helpers

// SortAsc creates an ascending sort
func SortAsc(field string) map[string]any {
	return map[string]any{
		field: map[string]any{
			"order": "asc",
		},
	}
}

// SortDesc creates a descending sort
func SortDesc(field string) map[string]any {
	return map[string]any{
		field: map[string]any{
			"order": "desc",
		},
	}
}

// Aggregation helpers

// TermsAggregation creates a terms aggregation
func TermsAggregation(field string, size int) map[string]any {
	return map[string]any{
		"terms": map[string]any{
			"field": field,
			"size":  size,
		},
	}
}

// DateHistogramAggregation creates a date histogram aggregation
func DateHistogramAggregation(field, interval string) map[string]any {
	return map[string]any{
		"date_histogram": map[string]any{
			"field":    field,
			"interval": interval,
		},
	}
}

// StatsAggregation creates a stats aggregation
func StatsAggregation(field string) map[string]any {
	return map[string]any{
		"stats": map[string]any{
			"field": field,
		},
	}
}

// Helper functions for building queries

// BuildSearchQuery builds a complete search query with common options
func BuildSearchQuery(query map[string]any, options ...SearchOption) map[string]any {
	searchQuery := map[string]any{
		"query": query,
	}

	// Apply options
	for _, option := range options {
		option(searchQuery)
	}

	return searchQuery
}

// SearchOption represents a search query option
type SearchOption func(map[string]any)

// WithIndices sets the target indices for the search (supports single or multiple indices)
func WithIndices(indices ...string) SearchOption {
	return func(query map[string]any) {
		query["indices"] = indices
	}
}

// WithSize sets the size parameter
func WithSize(size int) SearchOption {
	return func(query map[string]any) {
		query["size"] = size
	}
}

// WithFrom sets the from parameter
func WithFrom(from int) SearchOption {
	return func(query map[string]any) {
		query["from"] = from
	}
}

// WithSort adds sort parameters (can be called multiple times to add multiple sort fields)
func WithSort(sorts ...map[string]any) SearchOption {
	return func(query map[string]any) {
		if existing, ok := query["sort"]; ok {
			if existingSorts, ok := existing.([]map[string]any); ok {
				query["sort"] = append(existingSorts, sorts...)
			} else {
				// If existing sort is not the expected type, replace it
				query["sort"] = sorts
			}
		} else {
			query["sort"] = sorts
		}
	}
}

// WithAggregations sets the aggregations parameter
func WithAggregations(aggs map[string]any) SearchOption {
	return func(query map[string]any) {
		query["aggs"] = aggs
	}
}

// WithSource adds fields to include in results (can be called multiple times to add more fields)
func WithSource(includes ...string) SearchOption {
	return func(query map[string]any) {
		if existing, ok := query["_source"]; ok {
			switch v := existing.(type) {
			case string:
				// Convert single string to slice and append
				query["_source"] = append([]string{v}, includes...)
			case []string:
				// Append to existing slice
				query["_source"] = append(v, includes...)
			default:
				// If existing _source is not string or []string, replace it
				if len(includes) == 1 {
					query["_source"] = includes[0]
				} else {
					query["_source"] = includes
				}
			}
		} else {
			// No existing _source
			if len(includes) == 1 {
				query["_source"] = includes[0]
			} else {
				query["_source"] = includes
			}
		}
	}
}

// WithTimeout sets the timeout parameter
func WithTimeout(timeout string) SearchOption {
	return func(query map[string]any) {
		query["timeout"] = timeout
	}
}

// Common filter builders

// ByID creates a filter for finding by _id
func ByID(id any) map[string]any {
	return map[string]any{"ids": map[string]any{"values": []any{id}}}
}

// ByField creates a filter for a specific field
func ByField(field string, value any) map[string]any {
	return TermQuery(field, value)
}

// ByFields creates a filter for multiple fields (bool query with must clauses)
func ByFields(fields map[string]any) map[string]any {
	query := BoolQuery()
	for field, value := range fields {
		query = WithMust(query, TermQuery(field, value))
	}
	return query
}

// Common update builders (for update by query operations)

// SetScript creates a script for setting field values
func SetScript(fields map[string]any) map[string]any {
	source := "ctx._source.putAll(params)"
	return map[string]any{
		"source": source,
		"params": fields,
	}
}

// IncScript creates a script for incrementing field values
func IncScript(fields map[string]any) map[string]any {
	var statements []string
	for field := range fields {
		statements = append(statements, "ctx._source."+field+" += params."+field)
	}
	return map[string]any{
		"source": strings.Join(statements, "; "),
		"params": fields,
	}
}

// Common query builders (extending existing ones)

// ExistsQuery creates an exists query
func ExistsQuery(field string) map[string]any {
	return map[string]any{
		"exists": map[string]any{
			"field": field,
		},
	}
}

// WildcardQuery creates a wildcard query
func WildcardQuery(field, pattern string) map[string]any {
	return map[string]any{
		"wildcard": map[string]any{
			field: pattern,
		},
	}
}

// PrefixQuery creates a prefix query
func PrefixQuery(field, prefix string) map[string]any {
	return map[string]any{
		"prefix": map[string]any{
			field: prefix,
		},
	}
}

// FuzzyQuery creates a fuzzy query
func FuzzyQuery(field string, value any, fuzziness ...string) map[string]any {
	query := map[string]any{
		"fuzzy": map[string]any{
			field: map[string]any{
				"value": value,
			},
		},
	}

	if len(fuzziness) > 0 {
		query["fuzzy"].(map[string]any)[field].(map[string]any)["fuzziness"] = fuzziness[0]
	}

	return query
}

// Common aggregation builders (extending existing ones)

// AvgAggregation creates an average aggregation
func AvgAggregation(field string) map[string]any {
	return map[string]any{
		"avg": map[string]any{
			"field": field,
		},
	}
}

// MinAggregation creates a min aggregation
func MinAggregation(field string) map[string]any {
	return map[string]any{
		"min": map[string]any{
			"field": field,
		},
	}
}

// MaxAggregation creates a max aggregation
func MaxAggregation(field string) map[string]any {
	return map[string]any{
		"max": map[string]any{
			"field": field,
		},
	}
}

// SumAggregation creates a sum aggregation
func SumAggregation(field string) map[string]any {
	return map[string]any{
		"sum": map[string]any{
			"field": field,
		},
	}
}

// CardinalityAggregation creates a cardinality aggregation
func CardinalityAggregation(field string) map[string]any {
	return map[string]any{
		"cardinality": map[string]any{
			"field": field,
		},
	}
}

// TopHitsAggregation creates a top hits aggregation
func TopHitsAggregation(size int, sorts ...map[string]any) map[string]any {
	agg := map[string]any{
		"top_hits": map[string]any{
			"size": size,
		},
	}

	if len(sorts) > 0 {
		agg["top_hits"].(map[string]any)["sort"] = sorts
	}

	return agg
}

// Common sort builders (extending existing ones)

// SortByScore creates a sort by _score
func SortByScore(desc bool) map[string]any {
	order := "asc"
	if desc {
		order = "desc"
	}
	return map[string]any{
		"_score": map[string]any{
			"order": order,
		},
	}
}

// MultiSort creates multiple sorts
func MultiSort(sorts ...map[string]any) []map[string]any {
	return sorts
}

// Additional useful query builders

// MatchAllQuery creates a match_all query (useful for getting all documents)
func MatchAllQuery() map[string]any {
	return map[string]any{
		"match_all": map[string]any{},
	}
}

// MatchNoneQuery creates a match_none query (useful for empty results)
func MatchNoneQuery() map[string]any {
	return map[string]any{
		"match_none": map[string]any{},
	}
}

// MatchPhraseQuery creates a match_phrase query (exact phrase matching)
func MatchPhraseQuery(field string, phrase any) map[string]any {
	return map[string]any{
		"match_phrase": map[string]any{
			field: phrase,
		},
	}
}

// MultiMatchQuery creates a multi_match query (search across multiple fields)
func MultiMatchQuery(query any, fields ...string) map[string]any {
	return map[string]any{
		"multi_match": map[string]any{
			"query":  query,
			"fields": fields,
		},
	}
}

// TermsQuery creates a terms query (multiple exact values)
func TermsQuery(field string, values ...any) map[string]any {
	return map[string]any{
		"terms": map[string]any{
			field: values,
		},
	}
}

// Additional aggregation helpers

// FiltersAggregation creates a filters aggregation (multiple named filters)
func FiltersAggregation(filters map[string]map[string]any) map[string]any {
	return map[string]any{
		"filters": map[string]any{
			"filters": filters,
		},
	}
}

// Search convenience builders

// SimpleSearch creates a simple search query with common options
func SimpleSearch(queryText string, fields []string, size int) map[string]any {
	var query map[string]any

	if len(fields) > 1 {
		// Multi-field search
		query = MultiMatchQuery(queryText, fields...)
	} else if len(fields) == 1 {
		// Single field search
		query = MatchQuery(fields[0], queryText)
	} else {
		// Search all fields
		query = map[string]any{
			"query_string": map[string]any{
				"query": queryText,
			},
		}
	}

	return BuildSearchQuery(query, WithSize(size))
}

// PaginatedSearch creates a paginated search query
func PaginatedSearch(query map[string]any, page, pageSize int, sortField string, sortAsc bool) map[string]any {
	from := (page - 1) * pageSize

	var sort map[string]any
	if sortAsc {
		sort = SortAsc(sortField)
	} else {
		sort = SortDesc(sortField)
	}

	return BuildSearchQuery(
		query,
		WithSize(pageSize),
		WithFrom(from),
		WithSort(sort),
	)
}

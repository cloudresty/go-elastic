package elastic

// AggregationBuilder provides a fluent interface for building aggregations
type AggregationBuilder struct {
	agg map[string]any
}

// NewTermsAggregation creates a terms aggregation
func NewTermsAggregation(field string) *AggregationBuilder {
	return &AggregationBuilder{
		agg: map[string]any{
			"terms": map[string]any{
				"field": field,
			},
		},
	}
}

// NewDateHistogramAggregation creates a date histogram aggregation
func NewDateHistogramAggregation(field string, interval string) *AggregationBuilder {
	return &AggregationBuilder{
		agg: map[string]any{
			"date_histogram": map[string]any{
				"field":    field,
				"interval": interval,
			},
		},
	}
}

// NewRangeAggregation creates a range aggregation
func NewRangeAggregation(field string) *AggregationBuilder {
	return &AggregationBuilder{
		agg: map[string]any{
			"range": map[string]any{
				"field":  field,
				"ranges": []any{},
			},
		},
	}
}

// NewHistogramAggregation creates a histogram aggregation
func NewHistogramAggregation(field string, interval float64) *AggregationBuilder {
	return &AggregationBuilder{
		agg: map[string]any{
			"histogram": map[string]any{
				"field":    field,
				"interval": interval,
			},
		},
	}
}

// NewAvgAggregation creates an average aggregation
func NewAvgAggregation(field string) *AggregationBuilder {
	return &AggregationBuilder{
		agg: map[string]any{
			"avg": map[string]any{
				"field": field,
			},
		},
	}
}

// NewSumAggregation creates a sum aggregation
func NewSumAggregation(field string) *AggregationBuilder {
	return &AggregationBuilder{
		agg: map[string]any{
			"sum": map[string]any{
				"field": field,
			},
		},
	}
}

// NewMaxAggregation creates a max aggregation
func NewMaxAggregation(field string) *AggregationBuilder {
	return &AggregationBuilder{
		agg: map[string]any{
			"max": map[string]any{
				"field": field,
			},
		},
	}
}

// NewMinAggregation creates a min aggregation
func NewMinAggregation(field string) *AggregationBuilder {
	return &AggregationBuilder{
		agg: map[string]any{
			"min": map[string]any{
				"field": field,
			},
		},
	}
}

// NewStatsAggregation creates a stats aggregation
func NewStatsAggregation(field string) *AggregationBuilder {
	return &AggregationBuilder{
		agg: map[string]any{
			"stats": map[string]any{
				"field": field,
			},
		},
	}
}

// Size sets the size for terms aggregations
func (a *AggregationBuilder) Size(size int) *AggregationBuilder {
	if terms, ok := a.agg["terms"].(map[string]any); ok {
		terms["size"] = size
	}
	return a
}

// Order sets the order for terms aggregations
func (a *AggregationBuilder) Order(field string, direction string) *AggregationBuilder {
	if terms, ok := a.agg["terms"].(map[string]any); ok {
		terms["order"] = map[string]any{
			field: direction,
		}
	}
	return a
}

// MinDocCount sets the minimum document count for terms aggregations
func (a *AggregationBuilder) MinDocCount(count int) *AggregationBuilder {
	if terms, ok := a.agg["terms"].(map[string]any); ok {
		terms["min_doc_count"] = count
	}
	return a
}

// AddRange adds a range to a range aggregation
func (a *AggregationBuilder) AddRange(key string, from, to *float64) *AggregationBuilder {
	if rangeAgg, ok := a.agg["range"].(map[string]any); ok {
		if ranges, ok := rangeAgg["ranges"].([]any); ok {
			rangeEntry := map[string]any{"key": key}
			if from != nil {
				rangeEntry["from"] = *from
			}
			if to != nil {
				rangeEntry["to"] = *to
			}
			rangeAgg["ranges"] = append(ranges, rangeEntry)
		}
	}
	return a
}

// Format sets the format for date histogram aggregations
func (a *AggregationBuilder) Format(format string) *AggregationBuilder {
	if dateHist, ok := a.agg["date_histogram"].(map[string]any); ok {
		dateHist["format"] = format
	}
	return a
}

// TimeZone sets the timezone for date histogram aggregations
func (a *AggregationBuilder) TimeZone(tz string) *AggregationBuilder {
	if dateHist, ok := a.agg["date_histogram"].(map[string]any); ok {
		dateHist["time_zone"] = tz
	}
	return a
}

// SubAggregation adds a sub-aggregation
func (a *AggregationBuilder) SubAggregation(name string, subAgg *AggregationBuilder) *AggregationBuilder {
	if a.agg["aggs"] == nil {
		a.agg["aggs"] = map[string]any{}
	}
	if aggs, ok := a.agg["aggs"].(map[string]any); ok {
		aggs[name] = subAgg.Build()
	}
	return a
}

// Build returns the aggregation as a map
func (a *AggregationBuilder) Build() map[string]any {
	return a.agg
}

// WithAggregation creates a search option for aggregations
func WithAggregation(name string, agg *AggregationBuilder) SearchOption {
	return WithAggregations(map[string]any{
		name: agg.Build(),
	})
}

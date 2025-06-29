package elastic

import "context"

// DocumentsService bulk methods

// Bulk returns a BulkIndexer for chaining bulk operations on the specified index
func (s *DocumentsService) Bulk(indexName string) *BulkIndexer {
	return &BulkIndexer{
		client:     s.client,
		index:      indexName,
		operations: make([]*BulkOperation, 0),
	}
}

// BulkIndexer provides a fluent interface for building bulk operations
type BulkIndexer struct {
	client     *Client
	index      string
	operations []*BulkOperation
}

// Create adds a create operation to the bulk request (fails if document exists)
func (bi *BulkIndexer) Create(document any) *BulkIndexer {
	op := &BulkOperation{
		Action:   "create",
		Index:    bi.index,
		Document: document,
	}
	bi.operations = append(bi.operations, op)
	return bi
}

// CreateWithID adds a create operation with specific ID to the bulk request
func (bi *BulkIndexer) CreateWithID(id string, document any) *BulkIndexer {
	op := &BulkOperation{
		Action:   "create",
		Index:    bi.index,
		ID:       id,
		Document: document,
	}
	bi.operations = append(bi.operations, op)
	return bi
}

// Index adds an index operation to the bulk request (creates or replaces)
func (bi *BulkIndexer) Index(id string, document any) *BulkIndexer {
	op := &BulkOperation{
		Action:   "index",
		Index:    bi.index,
		ID:       id,
		Document: document,
	}
	bi.operations = append(bi.operations, op)
	return bi
}

// Update adds an update operation to the bulk request
func (bi *BulkIndexer) Update(id string, document any) *BulkIndexer {
	op := &BulkOperation{
		Action:   "update",
		Index:    bi.index,
		ID:       id,
		Document: document,
	}
	bi.operations = append(bi.operations, op)
	return bi
}

// UpdateWithScript adds an update operation with script to the bulk request
func (bi *BulkIndexer) UpdateWithScript(id string, script map[string]any) *BulkIndexer {
	op := &BulkOperation{
		Action: "update",
		Index:  bi.index,
		ID:     id,
		Script: script,
	}
	bi.operations = append(bi.operations, op)
	return bi
}

// Delete adds a delete operation to the bulk request
func (bi *BulkIndexer) Delete(id string) *BulkIndexer {
	op := &BulkOperation{
		Action: "delete",
		Index:  bi.index,
		ID:     id,
	}
	bi.operations = append(bi.operations, op)
	return bi
}

// Do executes the bulk request with all accumulated operations
func (bi *BulkIndexer) Do(ctx context.Context) (*BulkResponse, error) {
	bulkResource := &BulkResource{
		client: bi.client,
		index:  bi.index,
	}
	return bulkResource.Execute(ctx, bi.operations)
}

// Legacy methods for backward compatibility

// BulkRaw performs bulk operations using raw operation maps
func (s *DocumentsService) BulkRaw(ctx context.Context, operations []map[string]any) (*BulkResponse, error) {
	bulkResource := &BulkResource{
		client: s.client,
	}
	return bulkResource.ExecuteRaw(ctx, operations)
}

// ForIndex returns a BulkResource configured for a specific index
func (s *DocumentsService) ForIndex(indexName string) *BulkResource {
	return &BulkResource{
		client: s.client,
		index:  indexName,
	}
}

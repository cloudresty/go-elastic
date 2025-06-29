package elastic

import (
	"encoding/json"
	"time"

	"context"

	"github.com/cloudresty/emit"
)

// DocumentsService CRUD methods

// Get retrieves a document by ID
func (s *DocumentsService) Get(ctx context.Context, indexName, documentID string) (map[string]any, error) {
	doc := &Document{
		client: s.client,
		index:  indexName,
	}
	return doc.Get(ctx, documentID)
}

// MultiGet retrieves multiple documents by their IDs (uses Elasticsearch _mget API)
func (s *DocumentsService) MultiGet(ctx context.Context, indexName string, documentIDs []string) ([]map[string]any, error) {
	doc := &Document{
		client: s.client,
		index:  indexName,
	}
	return doc.GetMany(ctx, documentIDs)
}

// Create creates a new document with automatic ID generation
func (s *DocumentsService) Create(ctx context.Context, indexName string, document any) (*IndexResponse, error) {
	doc := &Document{
		client: s.client,
		index:  indexName,
	}
	return doc.Index(ctx, document)
}

// CreateWithID creates a new document with a specific ID (fails if document already exists)
func (s *DocumentsService) CreateWithID(ctx context.Context, indexName, documentID string, document any) (*IndexResponse, error) {
	doc := &Document{
		client: s.client,
		index:  indexName,
	}
	return doc.CreateWithID(ctx, documentID, document)
}

// Update updates a document
func (s *DocumentsService) Update(ctx context.Context, indexName, documentID string, document map[string]any) (*UpdateResponse, error) {
	doc := &Document{
		client: s.client,
		index:  indexName,
	}
	return doc.Update(ctx, documentID, document)
}

// Delete deletes a document by ID
func (s *DocumentsService) Delete(ctx context.Context, indexName, documentID string) (*DeleteResponse, error) {
	doc := &Document{
		client: s.client,
		index:  indexName,
	}
	return doc.Delete(ctx, documentID)
}

// Index creates or replaces a document with a specific ID (equivalent to PUT /<index>/_doc/<id>)
func (s *DocumentsService) Index(ctx context.Context, indexName, documentID string, document any) (*IndexResponse, error) {
	doc := &Document{
		client: s.client,
		index:  indexName,
	}
	return doc.IndexWithID(ctx, documentID, document)
}

// Exists checks if a document exists (more efficient than Get for existence checks)
func (s *DocumentsService) Exists(ctx context.Context, indexName, documentID string) (bool, error) {
	doc := &Document{
		client: s.client,
		index:  indexName,
	}
	return doc.Exists(ctx, documentID)
}

// UpdateByQuery updates all documents matching a query
func (s *DocumentsService) UpdateByQuery(ctx context.Context, indexName string, query map[string]any, script map[string]any) (map[string]any, error) {
	doc := &Document{
		client: s.client,
		index:  indexName,
	}
	return doc.UpdateByQuery(ctx, query, script)
}

// DeleteByQuery deletes all documents matching a query
func (s *DocumentsService) DeleteByQuery(ctx context.Context, indexName string, query map[string]any) (map[string]any, error) {
	doc := &Document{
		client: s.client,
		index:  indexName,
	}
	return doc.DeleteByQuery(ctx, query)
}

// GetIndex returns a Document resource for the given index for direct access
func (s *DocumentsService) GetIndex(indexName string) *Document {
	return &Document{
		client: s.client,
		index:  indexName,
	}
}

// enhanceDocument adds ID and metadata to a document based on client configuration
func (c *Client) enhanceDocument(doc any) map[string]any {
	var docMap map[string]any

	// Convert document to map
	if m, ok := doc.(map[string]any); ok {
		docMap = make(map[string]any)
		for k, v := range m {
			docMap[k] = v
		}
	} else {
		// Try to convert via JSON
		jsonBytes, err := json.Marshal(doc)
		if err != nil {
			emit.Error.StructuredFields("Failed to marshal document",
				emit.ZString("error", err.Error()))
			return map[string]any{}
		}
		if err := json.Unmarshal(jsonBytes, &docMap); err != nil {
			emit.Error.StructuredFields("Failed to unmarshal document",
				emit.ZString("error", err.Error()))
			return map[string]any{}
		}
	}

	// Add ID if not present and not in custom mode
	if c.config.IDMode != IDModeCustom {
		if _, exists := docMap["_id"]; !exists {
			switch c.config.IDMode {
			case IDModeULID:
				docMap["_id"] = generateULID()
			case IDModeElastic:
				// Let Elasticsearch generate its own random ID for optimal shard distribution
				// Don't set _id field - Elasticsearch will auto-generate
			default:
				// Default to Elasticsearch's ID generation for best performance
				// Don't set _id field - Elasticsearch will auto-generate
			}
		}
	}

	// Add timestamps
	now := time.Now()
	if _, exists := docMap["created_at"]; !exists {
		docMap["created_at"] = now
	}
	docMap["updated_at"] = now

	return docMap
}

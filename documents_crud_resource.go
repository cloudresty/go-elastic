package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/elastic/go-elasticsearch/v9/esapi"
)

// Document provides document-level operations for a specific index
type Document struct {
	client *Client
	index  string
}

// Index indexes a document with automatic ID generation
func (d *Document) Index(ctx context.Context, document any) (*IndexResponse, error) {
	return d.IndexWithID(ctx, "", document)
}

// IndexWithID indexes a document with a specific ID
func (d *Document) IndexWithID(ctx context.Context, documentID string, document any) (*IndexResponse, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second) //nolint:ineffassign
		defer cancel()
	}

	// Enhance document with metadata
	enhancedDoc := d.client.enhanceDocument(document)

	// Use the _id from enhanced document if no ID provided
	if documentID == "" {
		if id, exists := enhancedDoc["_id"]; exists {
			if idStr, ok := id.(string); ok {
				documentID = idStr
			}
		}
	}

	// Convert document to JSON
	docBytes, err := json.Marshal(enhancedDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal document: %w", err)
	}

	// Prepare the index request
	req := esapi.IndexRequest{
		Index:      d.index,
		DocumentID: documentID,
		Body:       bytes.NewReader(docBytes),
		Refresh:    "wait_for",
	}

	res, err := req.Do(ctx, d.client.client)
	if err != nil {
		return nil, fmt.Errorf("failed to execute index request: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("index request failed: %s - %s", res.Status(), string(body))
	}

	var indexResponse IndexResponse
	if err := json.NewDecoder(res.Body).Decode(&indexResponse); err != nil {
		return nil, fmt.Errorf("failed to decode index response: %w", err)
	}

	d.client.config.Logger.Info("Document indexed successfully - index: %s, document_id: %s, result: %s", d.index, indexResponse.ID, indexResponse.Result)

	return &indexResponse, nil
}

// Get retrieves a document by ID
func (d *Document) Get(ctx context.Context, documentID string) (map[string]any, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second) //nolint:ineffassign
		defer cancel()
	}

	req := esapi.GetRequest{
		Index:      d.index,
		DocumentID: documentID,
	}

	res, err := req.Do(ctx, d.client.client)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get request: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}()

	if res.IsError() {
		if res.StatusCode == 404 {
			return nil, fmt.Errorf("document not found")
		}
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("get request failed: %s - %s", res.Status(), string(body))
	}

	var getResponse struct {
		ID     string         `json:"_id"`
		Source map[string]any `json:"_source"`
		Found  bool           `json:"found"`
	}

	if err := json.NewDecoder(res.Body).Decode(&getResponse); err != nil {
		return nil, fmt.Errorf("failed to decode get response: %w", err)
	}

	if !getResponse.Found {
		return nil, fmt.Errorf("document not found")
	}

	d.client.config.Logger.Debug("Document retrieved successfully - index: %s, document_id: %s", d.index, documentID)

	return getResponse.Source, nil
}

// GetMany retrieves multiple documents by their IDs
func (d *Document) GetMany(ctx context.Context, documentIDs []string) ([]map[string]any, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second) //nolint:ineffassign
		defer cancel()
	}

	if len(documentIDs) == 0 {
		return []map[string]any{}, nil
	}

	// Build multi-get request
	mgetBody := map[string]any{
		"docs": make([]map[string]any, len(documentIDs)),
	}

	for i, id := range documentIDs {
		mgetBody["docs"].([]map[string]any)[i] = map[string]any{
			"_index": d.index,
			"_id":    id,
		}
	}

	bodyBytes, err := json.Marshal(mgetBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal mget request: %w", err)
	}

	req := esapi.MgetRequest{
		Body: bytes.NewReader(bodyBytes),
	}

	res, err := req.Do(ctx, d.client.client)
	if err != nil {
		return nil, fmt.Errorf("failed to execute mget request: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("mget request failed: %s - %s", res.Status(), string(body))
	}

	var mgetResponse struct {
		Docs []struct {
			ID     string         `json:"_id"`
			Source map[string]any `json:"_source"`
			Found  bool           `json:"found"`
		} `json:"docs"`
	}

	if err := json.NewDecoder(res.Body).Decode(&mgetResponse); err != nil {
		return nil, fmt.Errorf("failed to decode mget response: %w", err)
	}

	// Extract found documents
	var documents []map[string]any
	for _, doc := range mgetResponse.Docs {
		if doc.Found {
			documents = append(documents, doc.Source)
		}
	}

	d.client.config.Logger.Debug("Documents retrieved successfully - index: %s, requested: %d, found: %d", d.index, len(documentIDs), len(documents))

	return documents, nil
}

// Update updates a document
func (d *Document) Update(ctx context.Context, documentID string, doc map[string]any) (*UpdateResponse, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second) //nolint:ineffassign
		defer cancel()
	}

	// Wrap the document in an update request
	updateDoc := map[string]any{
		"doc": doc,
	}

	// Add updated_at timestamp
	if _, exists := doc["updated_at"]; !exists {
		updateDoc["doc"].(map[string]any)["updated_at"] = time.Now()
	}

	docBytes, err := json.Marshal(updateDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update document: %w", err)
	}

	req := esapi.UpdateRequest{
		Index:      d.index,
		DocumentID: documentID,
		Body:       bytes.NewReader(docBytes),
		Refresh:    "wait_for",
	}

	res, err := req.Do(ctx, d.client.client)
	if err != nil {
		return nil, fmt.Errorf("failed to execute update request: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("update request failed: %s - %s", res.Status(), string(body))
	}

	var updateResponse UpdateResponse
	if err := json.NewDecoder(res.Body).Decode(&updateResponse); err != nil {
		return nil, fmt.Errorf("failed to decode update response: %w", err)
	}

	d.client.config.Logger.Info("Document updated successfully - index: %s, document_id: %s, result: %s", d.index, documentID, updateResponse.Result)

	return &updateResponse, nil
}

// Delete deletes a document
func (d *Document) Delete(ctx context.Context, documentID string) (*DeleteResponse, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second) //nolint:ineffassign
		defer cancel()
	}

	req := esapi.DeleteRequest{
		Index:      d.index,
		DocumentID: documentID,
		Refresh:    "wait_for",
	}

	res, err := req.Do(ctx, d.client.client)
	if err != nil {
		return nil, fmt.Errorf("failed to execute delete request: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}()

	if res.IsError() {
		if res.StatusCode == 404 {
			return nil, fmt.Errorf("document not found")
		}
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("delete request failed: %s - %s", res.Status(), string(body))
	}

	var deleteResponse DeleteResponse
	if err := json.NewDecoder(res.Body).Decode(&deleteResponse); err != nil {
		return nil, fmt.Errorf("failed to decode delete response: %w", err)
	}

	d.client.config.Logger.Info("Document deleted successfully - index: %s, document_id: %s, result: %s", d.index, documentID, deleteResponse.Result)

	return &deleteResponse, nil
}

// Exists checks if a document exists using HEAD request (more efficient than GET)
func (d *Document) Exists(ctx context.Context, documentID string) (bool, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	req := esapi.ExistsRequest{
		Index:      d.index,
		DocumentID: documentID,
	}

	res, err := req.Do(ctx, d.client.client)
	if err != nil {
		d.client.config.Logger.Error("Failed to check document existence - index: %s, document_id: %s, error: %s", d.index, documentID, err.Error())
		return false, fmt.Errorf("failed to check document existence: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			d.client.config.Logger.Warn("Failed to close response body - error: %s", err.Error())
		}
	}()

	// HTTP 200 means document exists, HTTP 404 means it doesn't exist
	// Any other status code is an error
	switch res.StatusCode {
	case 200:
		d.client.config.Logger.Debug("Document exists - index: %s, document_id: %s", d.index, documentID)
		return true, nil
	case 404:
		d.client.config.Logger.Debug("Document does not exist - index: %s, document_id: %s", d.index, documentID)
		return false, nil
	default:
		bodyBytes, _ := io.ReadAll(res.Body)
		d.client.config.Logger.Error("Unexpected status when checking document existence - index: %s, document_id: %s, status: %s, response: %s", d.index, documentID, res.Status(), string(bodyBytes))
		return false, fmt.Errorf("unexpected status when checking document existence: %s - %s", res.Status(), string(bodyBytes))
	}
}

// CreateWithID creates a document with a specific ID using the _create endpoint (fails if document exists)
func (d *Document) CreateWithID(ctx context.Context, documentID string, document any) (*IndexResponse, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	// Enhance document with metadata
	enhancedDoc := d.client.enhanceDocument(document)

	// Convert document to JSON
	docBytes, err := json.Marshal(enhancedDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal document: %w", err)
	}

	// Use the _create endpoint which fails if document already exists
	req := esapi.CreateRequest{
		Index:      d.index,
		DocumentID: documentID,
		Body:       io.NopCloser(bytes.NewReader(docBytes)),
	}

	res, err := req.Do(ctx, d.client.client)
	if err != nil {
		d.client.config.Logger.Error("Failed to create document - index: %s, document_id: %s, error: %s", d.index, documentID, err.Error())
		return nil, fmt.Errorf("failed to create document: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			d.client.config.Logger.Warn("Failed to close response body - error: %s", err.Error())
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		d.client.config.Logger.Error("Failed to create document - index: %s, document_id: %s, status: %s, response: %s", d.index, documentID, res.Status(), string(bodyBytes))
		return nil, fmt.Errorf("create document failed: %s - %s", res.Status(), string(bodyBytes))
	}

	var indexResponse IndexResponse
	if err := json.NewDecoder(res.Body).Decode(&indexResponse); err != nil {
		return nil, fmt.Errorf("failed to decode create response: %w", err)
	}

	d.client.config.Logger.Info("Document created successfully - index: %s, document_id: %s, result: %s", d.index, documentID, indexResponse.Result)

	return &indexResponse, nil
}

// UpdateByQuery updates all documents matching a query using the _update_by_query API
func (d *Document) UpdateByQuery(ctx context.Context, query map[string]any, script map[string]any) (map[string]any, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second) // Longer timeout for bulk operations
		defer cancel()
	}

	// Build the request body
	body := map[string]any{
		"query": query,
	}
	if script != nil {
		body["script"] = script
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update by query body: %w", err)
	}

	req := esapi.UpdateByQueryRequest{
		Index: []string{d.index},
		Body:  io.NopCloser(bytes.NewReader(bodyBytes)),
	}

	res, err := req.Do(ctx, d.client.client)
	if err != nil {
		d.client.config.Logger.Error("Failed to update by query - index: %s, error: %s", d.index, err.Error())
		return nil, fmt.Errorf("failed to update by query: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			d.client.config.Logger.Warn("Failed to close response body - error: %s", err.Error())
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		d.client.config.Logger.Error("Update by query failed - index: %s, status: %s, response: %s", d.index, res.Status(), string(bodyBytes))
		return nil, fmt.Errorf("update by query failed: %s - %s", res.Status(), string(bodyBytes))
	}

	var result map[string]any
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode update by query response: %w", err)
	}

	d.client.config.Logger.Info("Update by query completed - index: %s", d.index)

	return result, nil
}

// DeleteByQuery deletes all documents matching a query using the _delete_by_query API
func (d *Document) DeleteByQuery(ctx context.Context, query map[string]any) (map[string]any, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second) // Longer timeout for bulk operations
		defer cancel()
	}

	// Build the request body
	body := map[string]any{
		"query": query,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal delete by query body: %w", err)
	}

	req := esapi.DeleteByQueryRequest{
		Index: []string{d.index},
		Body:  io.NopCloser(bytes.NewReader(bodyBytes)),
	}

	res, err := req.Do(ctx, d.client.client)
	if err != nil {
		d.client.config.Logger.Error("Failed to delete by query - index: %s, error: %s", d.index, err.Error())
		return nil, fmt.Errorf("failed to delete by query: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			d.client.config.Logger.Warn("Failed to close response body - error: %s", err.Error())
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		d.client.config.Logger.Error("Delete by query failed - index: %s, status: %s, response: %s", d.index, res.Status(), string(bodyBytes))
		return nil, fmt.Errorf("delete by query failed: %s - %s", res.Status(), string(bodyBytes))
	}

	var result map[string]any
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode delete by query response: %w", err)
	}

	d.client.config.Logger.Info("Delete by query completed - index: %s", d.index)

	return result, nil
}

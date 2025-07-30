package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v9/esapi"
)

// BulkResource provides bulk operations
type BulkResource struct {
	client *Client
	index  string // optional default index
}

// BulkOperation represents a single bulk operation
type BulkOperation struct {
	Action    string         `json:"action"`   // index, create, update, delete
	Index     string         `json:"index"`    // target index
	ID        string         `json:"id"`       // document ID
	Document  any            `json:"document"` // document data (can be any type)
	Source    map[string]any `json:"_source"`  // for updates
	Script    map[string]any `json:"script"`   // for script updates
	UpsertDoc map[string]any `json:"doc"`      // for upserts
}

// Index adds an index operation to the bulk request
func (br *BulkResource) Index(indexName, documentID string, document any) *BulkOperation {
	if indexName == "" && br.index != "" {
		indexName = br.index
	}

	return &BulkOperation{
		Action:   "index",
		Index:    indexName,
		ID:       documentID,
		Document: document,
	}
}

// Create adds a create operation to the bulk request
func (br *BulkResource) Create(indexName, documentID string, document any) *BulkOperation {
	if indexName == "" && br.index != "" {
		indexName = br.index
	}

	return &BulkOperation{
		Action:   "create",
		Index:    indexName,
		ID:       documentID,
		Document: document,
	}
}

// Update adds an update operation to the bulk request
func (br *BulkResource) Update(indexName, documentID string, doc any) *BulkOperation {
	if indexName == "" && br.index != "" {
		indexName = br.index
	}

	return &BulkOperation{
		Action:   "update",
		Index:    indexName,
		ID:       documentID,
		Document: doc,
	}
}

// UpdateWithScript adds an update operation with script to the bulk request
func (br *BulkResource) UpdateWithScript(indexName, documentID string, script map[string]any) *BulkOperation {
	if indexName == "" && br.index != "" {
		indexName = br.index
	}

	return &BulkOperation{
		Action: "update",
		Index:  indexName,
		ID:     documentID,
		Script: script,
	}
}

// Delete adds a delete operation to the bulk request
func (br *BulkResource) Delete(indexName, documentID string) *BulkOperation {
	if indexName == "" && br.index != "" {
		indexName = br.index
	}

	return &BulkOperation{
		Action: "delete",
		Index:  indexName,
		ID:     documentID,
	}
}

// Execute performs a bulk operation with the given operations
func (br *BulkResource) Execute(ctx context.Context, operations []*BulkOperation) (*BulkResponse, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	if len(operations) == 0 {
		return nil, fmt.Errorf("no operations provided")
	}

	// Build bulk request body
	var body strings.Builder
	for _, op := range operations {
		// Action line
		actionLine := map[string]map[string]any{
			op.Action: {
				"_index": op.Index,
			},
		}

		if op.ID != "" {
			actionLine[op.Action]["_id"] = op.ID
		}

		actionBytes, err := json.Marshal(actionLine)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal action line: %w", err)
		}
		body.Write(actionBytes)
		body.WriteString("\n")

		// Document line (if needed)
		switch op.Action {
		case "index", "create":
			if op.Document != nil {
				// Enhance document with metadata
				enhanced := br.client.enhanceDocument(op.Document)
				docBytes, err := json.Marshal(enhanced)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal document: %w", err)
				}
				body.Write(docBytes)
				body.WriteString("\n")
			}
		case "update":
			updateDoc := make(map[string]any)
			if op.UpsertDoc != nil {
				updateDoc["doc"] = op.UpsertDoc
				updateDoc["doc_as_upsert"] = true
			}
			if op.Script != nil {
				updateDoc["script"] = op.Script
			}

			docBytes, err := json.Marshal(updateDoc)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal update document: %w", err)
			}
			body.Write(docBytes)
			body.WriteString("\n")
		}
		// Delete operations only need the action line
	}

	req := esapi.BulkRequest{
		Body: strings.NewReader(body.String()),
	}

	res, err := req.Do(ctx, br.client.client)
	if err != nil {
		br.client.config.Logger.Error("Bulk operation failed - operations: %d, error: %s", len(operations), err.Error())
		return nil, fmt.Errorf("bulk request failed: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			br.client.config.Logger.Warn("Failed to close response body - error: %s", err.Error())
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		br.client.config.Logger.Error("Bulk operation failed - operations: %d, status: %s, response: %s", len(operations), res.Status(), string(bodyBytes))
		return nil, fmt.Errorf("bulk operation failed: %s - %s", res.Status(), string(bodyBytes))
	}

	var bulkResponse BulkResponse
	if err := json.NewDecoder(res.Body).Decode(&bulkResponse); err != nil {
		return nil, fmt.Errorf("failed to decode bulk response: %w", err)
	}

	br.client.config.Logger.Info("Bulk operation completed successfully - operations: %d, took: %d, errors: %t", len(operations), bulkResponse.Took, bulkResponse.Errors)

	return &bulkResponse, nil
}

// ExecuteRaw performs a bulk operation with raw operations (legacy compatibility)
func (br *BulkResource) ExecuteRaw(ctx context.Context, operations []map[string]any) (*BulkResponse, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	if len(operations) == 0 {
		return nil, fmt.Errorf("no operations provided")
	}

	// Build bulk request body
	var body strings.Builder
	for _, op := range operations {
		opBytes, err := json.Marshal(op)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal operation: %w", err)
		}
		body.Write(opBytes)
		body.WriteString("\n")
	}

	req := esapi.BulkRequest{
		Body: strings.NewReader(body.String()),
	}

	res, err := req.Do(ctx, br.client.client)
	if err != nil {
		return nil, fmt.Errorf("bulk request failed: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			br.client.config.Logger.Warn("Failed to close response body - error: %s", err.Error())
		}
	}()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("bulk operation failed: %s - %s", res.Status(), string(bodyBytes))
	}

	var bulkResponse BulkResponse
	if err := json.NewDecoder(res.Body).Decode(&bulkResponse); err != nil {
		return nil, fmt.Errorf("failed to decode bulk response: %w", err)
	}

	return &bulkResponse, nil
}

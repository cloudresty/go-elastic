package elastic

import (
	"testing"
)

func TestIDModeConfiguration(t *testing.T) {
	tests := []struct {
		name   string
		idMode IDMode
	}{
		{
			name:   "Elastic mode (default)",
			idMode: IDModeElastic,
		},
		{
			name:   "ULID mode",
			idMode: IDModeULID,
		},
		{
			name:   "Custom mode",
			idMode: IDModeCustom,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				Hosts:  []string{"localhost:9200"},
				IDMode: tt.idMode,
			}

			client := &Client{
				config: config,
			}

			doc := map[string]any{
				"name":  "test",
				"value": 123,
			}

			enhanced := client.enhanceDocument(doc)

			// Check the result based on mode
			switch tt.idMode {
			case IDModeElastic:
				if _, exists := enhanced["_id"]; exists {
					t.Errorf("Expected no _id field for elastic mode, but found: %v", enhanced["_id"])
				}
			case IDModeULID:
				if id, ok := enhanced["_id"].(string); !ok || len(id) != 26 {
					t.Errorf("Expected ULID string of length 26, got %T: %v", enhanced["_id"], enhanced["_id"])
				}
			case IDModeCustom:
				if _, exists := enhanced["_id"]; exists {
					t.Errorf("Expected no _id field for custom mode, but found: %v", enhanced["_id"])
				}
			}

			// Verify other fields are present
			if enhanced["name"] != "test" {
				t.Errorf("Expected name field to be preserved")
			}
			if enhanced["value"] != 123 {
				t.Errorf("Expected value field to be preserved")
			}
			if enhanced["created_at"] == nil {
				t.Errorf("Expected created_at field to be added")
			}
			if enhanced["updated_at"] == nil {
				t.Errorf("Expected updated_at field to be added")
			}
		})
	}
}

func TestUserProvidedID(t *testing.T) {
	config := &Config{
		Hosts:  []string{"localhost:9200"},
		IDMode: IDModeULID,
	}

	client := &Client{
		config: config,
	}

	// Test with user-provided ID
	doc := map[string]any{
		"_id":   "user-provided-id",
		"name":  "test",
		"value": 123,
	}

	enhanced := client.enhanceDocument(doc)

	// Should preserve user-provided ID
	if enhanced["_id"] != "user-provided-id" {
		t.Errorf("Expected user-provided ID to be preserved, got: %v", enhanced["_id"])
	}
}

func TestIDModeValidation(t *testing.T) {
	tests := []struct {
		mode  string
		valid bool
	}{
		{"ulid", true},
		{"custom", true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.mode, func(t *testing.T) {
			result := isValidIDMode(tt.mode)
			if result != tt.valid {
				t.Errorf("Expected isValidIDMode(%s) = %v, got %v", tt.mode, tt.valid, result)
			}
		})
	}
}

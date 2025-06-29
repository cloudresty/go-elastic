package elastic

import (
	"testing"
)

func TestBuildConnectionAddresses(t *testing.T) {
	tests := []struct {
		name     string
		hosts    []string
		expected []string
	}{
		{
			name:     "single host with port",
			hosts:    []string{"localhost:9200"},
			expected: []string{"http://localhost:9200"},
		},
		{
			name:     "single host without port (gets default)",
			hosts:    []string{"localhost"},
			expected: []string{"http://localhost:9200"},
		},
		{
			name:     "multiple hosts with ports",
			hosts:    []string{"host1:9201", "host2:9202", "host3:9203"},
			expected: []string{"http://host1:9201", "http://host2:9202", "http://host3:9203"},
		},
		{
			name:     "multiple hosts without ports (get default)",
			hosts:    []string{"host1", "host2", "host3"},
			expected: []string{"http://host1:9200", "http://host2:9200", "http://host3:9200"},
		},
		{
			name:     "mixed hosts (some with ports, some without)",
			hosts:    []string{"host1", "host2:9201", "host3"},
			expected: []string{"http://host1:9200", "http://host2:9201", "http://host3:9200"},
		},
		{
			name:     "TLS enabled",
			hosts:    []string{"secure-host:9200"},
			expected: []string{"https://secure-host:9200"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				Hosts:      tt.hosts,
				TLSEnabled: tt.name == "TLS enabled",
			}

			addresses := config.BuildConnectionAddresses()

			if len(addresses) != len(tt.expected) {
				t.Errorf("Expected %d addresses, got %d", len(tt.expected), len(addresses))
				return
			}

			for i, expected := range tt.expected {
				if addresses[i] != expected {
					t.Errorf("Expected address[%d] '%s', got '%s'", i, expected, addresses[i])
				}
			}
		})
	}
}

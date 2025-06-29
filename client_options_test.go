package elastic

import (
	"testing"
)

func TestClientOptions(t *testing.T) {
	// Test 1: Default client creation
	t.Run("default client", func(t *testing.T) {
		// This should not fail even if Elasticsearch is not running
		// because we're just testing the configuration logic

		// Simulate what NewClient does without actually connecting
		config, err := loadConfigWithPrefix("")
		if err != nil {
			t.Fatalf("Failed to load default config: %v", err)
		}

		if len(config.Hosts) == 0 || config.Hosts[0] != "localhost:9200" {
			t.Errorf("Expected default hosts ['localhost:9200'], got %v", config.Hosts)
		}
	})

	// Test 2: FromEnvWithPrefix option
	t.Run("with env prefix", func(t *testing.T) {
		opts := &clientOptions{}

		FromEnvWithPrefix("TEST_")(opts)

		if opts.config == nil {
			t.Fatal("Expected config to be created by FromEnvWithPrefix")
		}
		// The config should be loaded from environment with TEST_ prefix
		// We can't test specific values since they depend on environment variables
	})

	// Test 3: WithHosts option
	t.Run("with hosts", func(t *testing.T) {
		opts := &clientOptions{}

		// First apply WithHosts which should create a config
		WithHosts("customhost")(opts)

		if opts.config == nil {
			t.Fatal("Expected config to be created by WithHosts")
		}
		if len(opts.config.Hosts) != 1 || opts.config.Hosts[0] != "customhost" {
			t.Errorf("Expected hosts ['customhost'], got %v", opts.config.Hosts)
		}
	})

	// Test 3.1: WithHosts option with multiple hosts
	t.Run("with multiple hosts", func(t *testing.T) {
		opts := &clientOptions{}

		// Apply WithHosts with multiple hosts
		WithHosts("host1", "host2", "host3")(opts)

		if opts.config == nil {
			t.Fatal("Expected config to be created by WithHosts")
		}
		expectedHosts := []string{"host1", "host2", "host3"}
		if len(opts.config.Hosts) != 3 {
			t.Fatalf("Expected 3 hosts, got %d", len(opts.config.Hosts))
		}
		for i, expected := range expectedHosts {
			if opts.config.Hosts[i] != expected {
				t.Errorf("Expected host[%d] '%s', got '%s'", i, expected, opts.config.Hosts[i])
			}
		}
	})

	// Test 4: WithCredentials option
	t.Run("with credentials", func(t *testing.T) {
		opts := &clientOptions{}

		WithCredentials("user", "pass")(opts)

		if opts.config == nil {
			t.Fatal("Expected config to be created by WithCredentials")
		}
		if opts.config.Username != "user" {
			t.Errorf("Expected username 'user', got '%s'", opts.config.Username)
		}
		if opts.config.Password != "pass" {
			t.Errorf("Expected password 'pass', got '%s'", opts.config.Password)
		}
	})

	// Test 5: Multiple options
	t.Run("multiple options", func(t *testing.T) {
		opts := &clientOptions{}

		// Apply multiple options - using FromEnv first, then overrides
		FromEnv()(opts)
		WithHosts("multihost:9201")(opts)
		WithCredentials("multiuser", "multipass")(opts)

		if opts.config == nil {
			t.Fatal("Expected config to be created")
		}
		if len(opts.config.Hosts) == 0 || opts.config.Hosts[0] != "multihost:9201" {
			t.Errorf("Expected hosts ['multihost:9201'], got %v", opts.config.Hosts)
		}
		if opts.config.Username != "multiuser" {
			t.Errorf("Expected username 'multiuser', got '%s'", opts.config.Username)
		}
		if opts.config.Password != "multipass" {
			t.Errorf("Expected password 'multipass', got '%s'", opts.config.Password)
		}
	})

	// Test 6: WithConfig option
	t.Run("with config", func(t *testing.T) {
		customConfig := &Config{
			Hosts:    []string{"confighost:9202"},
			Username: "configuser",
		}

		opts := &clientOptions{}

		WithConfig(customConfig)(opts)

		if opts.config != customConfig {
			t.Error("Expected config to be set to custom config")
		}
		if len(opts.config.Hosts) == 0 || opts.config.Hosts[0] != "confighost:9202" {
			t.Errorf("Expected hosts ['confighost:9202'], got %v", opts.config.Hosts)
		}
	})
}

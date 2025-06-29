package elastic

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cloudresty/go-env"
)

// loadConfigWithPrefix loads Elasticsearch configuration with a custom environment prefix
// This is an internal function used by FromEnv() and FromEnvWithPrefix() functional options
func loadConfigWithPrefix(prefix string) (*Config, error) {
	// Create empty config struct - go-env will apply defaults from struct tags
	config := &Config{}

	bindOptions := env.DefaultBindingOptions()
	if prefix != "" {
		bindOptions.Prefix = prefix
	}

	// Bind environment variables and apply defaults from struct tags
	if err := env.Bind(config, bindOptions); err != nil {
		return nil, fmt.Errorf("failed to load environment config: %w", err)
	}

	// Parse RetryOnStatus from string if needed
	if err := parseRetryOnStatus(config, prefix); err != nil {
		return nil, fmt.Errorf("failed to parse retry status codes: %w", err)
	}

	// Validate the final configuration
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// parseRetryOnStatus parses the ELASTICSEARCH_RETRY_ON_STATUS environment variable
func parseRetryOnStatus(config *Config, prefix string) error {
	envVar := "ELASTICSEARCH_RETRY_ON_STATUS"
	if prefix != "" {
		envVar = prefix + envVar
	}

	retryStatusStr, exists := env.Lookup(envVar)
	if !exists || retryStatusStr == "" {
		// Set default retry status codes
		config.RetryOnStatus = []int{502, 503, 504, 429}
		return nil
	}

	statusCodes := strings.Split(retryStatusStr, ",")
	config.RetryOnStatus = make([]int, 0, len(statusCodes))

	for _, codeStr := range statusCodes {
		codeStr = strings.TrimSpace(codeStr)
		if codeStr == "" {
			continue
		}

		code, err := strconv.Atoi(codeStr)
		if err != nil {
			return fmt.Errorf("invalid status code '%s': %w", codeStr, err)
		}

		if code < 100 || code > 599 {
			return fmt.Errorf("invalid HTTP status code: %d", code)
		}

		config.RetryOnStatus = append(config.RetryOnStatus, code)
	}

	return nil
}

// validateConfig validates the Elasticsearch configuration
func validateConfig(config *Config) error {
	// Validate connection settings
	if config.CloudID == "" {
		// Only validate hosts if not using Cloud ID
		if len(config.Hosts) == 0 {
			return errors.New("ELASTICSEARCH_HOSTS must be set when not using ELASTICSEARCH_CLOUD_ID")
		}
		// Validate that all hosts have ports specified
		for _, host := range config.Hosts {
			if !strings.Contains(host, ":") {
				return fmt.Errorf("host '%s' must include a port (e.g., %s:9200)", host, host)
			}
		}
	}

	// Validate timeouts
	if config.ConnectTimeout <= 0 {
		return errors.New("connect timeout must be positive")
	}
	if config.RequestTimeout <= 0 {
		return errors.New("request timeout must be positive")
	}

	// Validate retry settings
	if config.MaxRetries < 0 {
		return errors.New("max retries cannot be negative")
	}

	// Validate reconnection settings
	if config.ReconnectDelay <= 0 {
		config.ReconnectDelay = 5 * time.Second
	}
	if config.MaxReconnectDelay <= 0 {
		config.MaxReconnectDelay = 1 * time.Minute
	}
	if config.ReconnectBackoff <= 1.0 {
		config.ReconnectBackoff = 2.0
	}
	if config.MaxReconnectAttempts < 0 {
		config.MaxReconnectAttempts = 10
	}

	// Validate health check settings
	if config.HealthCheckInterval <= 0 {
		config.HealthCheckInterval = 30 * time.Second
	}

	// Validate log level
	if !isValidLogLevel(config.LogLevel) {
		return fmt.Errorf("invalid log level: %s", config.LogLevel)
	}

	// Validate log format
	if !isValidLogFormat(config.LogFormat) {
		return fmt.Errorf("invalid log format: %s", config.LogFormat)
	}

	// Validate ID mode
	if !isValidIDMode(string(config.IDMode)) {
		return fmt.Errorf("invalid ID mode: %s", config.IDMode)
	}

	return nil
}

// isValidLogLevel checks if the log level is valid
func isValidLogLevel(level string) bool {
	validLevels := []string{"debug", "info", "warn", "error"}
	for _, valid := range validLevels {
		if level == valid {
			return true
		}
	}
	return false
}

// isValidLogFormat checks if the log format is valid
func isValidLogFormat(format string) bool {
	validFormats := []string{"json", "text"}
	for _, valid := range validFormats {
		if format == valid {
			return true
		}
	}
	return false
}

// isValidIDMode checks if the ID mode is valid
// Note: "elastic" (default) is recommended for optimal shard distribution
// "ulid" can cause shard hotspotting and should be used with caution
func isValidIDMode(mode string) bool {
	validModes := []string{"elastic", "ulid", "custom"}
	for _, valid := range validModes {
		if mode == valid {
			return true
		}
	}
	return false
}

// Environment variable names for reference
const (
	EnvElasticsearchHost                 = "ELASTICSEARCH_HOST"
	EnvElasticsearchPort                 = "ELASTICSEARCH_PORT"
	EnvElasticsearchUsername             = "ELASTICSEARCH_USERNAME"
	EnvElasticsearchPassword             = "ELASTICSEARCH_PASSWORD"
	EnvElasticsearchAPIKey               = "ELASTICSEARCH_API_KEY"
	EnvElasticsearchCloudID              = "ELASTICSEARCH_CLOUD_ID"
	EnvElasticsearchServiceToken         = "ELASTICSEARCH_SERVICE_TOKEN"
	EnvElasticsearchTLSEnabled           = "ELASTICSEARCH_TLS_ENABLED"
	EnvElasticsearchTLSInsecure          = "ELASTICSEARCH_TLS_INSECURE"
	EnvElasticsearchCompressionEnabled   = "ELASTICSEARCH_COMPRESSION_ENABLED"
	EnvElasticsearchRetryOnStatus        = "ELASTICSEARCH_RETRY_ON_STATUS"
	EnvElasticsearchMaxRetries           = "ELASTICSEARCH_MAX_RETRIES"
	EnvElasticsearchDiscoverNodesOnStart = "ELASTICSEARCH_DISCOVER_NODES_ON_START"
	EnvElasticsearchMaxIdleConns         = "ELASTICSEARCH_MAX_IDLE_CONNS"
	EnvElasticsearchMaxIdleConnsPerHost  = "ELASTICSEARCH_MAX_IDLE_CONNS_PER_HOST"
	EnvElasticsearchIdleConnTimeout      = "ELASTICSEARCH_IDLE_CONN_TIMEOUT"
	EnvElasticsearchMaxConnLifetime      = "ELASTICSEARCH_MAX_CONN_LIFETIME"
	EnvElasticsearchConnectTimeout       = "ELASTICSEARCH_CONNECT_TIMEOUT"
	EnvElasticsearchRequestTimeout       = "ELASTICSEARCH_REQUEST_TIMEOUT"
	EnvElasticsearchReconnectEnabled     = "ELASTICSEARCH_RECONNECT_ENABLED"
	EnvElasticsearchReconnectDelay       = "ELASTICSEARCH_RECONNECT_DELAY"
	EnvElasticsearchMaxReconnectDelay    = "ELASTICSEARCH_MAX_RECONNECT_DELAY"
	EnvElasticsearchReconnectBackoff     = "ELASTICSEARCH_RECONNECT_BACKOFF"
	EnvElasticsearchMaxReconnectAttempts = "ELASTICSEARCH_MAX_RECONNECT_ATTEMPTS"
	EnvElasticsearchHealthCheckEnabled   = "ELASTICSEARCH_HEALTH_CHECK_ENABLED"
	EnvElasticsearchHealthCheckInterval  = "ELASTICSEARCH_HEALTH_CHECK_INTERVAL"
	EnvElasticsearchAppName              = "ELASTICSEARCH_APP_NAME"
	EnvElasticsearchConnectionName       = "ELASTICSEARCH_CONNECTION_NAME"
	EnvElasticsearchIDMode               = "ELASTICSEARCH_ID_MODE"
	EnvElasticsearchLogLevel             = "ELASTICSEARCH_LOG_LEVEL"
	EnvElasticsearchLogFormat            = "ELASTICSEARCH_LOG_FORMAT"
)

package elastic

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cloudresty/emit"
	"github.com/elastic/go-elasticsearch/v9"
)

// IDMode defines the ID generation strategy for documents
type IDMode string

const (
	// IDModeElastic uses Elasticsearch's default random ID generation (default, recommended)
	// This ensures optimal shard distribution and write performance
	IDModeElastic IDMode = "elastic"
	// IDModeULID generates ULID strings as document IDs
	// WARNING: Can cause shard hotspotting in multi-shard indices due to time-ordering
	// Use only when you need sortable IDs and understand the trade-offs
	IDModeULID IDMode = "ulid"
	// IDModeCustom allows users to provide their own _id fields
	IDModeCustom IDMode = "custom"
)

// Client represents an Elasticsearch client with auto-reconnection and environment-first configuration
type Client struct {
	client         *elasticsearch.Client
	config         *Config
	mutex          sync.RWMutex
	isConnected    bool
	reconnectCount int64
	lastReconnect  time.Time
	healthTicker   *time.Ticker
	shutdownChan   chan struct{}
	shutdownOnce   sync.Once
}

// Config holds Elasticsearch connection configuration
type Config struct {
	// Connection settings
	Hosts    []string `env:"ELASTICSEARCH_HOSTS,default=localhost:9200"` // Single or multiple hosts with ports (comma-separated)
	Username string   `env:"ELASTICSEARCH_USERNAME"`
	Password string   `env:"ELASTICSEARCH_PASSWORD"`
	APIKey   string   `env:"ELASTICSEARCH_API_KEY"`

	// Cloud settings
	CloudID      string `env:"ELASTICSEARCH_CLOUD_ID"`
	ServiceToken string `env:"ELASTICSEARCH_SERVICE_TOKEN"`

	// TLS settings
	TLSEnabled  bool `env:"ELASTICSEARCH_TLS_ENABLED,default=false"`
	TLSInsecure bool `env:"ELASTICSEARCH_TLS_INSECURE,default=false"`

	// Performance settings
	CompressionEnabled   bool  `env:"ELASTICSEARCH_COMPRESSION_ENABLED,default=true"`
	RetryOnStatus        []int `env:"ELASTICSEARCH_RETRY_ON_STATUS"`
	MaxRetries           int   `env:"ELASTICSEARCH_MAX_RETRIES,default=3"`
	DiscoverNodesOnStart bool  `env:"ELASTICSEARCH_DISCOVER_NODES_ON_START,default=false"`

	// Connection pool settings
	MaxIdleConns        int           `env:"ELASTICSEARCH_MAX_IDLE_CONNS,default=100"`
	MaxIdleConnsPerHost int           `env:"ELASTICSEARCH_MAX_IDLE_CONNS_PER_HOST,default=10"`
	IdleConnTimeout     time.Duration `env:"ELASTICSEARCH_IDLE_CONN_TIMEOUT,default=90s"`
	MaxConnLifetime     time.Duration `env:"ELASTICSEARCH_MAX_CONN_LIFETIME,default=0s"` // 0 = no limit

	// Timeout settings
	ConnectTimeout time.Duration `env:"ELASTICSEARCH_CONNECT_TIMEOUT,default=10s"`
	RequestTimeout time.Duration `env:"ELASTICSEARCH_REQUEST_TIMEOUT,default=30s"`

	// Reconnection settings
	ReconnectEnabled     bool          `env:"ELASTICSEARCH_RECONNECT_ENABLED,default=true"`
	ReconnectDelay       time.Duration `env:"ELASTICSEARCH_RECONNECT_DELAY,default=5s"`
	MaxReconnectDelay    time.Duration `env:"ELASTICSEARCH_MAX_RECONNECT_DELAY,default=1m"`
	ReconnectBackoff     float64       `env:"ELASTICSEARCH_RECONNECT_BACKOFF,default=2.0"`
	MaxReconnectAttempts int           `env:"ELASTICSEARCH_MAX_RECONNECT_ATTEMPTS,default=10"`

	// Health check settings
	HealthCheckEnabled  bool          `env:"ELASTICSEARCH_HEALTH_CHECK_ENABLED,default=true"`
	HealthCheckInterval time.Duration `env:"ELASTICSEARCH_HEALTH_CHECK_INTERVAL,default=30s"`

	// Application settings
	AppName        string `env:"ELASTICSEARCH_APP_NAME,default=go-elastic-app"`
	ConnectionName string `env:"ELASTICSEARCH_CONNECTION_NAME"`

	// ID Generation settings
	IDMode IDMode `env:"ELASTICSEARCH_ID_MODE,default=elastic"`

	// Logging
	LogLevel  string `env:"ELASTICSEARCH_LOG_LEVEL,default=info"`
	LogFormat string `env:"ELASTICSEARCH_LOG_FORMAT,default=json"`
}

// BuildConnectionAddresses constructs Elasticsearch connection addresses from configuration
func (c *Config) BuildConnectionAddresses() []string {
	if c.CloudID != "" {
		// When using Cloud ID, addresses are handled by the Elasticsearch client
		return nil
	}

	scheme := "http"
	if c.TLSEnabled {
		scheme = "https"
	}

	var addresses []string

	// Use Hosts (all hosts must include ports)
	if len(c.Hosts) > 0 {
		for _, host := range c.Hosts {
			// Ensure host includes a port, default to 9200 if not specified
			if !strings.Contains(host, ":") {
				host = host + ":9200"
			}
			address := fmt.Sprintf("%s://%s", scheme, host)
			addresses = append(addresses, address)
		}
		return addresses
	}

	// Should not reach here with default config, but provide localhost fallback
	address := fmt.Sprintf("%s://localhost:9200", scheme)
	return []string{address}
}

// ConnectionStats represents connection statistics
type ConnectionStats struct {
	IsConnected   bool      `json:"is_connected"`
	Reconnects    int64     `json:"reconnects"`
	LastReconnect time.Time `json:"last_reconnect"`
}

// ClientOption represents a functional option for configuring the client
type ClientOption func(*clientOptions)

// clientOptions holds the configuration options for client creation
type clientOptions struct {
	config *Config
	prefix string
}

// WithConfig sets a custom configuration for the client
func WithConfig(config *Config) ClientOption {
	return func(opts *clientOptions) {
		opts.config = config
	}
}

// WithHosts sets custom hosts for the client (overrides environment)
// For single host, use: WithHosts("localhost")
// For multiple hosts, use: WithHosts("host1", "host2", "host3")
func WithHosts(hosts ...string) ClientOption {
	return func(opts *clientOptions) {
		if opts.config == nil {
			// Create a new config if none exists
			config, err := loadConfigWithPrefix("")
			if err != nil {
				// Use default config if loading fails
				config = &Config{}
			}
			opts.config = config
		}

		if len(hosts) > 0 {
			// Set all hosts
			opts.config.Hosts = hosts
		}
	}
}

// WithCredentials sets username and password for the client (overrides environment)
func WithCredentials(username, password string) ClientOption {
	return func(opts *clientOptions) {
		if opts.config == nil {
			// Create a new config if none exists
			config, err := loadConfigWithPrefix("")
			if err != nil {
				// Use default config if loading fails
				config = &Config{}
			}
			opts.config = config
		}
		opts.config.Username = username
		opts.config.Password = password
	}
}

// WithAPIKey sets API key for the client (overrides environment)
func WithAPIKey(apiKey string) ClientOption {
	return func(opts *clientOptions) {
		if opts.config == nil {
			// Create a new config if none exists
			config, err := loadConfigWithPrefix("")
			if err != nil {
				// Use default config if loading fails
				config = &Config{}
			}
			opts.config = config
		}
		opts.config.APIKey = apiKey
	}
}

// WithCloudID sets Elastic Cloud ID for the client (overrides environment)
func WithCloudID(cloudID string) ClientOption {
	return func(opts *clientOptions) {
		if opts.config == nil {
			// Create a new config if none exists
			config, err := loadConfigWithPrefix("")
			if err != nil {
				// Use default config if loading fails
				config = &Config{}
			}
			opts.config = config
		}
		opts.config.CloudID = cloudID
	}
}

// WithTLS enables or disables TLS for the client (overrides environment)
func WithTLS(enabled bool) ClientOption {
	return func(opts *clientOptions) {
		if opts.config == nil {
			// Create a new config if none exists
			config, err := loadConfigWithPrefix("")
			if err != nil {
				// Use default config if loading fails
				config = &Config{}
			}
			opts.config = config
		}
		opts.config.TLSEnabled = enabled
	}
}

// WithConnectionName sets a connection name for the client (useful for logging and identification)
func WithConnectionName(name string) ClientOption {
	return func(opts *clientOptions) {
		if opts.config == nil {
			// Create a new config if none exists
			config, err := loadConfigWithPrefix("")
			if err != nil {
				// Use default config if loading fails
				config = &Config{}
			}
			opts.config = config
		}
		opts.config.ConnectionName = name
	}
}

// FromEnv loads configuration from environment variables using the default
// "ELASTICSEARCH_" prefix. This is a functional option for NewClient.
// Example: client, err := elastic.NewClient(elastic.FromEnv())
func FromEnv() ClientOption {
	return FromEnvWithPrefix("")
}

// FromEnvWithPrefix loads configuration from environment variables using a
// custom prefix. For example, a prefix of "MYAPP_" would look for
// "MYAPP_ELASTICSEARCH_HOSTS".
// Example: client, err := elastic.NewClient(elastic.FromEnvWithPrefix("LOGS_"))
func FromEnvWithPrefix(prefix string) ClientOption {
	return func(opts *clientOptions) {
		// Load configuration from environment
		config, err := loadConfigWithPrefix(prefix)
		if err != nil {
			// If loading fails, create default config and continue
			// This allows other options to still configure the client
			config = &Config{}
		}
		opts.config = config
		opts.prefix = prefix
	}
}

// NewClient creates a new Elasticsearch client with functional options
func NewClient(options ...ClientOption) (*Client, error) {
	opts := &clientOptions{
		prefix: "", // default empty prefix
	}

	// Apply all options
	for _, option := range options {
		option(opts)
	}

	// If no config was provided via WithConfig, load from environment
	var config *Config
	var err error
	if opts.config != nil {
		config = opts.config
	} else {
		config, err = loadConfigWithPrefix("")
		if err != nil {
			return nil, err
		}
		opts.config = config

		// Apply options again in case they modify the loaded config
		for _, option := range options {
			option(opts)
		}
		config = opts.config
	}

	// Validate config
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Get the first host for logging
	firstHost := "localhost"
	logPort := 9200
	if len(config.Hosts) > 0 {
		firstHost = config.Hosts[0]
		// Extract port from host if present
		if strings.Contains(firstHost, ":") {
			parts := strings.Split(firstHost, ":")
			firstHost = parts[0]
			if len(parts) > 1 {
				if port, err := strconv.Atoi(parts[1]); err == nil {
					logPort = port
				}
			}
		}
	}

	emit.Info.StructuredFields("Creating new Elasticsearch client",
		emit.ZString("host", firstHost),
		emit.ZInt("port", logPort),
		emit.ZString("app_name", config.AppName),
		emit.ZBool("tls_enabled", config.TLSEnabled))

	client := &Client{
		config:       config,
		shutdownChan: make(chan struct{}),
	}

	if err := client.connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to Elasticsearch: %w", err)
	}

	if config.HealthCheckEnabled {
		client.startHealthCheck()
	}

	// Get the first host for logging
	logHost2 := "localhost"
	logPort2 := 9200
	if len(config.Hosts) > 0 {
		logHost2 = config.Hosts[0]
		// Extract port from host if present
		if strings.Contains(logHost2, ":") {
			parts := strings.Split(logHost2, ":")
			logHost2 = parts[0]
			if len(parts) > 1 {
				if port, err := strconv.Atoi(parts[1]); err == nil {
					logPort2 = port
				}
			}
		}
	}

	emit.Info.StructuredFields("Elasticsearch client initialized successfully",
		emit.ZString("host", logHost2),
		emit.ZInt("port", logPort2),
		emit.ZString("app_name", config.AppName))

	return client, nil
}

// connect establishes a connection to Elasticsearch
func (c *Client) connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	esConfig := c.buildClientConfig()

	client, err := elasticsearch.NewClient(esConfig)
	if err != nil {
		return fmt.Errorf("failed to create Elasticsearch client: %w", err)
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), c.config.ConnectTimeout)
	defer cancel()

	res, err := client.Info(client.Info.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to ping Elasticsearch: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	if res.IsError() {
		return fmt.Errorf("elasticsearch returned error: %s", res.String())
	}

	c.client = client
	c.isConnected = true
	c.lastReconnect = time.Now()

	return nil
}

// buildClientConfig constructs Elasticsearch client configuration
func (c *Client) buildClientConfig() elasticsearch.Config {
	config := elasticsearch.Config{
		Addresses: c.config.BuildConnectionAddresses(),
		Username:  c.config.Username,
		Password:  c.config.Password,
		APIKey:    c.config.APIKey,
		CloudID:   c.config.CloudID,

		// Transport settings
		Transport: &http.Transport{
			MaxIdleConns:          c.config.MaxIdleConns,
			MaxIdleConnsPerHost:   c.config.MaxIdleConnsPerHost,
			IdleConnTimeout:       c.config.IdleConnTimeout,
			ResponseHeaderTimeout: c.config.RequestTimeout,
			DisableCompression:    !c.config.CompressionEnabled,
		},

		// Retry settings
		RetryOnStatus: c.config.RetryOnStatus,
		MaxRetries:    c.config.MaxRetries,

		// Discovery settings
		DiscoverNodesOnStart: c.config.DiscoverNodesOnStart,
	}

	// Set default retry statuses if not configured
	if len(config.RetryOnStatus) == 0 {
		config.RetryOnStatus = []int{502, 503, 504, 429}
	}

	return config
}

// startHealthCheck starts the health check routine
func (c *Client) startHealthCheck() {
	if c.config.HealthCheckInterval <= 0 {
		return
	}

	c.healthTicker = time.NewTicker(c.config.HealthCheckInterval)

	go func() {
		for {
			select {
			case <-c.healthTicker.C:
				c.performHealthCheck()
			case <-c.shutdownChan:
				return
			}
		}
	}()

	emit.Info.StructuredFields("Health check started",
		emit.ZDuration("interval", c.config.HealthCheckInterval))
}

// performHealthCheck performs a health check
func (c *Client) performHealthCheck() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.Ping(ctx)
	if err != nil {
		emit.Warn.StructuredFields("Health check failed",
			emit.ZString("error", err.Error()))

		if c.config.ReconnectEnabled {
			c.attemptReconnect()
		}
	}
}

// attemptReconnect attempts to reconnect to Elasticsearch
func (c *Client) attemptReconnect() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.isConnected {
		return // Already connected
	}

	attempts := 0
	delay := c.config.ReconnectDelay

	for attempts < c.config.MaxReconnectAttempts {
		attempts++

		emit.Info.StructuredFields("Attempting to reconnect to Elasticsearch",
			emit.ZInt("attempt", attempts),
			emit.ZInt("max_attempts", c.config.MaxReconnectAttempts),
			emit.ZDuration("delay", delay))

		time.Sleep(delay)

		if err := c.connect(); err == nil {
			emit.Info.StructuredFields("Successfully reconnected to Elasticsearch",
				emit.ZInt("attempts", attempts))
			c.reconnectCount++
			return
		}

		// Exponential backoff
		delay = time.Duration(float64(delay) * c.config.ReconnectBackoff)
		if delay > c.config.MaxReconnectDelay {
			delay = c.config.MaxReconnectDelay
		}
	}

	emit.Error.StructuredFields("Failed to reconnect to Elasticsearch after maximum attempts",
		emit.ZInt("max_attempts", c.config.MaxReconnectAttempts))
}

// Close closes the client and stops background routines
func (c *Client) Close() error {
	c.shutdownOnce.Do(func() {
		close(c.shutdownChan)

		if c.healthTicker != nil {
			c.healthTicker.Stop()
		}

		emit.Info.Msg("Elasticsearch client closed")
	})

	return nil
}

// GetClient returns the underlying Elasticsearch client
func (c *Client) GetClient() *elasticsearch.Client {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.client
}

// Name returns the configured connection name for this client
// This is useful for logging and identifying clients in multi-client scenarios
func (c *Client) Name() string {
	return c.config.ConnectionName
}

// Resource-oriented API methods

// Indices returns an IndicesService for index operations
func (c *Client) Indices() *IndicesService {
	return &IndicesService{
		client: c,
	}
}

// Documents returns a DocumentsService for all document operations (CRUD, search, bulk)
func (c *Client) Documents() *DocumentsService {
	return &DocumentsService{
		client: c,
	}
}

// Cluster returns a ClusterService for cluster operations
func (c *Client) Cluster() *ClusterService {
	return &ClusterService{
		client: c,
	}
}

// Convenience methods for direct index access

// Search returns an Index instance for search operations
// This is a convenience method for search-focused workflows
func (c *Client) Search(indexName string) *Index {
	return &Index{
		client: c,
		name:   indexName,
	}
}

// Index returns an Index instance for direct index operations
// This provides direct access to index-specific operations
func (c *Client) Index(indexName string) *Index {
	return &Index{
		client: c,
		name:   indexName,
	}
}

// Service types for resource-oriented API

// IndicesService provides operations for managing Elasticsearch indices
type IndicesService struct {
	client *Client
}

// DocumentsService provides operations for managing Elasticsearch documents
// This includes CRUD operations, search, and bulk operations
type DocumentsService struct {
	client *Client
}

// ClusterService provides operations for cluster management
type ClusterService struct {
	client *Client
}

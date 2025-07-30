package elastic

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Shutdownable interface for resources that can be gracefully shutdown
type Shutdownable interface {
	Close() error
}

// ShutdownConfig holds configuration for graceful shutdown
type ShutdownConfig struct {
	Timeout          time.Duration // Maximum time to wait for shutdown
	GracePeriod      time.Duration // Grace period before forcing shutdown
	ForceKillTimeout time.Duration // Time to wait before force killing
}

// ShutdownManager manages graceful shutdown of Elasticsearch clients and other resources
type ShutdownManager struct {
	clients      []*Client
	resources    []Shutdownable
	shutdownChan chan os.Signal
	ctx          context.Context
	cancel       context.CancelFunc
	mutex        sync.Mutex
	config       *ShutdownConfig
	logger       Logger
}

// NewShutdownManager creates a new shutdown manager with default configuration
func NewShutdownManager(config *ShutdownConfig, logger Logger) *ShutdownManager {
	if config == nil {
		config = &ShutdownConfig{
			Timeout:          30 * time.Second,
			GracePeriod:      5 * time.Second,
			ForceKillTimeout: 10 * time.Second,
		}
	}

	if logger == nil {
		logger = &NopLogger{}
	}

	ctx, cancel := context.WithCancel(context.Background())

	logger.Info("Creating shutdown manager - timeout: %v, grace_period: %v", config.Timeout, config.GracePeriod)

	return &ShutdownManager{
		clients:      make([]*Client, 0),
		resources:    make([]Shutdownable, 0),
		shutdownChan: make(chan os.Signal, 1),
		ctx:          ctx,
		cancel:       cancel,
		config:       config,
		logger:       logger,
	}
}

// NewShutdownManagerWithConfig creates a shutdown manager with configuration
func NewShutdownManagerWithConfig(config *Config) *ShutdownManager {
	shutdownConfig := &ShutdownConfig{
		Timeout:          config.ConnectTimeout,
		GracePeriod:      5 * time.Second,
		ForceKillTimeout: 10 * time.Second,
	}

	config.Logger.Info("Creating shutdown manager with config - timeout: %v", shutdownConfig.Timeout)

	return NewShutdownManager(shutdownConfig, config.Logger)
}

// Register registers Elasticsearch clients for graceful shutdown
func (sm *ShutdownManager) Register(clients ...*Client) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sm.clients = append(sm.clients, clients...)

	sm.logger.Info("Registered clients for graceful shutdown - count: %d", len(clients))
}

// SetupSignalHandler sets up signal handlers for graceful shutdown
func (sm *ShutdownManager) SetupSignalHandler() {
	signal.Notify(sm.shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	sm.logger.Info("Signal handlers setup for graceful shutdown")
}

// Wait blocks until a shutdown signal is received and performs graceful shutdown
func (sm *ShutdownManager) Wait() {
	sig := <-sm.shutdownChan
	sm.logger.Info("Received shutdown signal - signal: %s", sig.String())

	sm.shutdown()
}

// Context returns the shutdown manager's context for background workers
func (sm *ShutdownManager) Context() context.Context {
	return sm.ctx
}

// RegisterResources registers shutdownable resources for graceful shutdown
func (sm *ShutdownManager) RegisterResources(resources ...Shutdownable) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sm.resources = append(sm.resources, resources...)

	sm.logger.Info("Registered resources for graceful shutdown - count: %d", len(resources))
}

// shutdown performs the actual shutdown logic
func (sm *ShutdownManager) shutdown() {
	start := time.Now()

	sm.logger.Info("Starting graceful shutdown - timeout: %v", sm.config.Timeout)

	// Cancel context to signal background workers to stop
	sm.cancel()

	// Create a timeout context for the shutdown process
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), sm.config.Timeout)
	defer shutdownCancel()

	// Channel to signal completion
	done := make(chan struct{})

	go func() {
		defer close(done)

		sm.mutex.Lock()
		clients := make([]*Client, len(sm.clients))
		copy(clients, sm.clients)
		resources := make([]Shutdownable, len(sm.resources))
		copy(resources, sm.resources)
		sm.mutex.Unlock()

		// Close Elasticsearch clients
		for i, client := range clients {
			if client != nil {
				sm.logger.Info("Closing Elasticsearch client - client_index: %d", i)

				if err := client.Close(); err != nil {
					sm.logger.Error("Error closing Elasticsearch client - client_index: %d, error: %s", i, err.Error())
				} else {
					sm.logger.Info("Elasticsearch client closed successfully - client_index: %d", i)
				}
			}
		}

		// Close other resources
		for i, resource := range resources {
			if resource != nil {
				sm.logger.Info("Closing resource - resource_index: %d", i)

				if err := resource.Close(); err != nil {
					sm.logger.Error("Error closing resource - resource_index: %d, error: %s", i, err.Error())
				} else {
					sm.logger.Info("Resource closed successfully - resource_index: %d", i)
				}
			}
		}

		// Wait for grace period to allow in-flight operations to complete
		if sm.config.GracePeriod > 0 {
			sm.logger.Info("Waiting grace period for in-flight operations - grace_period: %v", sm.config.GracePeriod)

			select {
			case <-time.After(sm.config.GracePeriod):
				sm.logger.Info("Grace period completed")
			case <-shutdownCtx.Done():
				sm.logger.Warn("Grace period interrupted by timeout")
			}
		}
	}()

	// Wait for shutdown to complete or timeout
	select {
	case <-done:
		elapsed := time.Since(start)
		sm.logger.Info("Graceful shutdown completed - elapsed: %v", elapsed)
	case <-shutdownCtx.Done():
		elapsed := time.Since(start)
		sm.logger.Warn("Graceful shutdown timed out - elapsed: %v, timeout: %v", elapsed, sm.config.Timeout)

		// Force close after timeout
		if sm.config.ForceKillTimeout > 0 {
			sm.logger.Warn("Waiting before force kill - force_kill_timeout: %v", sm.config.ForceKillTimeout)

			time.Sleep(sm.config.ForceKillTimeout)
		}

		sm.logger.Error("Force killing application")
		os.Exit(1)
	}
}

// SetTimeout updates the shutdown timeout
func (sm *ShutdownManager) SetTimeout(timeout time.Duration) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.config.Timeout = timeout

	sm.logger.Info("Shutdown timeout updated - timeout: %v", timeout)
}

// GetTimeout returns the current shutdown timeout
func (sm *ShutdownManager) GetTimeout() time.Duration {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	return sm.config.Timeout
}

// GetClientCount returns the number of registered clients
func (sm *ShutdownManager) GetClientCount() int {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	return len(sm.clients)
}

// GetResourceCount returns the number of registered resources
func (sm *ShutdownManager) GetResourceCount() int {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	return len(sm.resources)
}

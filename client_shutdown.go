package elastic

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/cloudresty/emit"
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
}

// NewShutdownManager creates a new shutdown manager with default configuration
func NewShutdownManager(config *ShutdownConfig) *ShutdownManager {
	if config == nil {
		config = &ShutdownConfig{
			Timeout:          30 * time.Second,
			GracePeriod:      5 * time.Second,
			ForceKillTimeout: 10 * time.Second,
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	emit.Info.StructuredFields("Creating shutdown manager",
		emit.ZDuration("timeout", config.Timeout),
		emit.ZDuration("grace_period", config.GracePeriod))

	return &ShutdownManager{
		clients:      make([]*Client, 0),
		resources:    make([]Shutdownable, 0),
		shutdownChan: make(chan os.Signal, 1),
		ctx:          ctx,
		cancel:       cancel,
		config:       config,
	}
}

// NewShutdownManagerWithConfig creates a shutdown manager with configuration
func NewShutdownManagerWithConfig(config *Config) *ShutdownManager {
	shutdownConfig := &ShutdownConfig{
		Timeout:          config.ConnectTimeout,
		GracePeriod:      5 * time.Second,
		ForceKillTimeout: 10 * time.Second,
	}

	emit.Info.StructuredFields("Creating shutdown manager with config",
		emit.ZDuration("timeout", shutdownConfig.Timeout))

	return NewShutdownManager(shutdownConfig)
}

// Register registers Elasticsearch clients for graceful shutdown
func (sm *ShutdownManager) Register(clients ...*Client) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sm.clients = append(sm.clients, clients...)

	emit.Info.StructuredFields("Registered clients for graceful shutdown",
		emit.ZInt("count", len(clients)))
}

// SetupSignalHandler sets up signal handlers for graceful shutdown
func (sm *ShutdownManager) SetupSignalHandler() {
	signal.Notify(sm.shutdownChan, syscall.SIGINT, syscall.SIGTERM)
	emit.Info.Msg("Signal handlers setup for graceful shutdown")
}

// Wait blocks until a shutdown signal is received and performs graceful shutdown
func (sm *ShutdownManager) Wait() {
	sig := <-sm.shutdownChan
	emit.Info.StructuredFields("Received shutdown signal",
		emit.ZString("signal", sig.String()))

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

	emit.Info.StructuredFields("Registered resources for graceful shutdown",
		emit.ZInt("count", len(resources)))
}

// shutdown performs the actual shutdown logic
func (sm *ShutdownManager) shutdown() {
	start := time.Now()

	emit.Info.StructuredFields("Starting graceful shutdown",
		emit.ZDuration("timeout", sm.config.Timeout))

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
				emit.Info.StructuredFields("Closing Elasticsearch client",
					emit.ZInt("client_index", i))

				if err := client.Close(); err != nil {
					emit.Error.StructuredFields("Error closing Elasticsearch client",
						emit.ZInt("client_index", i),
						emit.ZString("error", err.Error()))
				} else {
					emit.Info.StructuredFields("Elasticsearch client closed successfully",
						emit.ZInt("client_index", i))
				}
			}
		}

		// Close other resources
		for i, resource := range resources {
			if resource != nil {
				emit.Info.StructuredFields("Closing resource",
					emit.ZInt("resource_index", i))

				if err := resource.Close(); err != nil {
					emit.Error.StructuredFields("Error closing resource",
						emit.ZInt("resource_index", i),
						emit.ZString("error", err.Error()))
				} else {
					emit.Info.StructuredFields("Resource closed successfully",
						emit.ZInt("resource_index", i))
				}
			}
		}

		// Wait for grace period to allow in-flight operations to complete
		if sm.config.GracePeriod > 0 {
			emit.Info.StructuredFields("Waiting grace period for in-flight operations",
				emit.ZDuration("grace_period", sm.config.GracePeriod))

			select {
			case <-time.After(sm.config.GracePeriod):
				emit.Info.Msg("Grace period completed")
			case <-shutdownCtx.Done():
				emit.Warn.Msg("Grace period interrupted by timeout")
			}
		}
	}()

	// Wait for shutdown to complete or timeout
	select {
	case <-done:
		elapsed := time.Since(start)
		emit.Info.StructuredFields("Graceful shutdown completed",
			emit.ZDuration("elapsed", elapsed))
	case <-shutdownCtx.Done():
		elapsed := time.Since(start)
		emit.Warn.StructuredFields("Graceful shutdown timed out",
			emit.ZDuration("elapsed", elapsed),
			emit.ZDuration("timeout", sm.config.Timeout))

		// Force close after timeout
		if sm.config.ForceKillTimeout > 0 {
			emit.Warn.StructuredFields("Waiting before force kill",
				emit.ZDuration("force_kill_timeout", sm.config.ForceKillTimeout))

			time.Sleep(sm.config.ForceKillTimeout)
		}

		emit.Error.Msg("Force killing application")
		os.Exit(1)
	}
}

// SetTimeout updates the shutdown timeout
func (sm *ShutdownManager) SetTimeout(timeout time.Duration) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.config.Timeout = timeout

	emit.Info.StructuredFields("Shutdown timeout updated",
		emit.ZDuration("timeout", timeout))
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

package elastic

import (
	"context"
	"fmt"
	"log"
)

// Package-level convenience functions for client creation

// Connect creates a new Elasticsearch client using environment variables
func Connect() (*Client, error) {
	return NewClient(FromEnv())
}

// ConnectWithConfig creates a new Elasticsearch client with the provided configuration
func ConnectWithConfig(config *Config) (*Client, error) {
	return NewClient(WithConfig(config))
}

// MustConnect creates a new Elasticsearch client or panics on error
// Use this only in main functions or initialization code where panicking is acceptable
func MustConnect() *Client {
	client, err := NewClient(FromEnv())
	if err != nil {
		panic(err)
	}
	return client
}

// Ping tests the connection to Elasticsearch
func (c *Client) Ping(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	c.mutex.RLock()
	client := c.client
	c.mutex.RUnlock()

	if client == nil {
		return fmt.Errorf("client not connected")
	}

	res, err := client.Info(client.Info.WithContext(ctx))
	if err != nil {
		c.mutex.Lock()
		c.isConnected = false
		c.mutex.Unlock()
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	if res.IsError() {
		c.mutex.Lock()
		c.isConnected = false
		c.mutex.Unlock()
		return fmt.Errorf("ping failed: %s", res.String())
	}

	c.mutex.Lock()
	c.isConnected = true
	c.mutex.Unlock()

	return nil
}

// Stats returns connection statistics
func (c *Client) Stats() ConnectionStats {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return ConnectionStats{
		IsConnected:   c.isConnected,
		Reconnects:    c.reconnectCount,
		LastReconnect: c.lastReconnect,
	}
}

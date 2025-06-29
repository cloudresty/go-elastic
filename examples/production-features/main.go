package main

import (
	"context"
	"log"
	"time"

	"github.com/cloudresty/emit"
	elastic "github.com/cloudresty/go-elastic"
)

func main() {
	emit.Info.Msg("=== Production Features Demo ===")

	// Create clients with production configuration
	paymentsClient, err := elastic.NewClient(
		elastic.FromEnvWithPrefix("PAYMENTS_"),
		elastic.WithConnectionName("payments-cluster"),
	)
	if err != nil {
		log.Fatalf("Failed to create payments client: %v", err)
	}

	ordersClient, err := elastic.NewClient(
		elastic.FromEnvWithPrefix("ORDERS_"),
		elastic.WithConnectionName("orders-cluster"),
	)
	if err != nil {
		log.Fatalf("Failed to create orders client: %v", err)
	}

	// Set up graceful shutdown manager
	shutdownManager := elastic.NewShutdownManager(&elastic.ShutdownConfig{
		Timeout:          30 * time.Second,
		GracePeriod:      5 * time.Second,
		ForceKillTimeout: 10 * time.Second,
	})

	// Register clients for graceful shutdown
	shutdownManager.Register(paymentsClient, ordersClient)
	shutdownManager.SetupSignalHandler()

	// Start background workers in separate goroutines
	go func() {
		backgroundHealthChecker(shutdownManager.Context(), paymentsClient, "payments")
	}()

	go func() {
		backgroundHealthChecker(shutdownManager.Context(), ordersClient, "orders")
	}()

	// Demonstrate multi-client operations
	demonstrateMultiClientOperations(paymentsClient, ordersClient)

	// Wait for shutdown signal
	emit.Info.Msg("Application started. Press Ctrl+C to shutdown gracefully.")
	shutdownManager.Wait()

	emit.Info.Msg("Application shutdown completed")
}

func demonstrateMultiClientOperations(paymentsClient, ordersClient *elastic.Client) {
	emit.Info.Msg("Demonstrating multi-client operations")

	ctx := context.Background()

	// --- Demonstrating the improved connection API ---
	// 1. Get the client's static, configured name for logging
	paymentsName := paymentsClient.Name()
	ordersName := ordersClient.Name()
	emit.Info.StructuredFields("Working with clients",
		emit.ZString("payments_client", paymentsName),
		emit.ZString("orders_client", ordersName))

	// 2. Check real-time health/connectivity by performing an action
	if err := paymentsClient.Ping(ctx); err != nil {
		emit.Warn.StructuredFields("Client is not responsive",
			emit.ZString("client", paymentsName),
			emit.ZString("error", err.Error()))
	} else {
		emit.Info.StructuredFields("Client is responsive",
			emit.ZString("client", paymentsName))
	}

	if err := ordersClient.Ping(ctx); err != nil {
		emit.Warn.StructuredFields("Client is not responsive",
			emit.ZString("client", ordersName),
			emit.ZString("error", err.Error()))
	} else {
		emit.Info.StructuredFields("Client is responsive",
			emit.ZString("client", ordersName))
	}

	// 3. Get the dynamic, runtime metrics
	paymentsStats := paymentsClient.Stats()
	ordersStats := ordersClient.Stats()

	emit.Info.StructuredFields("Connection statistics",
		emit.ZString("payments_client", paymentsName),
		emit.ZInt64("payments_reconnects", paymentsStats.Reconnects),
		emit.ZString("orders_client", ordersName),
		emit.ZInt64("orders_reconnects", ordersStats.Reconnects))

	// Index a payment document
	paymentDoc := map[string]any{
		"amount":      100.50,
		"currency":    "USD",
		"customer_id": "cust_123",
		"status":      "completed",
		"timestamp":   time.Now(),
	}

	paymentResp, err := paymentsClient.Documents().Create(ctx, "payments", paymentDoc)
	if err != nil {
		emit.Error.StructuredFields("Failed to index payment",
			emit.ZString("error", err.Error()))
	} else {
		emit.Info.StructuredFields("Payment indexed successfully",
			emit.ZString("payment_id", paymentResp.ID))
	}

	// Index an order document
	orderDoc := map[string]any{
		"customer_id": "cust_123",
		"items": []map[string]any{
			{
				"product_id": "prod_456",
				"quantity":   2,
				"price":      50.25,
			},
		},
		"total":     100.50,
		"status":    "processing",
		"timestamp": time.Now(),
	}

	orderResp, err := ordersClient.Documents().Create(ctx, "orders", orderDoc)
	if err != nil {
		emit.Error.StructuredFields("Failed to index order",
			emit.ZString("error", err.Error()))
	} else {
		emit.Info.StructuredFields("Order indexed successfully",
			emit.ZString("order_id", orderResp.ID))
	}

	// Demonstrate health checks
	paymentsHealthy := paymentsClient.Ping(ctx) == nil
	ordersHealthy := ordersClient.Ping(ctx) == nil

	// Get detailed cluster health information
	var paymentsClusterName, ordersClusterName string
	if paymentsHealth, err := paymentsClient.Cluster().Health(ctx); err == nil {
		paymentsClusterName = paymentsHealth.ClusterName
	}
	if ordersHealth, err := ordersClient.Cluster().Health(ctx); err == nil {
		ordersClusterName = ordersHealth.ClusterName
	}

	emit.Info.StructuredFields("Health check results",
		emit.ZBool("payments_healthy", paymentsHealthy),
		emit.ZBool("orders_healthy", ordersHealthy),
		emit.ZString("payments_cluster", paymentsClusterName),
		emit.ZString("orders_cluster", ordersClusterName))
}

func backgroundHealthChecker(ctx context.Context, client *elastic.Client, serviceName string) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	emit.Info.StructuredFields("Starting background health checker",
		emit.ZString("service", serviceName))

	for {
		select {
		case <-ctx.Done():
			emit.Info.StructuredFields("Health checker shutting down",
				emit.ZString("service", serviceName))
			return
		case <-ticker.C:
			isHealthy := client.Ping(ctx) == nil
			stats := client.Stats()
			clientName := client.Name()
			if !isHealthy {
				emit.Warn.StructuredFields("Service health check failed",
					emit.ZString("service", serviceName),
					emit.ZString("client_name", clientName),
					emit.ZInt64("reconnects", stats.Reconnects))
			} else {
				emit.Debug.StructuredFields("Service health check passed",
					emit.ZString("service", serviceName),
					emit.ZString("client_name", clientName),
					emit.ZBool("connected", stats.IsConnected))
			}
		}
	}
}

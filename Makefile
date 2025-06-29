.PHONY: help test test-integration build clean lint fmt vet tidy examples docker-elasticsearch

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
fmt: ## Format Go code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

lint: ## Run golangci-lint
	golangci-lint run

tidy: ## Tidy go modules
	go mod tidy
	go mod verify

# Testing variables (can be overridden)
TEST_ELASTICSEARCH_HOST ?= localhost
TEST_ELASTICSEARCH_PORT ?= 9200
TEST_ELASTICSEARCH_USERNAME ?=
TEST_ELASTICSEARCH_PASSWORD ?=

# Testing
test: ## Run unit tests
	go test -v -race -short $$(go list ./... | grep -v '/examples/')

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	@mkdir -p test-results
	@if [ -n "$(TEST_ELASTICSEARCH_USERNAME)" ] && [ -n "$(TEST_ELASTICSEARCH_PASSWORD)" ]; then \
		echo "Using custom credentials: $(TEST_ELASTICSEARCH_USERNAME)@$(TEST_ELASTICSEARCH_HOST):$(TEST_ELASTICSEARCH_PORT)"; \
	else \
		echo "Using default connection: $(TEST_ELASTICSEARCH_HOST):$(TEST_ELASTICSEARCH_PORT)"; \
	fi
	@echo "Output will be saved to test-results/integration.txt"
	@env -i PATH="$(PATH)" HOME="$(HOME)" \
	ELASTICSEARCH_HOST=$(TEST_ELASTICSEARCH_HOST) \
	ELASTICSEARCH_PORT=$(TEST_ELASTICSEARCH_PORT) \
	ELASTICSEARCH_USERNAME=$${TEST_ELASTICSEARCH_USERNAME} \
	ELASTICSEARCH_PASSWORD=$${TEST_ELASTICSEARCH_PASSWORD} \
	go test -v -race $$(go list ./... | grep -v '/examples/') > test-results/integration.txt 2>&1
	@echo "Test completed. Check test-results/integration.txt for results."

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@env -i PATH="$(PATH)" HOME="$(HOME)" \
	ELASTICSEARCH_HOST=$(TEST_ELASTICSEARCH_HOST) \
	ELASTICSEARCH_PORT=$(TEST_ELASTICSEARCH_PORT) \
	ELASTICSEARCH_USERNAME=$${TEST_ELASTICSEARCH_USERNAME} \
	ELASTICSEARCH_PASSWORD=$${TEST_ELASTICSEARCH_PASSWORD} \
	go test -v -race -coverprofile=coverage.out $$(go list ./... | grep -v '/examples/')
	go tool cover -html=coverage.out -o coverage.html

# Build
build: ## Build examples
	@mkdir -p bin/
	go build -o bin/basic-client examples/basic-client/main.go
	go build -o bin/env-config examples/env-config/main.go
	go build -o bin/search-demo examples/search-demo/main.go
	go build -o bin/production-features examples/production-features/main.go
	go build -o bin/bulk-operations examples/bulk-operations/main.go
	go build -o bin/reconnection-test examples/reconnection-test/main.go
	go build -o bin/id-demo examples/id-demo/main.go

clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf test-results/
	rm -f coverage.out coverage.html

# Docker
docker-elasticsearch: ## Start Elasticsearch in Docker
	docker run -d --name elasticsearch-dev \
		-p 9200:9200 \
		-p 9300:9300 \
		-e "discovery.type=single-node" \
		-e "xpack.security.enabled=false" \
		-e "xpack.security.enrollment.enabled=false" \
		docker.elastic.co/elasticsearch/elasticsearch:9.0.1

docker-elasticsearch-auth: ## Start Elasticsearch with authentication in Docker
	docker run -d --name elasticsearch-dev-auth \
		-p 9200:9200 \
		-p 9300:9300 \
		-e "discovery.type=single-node" \
		-e "xpack.security.enabled=true" \
		-e "ELASTIC_PASSWORD=password" \
		-e "xpack.security.enrollment.enabled=false" \
		docker.elastic.co/elasticsearch/elasticsearch:9.0.1

docker-stop: ## Stop Elasticsearch Docker container
	docker stop elasticsearch-dev || true
	docker rm elasticsearch-dev || true
	docker stop elasticsearch-dev-auth || true
	docker rm elasticsearch-dev-auth || true

# Examples
run-basic-client: ## Run basic client example
	go run examples/basic-client/main.go

run-env-config: ## Run environment config example
	go run examples/env-config/main.go

run-search-demo: ## Run search demonstration example
	go run examples/search-demo/main.go

run-production-features: ## Run production features example
	go run examples/production-features/main.go

run-bulk-operations: ## Run bulk operations example
	go run examples/bulk-operations/main.go

run-reconnection-test: ## Run reconnection test example
	go run examples/reconnection-test/main.go

run-id-demo: ## Run ID generation demo
	go run examples/id-demo/main.go

# CI
ci: tidy fmt vet test ## Run CI pipeline

# Install tools
install-tools: ## Install development tools
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Security
security: ## Run security checks
	gosec ./...

# All tests
test-all: test test-integration ## Run all tests

test-results: ## Show latest test results
	@echo "=== Latest Test Results ==="
	@if [ -f test-results/integration.txt ]; then \
		echo "--- Integration Tests (test-results/integration.txt) ---"; \
		tail -n 20 test-results/integration.txt; \
		echo ""; \
	fi
	@echo "Use 'cat test-results/*.txt' to see full results"

test-status: ## Check if tests passed or failed
	@echo "=== Test Status Summary ==="
	@if [ -f test-results/integration.txt ]; then \
		if grep -q "PASS" test-results/integration.txt && ! grep -q "FAIL" test-results/integration.txt; then \
			echo "✅ Integration tests: PASSED"; \
		else \
			echo "❌ Integration tests: FAILED"; \
		fi; \
	else \
		echo "ℹ️  No integration test results found. Run 'make test-integration' first."; \
	fi

# Docker workflow
docker-elasticsearch-full: docker-elasticsearch ## Start Elasticsearch

# CI setup for testing
ci-setup-elasticsearch: ## Setup Elasticsearch for CI
	@echo "Setting up Elasticsearch for CI..."
	@$(MAKE) docker-elasticsearch
	@echo "Waiting for Elasticsearch to be ready..."
	@for i in $$(seq 1 60); do \
		if curl -s http://localhost:9200/_cluster/health >/dev/null 2>&1; then \
			echo "Elasticsearch is ready!"; \
			break; \
		fi; \
		echo "Waiting... ($$i/60)"; \
		sleep 2; \
	done
	@echo "Elasticsearch setup complete"

ci-test: ## Run tests in CI environment
	@echo "Running CI tests..."
	@$(MAKE) test-integration
	@$(MAKE) test-status

ci-cleanup: ## Cleanup CI environment
	@echo "Cleaning up CI environment..."
	@$(MAKE) docker-stop clean
	@echo "Cleanup complete"

# Complete local testing workflow
test-full-local: ## Complete local test workflow (setup + test + cleanup)
	@echo "=== Full Local Test Workflow ==="
	@$(MAKE) ci-setup-elasticsearch
	@$(MAKE) ci-test
	@echo "=== Test Workflow Complete ==="
	@echo "Tip: Run 'make ci-cleanup' to stop Elasticsearch and clean up"

test-clean: ## Clean test results and cache
	@echo "Cleaning test results and cache..."
	@rm -rf test-results/
	@go clean -testcache
	@echo "Test cleanup complete"

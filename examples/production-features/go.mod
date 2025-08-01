module github.com/cloudresty/go-elastic/examples/production-features

go 1.24.1

toolchain go1.24.5

require (
	github.com/cloudresty/emit v1.2.5
	github.com/cloudresty/go-elastic v0.0.0-00010101000000-000000000000
)

require (
	github.com/cloudresty/go-env v1.0.1 // indirect
	github.com/cloudresty/ulid v1.2.1 // indirect
	github.com/elastic/elastic-transport-go/v8 v8.7.0 // indirect
	github.com/elastic/go-elasticsearch/v9 v9.0.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.35.0 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	go.opentelemetry.io/otel/trace v1.35.0 // indirect
)

replace github.com/cloudresty/go-elastic => ../..

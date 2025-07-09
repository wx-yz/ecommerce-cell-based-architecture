module github.com/wx-yz/ecommerce-cell-based-architecture/frontend

go 1.21

require (
	cloud.google.com/go/profiler v0.4.0
	github.com/golang/protobuf v1.5.3
	github.com/google/uuid v1.4.0
	github.com/gorilla/mux v1.8.1
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.9.3
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.46.1
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.46.1
	go.opentelemetry.io/otel v1.21.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.21.0
	go.opentelemetry.io/otel/propagation v1.21.0
	go.opentelemetry.io/otel/sdk v1.21.0
	google.golang.org/grpc v1.60.1
	google.golang.org/protobuf v1.31.0
)
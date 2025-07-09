package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	port = 3550
)

type Product struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Picture     string   `json:"picture"`
	PriceUSD    Money    `json:"price_usd"`
	Categories  []string `json:"categories"`
}

type Money struct {
	CurrencyCode string `json:"currency_code"`
	Units        int64  `json:"units"`
	Nanos        int32  `json:"nanos"`
}

var (
	catalog []Product
	mu      sync.RWMutex
)

type productCatalogService struct{}

func main() {
	log := logrus.New()
	log.Level = logrus.InfoLevel
	log.Formatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	}

	log.Info("starting product catalog service")

	// Load product catalog
	if err := loadCatalog(); err != nil {
		log.Fatalf("failed to load catalog: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var srv *grpc.Server
	srv = grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	// Register service here when proto bindings are available

	log.Infof("product catalog service listening on port %d", port)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func loadCatalog() error {
	catalogJSON, err := ioutil.ReadFile("products.json")
	if err != nil {
		return fmt.Errorf("failed to read products.json: %v", err)
	}

	var products []Product
	if err := json.Unmarshal(catalogJSON, &products); err != nil {
		return fmt.Errorf("failed to unmarshal products.json: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	catalog = products
	return nil
}

func (s *productCatalogService) ListProducts(ctx context.Context) ([]Product, error) {
	mu.RLock()
	defer mu.RUnlock()
	return catalog, nil
}

func (s *productCatalogService) GetProduct(ctx context.Context, id string) (*Product, error) {
	mu.RLock()
	defer mu.RUnlock()
	for _, p := range catalog {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "product not found")
}

func (s *productCatalogService) SearchProducts(ctx context.Context, query string) ([]Product, error) {
	mu.RLock()
	defer mu.RUnlock()
	var results []Product
	for _, p := range catalog {
		if strings.Contains(strings.ToLower(p.Name), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(p.Description), strings.ToLower(query)) {
			results = append(results, p)
		}
	}
	return results, nil
}
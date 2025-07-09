package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	port = 50051
)

type Money struct {
	CurrencyCode string `json:"currency_code"`
	Units        int64  `json:"units"`
	Nanos        int32  `json:"nanos"`
}

type Address struct {
	StreetAddress string `json:"street_address"`
	City          string `json:"city"`
	State         string `json:"state"`
	Country       string `json:"country"`
	ZipCode       int32  `json:"zip_code"`
}

type CartItem struct {
	ProductID string `json:"product_id"`
	Quantity  int32  `json:"quantity"`
}

type shippingService struct{}

func main() {
	logger := logrus.New()
	logger.Level = logrus.InfoLevel
	logger.Formatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	}

	logger.Info("Starting shipping service")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	// Register service when proto bindings are available
	// pb.RegisterShippingServiceServer(server, &shippingService{})

	logger.Infof("Shipping service listening on port %d", port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *shippingService) GetQuote(ctx context.Context, address Address, items []CartItem) (*Money, error) {
	logger := logrus.WithFields(logrus.Fields{
		"method": "GetQuote",
		"city":   address.City,
		"state":  address.State,
		"items":  len(items),
	})

	logger.Info("Calculating shipping quote")

	// Calculate shipping cost based on items and destination
	baseCost := 5.0 // Base shipping cost
	itemCost := float64(len(items)) * 1.5 // Additional cost per item

	// Distance-based pricing (simplified)
	distanceMultiplier := 1.0
	switch address.State {
	case "CA", "OR", "WA":
		distanceMultiplier = 1.0 // West Coast
	case "NY", "NJ", "CT":
		distanceMultiplier = 1.2 // East Coast
	case "TX", "FL":
		distanceMultiplier = 1.1 // South
	default:
		distanceMultiplier = 1.15 // Midwest/Other
	}

	// International shipping
	if address.Country != "United States" && address.Country != "US" {
		distanceMultiplier *= 2.0
	}

	totalCost := (baseCost + itemCost) * distanceMultiplier

	// Convert to Money structure
	units := int64(math.Floor(totalCost))
	nanos := int32((totalCost - float64(units)) * 1e9)

	result := &Money{
		CurrencyCode: "USD",
		Units:        units,
		Nanos:        nanos,
	}

	logger.WithFields(logrus.Fields{
		"cost_usd": totalCost,
		"units":    units,
		"nanos":    nanos,
	}).Info("Shipping quote calculated")

	return result, nil
}

func (s *shippingService) ShipOrder(ctx context.Context, address Address, items []CartItem) (string, error) {
	logger := logrus.WithFields(logrus.Fields{
		"method": "ShipOrder",
		"city":   address.City,
		"state":  address.State,
		"items":  len(items),
	})

	logger.Info("Processing shipping order")

	// Validate address
	if address.StreetAddress == "" || address.City == "" || address.State == "" {
		return "", status.Errorf(codes.InvalidArgument, "incomplete address")
	}

	// Generate tracking ID
	trackingID := generateTrackingID()

	// Simulate shipping processing
	time.Sleep(50 * time.Millisecond)

	logger.WithFields(logrus.Fields{
		"tracking_id": trackingID,
	}).Info("Order shipped successfully")

	return trackingID, nil
}

func generateTrackingID() string {
	// Generate a realistic tracking ID
	uuid := uuid.New()
	return fmt.Sprintf("SH-%s", uuid.String()[:8])
}

// Health check endpoint
func (s *shippingService) healthCheck() error {
	return nil
}
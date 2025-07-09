package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/google/uuid"
)

type CheckoutService struct {
	port string
}

// Request/Response models
type PlaceOrderRequest struct {
	UserID       string                 `json:"user_id"`
	UserCurrency string                 `json:"user_currency"`
	Address      Address                `json:"address"`
	Email        string                 `json:"email"`
	CreditCard   CreditCardInfo         `json:"credit_card"`
	Cart         []CartItem             `json:"cart"`
}

type Address struct {
	StreetAddress string `json:"street_address"`
	City          string `json:"city"`
	State         string `json:"state"`
	Country       string `json:"country"`
	ZipCode       string `json:"zip_code"`
}

type CreditCardInfo struct {
	Number          string `json:"number"`
	CVV             int32  `json:"cvv"`
	ExpirationYear  int32  `json:"expiration_year"`
	ExpirationMonth int32  `json:"expiration_month"`
}

type CartItem struct {
	ProductID string `json:"product_id"`
	Quantity  int32  `json:"quantity"`
}

type OrderResult struct {
	OrderID          string    `json:"order_id"`
	ShippingTrackingID string  `json:"shipping_tracking_id"`
	ShippingCost     Money     `json:"shipping_cost"`
	ShippingAddress  Address   `json:"shipping_address"`
	Items            []OrderItem `json:"items"`
	CreatedAt        time.Time `json:"created_at"`
	Status           string    `json:"status"`
}

type OrderItem struct {
	Item CartItem `json:"item"`
	Cost Money    `json:"cost"`
}

type Money struct {
	CurrencyCode string `json:"currency_code"`
	Units        int64  `json:"units"`
	Nanos        int32  `json:"nanos"`
}

func NewCheckoutService() *CheckoutService {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	return &CheckoutService{
		port: port,
	}
}

func (s *CheckoutService) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var req PlaceOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Generate order ID
	orderID := uuid.New().String()
	
	// Mock processing - in real implementation this would:
	// 1. Get quote from shipping service
	// 2. Charge credit card via payment service
	// 3. Ship order via shipping service
	// 4. Send confirmation email
	
	// For now, create a mock response
	result := OrderResult{
		OrderID:          orderID,
		ShippingTrackingID: "SP-" + orderID[:8],
		ShippingCost: Money{
			CurrencyCode: req.UserCurrency,
			Units:        5,
			Nanos:        0,
		},
		ShippingAddress: req.Address,
		Items:          []OrderItem{},
		CreatedAt:      time.Now(),
		Status:         "confirmed",
	}

	// Convert cart items to order items
	for _, item := range req.Cart {
		orderItem := OrderItem{
			Item: item,
			Cost: Money{
				CurrencyCode: req.UserCurrency,
				Units:        10, // Mock price
				Nanos:        0,
			},
		}
		result.Items = append(result.Items, orderItem)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *CheckoutService) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func (s *CheckoutService) Start() error {
	r := mux.NewRouter()
	
	// Routes
	r.HandleFunc("/api/checkout/place-order", s.PlaceOrder).Methods("POST")
	r.HandleFunc("/health", s.Health).Methods("GET")
	
	log.Printf("Checkout service listening on port %s", s.port)
	return http.ListenAndServe(":"+s.port, r)
}

func main() {
	service := NewCheckoutService()
	log.Fatal(service.Start())
}
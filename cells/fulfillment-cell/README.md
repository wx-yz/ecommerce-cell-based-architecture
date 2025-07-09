# Fulfillment Cell

## Overview
The Fulfillment Cell handles shipping and customer communications.

## Components
- **shippingservice**: Go service for shipping calculations
- **emailservice**: Python service for email notifications

## Responsibilities
- Shipping cost calculations
- Shipping address validation
- Order confirmation emails
- Customer notifications

## API Endpoints
### Shipping Service
- `GetQuote(address, items)`: Get shipping quote
- `ShipOrder(address, items)`: Ship order

### Email Service
- `SendOrderConfirmation(email, order)`: Send confirmation email

## Dependencies
- No external cell dependencies

## Deployment Configuration
- **Runtime**: Go (shipping), Python (email)
- **Port**: 50051 (shipping), 8080 (email)

## Choreo Deployment
Deployed as a single Choreo project with two components.
# Order Cell

## Overview
The Order Cell handles order processing and orchestration.

## Components
- **checkoutservice**: Go service for order orchestration

## Responsibilities
- Order processing workflow
- Payment coordination
- Shipping coordination
- Email notification coordination

## API Endpoints
- `PlaceOrder(user_id, currency, address, email, credit_card)`: Process order

## Dependencies
- Cart Cell: Get cart items
- Payment Cell: Process payment
- Fulfillment Cell: Handle shipping and email
- Platform Cell: Currency conversion

## Deployment Configuration
- **Runtime**: Go
- **Port**: 5050

## Choreo Deployment
Deployed as a single Choreo project with checkout service component.
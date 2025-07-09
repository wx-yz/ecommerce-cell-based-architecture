# Payment Cell

## Overview
The Payment Cell handles payment processing.

## Components
- **paymentservice**: Node.js service for payment processing

## Responsibilities
- Credit card processing
- Payment validation
- Transaction management

## API Endpoints
- `Charge(amount, credit_card)`: Process payment

## Dependencies
- No external cell dependencies
- Mock payment processing

## Deployment Configuration
- **Runtime**: Node.js
- **Port**: 50051

## Choreo Deployment
Deployed as a single Choreo project with payment service component.
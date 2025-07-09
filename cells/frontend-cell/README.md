# Frontend Cell

## Overview
The Frontend Cell represents the User Experience domain in our cell-based architecture. It contains the web frontend and load generation components.

## Components
- **frontend**: Go-based web application serving the e-commerce UI
- **loadgenerator**: Python/Locust-based load testing tool

## Responsibilities
- User interface for browsing products
- Shopping cart management interface
- Checkout process interface
- User session management
- Load testing and performance validation

## API Dependencies
This cell communicates with the following cells:
- **Product Cell**: Product catalog and recommendations
- **Cart Cell**: Shopping cart operations
- **Order Cell**: Checkout and order placement
- **Payment Cell**: Payment processing
- **Marketing Cell**: Advertisement display
- **Platform Cell**: Currency conversion

## Deployment Configuration
- **Runtime**: Node.js for web server
- **Port**: 8080 (configurable via PORT environment variable)
- **Static Assets**: Served from `/src/frontend/static/`
- **Templates**: Go templates in `/src/frontend/templates/`

## Environment Variables
- `PORT`: Server port (default: 8080)
- `PRODUCT_CATALOG_SERVICE_ADDR`: Product service endpoint
- `CART_SERVICE_ADDR`: Cart service endpoint
- `RECOMMENDATION_SERVICE_ADDR`: Recommendation service endpoint
- `CHECKOUT_SERVICE_ADDR`: Checkout service endpoint
- `CURRENCY_SERVICE_ADDR`: Currency service endpoint
- `AD_SERVICE_ADDR`: Ad service endpoint

## Choreo Connections
The frontend cell requires connections to:
1. Product Cell API (REST)
2. Cart Cell API (REST)
3. Order Cell API (REST)
4. Platform Cell API (REST)
5. Marketing Cell API (REST)

## Files in this Cell
- `frontend/`: Web application source code
- `loadgenerator/`: Load testing configuration
- `choreo/`: Choreo-specific configuration files
- `README.md`: This file
# Product Cell

## Overview
The Product Cell represents the Product Catalog and Recommendations domain. It manages product information and provides personalized recommendations.

## Components
- **productcatalogservice**: Go-based service for product catalog management
- **recommendationservice**: Python-based service for product recommendations

## Responsibilities
- Product catalog management
- Product search functionality
- Product recommendations based on user behavior
- Product metadata management

## API Endpoints
### Product Catalog Service
- `ListProducts()`: Get all products
- `GetProduct(id)`: Get specific product details
- `SearchProducts(query)`: Search products by query

### Recommendation Service
- `ListRecommendations(user_id, product_ids)`: Get recommended products

## Dependencies
- No external cell dependencies
- Uses internal product data
- Recommendation engine uses ML models

## Deployment Configuration
- **Runtime**: Go for catalog service, Python for recommendation service
- **Port**: 3550 (catalog), 8080 (recommendations)
- **Data**: JSON file for product catalog, ML models for recommendations

## Environment Variables
- `PORT`: Server port
- `PRODUCT_CATALOG_SERVICE_ADDR`: Internal catalog service address
- `DISABLE_TRACING`: Disable OpenTelemetry tracing
- `DISABLE_PROFILER`: Disable profiler

## Choreo Deployment
This cell is deployed as a single Choreo project with two components:
1. Product Catalog Service (Go)
2. Recommendation Service (Python)

## Files in this Cell
- `productcatalogservice/`: Product catalog service source code
- `recommendationservice/`: Recommendation service source code
- `choreo/`: Choreo-specific configuration files
- `README.md`: This file
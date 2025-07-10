# Choreo Deployment Strategy

## Overview

This document outlines the deployment strategy for the cell-based e-commerce application on the Choreo platform.

## Choreo Platform Mapping

### Cell-to-Project Mapping

Each cell maps to a separate Choreo project, following the cell-based architecture principle:
| Cell | Choreo Project | Components |
|------|----------------|-----------|
| Frontend Cell | `ecommerce-frontend` | frontend, loadgenerator |
| Product Cell | `ecommerce-product` | productcatalogservice, recommendationservice |
| Cart Cell | `ecommerce-cart` | cartservice |
| Order Cell | `ecommerce-order` | checkoutservice |
| Payment Cell | `ecommerce-payment` | paymentservice |
| Fulfillment Cell | `ecommerce-fulfillment` | shippingservice, emailservice |
| Marketing Cell | `ecommerce-marketing` | adservice |
| Platform Cell | `ecommerce-platform` | currencyservice |

### Project Configuration\nEach project will be configured with:
- **Region**: US (primary), EU (secondary)
- **Environment**: Development, Testing, Production
- **Git Repository**: This monorepo with cell-specific paths
- **Build Pipeline**: Automated CI/CD per cell

## Deployment Order

### Phase 1: Foundation Services
1. **Platform Cell** (currencyservice)
   - No dependencies\n
   - Provides utility services
   - Deploy first to ensure availability
2. **Product Cell** (productcatalogservice, recommendationservice)
   - No dependencies
   - Self-contained domain
   - Critical for frontend functionality
3. **Cart Cell** (cartservice)
   - No dependencies
   - Self-contained domain
   - Requires Redis configuration

### Phase 2: Business Logic Services
4. **Payment Cell** (paymentservice)
   - No dependencies\n
   - Self-contained domain
   - Mock payment processing
5. **Fulfillment Cell** (shippingservice, emailservice)
   - No dependencies
   - Self-contained domain
   - Email templates configuration
6. **Marketing Cell** (adservice)
   - No dependencies
   - Self-contained domain
   - Ad content configuration

### Phase 3: Orchestration Services
7. **Order Cell** (checkoutservice)
   - Depends on: Cart, Payment, Fulfillment, Platform
   - Orchestrates business workflow
   - Deploy after dependencies are stable

### Phase 4: User Interface
8. **Frontend Cell** (frontend, loadgenerator)
   - Depends on: All other cells
   - User-facing interface
   - Deploy last to ensure all APIs are available

## Choreo Project Setup

### 1. Create Projects
```bash
# Create each project using Choreo CLI or API
for project in frontend product cart order payment fulfillment marketing platform; do
    choreo project create \"ecommerce-${project}\" \
        --description \"E-commerce ${project} cell\" \
        --region \"US\" \
        --repository \"https://github.com/wx-yz/ecommerce-cell-based-architecture.git\" \
        --branch \"main\"
done
```
### 2. Configure Components
Each project will have components configured with:
- **Source Path**: `cells/{cell-name}/{component-name}/`
- **Build Pack**: Language-specific (Go, Node.js, Python, Java, .NET)
- **Port**: Component-specific port configuration
- **Environment Variables**: Cell-specific configuration
### 3. Environment Configuration\nEach project will have three environments:
- **Development**: For feature development and testing
- **Testing**: For integration testing and QA
- **Production**: For live traffic

## Choreo Connections
### Connection Types
#### 1. Service-to-Service Connections
- **Type**: Internal API connections
- **Protocol**: HTTP/gRPC
- **Authentication**: Service-to-service tokens
- **Load Balancing**: Choreo managed
#### 2. External Connections
- **Type**: External API connections
- **Usage**: Platform cell → Currency APIs
- **Authentication**: API keys/tokens
- **Rate Limiting**: Choreo managed

### Connection Configuration
#### Frontend Cell Connections
```yaml
connections:
  - name: \"product-api\"
    target: \"ecommerce-product\"
    protocol: \"http\"
    path: \"/api/v1/products\"
  - name: \"cart-api\"
    target: \"ecommerce-cart\"
    protocol: \"http\"
    path: \"/api/v1/cart\"
  - name: \"order-api\"
    target: \"ecommerce-order\"
    protocol: \"http\"
    path: \"/api/v1/checkout\"
  - name: \"platform-api\"
    target: \"ecommerce-platform\"
    protocol: \"http\"
    path: \"/api/v1/currency\"
  - name: \"marketing-api\"
    target: \"ecommerce-marketing\"
    protocol: \"http\"
    path: \"/api/v1/ads\"
```
#### Order Cell Connections
```yaml
connections:
  - name: \"cart-api\"
    target: \"ecommerce-cart\"
    protocol: \"grpc\"
    port: 7070
  - name: \"payment-api\"
    target: \"ecommerce-payment\"
    protocol: \"grpc\"
    port: 50051
  - name: \"fulfillment-shipping-api\"
    target: \"ecommerce-fulfillment\"
    protocol: \"grpc\"
    port: 50051
  - name: \"fulfillment-email-api\"
    target: \"ecommerce-fulfillment\"
    protocol: \"grpc\"
    port: 8080
  - name: \"platform-api\"
    target: \"ecommerce-platform\"
    protocol: \"grpc\"
    port: 7000
```
## Environment Variables
### Frontend Cell
```yaml
env:
  - name: \"PORT\"
    value: \"8080\"
  - name: \"PRODUCT_CATALOG_SERVICE_ADDR\"
    value: \"${CHOREO_CONNECTION_product-api}\"
  - name: \"CART_SERVICE_ADDR\"
    value: \"${CHOREO_CONNECTION_cart-api}\"
  - name: \"RECOMMENDATION_SERVICE_ADDR\"
    value: \"${CHOREO_CONNECTION_product-api}\"
  - name: \"CHECKOUT_SERVICE_ADDR\"
    value: \"${CHOREO_CONNECTION_order-api}\"
  - name: \"CURRENCY_SERVICE_ADDR\"
    value: \"${CHOREO_CONNECTION_platform-api}\"
  - name: \"AD_SERVICE_ADDR\"
    value: \"${CHOREO_CONNECTION_marketing-api}\"
```
### Order Cell
```yaml
env:
  - name: \"PORT\"
    value: \"5050\"
  - name: \"CART_SERVICE_ADDR\"
    value: \"${CHOREO_CONNECTION_cart-api}\"
  - name: \"PAYMENT_SERVICE_ADDR\"
    value: \"${CHOREO_CONNECTION_payment-api}\"
  - name: \"SHIPPING_SERVICE_ADDR\"
    value: \"${CHOREO_CONNECTION_fulfillment-shipping-api}\"
  - name: \"EMAIL_SERVICE_ADDR\"
    value: \"${CHOREO_CONNECTION_fulfillment-email-api}\"
  - name: \"CURRENCY_SERVICE_ADDR\"
    value: \"${CHOREO_CONNECTION_platform-api}\"
```
## Configurations and Secrets
### Configuration Management
- **Application Config**: Use Choreo config maps
- **Environment-specific**: Different values per environment
- **Secret Management**: Use Choreo secrets for sensitive data
### Example Configurations
#### Cart Cell Redis Configuration
```yaml
configurations:
  - name: \"redis-config\"
    type: \"config-map\"
    mount_type: \"env variable\"
    env_variables: |
      {
        \"REDIS_ADDR\": \"redis:6379\",
        \"REDIS_DB\": \"0\"
      }
```
#### Platform Cell Currency API Secret
```yaml
configurations:
  - name: \"currency-api-secret\"
    type: \"secret\"
    mount_type: \"env variable\"
    env_variables: |
      {
        \"CURRENCY_API_KEY\": \"${CURRENCY_API_KEY}\"
      }
```
## Monitoring and Observability
### Choreo Built-in Monitoring
- **Metrics**: Request rates, response times, error rates
- **Logs**: Application logs with cell context
- **Tracing**: Distributed tracing across cells
- **Health Checks**: Endpoint health monitoring
### Custom Monitoring
- **Business Metrics**: Order completion rates, cart abandonment
- **Performance Metrics**: Page load times, API response times
- **Error Tracking**: Application errors with cell context
## Security Configuration
### Authentication
- **Service-to-Service**: Choreo managed service tokens
- **External APIs**: API keys in secrets
- **User Authentication**: Session-based (frontend cell)
### Authorization\n- **API Gateway**: Choreo managed API gateway
- **Rate Limiting**: Per-cell rate limiting
- **Access Control**: Role-based access per cell
## Scaling Configuration
### Auto-scaling
- **Horizontal**: Scale replicas based on CPU/memory
- **Vertical**: Scale resources based on demand
- **Per-cell**: Independent scaling per cell
### Resource Limits
```yaml
resources:
  limits:
    cpu: \"1000m\"
    memory: \"512Mi\"
  requests:
    cpu: \"100m\"
    memory: \"128Mi\"
```
## Disaster Recovery
### Backup Strategy
- **Configuration Backup**: Export Choreo configurations
- **Data Backup**: Per-cell data backup strategy
- **Cross-region**: Deploy to multiple regions
### Recovery Process
1. **Identify Failed Cell**: Monitor alerts
2. **Isolate Impact**: Cell boundaries limit blast radius
3. **Restore Service**: Redeploy affected cell
4. **Verify Integration**: Test inter-cell connections
## CI/CD Pipeline
### Build Pipeline
```yaml
stages:
  - build:
      - compile/build application
      - run unit tests
      - security scanning
  - deploy:
      - deploy to development
      - integration tests
      - deploy to testing
      - acceptance tests
      - deploy to production
```
### Deployment Strategy
- **Blue-Green**: Zero-downtime deployments
- **Canary**: Gradual rollout for production
- **Rollback**: Quick rollback on issues
## Troubleshooting
### Common Issues
1. **Connection Failures**: Check Choreo connections
2. **Service Discovery**: Verify service endpoints
3. **Configuration Issues**: Check environment variables
4. **Resource Limits**: Monitor resource usage
### Debug Process
1. **Check Logs**: View application logs in Choreo
2. **Monitor Metrics**: Check performance metrics
3. **Trace Requests**: Use distributed tracing
4. **Test Connections**: Verify inter-cell connectivity

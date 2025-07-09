# Cart Cell

## Overview
The Cart Cell manages shopping cart operations and persistence.

## Components
- **cartservice**: C# service for cart operations

## Responsibilities
- Shopping cart management
- Cart item persistence
- Cart operations (add, remove, empty)

## API Endpoints
- `AddItem(user_id, item)`: Add item to cart
- `GetCart(user_id)`: Get user's cart
- `EmptyCart(user_id)`: Empty user's cart

## Dependencies
- Redis for cart storage
- No external cell dependencies

## Deployment Configuration
- **Runtime**: .NET Core
- **Port**: 7070
- **Storage**: Redis

## Choreo Deployment
Deployed as a single Choreo project with cart service component.
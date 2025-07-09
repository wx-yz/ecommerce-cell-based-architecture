# Platform Cell

## Overview
The Platform Cell provides shared utilities and cross-cutting concerns.

## Components
- **currencyservice**: Node.js service for currency conversion

## Responsibilities
- Currency conversion
- Supported currencies management
- Real-time exchange rates

## API Endpoints
- `GetSupportedCurrencies()`: Get supported currencies
- `Convert(from, to_code)`: Convert currency

## Dependencies
- External currency API (European Central Bank)

## Deployment Configuration
- **Runtime**: Node.js
- **Port**: 7000

## Choreo Deployment
Deployed as a single Choreo project with currency service component.
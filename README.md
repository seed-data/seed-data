# seed-data
This is an experiment in using docker & docker-compose to build a real time monitoring platform for fun and profit.

The goal of the project is to create a stand-alone application that can display historical and real-time options pricing data.

The project is divided into three sub-projects:
- web-frontend: React UI to display in your browser
- web-api: Python Flask backend API for serving JSON data
- worker: Background worker for importing new and historical data


# Architecture
The system is designed as a read-only view into read time and historical options data.  The web-frontend uses React and HTML5 to display the data to the users.  The web-api serves the JSON data to the front-end with Flask python.  Finally the backend worker processes new data and imports historical data from the archives files.


## Data Models
Each source data is setup as a data model in the worker/models package.  
The primary models are:
- Symbols: name, stock symbol, exchange, etc of a trade-able entity
- Contracts: an options contract with a given strike, expiration date, and symbol
- ContractsDaily: Intra-day price data for each contract with price, volume, open interest, etc on only the current/last day.
- ContractsEndOfDay: Daily end-of-day data for each contract with price, volume, open interest, etc.


## Data Sources
Realtime data is imported from the Tradier Brokerage API: https://developer.tradier.com/documentation

Historical data is imported from a dump file in S3.  This is only used to populate seed-data with useful data for historical comparisons.

All data adapters from external sources are in the worker/adapters package.


## Fetching Data
TODO: complete me


## Multi-instance workers and fail-overs
TODO: complete me


# Dashboards

## Runtime health checks
TODO

## Worker status
TODO

## Configuration
TODO

# Testing

## Unit tests
TODO

## Integration tests
TODO

# Deployment
## Local deployment
TODO

## AWS deployment
TODO

## Logs
TODO

# BookCabin Technical Test

A Go-based backend service that aggregates flight data from multiple airline providers (mocked), allowing users to search, filter, and sort flights.

## Table of Contents
- [Project Overview](#project-overview)
- [Technologies](#technologies)
- [Architecture & Design Choices](#architecture--design-choices)
- [Prerequisites](#prerequisites)
- [Installation & Setup](#installation--setup)
- [Running the Application](#running-the-application)
- [Running Tests](#running-tests)
- [API Documentation](#api-documentation)

## Project Overview

This service acts as a Flight Aggregator. It queries multiple distinct airline providers (Lion Air, Garuda Indonesia, Batik Air, AirAsia) in parallel, aggregates the results, filters them based on user criteria, and returns a consolidated list of flight options.

## Technologies

*   **Language**: Go (v1.25+)
*   **Web Framework**: [Gin](https://github.com/gin-gonic/gin)
*   **Logging**: [Zap](https://github.com/uber-go/zap)
*   **Configuration**: [envconfig](https://github.com/kelseyhightower/envconfig)
*   **Caching**: In-memory cache ([go-cache](https://github.com/patrickmn/go-cache))
*   **Testing**: [Testify](https://github.com/stretchr/testify)

## Architecture & Design Choices

The project follows a **Clean / Layered Architecture** to separate concerns and improve maintainability and testability.

### 1. Structure
*   `cmd/api`: Entry point of the application.
*   `internal/transport/api`: layer responsible for handling HTTP requests (Gin handlers) and decoding/encoding JSON.
*   `internal/service`: Business logic layer. It orchestrates the flow: checking cache, calling providers, filtering and sorting results.
*   `internal/provider`: Integration layer for external airline APIs. Each airline has its own implementation satisfying the `AirlineInterface`.
*   `internal/domain`: Core domain models and interfaces.
*   `pkg`: Shared utilities (Logging, Error handling, Caching).

### 2. Design Patterns & Decisions

*   **Concurrency (Fan-Out/Fan-In)**: The `AirlineProvider` uses `sync.WaitGroup` and Goroutines to query all airline providers simultaneously. This significantly reduces the total response time compared to sequential requests.
*   **Interface Segregation**: 
    *   `AirlineInterface`: Defines the contract for fetching flights from an airline.
    *   `FlightInterface`: Defines the contract for the service layer.
    *   `AirlineAggregator`: Wraps the complexity of multiple providers, making the service layer easier to test by mocking the aggregator.
*   **Caching Strategy**: Search results are cached based on a composite key of the search parameters (Origin, Destination, Date, etc.). This allows identical queries to return instantly, reducing load on providers. 
    *   *Note*: In a real-world distributed system, Redis would be preferred over in-memory cache.
*   **Smart Sorting/Ranking**: A "Best Value" score is calculated combining Price and Duration to give users the optimal trade-off.
*   **Mocking**: The provider layer currently loads data from local JSON files to simulate external API calls, with random delays added to mimic network latency.

## Prerequisites

*   **Go**: Version 1.25 or later installed.
*   **Make**: (Optional) for running Makefile commands.

## Installation & Setup

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/azcov/bookcabin_test.git
    cd bookcabin_test
    ```

2.  **Download dependencies**:
    ```bash
    go mod download
    ```

3.  **Environment Variables**:
    Copy the example environment file:
    ```bash
    cp .env.example .env
    ```
    Review the `.env` file to configure port or log levels if necessary.

## Running the Application

To run the server locally:

```bash
go run cmd/api/main.go
```

The server will start on port `8080` (default).

## Running Tests

To run all unit tests:

```bash
go test ./internal/...
```

To run detailed tests with coverage:

```bash
go test -v -cover ./internal/...
```

## API Documentation

### Search Flights
**Endpoint**: `POST /v1/flights/search`

**Request Body**:
```json
{
    "origin": "CGK",
    "destination": "DPS",
    "departureDate": "2025-12-15",
    "passengers": 1,
    "cabinClass": "Economy",
    "sort": {
        "key": "price",
        "order": "asc"
    },
    "filters": [
        { "key": "max_price", "value": 2000000 }
    ]
}
```

**Response**:
```json
{
    "metadata": {
        "total_results": 5,
        "search_time_ms": 120,
        "cache_hit": false
    },
    "flights": [
        {
            "id": "JT740_Lion Air",
            "airline": { "name": "Lion Air", "code": "JT" },
            "price": { "amount": 950000, "currency": "IDR" },
            // ... other details
        }
    ]
}
```
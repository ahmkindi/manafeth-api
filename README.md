# Trade Data Warehouse API 🚀

A high-performance REST API built with Go (Fiber) and PostgreSQL for querying international trade data with flexible aggregation capabilities. Optimized for AI agent consumption.

## ✨ Features

- **Flexible Aggregation**: Dynamic grouping by product, country, port, year, and trade type
- **High Performance**: Connection pooling, query optimization, caching
- **AI-Friendly**: Structured JSON requests/responses perfect for LLM agents
- **Production Ready**: Security headers, rate limiting, error handling, logging
- **Scalable**: Star schema design for fast analytical queries

## 📋 Prerequisites

- Go 1.21+
- PostgreSQL 14+
- Docker & Docker Compose (optional)

## 🚀 Quick Start

### Option 1: Using Docker Compose (Recommended)

```bash
# Clone the repository
git clone <your-repo>
cd trade-api

# Start services
docker-compose up -d

# Check health
curl http://localhost:3000/health
```

### Option 2: Local Development

1. **Install dependencies**
```bash
go mod download
```

2. **Setup PostgreSQL**
```bash
# Create database
createdb trade_db

# Run schema creation (from your SQL file)
psql -d trade_db -f schema.sql
```

3. **Configure environment**
```bash
cp .env.example .env
# Edit .env with your database credentials
```

4. **Run the application**
```bash
go run main.go
# or
make run
```

The API will be available at `http://localhost:3000`

## 📁 Project Structure

```
trade-api/
├── main.go                 # Application entry point
├── config/
│   └── database.go        # Database connection & pooling
├── models/
│   └── models.go          # Data structures
├── handlers/
│   ├── dimensions.go      # Dimension endpoints (products, countries, ports)
│   └── trade.go          # Trade query endpoints
├── middleware/
│   └── middleware.go      # Cache & other middleware
├── utils/
│   └── query_builder.go   # Dynamic SQL query builder
├── docker-compose.yml     # Docker orchestration
├── Dockerfile            # Application container
├── Makefile              # Build commands
└── .env.example          # Environment template
```

## 🔧 Configuration

### Environment Variables

Create a `.env` file with:

```env
PORT=3000

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=trade_db
DB_SSLMODE=disable
```

### Connection Pool Settings

Configured in `config/database.go`:
- Max connections: 25
- Min connections: 5
- Max lifetime: 1 hour
- Health check: 1 minute

## 📚 API Documentation

### Base URL
```
http://localhost:3000/api/v1
```

### Key Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/dimensions/products` | List/search products |
| GET | `/dimensions/countries` | List/search countries |
| GET | `/dimensions/ports` | List/search ports |
| GET | `/trade/summary` | Yearly trade summary |
| GET | `/trade/balance` | Trade balance calculation |
| POST | `/trade/aggregate` | **Main query endpoint** |

### Example: Aggregate Query

```bash
curl -X POST http://localhost:3000/api/v1/trade/aggregate \
  -H "Content-Type: application/json" \
  -d '{
    "date_range": {
      "start_year": 2020,
      "end_year": 2023
    },
    "trade_types": ["Import"],
    "group_by": ["country", "year"],
    "filters": {
      "port_types": ["Sea"]
    },
    "pagination": {
      "page": 1,
      "limit": 25
    },
    "sorting": {
      "sort_by": "total_value",
      "sort_order": "desc"
    }
  }'
```

**Response:**
```json
{
  "data": [
    {
      "country_id": 5,
      "country_name_en": "United Arab Emirates",
      "year": 2023,
      "total_value": 5000000000
    }
  ],
  "pagination": {
    "current_page": 1,
    "page_size": 25,
    "total_count": 150,
    "total_pages": 6
  }
}
```

## 🤖 AI Agent Examples

### Query 1: "What were the top 10 imported products in 2022?"

```json
{
  "date_range": {"start_year": 2022, "end_year": 2022},
  "trade_types": ["Import"],
  "group_by": ["product"],
  "pagination": {"limit": 10},
  "sorting": {"sort_by": "total_value", "sort_order": "desc"}
}
```

### Query 2: "Total trade balance from 2020-2023"

```bash
curl "http://localhost:3000/api/v1/trade/balance?start_year=2020&end_year=2023"
```

### Query 3: "Imports by country for the last 5 years"

```json
{
  "date_range": {"start_year": 2019, "end_year": 2023},
  "trade_types": ["Import"],
  "group_by": ["country", "year"],
  "pagination": {"limit": 100}
}
```

## 📊 Performance Optimization

1. **Caching**: Dimension endpoints cached for 5 minutes
2. **Connection Pooling**: Optimized pool settings
3. **Indexes**: Ensure composite indexes on fact tables:
   ```sql
   CREATE INDEX idx_product_port_year ON fact_trade_by_product_port(year, trade_type, product_id, port_id);
   CREATE INDEX idx_country_port_year ON fact_trade_by_country_port(year, trade_type, country_id, port_id);
   ```
4. **Query Optimization**: Dynamic query building avoids unnecessary joins

## 🛡️ Security Features

- Helmet middleware (security headers)
- Rate limiting (100 req/min)
- SQL injection protection (parameterized queries)
- Input validation
- Panic recovery
- Request size limits

## 📝 Development

### Building

```bash
# Development build
go build -o bin/trade-api .

# Production build (optimized)
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
```
```

## 🚢 Deployment

### Docker Production

```bash
# Build image
docker build -t trade-api:latest .

# Run container
docker run -d \
  -p 3000:3000 \
  -e DB_HOST=your-db-host \
  -e DB_PASSWORD=your-password \
  --name trade-api \
  trade-api:latest
```

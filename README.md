# SpiderLite 🕷️

SpiderLite is a lightweight web crawler with a built-in API server for querying crawled data. It respects robots.txt rules and stores crawl results in a SQLite database.

## Features

- Web crawling with configurable depth
- robots.txt compliance
- SQLite storage for crawl results
- RESTful API to query crawled data
- Datadog integration for metrics
- Docker support

## Prerequisites

- Go 1.23 or higher
- Docker and Docker Compose
- Datadog account and API key

## Quick Start

1. Clone the repository:
```bash
git clone https://github.com/yourusername/spiderlite.git
cd spiderlite
```

2. Create a `.env` file:
```bash
cp .env.example .env
# Edit .env with your Datadog API key
```

3. Start the application:
```bash
docker-compose up -d
```

4. Start crawling a website:
```bash
curl -X POST "http://localhost:8080/crawl?url=https://example.com"
```

## API Endpoints

### Start a Crawl
```bash
POST /crawl?url=https://example.com
```
Response:
```json
{
  "status": "started",
  "message": "Crawl started for https://example.com"
}
```

### Get All Pages
```bash
GET /pages
```
Response:
```json
{
  "count": 1,
  "pages": [
    {
      "URL": "https://example.com",
      "StatusCode": 200,
      "CrawledAt": "2024-01-01T12:34:56Z"
    }
  ]
}
```

### Get Pages by Status Code
```bash
GET /pages/status?code=200
```
Response:
```json
{
  "status": 200,
  "count": 1,
  "pages": [
    {
      "URL": "https://example.com",
      "StatusCode": 200,
      "CrawledAt": "2024-01-01T12:34:56Z"
    }
  ]
}
```

### Debug Information
```bash
GET /debug
```
Response:
```json
{
  "database_path": "/data/crawler.db",
  "table_count": 1,
  "tables": {
    "pages": 1
  }
}
```

## Configuration

Environment variables (in `.env`):

```env
DD_API_KEY=your_api_key_here    # Datadog API key
DD_ENV=dev                      # Environment (dev, prod, etc.)
DD_SERVICE=spiderlite          # Service name in Datadog
DB_PATH=/data/crawler.db       # SQLite database path
```

## Monitoring

### Datadog Metrics

- `spiderlite.crawler.pages_processed`: Counter of processed pages
- `spiderlite.crawler.page_process_time`: Timing of page processing
- `spiderlite.crawler.errors`: Counter of crawl errors
- `spiderlite.api.requests`: Counter of API requests

### Datadog Logs

Logs are automatically collected and include:
- Crawl operations
- Page processing results
- API requests
- Error details

## Development

### Prerequisites

- Go 1.23+
- Docker and Docker Compose
- Make (optional)

### Local Development

1. Install dependencies:
```bash
go mod download
```

2. Run tests:
```bash
go test ./...
```

3. Build:
```bash
go build -o spiderlite ./cmd/crawler
go build -o spiderlite-server ./cmd/server
```

### Docker Development

```bash
# Build and start services
docker-compose up --build -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT

## Acknowledgments

- [robotstxt](https://github.com/temoto/robotstxt) - For robots.txt parsing
- [go-sqlite3](https://github.com/mattn/go-sqlite3) - SQLite driver for Go

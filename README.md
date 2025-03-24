# SpiderLite

SpiderLite is a lightweight web crawler with a built-in API server for querying crawled data. It respects robots.txt rules and stores crawl results in a SQLite database.

## Features

- Web crawling with configurable depth
- robots.txt compliance
- SQLite storage for crawl results
- RESTful API to query crawled data
- Docker support

## Installation

### Prerequisites

- Go 1.21 or higher
- SQLite3
- Make (optional, for using Makefile commands)

### Building from source

```bash
# Clone the repository
git clone https://github.com/yourusername/spiderlite.git
cd spiderlite

# Install dependencies
go mod download

# Build both crawler and server
make build
```

## Usage

### Web Crawler

To start crawling a website:

```bash
# Using make
make run URL=https://example.com

# Or directly
./spiderlite https://example.com
```

### API Server

To start the API server:

```bash
# Using make
make server

# Or directly
./spiderlite-server -addr=:8080 -db=crawler.db
```

### Docker

Build and run using Docker:

```bash
# Build the image
docker build -t spiderlite .

# Run the crawler
docker run spiderlite https://example.com

# Run the server
docker run -p 8080:8080 spiderlite-server
```

## API Endpoints

### Get all crawled pages

GET /pages

Response:
```json
[
  {
    "URL": "https://example.com",
    "StatusCode": 200,
    "CrawledAt": "2024-01-01T12:34:56Z"
  }
]
```

### Get pages by status code

GET /pages/status?code=200

Response:
```json
[
  {
    "URL": "https://example.com",
    "StatusCode": 200,
    "CrawledAt": "2024-01-01T12:34:56Z"
  }
]
```

## Project Structure

spiderlite/
├── cmd/
│ ├── crawler/ # Crawler executable
│ └── server/ # API server executable
├── internal/
│ ├── crawler/ # Crawler logic
│ ├── database/ # Database operations
│ ├── parser/ # HTML parsing
│ └── server/ # HTTP server
├── Dockerfile
├── Makefile
└── README.md


## Development

### Available Make Commands

```bash
make build    # Build both crawler and server
make run      # Build and run the crawler
make server   # Build and run the API server
make clean    # Clean up built binaries
```

### Adding New Features

1. Create a new branch
2. Add your feature
3. Write tests
4. Submit a pull request

## License

MIT

## Acknowledgments

- [robotstxt](https://github.com/temoto/robotstxt) - For robots.txt parsing
- [go-sqlite3](https://github.com/mattn/go-sqlite3) - SQLite driver for Go

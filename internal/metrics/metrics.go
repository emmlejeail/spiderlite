package metrics

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
)

type Metrics struct {
	client *statsd.Client
}

func New() (*Metrics, error) {
	// Get host from environment or use default
	host := os.Getenv("DD_AGENT_HOST")
	if host == "" {
		host = "datadog-agent"
	}

	// Configure StatsD client
	client, err := statsd.New(fmt.Sprintf("%s:8125", host),
		statsd.WithNamespace("spiderlite."),
		statsd.WithTags([]string{
			fmt.Sprintf("env:%s", os.Getenv("DD_ENV")),
			fmt.Sprintf("service:%s", os.Getenv("DD_SERVICE")),
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create statsd client: %v", err)
	}

	// Test connection
	if err := client.Gauge("startup", 1, nil, 1); err != nil {
		log.Printf("Warning: Failed to send test metric: %v", err)
	} else {
		log.Printf("Successfully connected to Datadog agent at %s", host)
	}

	return &Metrics{client: client}, nil
}

func (m *Metrics) Close() error {
	return m.client.Close()
}

// Métriques pour le crawler
func (m *Metrics) IncrementPagesProcessed(statusCode int, host string) {
	tags := []string{
		fmt.Sprintf("status:%d", statusCode),
		fmt.Sprintf("host:%s", host),
	}
	if err := m.client.Incr("crawler.pages_processed", tags, 1); err != nil {
		log.Printf("Failed to send metric crawler.pages_processed: %v", err)
	}
}

func (m *Metrics) TimeCrawl(duration time.Duration, host string) {
	tags := []string{fmt.Sprintf("host:%s", host)}
	if err := m.client.Timing("crawler.page_process_time", duration, tags, 1); err != nil {
		log.Printf("Failed to send metric crawler.page_process_time: %v", err)
	}
}

func (m *Metrics) GaugeLinksFound(count int, host string) {
	tags := []string{"host:" + host}
	m.client.Gauge("crawler.links_found", float64(count), tags, 1)
}

// Métriques pour l'API
func (m *Metrics) IncrementAPIRequests(endpoint, method string, statusCode int) {
	tags := []string{
		"endpoint:" + endpoint,
		"method:" + method,
		"status:" + strconv.Itoa(statusCode),
	}
	m.client.Incr("api.requests", tags, 1)
}

func (m *Metrics) TimeAPIRequest(endpoint string, duration time.Duration) {
	tags := []string{"endpoint:" + endpoint}
	m.client.Timing("api.request_duration", duration, tags, 1)
}

func (m *Metrics) IncrementCrawlErrors() {
	if err := m.client.Incr("crawler.errors", nil, 1); err != nil {
		log.Printf("Failed to send metric crawler.errors: %v", err)
	}
}

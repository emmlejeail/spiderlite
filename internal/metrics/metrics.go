package metrics

import (
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
)

type Metrics struct {
	client *statsd.Client
}

func New() (*Metrics, error) {
	// Initialiser le client StatsD
	client, err := statsd.New("127.0.0.1:8125",
		statsd.WithNamespace("spiderlite."),
		statsd.WithTags([]string{"env:prod"}),
	)
	if err != nil {
		return nil, err
	}

	return &Metrics{client: client}, nil
}

func (m *Metrics) Close() error {
	return m.client.Close()
}

// Métriques pour le crawler
func (m *Metrics) IncrementPagesProcessed(statusCode int, host string) {
	tags := []string{
		"status:" + string(statusCode),
		"host:" + host,
	}
	m.client.Incr("crawler.pages_processed", tags, 1)
}

func (m *Metrics) TimeCrawl(duration time.Duration, host string) {
	tags := []string{"host:" + host}
	m.client.Timing("crawler.page_process_time", duration, tags, 1)
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
		"status:" + string(statusCode),
	}
	m.client.Incr("api.requests", tags, 1)
}

func (m *Metrics) TimeAPIRequest(endpoint string, duration time.Duration) {
	tags := []string{"endpoint:" + endpoint}
	m.client.Timing("api.request_duration", duration, tags, 1)
}

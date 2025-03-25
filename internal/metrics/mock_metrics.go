package metrics

import "time"

// NoopMetrics is a mock metrics client that does nothing
type NoopMetrics struct{}

func NewNoopMetrics() *NoopMetrics {
	return &NoopMetrics{}
}

func (m *NoopMetrics) IncrementPagesProcessed(statusCode int, host string)          {}
func (m *NoopMetrics) IncrementCrawlErrors()                                        {}
func (m *NoopMetrics) TimeCrawl(duration time.Duration, host string)                {}
func (m *NoopMetrics) IncrementAPIRequests(endpoint, method string, statusCode int) {}
func (m *NoopMetrics) TimeAPIRequest(endpoint string, duration time.Duration)       {}
func (m *NoopMetrics) Close() error                                                 { return nil }

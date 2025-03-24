package crawler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"spiderlite/internal/database"
	"spiderlite/internal/metrics"
	"spiderlite/internal/parser"
)

type Crawler struct {
	db      *database.DB
	metrics *metrics.Metrics
	visited map[string]bool
	client  *http.Client
}

func New(db *database.DB, metrics *metrics.Metrics) *Crawler {
	return &Crawler{
		db:      db,
		metrics: metrics,
		visited: make(map[string]bool),
		client:  &http.Client{},
	}
}

func (c *Crawler) Start(startURL *url.URL) error {
	robots, err := NewRobotsChecker(startURL)
	if err != nil {
		return err
	}

	if !robots.IsAllowed(startURL.Path) {
		return fmt.Errorf("URL disallowed by robots.txt: %s", startURL)
	}

	return c.crawl(startURL, robots)
}

func (c *Crawler) crawl(u *url.URL, robots *RobotsChecker) error {
	start := time.Now()
	defer func() {
		c.metrics.TimeCrawl(time.Since(start), u.Host)
	}()

	if c.visited[u.String()] {
		return nil
	}
	c.visited[u.String()] = true

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return c.storeError(u, err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return c.storeError(u, err)
	}
	defer resp.Body.Close()

	if err := c.db.StorePage(database.PageData{
		URL:        u.String(),
		StatusCode: resp.StatusCode,
		CrawledAt:  time.Now(),
	}); err != nil {
		return err
	}

	c.metrics.IncrementPagesProcessed(resp.StatusCode, u.Host)

	if resp.StatusCode != 200 {
		return nil
	}

	links, err := parser.ExtractLinks(resp.Body, u)
	if err != nil {
		return err
	}

	c.metrics.GaugeLinksFound(len(links), u.Host)

	for _, link := range links {
		if !robots.IsAllowed(link.Path) {
			continue
		}
		if link.Host != u.Host {
			continue
		}
		if err := c.crawl(link, robots); err != nil {
			return err
		}
	}

	return nil
}

func (c *Crawler) storeError(u *url.URL, err error) error {
	return c.db.StorePage(database.PageData{
		URL:        u.String(),
		StatusCode: 0,
		CrawledAt:  time.Now(),
	})
}

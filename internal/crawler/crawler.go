package crawler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"spiderlite/internal/database"
	"spiderlite/internal/metrics"
	"spiderlite/internal/parser"
)

type Crawler struct {
	db      *database.DB
	metrics metrics.MetricsClient
	visited map[string]bool
	client  *http.Client
}

func New(db *database.DB, m metrics.MetricsClient) *Crawler {
	return &Crawler{
		db:      db,
		metrics: m,
		visited: make(map[string]bool),
		client:  &http.Client{},
	}
}

func (c *Crawler) Start(startURL *url.URL) error {
	log.Printf("Starting crawl for: %s", startURL.String())

	robots, err := NewRobotsChecker(startURL)
	if err != nil {
		log.Printf("Robots.txt error: %v", err)
		// Continue anyway
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

	log.Printf("Crawling: %s", u.String())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		log.Printf("Request error for %s: %v", u.String(), err)
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.metrics.IncrementCrawlErrors()
		log.Printf("HTTP error for %s: %v", u.String(), err)
		// Store error page
		if err := c.db.StorePage(database.PageData{
			URL:        u.String(),
			StatusCode: 0,
			CrawledAt:  time.Now(),
		}); err != nil {
			log.Printf("Failed to store error page: %v", err)
		}
		return err
	}
	defer resp.Body.Close()

	// Store successful page
	pageData := database.PageData{
		URL:        u.String(),
		StatusCode: resp.StatusCode,
		CrawledAt:  time.Now(),
	}

	if err := c.db.StorePage(pageData); err != nil {
		log.Printf("Failed to store page %s: %v", u.String(), err)
		return err
	}
	log.Printf("Successfully stored page: %s with status: %d", u.String(), resp.StatusCode)

	// Increment pages processed with status code
	c.metrics.IncrementPagesProcessed(resp.StatusCode, u.Host)

	if resp.StatusCode != 200 {
		log.Printf("Non-200 status code for %s: %d", u.String(), resp.StatusCode)
		return nil
	}

	links, err := parser.ExtractLinks(resp.Body, u)
	if err != nil {
		log.Printf("Link extraction error for %s: %v", u.String(), err)
		return err
	}

	log.Printf("Found %d links on %s", len(links), u.String())

	for _, link := range links {
		if !robots.IsAllowed(link.Path) {
			log.Printf("Skipping disallowed URL: %s", link.String())
			continue
		}
		if link.Host != u.Host {
			log.Printf("Skipping external URL: %s", link.String())
			continue
		}
		if err := c.crawl(link, robots); err != nil {
			log.Printf("Error crawling %s: %v", link.String(), err)
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

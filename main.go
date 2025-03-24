package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"io"
	"strings"

	"github.com/temoto/robotstxt"
	"golang.org/x/net/html"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <url>", os.Args[0])
	}
	startURL := os.Args[1]

	parsedURL, err := url.Parse(startURL)
	if err != nil {
		log.Fatalf("Invalid URL: %v", err)
	}

	robotsURL := fmt.Sprintf("%s://%s/robots.txt", parsedURL.Scheme, parsedURL.Host)
	robots, err := fetchRobots(robotsURL)
	if err != nil {
		log.Printf("Warning: Failed to fetch robots.txt: %v", err)
	}

	// Use User-agent: *
	allowed := func(path string) bool { return true }
	if robots != nil {
		grp := robots.FindGroup("*")
		allowed = grp.Test
	}

	if !allowed(parsedURL.Path) {
		log.Fatalf("Disallowed by robots.txt: %s", parsedURL.Path)
	}

	log.Printf("Starting crawl at %s\n", startURL)
	visited := make(map[string]bool)
	crawl(parsedURL, allowed, visited)
}

func fetchRobots(robotsURL string) (*robotstxt.RobotsData, error) {
	resp, err := http.Get(robotsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := robotstxt.FromResponse(resp)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func crawl(u *url.URL, allowed func(string) bool, visited map[string]bool) {
	if visited[u.String()] {
		return
	}
	visited[u.String()] = true

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		log.Printf("Request error: %v", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("HTTP error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Non-200: %s (%d)", u.String(), resp.StatusCode)
		return
	}

	log.Printf("Crawled: %s", u.String())

	links, err := extractLinks(resp.Body, u)
	if err != nil {
		log.Printf("Link extraction failed: %v", err)
		return
	}

	for _, link := range links {
		if !allowed(link.Path) {
			log.Printf("Disallowed: %s", link.String())
			continue
		}
		if link.Host != u.Host {
			continue // skip external links
		}
		crawl(link, allowed, visited) // one-level recursive crawl
	}
}

func extractLinks(body io.Reader, base *url.URL) ([]*url.URL, error) {
	tokens := html.NewTokenizer(body)
	links := []*url.URL{}

	for {
		tt := tokens.Next()
		if tt == html.ErrorToken {
			break
		}
		token := tokens.Token()
		if token.Type == html.StartTagToken && token.DataAtom.String() == "a" {
			for _, attr := range token.Attr {
				if attr.Key == "href" {
					href := strings.TrimSpace(attr.Val)
					link, err := base.Parse(href)
					if err == nil {
						links = append(links, link)
					}
				}
			}
		}
	}
	return links, nil
}

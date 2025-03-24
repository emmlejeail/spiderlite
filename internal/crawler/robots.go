package crawler

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/temoto/robotstxt"
)

type RobotsChecker struct {
	robots *robotstxt.RobotsData
}

func NewRobotsChecker(baseURL *url.URL) (*RobotsChecker, error) {
	robotsURL := fmt.Sprintf("%s://%s/robots.txt", baseURL.Scheme, baseURL.Host)
	robots, err := fetchRobots(robotsURL)
	if err != nil {
		return nil, err
	}
	return &RobotsChecker{robots: robots}, nil
}

func fetchRobots(robotsURL string) (*robotstxt.RobotsData, error) {
	resp, err := http.Get(robotsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return robotstxt.FromResponse(resp)
}

func (r *RobotsChecker) IsAllowed(path string) bool {
	if r.robots == nil {
		return true
	}
	return r.robots.FindGroup("*").Test(path)
}

package parser

import (
	"io"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func ExtractLinks(body io.Reader, base *url.URL) ([]*url.URL, error) {
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

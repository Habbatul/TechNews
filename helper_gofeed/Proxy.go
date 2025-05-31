package helper_gofeed

import (
	"github.com/mmcdole/gofeed"
	"net/http"
	"net/url"
	"os"
	"time"
)

func NewParserWithHTTPProxy() (*gofeed.Parser, error) {
	proxyStr := os.Getenv("PROXY_URL")
	if proxyStr == "" {
		//tangani didepan kalau nil pakek parser default gofeed
		return nil, nil
	}

	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   50 * time.Second,
	}

	parser := gofeed.NewParser()
	parser.Client = client
	return parser, nil
}

func ParseURLWithProxy(feedURL string, parser *gofeed.Parser) (*gofeed.Feed, error) {
	req, err := http.NewRequest("GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := parser.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	feed, err := parser.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	return feed, nil
}

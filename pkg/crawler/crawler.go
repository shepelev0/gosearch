package crawler

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"strings"
)

// Result — результат сканирования одной страницы
type Result struct {
	URL  string
	Body string
}

// Crawler — структура для сканирования сайтов
type Crawler struct{}

func New() *Crawler {
	return &Crawler{}
}

func (c *Crawler) Crawl(url string, depth int) ([]Result, error) {
	visited := make(map[string]bool)
	var results []Result

	var crawl func(string, int)
	crawl = func(u string, d int) {
		if d <= 0 || visited[u] {
			return
		}
		visited[u] = true

		body, links, err := fetch(u)
		if err != nil {
			fmt.Printf("Ошибка при получении %s: %v\n", u, err)
			return
		}

		results = append(results, Result{URL: u, Body: body})
		for _, link := range links {
			crawl(link, d-1)
		}
	}

	crawl(url, depth)
	return results, nil
}

func fetch(url string) (string, []string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}
	body := string(bodyBytes)
	links := extractLinks(body, url)
	return body, links, nil
}

func extractLinks(htmlBody, baseURL string) []string {
	doc, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return nil
	}
	var links []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" && strings.HasPrefix(a.Val, "/") {
					links = append(links, baseURL+a.Val)
				} else if a.Key == "href" && (strings.HasPrefix(a.Val, "http") || strings.HasPrefix(a.Val, "https")) {
					links = append(links, a.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return links
}

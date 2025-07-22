// launches concurrent scrapers on given business urls.
package web

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// attempt to find a job-related page (careers, jobs, hiring) on the given site.
func ScrapeWebsite(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}

	var jobPage string
	var visit func(*html.Node)
	visit = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					href := strings.ToLower(attr.Val)
					if strings.Contains(href, "careers") || strings.Contains(href, "jobs") || strings.Contains(href, "hiring") {
						jobPage = attr.Val
						return
					}
				}
			}
		}
		for c := n.FirstChild; c != nil && jobPage == ""; c = c.NextSibling {
			visit(c)
		}
	}
	visit(doc)

	if jobPage == "" {
		return "", errors.New("no job page found")
	}

	// Resolve relative URL
	if strings.HasPrefix(jobPage, "/") {
		return strings.TrimSuffix(url, "/") + jobPage, nil
	}

	return jobPage, nil
}


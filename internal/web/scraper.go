// launches concurrent scrapers on given business urls.
package web

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// take a root website URL and tries to find a careers/job page.
func ScrapeWebsite(rootURL string) (string, error) {
	// fetch url root and checks if responds 
	resp, err := http.Get(rootURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch %s: %w", rootURL, err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body: %w", err)
	}
	body := string(bodyBytes)

	// if the root page itself looks like a job page (unlikely to be true)
	if IsJobPage(rootURL, body) {
		return rootURL, nil
	}

	// , parse HTML and scan links
	pageLinks := extractLinks(body, rootURL)
	for _, link := range pageLinks {
		// quick keyword check before fetching
		for _, kw := range JobPageKeywords {
			if strings.Contains(strings.ToLower(link), kw) {
				// fetch link and confirm it’s a job page
				jobURL, ok := checkLink(link)
				if ok {
					return jobURL, nil
				}
			}
		}
	}

	return "", nil // nothing found
}

// fetch a link and applies IsJobPage
func checkLink(link string) (string, bool) {
	resp, err := http.Get(link)
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", false
	}
	body := string(bodyBytes)
	fmt.Printf("Checking candidate link: %s\n", link)

	if IsJobPage(link, body) {
		return link, true
	}
	return "", false
}

// this is only parsing <a href=".."> links from the HTML body
func extractLinks(body, base string) []string {
	var links []string
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return links
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					href := resolveURL(base, attr.Val)
					if href != "" {
						links = append(links, href)
					}
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

// resolve relative/absolute URLs
func resolveURL(base, href string) string {
	parsedBase, err := url.Parse(base)
	if err != nil {
		return ""
	}
	parsedHref, err := url.Parse(href)
	if err != nil {
		return ""
	}
	return parsedBase.ResolveReference(parsedHref).String()
}


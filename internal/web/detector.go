// applies heuristic to locate job/career pages.
package web

import (
	"strings"

	"golang.org/x/net/html"
)

var JobPageKeywords = []string{
	"careers", "jobs", "join-us", "employment",
	"opportunities", "work-with-us", "hiring",
}

func IsJobPage(url, body string) bool {
	urlLower := strings.ToLower(url)
	for _, kw := range JobPageKeywords {
		if strings.Contains(urlLower, kw) {
			return true
		}
	}

	// filter out scripts/styles, then scan for keywords
	text := extractVisibleText(body)
	for _, kw := range JobPageKeywords {
		if strings.Contains(text, kw) {
			return true
		}
	}
	return false
}

func MatchesJobTitle(body string, titles []string) bool {
	// if no title is provided, assume true and search for any job page 
	if len(titles) == 0 {
		return true
	}
	text := extractVisibleText(body)

	// build quick word index
	words := strings.FieldsFunc(text, func(r rune) bool {
		return r < 'a' || r > 'z' // split on anything non a–z
	})
	wordSet := make(map[string]struct{}, len(words))
	for _, w := range words {
		wordSet[w] = struct{}{}
	}

	for _, t := range titles {
		tLower := strings.ToLower(t)
		if _, ok := wordSet[tLower]; ok && hasNearbyContext(words, tLower) {
			return true
		}
	}
	return false
}

// context words to avoid false hits
var contextWords = []string{
	"apply", "opening", "position", "role", "responsibilities",
	"full-time", "part-time", "hiring", "join", "career", "benefits",
}

// DOM text extractor
func extractVisibleText(htmlBody string) string {
	doc, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return strings.ToLower(htmlBody) // fallback
	}

	var sb strings.Builder
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.TextNode {
			txt := strings.TrimSpace(n.Data)
			if txt != "" {
				sb.WriteString(txt)
				sb.WriteByte(' ')
			}
		}
		if n.Type == html.ElementNode {
			switch strings.ToLower(n.Data) {
			case "script", "style", "noscript":
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)

	return strings.ToLower(sb.String())
}

func hasNearbyContext(words []string, token string) bool {
	for i, w := range words {
		if w == token {
			start := max(0, i-8)
			end := min(len(words), i+8)
			for _, ctx := range contextWords {
				for j := start; j < end; j++ {
					if words[j] == ctx {
						return true
					}
				}
			}
		}
	}
	return false
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}


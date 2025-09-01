// applies heuristic to locate job/career pages. 
package web

import "strings"

var JobPageKeywords = []string{"careers", "jobs", "join-us", "employment", "opportunities", "work-with-us", "hiring"}

func IsJobPage(url string, body string) bool {
	urlLower := strings.ToLower(url)
	for _, kw := range JobPageKeywords {
		if strings.Contains(urlLower, kw) {
			return true
		}
	}
	bodyLower := strings.ToLower(body)
	for _, kw := range JobPageKeywords {
		if strings.Contains(bodyLower, kw) {
			return true
		}
	}
	return false
}

func MatchesJobTitle(body string, titles []string) bool {
	bodyLower := strings.ToLower(body)
	for _, t := range titles {
		if strings.Contains(bodyLower, strings.ToLower(t)) {
			return true
		}
	}
	return false
}

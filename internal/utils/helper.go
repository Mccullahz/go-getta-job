// contains reusable helpers to be used across the application 
package utils

import (

)

func NormalizeURL(url string) string {
	// sanitize and normalize URL doesnt need to be much here, just making sure urls are easy to work with
	if url == "" {
		return ""
	}
	if url[len(url)-1] == '/' {
		return url[:len(url)-1]
	}
	return url
}

func IsValidZip(zip string) bool {
    if len(zip) != 5 {
        return false
    }
    for _, ch := range zip {
        if ch < '0' || ch > '9' {
            return false
        }
    }
    return true
}

func IsValidRadius(r string) bool {
    if r == "" {
        return false
    }
    for _, ch := range r {
        if ch < '0' || ch > '9' {
            return false
        }
    }
    return true
}

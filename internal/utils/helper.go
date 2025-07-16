// contains reusable helpers 
package utils
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


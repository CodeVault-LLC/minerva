package utils

import "strings"

func NormalizeURL(url string) string {
	var urlNormalized string

	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		urlNormalized = url
	} else {
		urlNormalized = "https://" + url
	}

	if strings.HasSuffix(urlNormalized, "/") {
		urlNormalized = urlNormalized[:len(urlNormalized)-1]
	}

	if strings.Contains(urlNormalized, "www.") {
		urlNormalized = strings.Replace(urlNormalized, "www.", "", 1)
	}

	return urlNormalized
}

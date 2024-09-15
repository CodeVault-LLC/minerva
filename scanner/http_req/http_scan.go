package http_req

import (
	"net/http"
)

func ScanHTTPHeaders(url string) (map[string][]string, error) {
	resp, err := http.Head(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return resp.Header, nil
}

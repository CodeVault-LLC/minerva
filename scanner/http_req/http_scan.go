package http_req

import (
	"fmt"
	"net/http"
)

func ScanHTTPHeaders(url string) ([]string, error) {
	resp, err := http.Head(url)
	if err != nil {
		fmt.Println("Error scanning HTTP headers", err)
		return nil, err
	}

	defer resp.Body.Close()

	var headers []string
	for key, value := range resp.Header {
		headers = append(headers, fmt.Sprintf("%s: %s", key, value))
	}

	return headers, nil
}

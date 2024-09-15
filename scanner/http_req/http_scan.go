package http_req

import (
	"encoding/json"
	"net/http"
)

func ScanHTTPHeaders(url string) (json.RawMessage, error) {
	resp, err := http.Head(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	headers, err := json.Marshal(resp.Header)
	if err != nil {
		return nil, err
	}

	return headers, nil
}

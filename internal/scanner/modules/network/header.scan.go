package network

import (
	"fmt"
	"net/http"
)

type HTTPResponse struct {
	StatusCode int
	Headers    http.Header
}

func getHeaders(url string) (HTTPResponse, error) {
	resp, err := http.Head(url)
	if err != nil {
		fmt.Println("Error scanning HTTP headers", err)
		return HTTPResponse{}, err
	}

	defer resp.Body.Close()

	return HTTPResponse{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
	}, nil
}

package network

import (
	"fmt"
	"net/http"

	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
)

type HeaderModule struct{}

func (m *HeaderModule) Run(job entities.JobModel) (interface{}, error) {
	headers, err := getHeaders(job.URL)
	if err != nil {
		return nil, err
	} else {
		return headers, nil
	}
}

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

func (m *HeaderModule) Name() string {
	return "Header"
}

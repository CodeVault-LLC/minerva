package network

import (
	"fmt"
	"net/http"

	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"go.uber.org/zap"
)

type HeaderModule struct{}

func (m *HeaderModule) Run(job entities.JobModel) (interface{}, error) {
	headers, err := getHeaders(job.URL)
	if err != nil {
		logger.Log.Error("Error getting headers", zap.Error(err))
		return nil, err
	}

	httpHeaders := make([]string, 0)
	for key, value := range headers.Headers {
		httpHeaders = append(httpHeaders, fmt.Sprintf("%s: %s", key, value))
	}

	return httpHeaders, nil
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

package network

import (
	"net"

	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"go.uber.org/zap"
)

// IPLookupModule implements LookupModule
type IPLookupModule struct{}

func (m *IPLookupModule) Run(job entities.JobModel) (interface{}, error) {
	url := utils.StripProtocol(job.URL)

	ipAddresses, err := net.LookupIP(url)
	if err != nil {
		logger.Log.Error("Error in IP Lookup", zap.Error(err))
		return nil, err
	}

	ips := make([]string, len(ipAddresses))
	for i, ip := range ipAddresses {
		ips[i] = ip.String()
	}

	return ips, nil
}

func (m *IPLookupModule) Name() string {
	return "IPLookup"
}

// IPRangeLookupModule implements LookupModule
type IPRangeLookupModule struct{}

func (m *IPRangeLookupModule) Run(job entities.JobModel) (interface{}, error) {
	url := utils.StripProtocol(job.URL)

	ips, err := net.LookupHost(url)
	if err != nil {
		return nil, err
	}

	return ips, nil
}

func (m *IPRangeLookupModule) Name() string {
	return "IPRangeLookup"
}

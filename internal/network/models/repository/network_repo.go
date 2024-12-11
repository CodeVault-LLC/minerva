package repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/codevault-llc/minerva/internal/database"
	"github.com/codevault-llc/minerva/internal/network/models/entities"
	"github.com/codevault-llc/minerva/pkg/logger"
	"go.uber.org/zap"
)

type NetworkRepo struct {
	database *database.Database
}

var NetworkRepository *NetworkRepo

// NewNetworkRepository creates a new NetworkRepository
func NewNetworkRepository(database *database.Database) *NetworkRepo {
	return &NetworkRepo{
		database: database,
	}
}

// NetworkRepositoryInterface is the interface for the NetworkRepository
func (n *NetworkRepo) Create(scanId string, network entities.NetworkModel) error {
	query := "INSERT INTO minerva.networks (id, scan_id, ip_addresses, ip_ranges, http_headers, dns, whois, certificates, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	ipAddresses, err := json.Marshal(network.IpAddresses)
	if err != nil {
		logger.Log.Error("Failed to marshal ip_addresses", zap.Error(err))
		return err
	}

	ipRanges, err := json.Marshal(network.IpRanges)
	if err != nil {
		logger.Log.Error("Failed to marshal ip_ranges", zap.Error(err))
		return err
	}

	httpHeaders, err := json.Marshal(network.HttpHeaders)
	if err != nil {
		logger.Log.Error("Failed to marshal http_headers", zap.Error(err))
		return err
	}

	certificates, err := json.Marshal(network.CertificateModel)
	if err != nil {
		logger.Log.Error("Failed to marshal certificates", zap.Error(err))
		return err
	}

	dnsModel, err := json.Marshal(network.DnsModel)
	if err != nil {
		logger.Log.Error("Failed to marshal dns", zap.Error(err))
		return err
	}

	whoisModel, err := json.Marshal(network.WhoisModel)
	if err != nil {
		logger.Log.Error("Failed to marshal whois", zap.Error(err))
		return err
	}

	queryResult := n.database.GetDatabase().Query(query, network.Id, scanId, string(ipAddresses), string(ipRanges), string(httpHeaders), string(dnsModel), string(whoisModel), string(certificates), network.CreatedAt, network.UpdatedAt)
	err = queryResult.Exec()
	if err != nil {
		logger.Log.Error("Failed to execute query", zap.Error(err))
		return err
	}

	return nil
}

type combinedNetwork struct {
	entities.NetworkModel
}

func (n *NetworkRepo) GetScanNetwork(id string) (entities.NetworkModel, error) {
	ctx := context.Background()
	var combinedNetworks entities.NetworkModel

	query := "SELECT id, scan_id, ip_addresses, ip_ranges, http_headers, whois, dns, certificates, created_at, updated_at FROM networks WHERE id = ?"
	err := n.database.Select(ctx, query, &combinedNetworks, id)
	if err != nil {
		logger.Log.Error("Failed to fetch scan result", zap.Error(err))
		return entities.NetworkModel{}, err
	}

	if combinedNetworks.Id == "" {
		return entities.NetworkModel{}, errors.New("no scan result found")
	}

	return combinedNetworks, nil
}

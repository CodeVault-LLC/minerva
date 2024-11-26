package repository

import (
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
func (n *NetworkRepo) Create(network entities.NetworkModel) (uint, error) {
	
}

type combinedNetwork struct {
	entities.NetworkModel
	entities.DnsModel
	entities.WhoisModel
	entities.CertificateModel
}

func (n *NetworkRepo) GetScanNetwork(id uint) (combinedNetwork, error) {
	query := "SELECT * FROM networks LEFT JOIN dns ON networks.id = dns.network_id LEFT JOIN whois ON networks.id = whois.network_id LEFT JOIN certificates ON networks.id = certificates.network_id WHERE scan_id = $1"
	stmt, err := n.db.Preparex(query)
	if err != nil {
		return combinedNetwork{}, err
	}

	var combinedNetworks combinedNetwork
	err = stmt.Get(&combinedNetworks, id)
	if err != nil {
		return combinedNetwork{}, err
	}

	return combinedNetworks, nil
}

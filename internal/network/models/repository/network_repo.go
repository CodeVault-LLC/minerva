package repository

import (
	"crypto/x509"

	"github.com/codevault-llc/minerva/internal/database"
	"github.com/codevault-llc/minerva/internal/network/models/entities"
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
func (n *NetworkRepo) Create(scanId uint, network entities.NetworkModel, whois entities.WhoisModel, dns entities.DnsModel, certificates []x509.Certificate) (uint, error) {
	query := "INSERT INTO networks (scan_id, network, whois, dns, certificates) VALUES (?, ?, ?, ?, ?)"

	queryResult := n.database.GetDatabase().Query(query, scanId, network, whois, dns, certificates)
	err := queryResult.Exec()
	if err != nil {
		return 0, err
	}

	var networkId uint
	err = queryResult.Scan(&networkId)
	if err != nil {
		return 0, err
	}

	return networkId, nil
}

type combinedNetwork struct {
	entities.NetworkModel
	entities.DnsModel
	entities.WhoisModel
	entities.CertificateModel
}

func (n *NetworkRepo) GetScanNetwork(id uint) (combinedNetwork, error) {
	query := "SELECT * FROM networks WHERE id = ?"

	queryResult := n.database.GetDatabase().Query(query, id)
	err := queryResult.Exec()
	if err != nil {
		return combinedNetwork{}, err
	}

	var combinedNetworks combinedNetwork
	err = queryResult.Scan(&combinedNetworks)
	if err != nil {
		return combinedNetwork{}, err
	}

	return combinedNetworks, nil
}

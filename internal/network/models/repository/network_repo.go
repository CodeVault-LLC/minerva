package repository

import (
	"github.com/codevault-llc/minerva/internal/database"
	"github.com/codevault-llc/minerva/internal/network/models/entities"
	"github.com/codevault-llc/minerva/pkg/logger"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type NetworkRepo struct {
	db *sqlx.DB
}

var NetworkRepository *NetworkRepo

// NewNetworkRepository creates a new NetworkRepository
func NewNetworkRepository(db *sqlx.DB) *NetworkRepo {
	return &NetworkRepo{
		db: db,
	}
}

// NetworkRepositoryInterface is the interface for the NetworkRepository
func (n *NetworkRepo) Create(network entities.NetworkModel) (uint, error) {
	tx, err := n.db.Beginx()
	if err != nil {
		return 0, err
	}

	// Get query and values from StructToQuery
	query, values, err := database.StructToQuery(network, "networks")
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	logger.Log.Info("Query", zap.String("query", query), zap.Any("values", values))

	// Insert using InsertStruct with query and values
	networkId, err := database.InsertStruct(tx, query, values)
	if err != nil {
		logger.Log.Error("Failed to insert network", zap.Error(err))

		err := tx.Rollback()
		if err != nil {
			logger.Log.Error("Failed to rollback transaction", zap.Error(err))
		}
		return 0, err
	}

	// Commit transaction
	err = tx.Commit()
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

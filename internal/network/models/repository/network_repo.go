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
		tx.Rollback()
		return 0, err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return networkId, nil
}

func (n *NetworkRepo) GetScanNetwork(id uint) (entities.NetworkModel, error) {
	var network entities.NetworkModel
	if err := n.db.Get(&network, "SELECT * FROM networks WHERE scan_id = $1", id); err != nil {
		return entities.NetworkModel{}, err
	}

	return network, nil
}

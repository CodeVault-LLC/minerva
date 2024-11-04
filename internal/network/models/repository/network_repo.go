package repository

import (
	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/internal/network/models/entities"
	"github.com/jmoiron/sqlx"
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
func (n *NetworkRepo) Create(network entities.NetworkModel) (entities.NetworkModel, error) {
	tx, err := n.db.Beginx()
	if err != nil {
		return entities.NetworkModel{}, err
	}

	query, err := database.StructToQuery(network, "network")
	if err != nil {
		return entities.NetworkModel{}, err
	}

	_, err = database.InsertStruct(tx, query, network)
	if err != nil {
		return entities.NetworkModel{}, err
	}

	err = tx.Commit()
	if err != nil {
		return entities.NetworkModel{}, err
	}

	return network, nil
}

func (n *NetworkRepo) GetScanNetwork(id uint) (entities.NetworkModel, error) {
	var network entities.NetworkModel
	if err := n.db.Get(&network, "SELECT * FROM network WHERE scan_id = $1", id); err != nil {
		return entities.NetworkModel{}, err
	}

	return network, nil
}

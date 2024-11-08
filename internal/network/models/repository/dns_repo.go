package repository

import (
	"github.com/codevault-llc/minerva/internal/database"
	"github.com/codevault-llc/minerva/internal/network/models/entities"
	"github.com/jmoiron/sqlx"
)

type DnsRepo struct {
	db *sqlx.DB
}

func NewDnsRepository(db *sqlx.DB) *DnsRepo {
	return &DnsRepo{
		db: db,
	}
}

var DnsRepository *DnsRepo

func (repository *DnsRepo) SaveDnsResult(dns entities.DnsModel) error {
	tx, err := repository.db.Beginx()
	if err != nil {
		return err
	}

	query, values, err := database.StructToQuery(dns, "dns")
	if err != nil {
		return err
	}

	_, err = database.InsertStruct(tx, query, values)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

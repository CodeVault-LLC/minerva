package repository

import (
	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/internal/network/models/entities"
	"github.com/jmoiron/sqlx"
)

type WhoisRepo struct {
	db *sqlx.DB
}

func NewWhoisRepository(db *sqlx.DB) *WhoisRepo {
	return &WhoisRepo{
		db: db,
	}
}

var WhoisRepository *WhoisRepo

func (repository *WhoisRepo) SaveWhoisResult(whois entities.WhoisModel) error {
	tx, err := repository.db.Beginx()
	if err != nil {
		return err
	}

	query, values, err := database.StructToQuery(whois, "whois")
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

package repository

import (
	"context"

	"github.com/codevault-llc/minerva/internal/contents/models/entities"
	"github.com/codevault-llc/minerva/internal/database"
	"github.com/codevault-llc/minerva/pkg/logger"
	"github.com/codevault-llc/minerva/pkg/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type FingerprintRepo struct {
	database *database.Database
}

func NewFingerprintRepo(database *database.Database) *FingerprintRepo {
	return &FingerprintRepo{database: database}
}

var FingerprintRepository *FingerprintRepo

func (repository *FingerprintRepo) SaveFingerprintResult(fingerprint entities.FingerprintModel) error {
	finding := entities.FingerprintModel{
		Id:            uuid.New().String(),
		ContentId:     fingerprint.ContentId,
		FingerprintId: fingerprint.Id,

		CreatedAt: utils.GetCurrentTime(),
		UpdatedAt: utils.GetCurrentTime(),
	}

	ctx := context.Background()

	query := "INSERT INTO fingerprint (id, content_id, fingerprint_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	_, err := repository.database.Db.Query(ctx, query, finding.Id, finding.ContentId, finding.FingerprintId, finding.CreatedAt, finding.UpdatedAt)
	if err != nil {
		logger.Log.Error("Failed to save finding", zap.Error(err))
		return err
	}

	return nil
}

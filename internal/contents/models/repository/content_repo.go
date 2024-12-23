package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/codevault-llc/minerva/internal/contents/models/entities"
	"github.com/codevault-llc/minerva/internal/database"
	"github.com/codevault-llc/minerva/pkg/logger"
	"go.uber.org/zap"
)

type ContentRepo struct {
	database *database.Database
}

func NewContentRepo(database *database.Database) *ContentRepo {
	return &ContentRepo{database: database}
}

var ContentRepository *ContentRepo

func (repository *ContentRepo) SaveContentResult(content entities.ContentModel) error {
	query := "INSERT INTO content (id, scan_id, file_size, file_type, source, md5, sha1, sha256, duration, tags) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	ctx := context.Background()

	err := repository.database.Db.Exec(ctx, query, content.Id, content.ScanId, content.FileSize, content.FileType, content.Source, content.Md5, content.Sha1, content.Sha256, content.Duration, content.Tags)
	if err != nil {
		logger.Log.Error("Failed to save content result", zap.Error(err))
		return err
	}

	return nil
}

func (repository *ContentRepo) FindContentByMd5(md5 string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var id string

	query := "SELECT id FROM content WHERE md5 = ?"
	err := repository.database.Db.QueryRow(ctx, query, md5).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}

		logger.Log.Error("Failed to fetch content by md5", zap.Error(err))
		return "", err
	}

	return id, nil
}

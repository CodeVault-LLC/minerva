package helper

import (
	"context"

	"github.com/codevault-llc/humblebrag-api/internal/database/models"
)

func AddLicenseToContext(ctx context.Context, license models.LicenseModel) context.Context {
	return context.WithValue(ctx, "license", license)
}

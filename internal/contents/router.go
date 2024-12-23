package contents

import (
	"github.com/codevault-llc/minerva/internal/contents/models/repository"
	"github.com/codevault-llc/minerva/internal/contents/models/viewmodels"
	"github.com/codevault-llc/minerva/pkg/logger"
	"github.com/codevault-llc/minerva/pkg/responder"
	"github.com/codevault-llc/minerva/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func RegisterContentRoutes(router fiber.Router) error {
	router.Get("/findings/:scanID/", getScanFindings)

	return nil
}

// @Summary Get scan findings
// @Description Get scan findings
// @Tags scans
// @Accept json
// @Produce json
// @Param scanID path string true "Scan ID"
// @Success 200 {array} responder.APIResponse{data=models.FindingResponse}
// @Failure 400 {object} responder.APIResponse{error=responder.APIError}
// @Failure 404 {object} responder.APIResponse{error=responder.APIError}
// @Router /scans/{scanID}/findings [get]
func getScanFindings(c *fiber.Ctx) error {
	scanID := c.Params("scanID")

	scanUint, err := utils.ParseUint(scanID)
	if err != nil {
		return responder.CreateError(responder.ErrInvalidRequest).Error
	}

	findings, err := repository.FindingRepository.GetScanFindings(uint(scanUint))
	if err != nil {
		return responder.CreateError(responder.ErrDatabaseQueryFailed).Error
	}

	if len(findings) == 0 {
		logger.Log.Info("No findings found for scan", zap.Uint("scanID", uint(scanUint)))
		return responder.CreateError(responder.ErrResourceNotFound).Error
	}

	responder.WriteJSONResponse(c, responder.CreateSuccessResponse(viewmodels.ConvertFindings(findings), "Successfully retrieved scan findings"))
	return nil
}

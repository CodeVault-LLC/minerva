package contents

import (
	"github.com/codevault-llc/humblebrag-api/internal/contents/models/repository"
	"github.com/codevault-llc/humblebrag-api/internal/contents/models/viewmodels"
	"github.com/codevault-llc/humblebrag-api/pkg/responder"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func RegisterContentRoutes(router fiber.Router) error {
	router.Get("/contents/:scanID/", getScanContents)
	router.Get("/contents/:scanID/findings", getScanFindings)
	router.Get("/contents/:scanID/:contentID", getScanContent)

	return nil
}

// @Summary Get scan contents
// @Description Get scan contents
// @Tags scans
// @Accept json
// @Produce json
// @Param scanID path string true "Scan ID"
// @Success 200 {array} responder.APIResponse{data=models.ContentResponse}
// @Failure 400 {object} responder.APIResponse{error=responder.APIError}
// @Failure 404 {object} responder.APIResponse{error=responder.APIError}
// @Router /scans/{scanID}/contents [get]
func getScanContents(c *fiber.Ctx) error {
	scanID := c.Params("scanID")

	scanUint, err := utils.ParseUint(scanID)
	if err != nil {
		return responder.CreateError(responder.ErrInvalidRequest).Error
	}

	contents, err := repository.ContentRepository.GetScanContents(uint(scanUint))
	if err != nil {
		return responder.CreateError(responder.ErrDatabaseQueryFailed).Error
	}

	if len(contents) == 0 {
		return responder.CreateError(responder.ErrResourceNotFound).Error
	}

	responder.WriteJSONResponse(c, responder.CreateSuccessResponse(contents, "Successfully retrieved scan contents"))
	return nil
}

// @Summary Get scan content
// @Description Get scan content
// @Tags scans
// @Accept json
// @Produce json
// @Param scanID path string true "Scan ID"
// @Param contentID path string true "Content ID"
// @Success 200 {object} responder.APIResponse{data=models.ContentResponse}
// @Failure 400 {object} responder.APIResponse{error=responder.APIError}
// @Failure 404 {object} responder.APIResponse{error=responder.APIError}
func getScanContent(c *fiber.Ctx) error {
	scanID := c.Params("scanID")
	contentID := c.Params("contentID")

	_, err := utils.ParseUint(scanID)
	if err != nil {
		return responder.CreateError(responder.ErrInvalidRequest).Error
	}

	contentUint, err := utils.ParseUint(contentID)
	if err != nil {
		return responder.CreateError(responder.ErrInvalidRequest).Error
	}

	content, err := repository.ContentRepository.GetScanContent(uint(contentUint))
	if err != nil {
		return responder.CreateError(responder.ErrDatabaseQueryFailed).Error
	}

	if content.Id == 0 {
		return responder.CreateError(responder.ErrResourceNotFound).Error
	}

	responder.WriteJSONResponse(c, responder.CreateSuccessResponse(content, "Successfully retrieved scan content"))
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
		return responder.CreateError(responder.ErrResourceNotFound).Error
	}

	responder.WriteJSONResponse(c, responder.CreateSuccessResponse(viewmodels.ConvertFindings(findings), "Successfully retrieved scan findings"))
	return nil
}

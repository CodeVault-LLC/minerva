package network

import (
	"github.com/codevault-llc/minerva/internal/network/models/repository"
	"github.com/codevault-llc/minerva/internal/network/models/viewmodels"
	"github.com/codevault-llc/minerva/pkg/logger"
	"github.com/codevault-llc/minerva/pkg/responder"
	"github.com/codevault-llc/minerva/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func RegisterNetworkRouter(router fiber.Router) error {
	router.Get("/network/:scanID/", getScanNetwork)

	return nil
}

// @Summary Get scan network
// @Description Get scan network
// @Tags scans
// @Accept json
// @Produce json
// @Param scanID path string true "Scan ID"
// @Success 200 {array} responder.APIResponse{data=models.NetworkResponse}
// @Failure 400 {object} responder.APIResponse{error=responder.APIError}
// @Failure 404 {object} responder.APIResponse{error=responder.APIError}
// @Router /scans/{scanID}/network [get]
func getScanNetwork(c *fiber.Ctx) error {
	scanID := c.Params("scanID")

	scanUint, err := utils.ParseUint(scanID)
	if err != nil {
		return responder.CreateError(responder.ErrInvalidRequest).Error
	}

	network, err := repository.NetworkRepository.GetScanNetwork(uint(scanUint))
	if err != nil {
		logger.Log.Error("Failed to get scan network", zap.Error(err))
		return responder.CreateError(responder.ErrDatabaseQueryFailed).Error
	}

	if network.Id == 0 {
		return responder.CreateError(responder.ErrResourceNotFound).Error
	}

	responder.WriteJSONResponse(c, responder.CreateSuccessResponse(viewmodels.ConvertNetwork(network), "Successfully retrieved scan network"))
	return nil
}

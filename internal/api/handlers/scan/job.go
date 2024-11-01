package scan

import (
	"fmt"

	"github.com/codevault-llc/humblebrag-api/internal/core"
	"github.com/codevault-llc/humblebrag-api/internal/models/viewmodels"
	"github.com/codevault-llc/humblebrag-api/pkg/responder"
	"github.com/gofiber/fiber/v2"
)

func RegisterJobRoutes(router fiber.Router) error {
	router.Get("/jobs/:jobID", GetJob)
	return nil
}

func GetJob(c *fiber.Ctx) error {
	jobID := c.Params("jobID")

	fmt.Println("Job ID:", jobID)

	job := core.Scheduler.GetArchivedJob(jobID)
	if job == nil {
		return responder.CreateError(responder.ErrResourceNotFound).Error
	}

	responder.WriteJSONResponse(c, responder.CreateSuccessResponse(viewmodels.ConvertJob(*job), "Job retrieved successfully"))
	return nil
}

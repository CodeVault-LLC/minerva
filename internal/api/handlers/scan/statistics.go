package scan

import (
	"encoding/json"
	"net/http"

	"github.com/codevault-llc/humblebrag-api/helper"
	"github.com/codevault-llc/humblebrag-api/internal/service"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/gorilla/mux"
)

func RegisterStatisticsRoutes(api *mux.Router) {
	api.HandleFunc("/scans/statistics", GetScanStatistics).Methods("GET")
}

type UserStatisticsResponse struct {
	TotalScans          int64 `json:"totalScans"`
	TotalDomainsScanned int64 `json:"totalDomainsScanned"`
	LastScansIn24Hours  int64 `json:"lastScansIn24Hours"`

	MostScannedDomains []models.ScanAPIResponse `json:"mostScannedWebsites"`
}

// @Summary Get scan statistics
// @Description Get scan statistics
// @Tags scans
// @Accept json
// @Produce json
// @Success 200 {object} UserStatisticsResponse
// @Failure 400 {object} types.Error
// @Failure 404 {object} types.Error
// @Router /scans/statistics [get]
func GetScanStatistics(w http.ResponseWriter, r *http.Request) {
	subscription := r.Context().Value("license").(models.LicenseModel)
	if subscription.ID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	totalScans, err := service.GetTotalScans()
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to get total scans")
		return
	}

	totalDomainsScanned, err := service.GetTotalDomainsScanned()
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to get total domains scanned")
		return
	}

	lastScansIn24Hours, err := service.GetRecentScans()
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to get last scans in 24 hours")
		return
	}

	mostScannedDomains, err := service.GetMostScannedDomains()
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to get most scanned domains")
		return
	}

	statistics := UserStatisticsResponse{
		TotalScans:          totalScans,
		TotalDomainsScanned: totalDomainsScanned,
		LastScansIn24Hours:  lastScansIn24Hours,

		MostScannedDomains: models.ConvertScans(mostScannedDomains),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(statistics); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to encode statistics response")
	}
}

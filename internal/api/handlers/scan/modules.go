package scan

import (
	"net/http"

	"github.com/codevault-llc/humblebrag-api/helper"
	"github.com/codevault-llc/humblebrag-api/internal/service"
	"github.com/gorilla/mux"
)

func RegisterModulesRoutes(api *mux.Router) {
	api.HandleFunc("/scans/{scanID}/findings", getScanFindings).Methods("GET")
	api.HandleFunc("/scans/{scanID}/contents", getScanContents).Methods("GET")
	api.HandleFunc("/scans/{scanID}/network", getScanNetwork).Methods("GET")
}

// @Summary Get scan findings
// @Description Get scan findings
// @Tags scans
// @Accept json
// @Produce json
// @Param scanID path string true "Scan ID"
// @Success 200 {array} models.FindingResponse
// @Failure 400 {object} types.Error
// @Failure 404 {object} types.Error
// @Router /scans/{scanID}/findings [get]
func getScanFindings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scanID := vars["scanID"]

	findings, err := service.GetScanFindings(scanID)
	if err != nil {
		helper.RespondWithError(w, 500, "Failed to get scan findings")
		return
	}

	helper.RespondWithJSON(w, 200, findings)
}

// @Summary Get scan contents
// @Description Get scan contents
// @Tags scans
// @Accept json
// @Produce json
// @Param scanID path string true "Scan ID"
// @Success 200 {array} models.ContentResponse
// @Failure 400 {object} types.Error
// @Failure 404 {object} types.Error
// @Router /scans/{scanID}/contents [get]
func getScanContents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scanID := vars["scanID"]

	contents, err := service.GetScanContent(scanID)
	if err != nil {
		helper.RespondWithError(w, 500, "Failed to get scan contents")
		return
	}

	helper.RespondWithJSON(w, 200, contents)
}

// @Summary Get scan network
// @Description Get scan network
// @Tags scans
// @Accept json
// @Produce json
// @Param scanID path string true "Scan ID"
// @Success 200 {array} models.NetworkResponse
// @Failure 400 {object} types.Error
// @Failure 404 {object} types.Error
// @Router /scans/{scanID}/network [get]
func getScanNetwork(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scanID := vars["scanID"]

	network, err := service.GetScanNetwork(scanID)
	if err != nil {
		helper.RespondWithError(w, 500, "Failed to get scan network")
		return
	}

	helper.RespondWithJSON(w, 200, network)
}

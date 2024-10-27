package scan

import (
	"net/http"
	"strconv"

	"github.com/codevault-llc/humblebrag-api/internal/service"
	"github.com/codevault-llc/humblebrag-api/pkg/responder"
	"github.com/gorilla/mux"
)

func RegisterModulesRoutes(api *mux.Router) {
	api.HandleFunc("/scans/{scanID}/findings", getScanFindings).Methods("GET")
	api.HandleFunc("/scans/{scanID}/contents", getScanContents).Methods("GET")
	api.HandleFunc("/scans/{scanID}/network", getScanNetwork).Methods("GET")
	api.HandleFunc("/scans/{scanID}/metadata", getScanMetadata).Methods("GET")
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
func getScanFindings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scanID := vars["scanID"]

	_, err := strconv.ParseUint(scanID, 10, 64)
	if err != nil {
		responder.WriteJSONResponse(w, responder.CreateError(responder.ErrInvalidRequest))
		return
	}

	findings, err := service.GetScanFindings(scanID)
	if err != nil {
		responder.WriteJSONResponse(w, responder.CreateError(responder.ErrDatabaseQueryFailed))
		return
	}

	responder.WriteJSONResponse(w, responder.CreateSuccessResponse(findings, "Successfully retrieved scan findings"))
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
func getScanContents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scanID := vars["scanID"]

	scanIDUint, err := strconv.ParseUint(scanID, 10, 64)
	if err != nil {
		responder.WriteJSONResponse(w, responder.CreateError(responder.ErrInvalidRequest))
		return
	}

	contents, err := service.GetScanContent(uint(scanIDUint))
	if err != nil {
		responder.WriteJSONResponse(w, responder.CreateError(responder.ErrDatabaseQueryFailed))
		return
	}

	responder.WriteJSONResponse(w, responder.CreateSuccessResponse(contents, "Successfully retrieved scan contents"))
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
func getScanNetwork(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scanID := vars["scanID"]

	network, err := service.GetScanNetwork(scanID)
	if err != nil {
		responder.WriteJSONResponse(w, responder.CreateError(responder.ErrDatabaseQueryFailed))
		return
	}

	responder.WriteJSONResponse(w, responder.CreateSuccessResponse(network, "Successfully retrieved scan network"))
}

// @Summary Get scan metadata
// @Description Get scan metadata
// @Tags scans
// @Accept json
// @Produce json
// @Param scanID path string true "Scan ID"
// @Success 200 {object} responder.APIResponse{data=models.MetadataResponse}
// @Failure 400 {object} responder.APIResponse{error=responder.APIError}
// @Failure 404 {object} responder.APIResponse{error=responder.APIError}
// @Router /scans/{scanID}/metadata [get]
func getScanMetadata(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scanID := vars["scanID"]

	metadata, err := service.GetScanMetadataByScanID(scanID)
	if err != nil {
		responder.WriteJSONResponse(w, responder.CreateError(responder.ErrDatabaseQueryFailed))
		return
	}

	responder.WriteJSONResponse(w, responder.CreateSuccessResponse(metadata, "Successfully retrieved scan metadata"))
}

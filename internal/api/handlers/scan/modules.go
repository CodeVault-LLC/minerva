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

package scan

import (
	"encoding/json"
	"net/http"

	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/scanner/secrets"
	"github.com/codevault-llc/humblebrag-api/utils"
	"github.com/gorilla/mux"
)

func ScanRouter(router *mux.Router) {
	router.HandleFunc("/scan", CreateScan).Methods("POST")
}

func CreateScan(w http.ResponseWriter, r *http.Request) {
	var scan models.ScanRequest
	err := json.NewDecoder(r.Body).Decode(&scan)
	if err != nil {
		utils.RespondWithError(w, 400, "Invalid request")
		return
	}

	type ScanResponse struct {
		Secrets []utils.RegexReturn `json:"secrets"`
	}

	var newScan ScanResponse
	secrets := secrets.ScanSecrets(scan.Scripts)
	newScan.Secrets = secrets

	utils.RespondWithJSON(w, 200, newScan)
}

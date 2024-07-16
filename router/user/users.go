package user

import (
	"net/http"

	"github.com/codevault-llc/humblebrag-api/constants"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/utils"
	"github.com/gorilla/mux"
)

func UserRouter(router *mux.Router) {
	router.HandleFunc("/me", GetCurrentUser).Methods("GET")
}

func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	user, ok := constants.SessionManager.Get(r.Context(), "user").(models.User)

	if !ok {
		utils.RespondWithError(w, 500, "Error getting user from context")
		return
	}

	response := utils.ConvertUser(user)

	utils.RespondWithJSON(w, 200, response)
}

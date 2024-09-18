package controller

import (
	"net/http"
	"strconv"

	"github.com/codevault-llc/humblebrag-api/helper"
	"github.com/codevault-llc/humblebrag-api/service"
	"github.com/gorilla/mux"
)

func RegisterNotificationRoutes(router *mux.Router) {
	router.HandleFunc("/notifications/{notificationId}", markNotification).Methods("POST")
}

func markNotification(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	notificationIDStr := vars["notificationId"]

	notificationID, err := strconv.Atoi(notificationIDStr)
	if err != nil {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	if err := service.MarkNotificationAsRead(uint(notificationID)); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to mark notification as read")
		return
	}

	helper.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Notification marked as read"})
}

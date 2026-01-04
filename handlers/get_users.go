package handlers

import (
	"encoding/json"
	"net/http"

	"golang-http-patch/database"
	"golang-http-patch/models"
)

// GetUsers handles GET /users - Get all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	result := database.DB.Find(&users)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

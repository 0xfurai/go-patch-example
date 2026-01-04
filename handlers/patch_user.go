package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"golang-http-patch/database"
	"golang-http-patch/models"
	"golang-http-patch/patch"
	"golang-http-patch/validation"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// PatchUser handles PATCH /users/{id} - Partially update a user
func PatchUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Check if user exists
	var user models.User
	result := database.DB.First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		}
		return
	}

	var dto models.PatchUserDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate DTO
	if !validation.ValidateStruct(w, dto) {
		return
	}

	// Build updates map only for provided fields
	updates := make(map[string]interface{})
	patch.SetUpdate(updates, "name", dto.Name)
	patch.SetUpdate(updates, "age", dto.Age)
	patch.SetUpdate(updates, "phone", dto.Phone)
	patch.SetUpdate(updates, "active", dto.Active)
	patch.SetUpdate(updates, "bio", dto.Bio)
	patch.SetUpdate(updates, "role", dto.Role)
	patch.SetUpdate(updates, "score", dto.Score)
	// Email is immutable and cannot be updated

	if len(updates) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	result = database.DB.Model(&user).Updates(updates)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Reload user to get updated data
	database.DB.First(&user, id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

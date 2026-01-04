package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"golang-http-patch/database"
	"golang-http-patch/models"
	"golang-http-patch/validation"

	"github.com/gorilla/mux"
	"gorm.io/gorm/clause"
)

// UpdateUser handles PUT /users/{id} - Update a user (full update)
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var dto models.UpdateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !validation.ValidateStruct(w, dto) {
		return
	}

	var updatedUser models.User
	result := database.DB.Model(&updatedUser).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"name":   dto.Name,
			"age":    dto.Age,
			"phone":  dto.Phone,
			"active": dto.Active,
			"bio":    dto.Bio,
			"role":   dto.Role,
			"score":  dto.Score,
			// Email is immutable and cannot be updated
		})
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

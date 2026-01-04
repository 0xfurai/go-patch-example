package handlers

import (
	"encoding/json"
	"net/http"

	"golang-http-patch/database"
	"golang-http-patch/models"
	"golang-http-patch/validation"
)

// CreateUser handles POST /users - Create a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var dto models.CreateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate DTO
	if !validation.ValidateStruct(w, dto) {
		return
	}

	user := models.User{
		Name:   dto.Name,
		Email:  dto.Email,
		Age:    dto.Age,
		Phone:  dto.Phone,
		Active: dto.Active,
		Bio:    dto.Bio,
		Role:   dto.Role,
		Score:  dto.Score,
	}
	// Set defaults if not provided
	if dto.Role == "" {
		user.Role = "user"
	}
	// Active defaults to true (handled by GORM default:true in schema)

	result := database.DB.Create(&user)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

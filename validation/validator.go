package validation

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

// InitValidator initializes the validator instance
func InitValidator() {
	Validate = validator.New()
	RegisterPatchValidators(Validate)
}

// ValidateStruct validates a struct and returns validation errors as JSON
func ValidateStruct(w http.ResponseWriter, s interface{}) bool {
	if err := Validate.Struct(s); err != nil {
		var errors []map[string]string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, map[string]string{
				"field":   err.Field(),
				"tag":     err.Tag(),
				"message": GetValidationMessage(err),
			})
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  "Validation failed",
			"errors": errors,
		})
		return false
	}
	return true
}

// GetValidationMessage returns a user-friendly validation message
func GetValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + " is required"
	case "email":
		return err.Field() + " must be a valid email address"
	case "min":
		return err.Field() + " must be at least " + err.Param() + " characters"
	case "max":
		return err.Field() + " must be at most " + err.Param() + " characters"
	case "opt":
		if err.Param() != "" {
			return err.Field() + " failed validation: " + err.Param()
		}
		return err.Field() + " is invalid"
	case "nonull":
		return err.Field() + " cannot be null"
	case "gte":
		return err.Field() + " must be greater than or equal to " + err.Param()
	case "lte":
		return err.Field() + " must be less than or equal to " + err.Param()
	case "gt":
		return err.Field() + " must be greater than " + err.Param()
	case "lt":
		return err.Field() + " must be less than " + err.Param()
	default:
		return err.Field() + " is invalid"
	}
}

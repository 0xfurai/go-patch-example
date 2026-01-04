package validation

import (
	"reflect"
	"strings"

	"golang-http-patch/patch"

	"github.com/go-playground/validator/v10"
)

// RegisterPatchValidators registers custom validators for patch.Optional types
func RegisterPatchValidators(v *validator.Validate) {
	// Note: We don't register a custom type function because it interferes with
	// the validator's ability to call our custom "opt" tag validator

	// opt <inner rules>
	// - unset => valid
	// - null  => valid (use nonull if you want to forbid null)
	// - value => validate.Var(value, <inner rules>)
	// Rules are semicolon-separated, e.g. opt=min=2;max=100
	v.RegisterValidation("opt", func(fl validator.FieldLevel) bool {
		// Get the field value
		fieldVal := fl.Field()

		// Check if it's a pointer to Optional
		if fieldVal.Kind() == reflect.Ptr {
			if fieldVal.IsNil() {
				return true
			}
			fieldVal = fieldVal.Elem()
		}

		oa, ok := fieldVal.Interface().(patch.OptionalAny)
		if !ok {
			// not an Optional -> don't block other structs accidentally
			return true
		}

		// If unset or null, validation passes
		if !oa.IsSet() || oa.IsNull() {
			return true
		}
		// Extract the value and validate it with the inner rules
		val, ok := oa.Any()
		if !ok {
			return true
		}

		// If no validation rules specified, it's valid
		param := fl.Param()
		if param == "" {
			return true
		}

		// Replace semicolons with commas for validator syntax
		// This allows us to use simple semicolon-separated rules: opt=min=2;max=100
		rules := strings.ReplaceAll(param, ";", ",")

		// Create a new validator instance to avoid infinite recursion
		// We can't use the parent validator because it would re-trigger validation on Optional types
		tmpValidator := validator.New()

		// Validate the inner value with the provided rules
		return tmpValidator.Var(val, rules) == nil
	})

	// nonull
	// - unset => valid
	// - null  => invalid
	// - value => valid
	v.RegisterValidation("nonull", func(fl validator.FieldLevel) bool {
		oa, ok := fl.Field().Interface().(patch.OptionalAny)
		if !ok {
			return true
		}
		return !oa.IsSet() || !oa.IsNull()
	})
}

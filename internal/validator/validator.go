package validator

import (
	"slices"
	"strings"
	"unicode/utf8"
)

// Validator is a struct that contains a map of field errors.
type Validator struct {
	FieldErrors map[string]string
}

// Valid returns true if the Validator has no field errors.
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// AddFieldError adds an error message to the map of field errors.
func (v *Validator) AddFieldError(field, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, ok := v.FieldErrors[field]; !ok {
		v.FieldErrors[field] = message
	}
}

// CheckField checks if a condition is met and adds an error message to the map of field errors if it is not.
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// NotBlank returns true if the value is not empty.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxLength returns true if the value is less than or equal to the max length.
func MaxLength(value string, max int) bool {
	return utf8.RuneCountInString(value) <= max
}

// PermittedValue returns true if the value is in the list of permitted values.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

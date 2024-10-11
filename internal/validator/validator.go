package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Validator is a struct that contains a map of field errors.
type Validator struct {
	FieldErrors   map[string]string
	GeneralErrors []string
}

// Valid returns true if the Validator has no field errors.
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.GeneralErrors) == 0
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

// AddGeneralError adds an error message to the slice of general errors.
func (v *Validator) AddGeneralError(message string) {
	v.GeneralErrors = append(v.GeneralErrors, message)
}

// CheckField checks if a condition is met and adds an error message to the map of field errors if it is not.
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// Matches returns true if the value matches the regular expression.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// MaxLength returns true if the value is less than or equal to the max length.
func MaxLength(value string, max int) bool {
	return utf8.RuneCountInString(value) <= max
}

// MinLength returns true if the value is greater than or equal to the min length.
func MinLength(value string, min int) bool {
	return utf8.RuneCountInString(value) >= min
}

// NotBlank returns true if the value is not empty.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// PermittedValue returns true if the value is in the list of permitted values.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

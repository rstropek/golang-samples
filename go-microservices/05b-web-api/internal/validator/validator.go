package validator

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Define a new Validator type which contains a map of validation errors.
type Validator struct {
    Errors map[string]string
}

// New is a helper which creates a new Validator instance with an empty errors map.
func New() *Validator {
    return &Validator{Errors: make(map[string]string)}
}

// Valid returns true if the errors map doesn't contain any entries.
func (v *Validator) Valid() bool {
    return len(v.Errors) == 0
}

// AddError adds an error message to the map (so long as no entry already exists for
// the given key).
func (v *Validator) AddError(key, message string) {
    if _, exists := v.Errors[key]; !exists {
        v.Errors[key] = message
    }
}

// Check adds an error message to the map only if a validation check is not 'ok'.
func (v *Validator) Check(ok bool, key, message string) {
    if !ok {
        v.AddError(key, message)
    }
}

func In(value string, list ...string) bool {
    for i := range list {
        if value == list[i] {
            return true
        }
    }
    return false
}

func IsEmptyUuid(id uuid.UUID) bool {
	return id == uuid.Nil
}

func IsNotEmptyString(text string) bool {
	return len(text) > 0
}

func HasLen(text string, requiredLength int) bool {
	return len(text) == requiredLength
}

func IsGreaterThan(value decimal.Decimal, minValue decimal.Decimal) bool {
	return value.GreaterThan(minValue)
}

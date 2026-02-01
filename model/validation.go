package model

// Default tag type.
const VLD string = "vld"

// Validation represents a validation rule for a struct field.
type Validation struct {
	Key         string
	Type        ValidatorType
	Requirement string
	Groups      []*Group
	Default     string
	// Inner Struct validation
	InnerValidation []Validation
}

// ValidatorMap is a map of validation keys to Validation objects.
// It is used to store and manage multiple validation rules for different struct fields.
type ValidatorMap map[string]Validation

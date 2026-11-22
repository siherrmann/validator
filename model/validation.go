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
	OmitEmpty   bool
	// Inner Struct validation
	InnerValidation []Validation
	// InnerValidationRef is set by GetValidationsFromStruct instead of InnerValidation for a field
	// whose type refers back to a struct that contains it (eg. `Children []Node` in Node).
	// It points to the validations of that struct, so recursive types do not need endless validations.
	InnerValidationRef *[]Validation
	// Parsed AST representation of the requirement to avoid redundant parsing
	RequirementAST *AstValue
	// NoFieldUsage is set by GetValidationsFromStruct if neither the requirement nor any inner validation
	// uses `rel` or a `$` reference. The validator then skips tracking the levels of the input.
	// It is false for validations built by hand, which is always correct, just a bit slower.
	NoFieldUsage bool
}

// UsesFields reports if the validation (or one of its inner validations) may use `rel` or a `$` reference.
func (v *Validation) UsesFields() bool {
	return !v.NoFieldUsage
}

// Inner returns the inner validations of a struct field or of the items of a list of structs.
func (v *Validation) Inner() []Validation {
	if v.InnerValidationRef != nil {
		return *v.InnerValidationRef
	}
	return v.InnerValidation
}

// ValidatorMap is a map of validation keys to Validation objects.
// It is used to store and manage multiple validation rules for different struct fields.
type ValidatorMap map[string]Validation

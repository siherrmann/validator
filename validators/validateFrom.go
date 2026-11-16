package validators

import (
	"fmt"
	"slices"

	"github.com/siherrmann/validator/model"
)

func ValidateFrom[T any](v T, ast *model.AstValue) error {
	from, err := fromWithList(v, ast, false)
	if err != nil {
		return fmt.Errorf("error checking from: %v", err)
	}
	if !from {
		return fmt.Errorf("from %v does not contain value", ast.ConditionValue)
	}

	return nil
}

func ValidateNotFrom[T any](v T, ast *model.AstValue) error {
	notFrom, err := fromWithList(v, ast, true)
	if err != nil {
		return fmt.Errorf("error checking not from: %v", err)
	}
	if !notFrom {
		return fmt.Errorf("from %v does contain value", ast.ConditionValue)
	}

	return nil
}

// ValidateFromValues is `frm`/`nfr` with the value of a referenced field (eg. `frm$/ports.port`)
// instead of a comma separated list. The referenced value can be a single value, a list or a map (its keys).
// Like for `frm`, a list input checks every item and a map input checks every key.
// Values of different types are compared by converting the referenced value to the type of the input (eg. 80 and 80.0).
func ValidateFromValues(v any, from any, not bool, name string) error {
	if fromValues(v, from, not) {
		return nil
	}
	if not {
		return fmt.Errorf("value is part of %v", name)
	}
	return fmt.Errorf("value is not part of %v", name)
}

// fromWithList checks a string value against the list split by the parser,
// so the condition value does not have to be split and converted on every validation.
// All other values (and ASTs built without the parser) are checked with From.
func fromWithList(v any, ast *model.AstValue, not bool) (bool, error) {
	if s, ok := v.(string); ok && ast.List != nil {
		return !not == slices.Contains(ast.List, s), nil
	}
	return From(v, ast.ConditionValue, not)
}

package model

import (
	"fmt"
	"strings"
)

// FieldRef is the parsed path to another field, used by `rel` conditions and `$` references.
// The path is relative to the level of the validated field: `name` is a field on the same level,
// `../name` is a field one level up (repeatable) and `/name` is a field of the top level.
// Inner fields are reached with `.` (eg. `ports.port`), on a list this collects the field of every item.
//
// `${name}` is not a field, but a value injected with [validator.Inject] for a single validation call.
// Then Injected is set and Path holds only the name.
type FieldRef struct {
	Raw      string
	Root     bool
	Up       int
	Path     []string
	Injected bool
}

// ParseFieldRef parses a field path like `name`, `../name`, `/ports.port` or an injected value like `${name}`.
func ParseFieldRef(raw string) (*FieldRef, error) {
	if strings.HasPrefix(raw, "${") {
		name, ok := strings.CutSuffix(raw[2:], "}")
		if !ok || name == "" || strings.ContainsAny(name, "{}") {
			return nil, fmt.Errorf("invalid injected value %q, expected ${name}", raw)
		}
		return &FieldRef{Raw: raw, Path: []string{name}, Injected: true}, nil
	}

	ref := &FieldRef{Raw: raw}
	s := raw
	if strings.HasPrefix(s, "/") {
		ref.Root = true
		s = s[1:]
	}
	for strings.HasPrefix(s, "../") {
		if ref.Root {
			return nil, fmt.Errorf("invalid field path %q, `/` and `../` can not be combined", raw)
		}
		ref.Up++
		s = s[3:]
	}
	if s == "" {
		return nil, fmt.Errorf("invalid field path %q, field name missing", raw)
	}
	ref.Path = strings.Split(s, ".")
	for _, segment := range ref.Path {
		if segment == "" || segment == ".." {
			return nil, fmt.Errorf("invalid field path %q", raw)
		}
	}
	return ref, nil
}

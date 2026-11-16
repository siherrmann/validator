package model

import (
	"fmt"
	"strings"
)

// FieldRef is the parsed path to another field, used by `rel` conditions and `$` references.
// The path is relative to the level of the validated field: `name` is a field on the same level,
// `../name` is a field one level up (repeatable) and `/name` is a field of the top level.
// Inner fields are reached with `.` (eg. `ports.port`), on a list this collects the field of every item.
type FieldRef struct {
	Raw  string
	Root bool
	Up   int
	Path []string
}

// ParseFieldRef parses a field path like `name`, `../name`, `/ports.port`.
func ParseFieldRef(raw string) (*FieldRef, error) {
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

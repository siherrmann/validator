package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JsonMap is a type that represents a map of string keys to any values.
type JsonMap map[string]any

// JsonMap implements the sql.Valuer and sql.Scanner interfaces for JSON data.
func (c JsonMap) Value() (driver.Value, error) {
	return c.Marshal()
}

// Scan implements the sql.Scanner interface for JsonMap.
func (c *JsonMap) Scan(value any) error {
	return c.Unmarshal(value)
}

// Marshal marshals the JsonMap into a JSON byte slice.
func (r JsonMap) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Unmarshal unmarshals a JSON value into the JsonMap.
func (r *JsonMap) Unmarshal(value any) error {
	if s, ok := value.(JsonMap); ok {
		*r = JsonMap(s)
	} else {
		b, ok := value.([]byte)
		if !ok {
			return errors.New("type assertion to []byte failed")
		}
		return json.Unmarshal(b, r)
	}
	return nil
}

// Has checks if the JsonMap contains a specific key.
func (r *JsonMap) Has(key string) bool {
	if _, ok := (*r)[key]; ok {
		return true
	}
	return false
}

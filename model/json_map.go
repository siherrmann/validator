package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JsonMap map[string]interface{}

func (c JsonMap) Value() (driver.Value, error) {
	return c.Marshal()
}

func (c *JsonMap) Scan(value interface{}) error {
	return c.Unmarshal(value)
}

func (r JsonMap) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *JsonMap) Unmarshal(value interface{}) error {
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

func (r *JsonMap) Has(key string) bool {
	if _, ok := (*r)[key]; ok {
		return true
	}
	return false
}

package model

type Validation struct {
	Key        string
	Type       ValidatorType
	Conditions string
}

type ValidatorMap map[string]Validation

type ValidatorType string

const (
	String ValidatorType = "string"
	Int    ValidatorType = "int"
	Float  ValidatorType = "float"
	Bool   ValidatorType = "bool"
	Map    ValidatorType = "map"
	Time   ValidatorType = "time"
)

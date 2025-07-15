package main

import (
	"fmt"

	"github.com/siherrmann/validator/model"
	"github.com/siherrmann/validator/validators"
)

func main() {
	err := validators.ValidateEqual("test", &model.AstValue{ConditionValue: "testing"})
	fmt.Printf("Validation result: %v\n", err)

	err = validators.ValidateEqual(42, &model.AstValue{ConditionValue: "43"})
	fmt.Printf("Validation result: %v\n", err)

	err = validators.ValidateEqual(3.14, &model.AstValue{ConditionValue: "3.15"})
	fmt.Printf("Validation result: %v\n", err)

	err = validators.ValidateEqual(true, &model.AstValue{ConditionValue: "false"})
	fmt.Printf("Validation result: %v\n", err)

	err = validators.ValidateEqual([]int{1, 2, 3}, &model.AstValue{ConditionValue: "2"})
	fmt.Printf("Validation result: %v\n", err)

	err = validators.ValidateEqual(map[string]int{"a": 1, "b": 2}, &model.AstValue{ConditionValue: "1"})
	fmt.Printf("Validation result: %v\n", err)
}

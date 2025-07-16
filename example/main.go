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

	err = validators.ValidateMin(2, &model.AstValue{ConditionValue: "3"})
	fmt.Printf("Validation result: %v\n", err)
	err = validators.ValidateMin(5.1, &model.AstValue{ConditionValue: "5.2"})
	fmt.Printf("Validation result: %v\n", err)
	err = validators.ValidateMin("test", &model.AstValue{ConditionValue: "5"})
	fmt.Printf("Validation result: %v\n", err)
	err = validators.ValidateMin([]int{1, 2}, &model.AstValue{ConditionValue: "3"})
	fmt.Printf("Validation result: %v\n", err)

	err = validators.ValidateMax(10, &model.AstValue{ConditionValue: "5"})
	fmt.Printf("Validation result: %v\n", err)
	err = validators.ValidateMax(3.14, &model.AstValue{ConditionValue: "3.13"})
	fmt.Printf("Validation result: %v\n", err)
	err = validators.ValidateMax("test", &model.AstValue{ConditionValue: "3"})
	fmt.Printf("Validation result: %v\n", err)

	err = validators.ValidateContains("banana", &model.AstValue{ConditionValue: "ane"})
	fmt.Printf("Validation result: %v\n", err)
	err = validators.ValidateContains([]int{1, 2, 3}, &model.AstValue{ConditionValue: "4"})
	fmt.Printf("Validation result: %v\n", err)
	err = validators.ValidateContains([]string{"apple", "banana"}, &model.AstValue{ConditionValue: "orange"})
	fmt.Printf("Validation result: %v\n", err)
	err = validators.ValidateContains(map[string]int{"apple": 1, "banana": 2}, &model.AstValue{ConditionValue: "orange"})
	fmt.Printf("Validation result: %v\n", err)

	err = validators.ValidateNotContains("banana", &model.AstValue{ConditionValue: "ana"})
	fmt.Printf("Validation result: %v\n", err)
	err = validators.ValidateNotContains([]int{1, 2, 3}, &model.AstValue{ConditionValue: "2"})
	fmt.Printf("Validation result: %v\n", err)
	err = validators.ValidateNotContains([]string{"apple", "banana"}, &model.AstValue{ConditionValue: "apple"})
	fmt.Printf("Validation result: %v\n", err)
	err = validators.ValidateNotContains(map[string]int{"apple": 1, "banana": 2}, &model.AstValue{ConditionValue: "banana"})
	fmt.Printf("Validation result: %v\n", err)

	err = validators.ValidateFrom("banana", &model.AstValue{ConditionValue: "orange,apple"})
	fmt.Printf("Validation result from: %v\n", err)
	err = validators.ValidateFrom([]int{1, 2, 3}, &model.AstValue{ConditionValue: "2,3"})
	fmt.Printf("Validation result from: %v\n", err)
	err = validators.ValidateFrom([]string{"$", "$"}, &model.AstValue{ConditionValue: "@,$"})
	fmt.Printf("Validation result from: %v\n", err)

	err = validators.ValidateNotFrom("banana", &model.AstValue{ConditionValue: "orange,banana"})
	fmt.Printf("Validation result: %v\n", err)
	err = validators.ValidateNotFrom([]int{1, 2, 3}, &model.AstValue{ConditionValue: "2,3"})
	fmt.Printf("Validation result: %v\n", err)
	err = validators.ValidateNotFrom([]string{"apple", "banana"}, &model.AstValue{ConditionValue: "orange,banana"})
	fmt.Printf("Validation result: %v\n", err)

	err = validators.ValidateRegex("test123wrong", &model.AstValue{ConditionValue: "^[a-z]+[0-9]+$"})
	fmt.Printf("Validation result: %v\n", err)
	err = validators.ValidateRegex([]string{"test123wrong", "example456"}, &model.AstValue{ConditionValue: "^[a-z]+[0-9]+$"})
	fmt.Printf("Validation result: %v\n", err)
	err = validators.ValidateRegex(map[string]int{"key1wrong": 1, "key2": 2}, &model.AstValue{ConditionValue: "^[a-z]+[0-9]+$"})
	fmt.Printf("Validation result: %v\n", err)
	err = validators.ValidateRegex(52, &model.AstValue{ConditionValue: "^[0-4]+$"})
	fmt.Printf("Validation result: %v\n", err)
	err = validators.ValidateRegex(5.14, &model.AstValue{ConditionValue: "^[0-4]+[.][0-4]+$"})
	fmt.Printf("Validation result: %v\n", err)
	err = validators.ValidateRegex(true, &model.AstValue{ConditionValue: "^false$"})
	fmt.Printf("Validation result: %v\n", err)
}

package validator

import (
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"
	"testing"

	"github.com/siherrmann/validator/model"
	"github.com/siherrmann/validator/parser"
)

type TestStructNoValidation struct {
	String string   `vld:"-"`
	Int    int      `vld:"-"`
	Float  float64  `vld:"-"`
	Array  []string `vld:"-"`
}

type TestStructValidation struct {
	String string   `json:"string" vld:"rex^[a-zA-Z0-9]+$, gr1min3"`
	Int    int      `json:"int" vld:"equ2 || equ3, gr1min3"`
	Float  float64  `json:"float" vld:"equ2 || equ3, gr1min3"`
	Array  []string `json:"array" vld:"min3, gr1min3"`
}

// BenchStructFourFields mirrors the struct used in the module's documented allocation
// baseline: one simple tag, one grouped/OR tag, one float min/max tag, and one field
// with no validation at all.
type BenchStructFourFields struct {
	Field1 string  `vld:"min3 max20"`
	Field2 int     `vld:"(min18 && max99) || equ0"`
	Field3 float64 `vld:"min0.0 max100.0"`
	Field4 string  `vld:"-"`
}

func BenchmarkNoValidation(b *testing.B) {
	noValidation := &TestStructNoValidation{
		String: "test",
		Int:    2,
		Float:  3.0,
		Array:  []string{"", "", ""},
	}

	b.ReportAllocs()
	for b.Loop() {
		err := Validate(noValidation)
		if err != nil {
			b.Fatalf("error no validation %v", err)
		}
	}
}

func BenchmarkUnmarshalAndManualValidation(b *testing.B) {
	jsonString := []byte(`{"string": "test", "int": 2, "float": 3.0, "array": ["", "", ""]}`)

	b.ReportAllocs()
	for b.Loop() {
		// Unmarshal JSON string to struct
		manualValidation := &TestStructValidation{}
		err := json.Unmarshal(jsonString, manualValidation)
		if err != nil {
			b.Fatalf("error unmarshal %v", err)
		}

		// manual validation with minimal regex (only string)
		var errCount int
		match, err := regexp.MatchString("^[a-zA-Z0-9]+$", manualValidation.String)
		if err != nil || !match {
			errCount++
		}
		if manualValidation.Int != 2 && manualValidation.Int != 3 {
			errCount++
		}
		if manualValidation.Float != 2 && manualValidation.Float != 3 {
			errCount++
		}
		if len(manualValidation.Array) < 3 {
			errCount++
		}
		if errCount > 1 {
			b.Fatalf("error manual validation, %d checks failed", errCount)
		}
	}
}

func BenchmarkUnmarshalAndValidate(b *testing.B) {
	jsonString := []byte(`{"string": "test", "int": 2, "float": 3.0, "array": ["", "", ""]}`)

	b.ReportAllocs()
	for b.Loop() {
		// unmarshal and validate
		unmarshalAndValidate := &TestStructValidation{}
		req, err := http.NewRequest("POST", "/", bytes.NewBuffer(jsonString))
		if err != nil {
			b.Fatalf("error creating request %v", err)
		}

		err = UnmarshalAndValidate(req, unmarshalAndValidate)
		if err != nil {
			b.Fatalf("error unmarshal and validate %v", err)
		}
	}
}

// BenchmarkGetValidationsFromStruct measures the uncached, reflection-based extraction
// path directly — this is what Validate paid on every call before the per-type cache
// was added, and what GetValidationsFromStruct still does whenever called directly.
func BenchmarkGetValidationsFromStruct(b *testing.B) {
	s := &BenchStructFourFields{}

	b.ReportAllocs()
	for b.Loop() {
		_, err := GetValidationsFromStruct(s, "vld")
		if err != nil {
			b.Fatalf("error getting validations from struct %v", err)
		}
	}
}

// BenchmarkValidateCachedType measures Validate once the (type, tagType) extraction is
// already warm in validationCache — i.e. the steady-state cost for any type after its
// first call.
func BenchmarkValidateCachedType(b *testing.B) {
	s := &BenchStructFourFields{Field1: "hello", Field2: 50, Field3: 42.5, Field4: "ignored"}

	// Warm the cache before timing.
	if err := Validate(s); err != nil {
		b.Fatalf("error validating %v", err)
	}

	b.ReportAllocs()
	for b.Loop() {
		if err := Validate(s); err != nil {
			b.Fatalf("error validating %v", err)
		}
	}
}

// BenchmarkParseValidationSimple parses a plain, ungrouped requirement.
func BenchmarkParseValidationSimple(b *testing.B) {
	p := parser.NewParser()

	b.ReportAllocs()
	for b.Loop() {
		if _, err := p.ParseValidation("min3 max20"); err != nil {
			b.Fatalf("error parsing validation %v", err)
		}
	}
}

// BenchmarkParseValidationGrouped parses a requirement with a parenthesised group and
// an OR, exercising nested parseGroup/parseCondition calls.
func BenchmarkParseValidationGrouped(b *testing.B) {
	p := parser.NewParser()

	b.ReportAllocs()
	for b.Loop() {
		if _, err := p.ParseValidation("(min18 && max99) || equ0"); err != nil {
			b.Fatalf("error parsing validation %v", err)
		}
	}
}

// BenchmarkLexing measures lexing alone, with no AST built on top, by draining a fresh
// Lexer to EOF. This isolates the allocation cost the parser pays before it constructs
// a single AstValue.
func BenchmarkLexing(b *testing.B) {
	const tag = "(min18 && max99) || equ0"

	b.ReportAllocs()
	for b.Loop() {
		l := parser.NewLexer(tag)
		for {
			tok := l.NextToken()
			if tok.Type == model.LexerEOF {
				break
			}
		}
	}
}

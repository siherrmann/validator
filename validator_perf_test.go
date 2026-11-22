package validator

import (
	"bytes"
	"encoding/json/v2"
	"errors"
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
		manualValidation := &TestStructValidation{}
		err := json.Unmarshal(jsonString, manualValidation)
		if err != nil {
			b.Fatalf("error unmarshal %v", err)
		}

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

func BenchmarkParseValidationSimple(b *testing.B) {
	p := parser.NewParser()

	b.ReportAllocs()
	for b.Loop() {
		if _, err := p.ParseValidation("min3 max20"); err != nil {
			b.Fatalf("error parsing validation %v", err)
		}
	}
}

func BenchmarkParseValidationGrouped(b *testing.B) {
	p := parser.NewParser()

	b.ReportAllocs()
	for b.Loop() {
		if _, err := p.ParseValidation("(min18 && max99) || equ0"); err != nil {
			b.Fatalf("error parsing validation %v", err)
		}
	}
}

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

func BenchmarkValidateRelations(b *testing.B) {
	r := newRelValidator()
	spec := validRelServiceSpec()
	if err := r.Validate(spec); err != nil {
		b.Fatalf("error validating %v", err)
	}

	b.ReportAllocs()
	for b.Loop() {
		if err := r.Validate(spec); err != nil {
			b.Fatalf("error validating %v", err)
		}
	}
}

// BenchmarkManualRelations is the hand written equivalent of BenchmarkValidateRelations as baseline.
func BenchmarkManualRelations(b *testing.B) {
	spec := validRelServiceSpec()
	envKey := regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)
	if err := manualRelRules(spec, envKey); err != nil {
		b.Fatalf("error validating %v", err)
	}

	b.ReportAllocs()
	for b.Loop() {
		if err := manualRelRules(spec, envKey); err != nil {
			b.Fatalf("error validating %v", err)
		}
	}
}

func manualRelRules(s *relServiceSpec, envKey *regexp.Regexp) error {
	guestPath := func(p string) error {
		if len(p) == 0 || p[0] != '/' {
			return errors.New("path must be absolute")
		}
		return nil
	}
	switch s.Image.Type {
	case "oci":
		if s.Image.Ref == "" {
			return errors.New("image ref is required for oci images")
		}
	case "go", "dockerfile":
		if s.Image.Path == "" {
			return errors.New("image path is required")
		}
	default:
		return errors.New("invalid image type")
	}
	if a := s.Autoscale; a != nil {
		if s.Stateful {
			return errors.New("stateful services can not be autoscaled")
		}
		if a.MaxReplicas < a.MinReplicas {
			return errors.New("max_replicas must not be lower than min_replicas")
		}
		if a.TargetCPUPercent == 0 && a.TargetMemoryPercent == 0 {
			return errors.New("autoscale requires a cpu or memory target")
		}
	}
	ports := map[int]bool{}
	for _, p := range s.Ports {
		if p.Port < 1 || p.Port > 65535 {
			return errors.New("invalid port")
		}
		ports[p.Port] = true
	}
	if h := s.HealthCheck; h != nil && h.Type == "http" && h.Path == "" {
		return errors.New("http health checks need a path")
	}
	if in := s.Ingress; in != nil && !ports[in.Port] {
		return errors.New("ingress port is not a service port")
	}
	for k := range s.Env {
		if !envKey.MatchString(k) {
			return errors.New("invalid env name")
		}
	}
	for _, k := range s.Secrets {
		if !envKey.MatchString(k) {
			return errors.New("invalid secret name")
		}
		if _, ok := s.Env[k]; ok {
			return errors.New("secret collides with env")
		}
	}
	if len(s.Volumes) > 0 && !s.Stateful {
		return errors.New("volumes require a stateful service")
	}
	names := map[string]bool{}
	for _, v := range s.Volumes {
		if names[v.Name] {
			return errors.New("duplicate volume")
		}
		names[v.Name] = true
		if err := guestPath(v.MountPath); err != nil {
			return err
		}
	}
	for _, f := range s.Files {
		if err := guestPath(f.Path); err != nil {
			return err
		}
	}
	return nil
}

func BenchmarkValidateInjected(b *testing.B) {
	r := NewValidator()
	request := validInjectRequest()
	if err := r.ValidateContext(injectContext("user"), request); err != nil {
		b.Fatalf("error validating %v", err)
	}

	b.ReportAllocs()
	for b.Loop() {
		if err := r.ValidateContext(injectContext("user"), request); err != nil {
			b.Fatalf("error validating %v", err)
		}
	}
}

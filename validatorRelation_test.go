package validator

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/siherrmann/validator/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The service spec structs express all rules of a hand written `validateRules` function with tags only:
// conditional requirements (rel and =>), comparisons with other fields ($), uniqueness (unq),
// groups and custom functions (fun).
type relImage struct {
	Type string `json:"type" vld:"frmoci,go,dockerfile"`
	Ref  string `json:"ref,omitempty" vld:"rel(type equoci) => min1"`
	Path string `json:"path,omitempty" vld:"rel(type frmgo,dockerfile) => min1"`
}

type relAutoscale struct {
	MinReplicas         int `json:"min_replicas" vld:"min1"`
	MaxReplicas         int `json:"max_replicas" vld:"min$min_replicas"`
	TargetCPUPercent    int `json:"target_cpu_percent" vld:"min1 && max100, gr1min1"`
	TargetMemoryPercent int `json:"target_memory_percent" vld:"min1 && max100, gr1min1"`
}

type relPort struct {
	Port int `json:"port" vld:"min1 && max65535"`
}

type relHealthCheck struct {
	Type string `json:"type" vld:"frmhttp,tcp"`
	Path string `json:"path,omitempty" vld:"rel(type equhttp) => min1"`
}

type relIngress struct {
	Port int `json:"port" vld:"frm$../ports.port"`
}

type relVolume struct {
	Name      string `json:"name" vld:"min1"`
	MountPath string `json:"mount_path" vld:"funguestPath"`
}

type relFile struct {
	Path string `json:"path" vld:"funguestPath"`
}

type relServiceSpec struct {
	Image       relImage          `json:"image" vld:"-"`
	Stateful    bool              `json:"stateful"`
	Autoscale   *relAutoscale     `json:"autoscale,omitempty" vld:"rel(stateful equtrue) => max0"`
	Ports       []relPort         `json:"ports" vld:"-"`
	HealthCheck *relHealthCheck   `json:"health_check,omitempty" vld:"-"`
	Ingress     *relIngress       `json:"ingress,omitempty" vld:"-"`
	Env         map[string]string `json:"env" vld:"rex^[A-Z_][A-Z0-9_]*$"`
	Secrets     []string          `json:"secrets" vld:"rex^[A-Z_][A-Z0-9_]*$ && nfr$env"`
	Volumes     []relVolume       `json:"volumes" vld:"(rel(stateful equfalse) => max0) && unqname"`
	Files       []relFile         `json:"files" vld:"-"`
}

func newRelValidator() *Validator {
	r := NewValidator()
	r.AddValidationFunc(func(input any, astValue *model.AstValue) error {
		if s, _ := input.(string); !strings.HasPrefix(s, "/") {
			return fmt.Errorf("path must be absolute")
		}
		return nil
	}, "guestPath")
	return r
}

func validRelServiceSpec() *relServiceSpec {
	return &relServiceSpec{
		Image:       relImage{Type: "oci", Ref: "nginx"},
		Stateful:    true,
		Ports:       []relPort{{Port: 80}, {Port: 443}},
		HealthCheck: &relHealthCheck{Type: "http", Path: "/health"},
		Ingress:     &relIngress{Port: 443},
		Env:         map[string]string{"FOO": "1"},
		Secrets:     []string{"BAR"},
		Volumes:     []relVolume{{Name: "data", MountPath: "/data"}, {Name: "logs", MountPath: "/logs"}},
		Files:       []relFile{{Path: "/etc/config"}},
	}
}

func TestValidateRelations(t *testing.T) {
	tests := []struct {
		name        string
		modify      func(s *relServiceSpec)
		errContains string
	}{
		{name: "Valid spec", modify: func(s *relServiceSpec) {}},
		{
			name:        "Oci image without ref",
			modify:      func(s *relServiceSpec) { s.Image.Ref = "" },
			errContains: "field ref invalid",
		},
		{
			name:        "Go image without path",
			modify:      func(s *relServiceSpec) { s.Image = relImage{Type: "go"} },
			errContains: "field path invalid",
		},
		{
			name:   "Go image with path and without ref",
			modify: func(s *relServiceSpec) { s.Image = relImage{Type: "go", Path: "./cmd"} },
		},
		{
			name:        "Stateful service with autoscale",
			modify:      func(s *relServiceSpec) { s.Autoscale = &relAutoscale{1, 2, 50, 0} },
			errContains: "field autoscale invalid",
		},
		{
			name: "Stateless service with autoscale",
			modify: func(s *relServiceSpec) {
				s.Stateful, s.Volumes, s.Autoscale = false, nil, &relAutoscale{1, 2, 50, 0}
			},
		},
		{
			name: "Max replicas lower than min replicas",
			modify: func(s *relServiceSpec) {
				s.Stateful, s.Volumes, s.Autoscale = false, nil, &relAutoscale{3, 2, 50, 0}
			},
			errContains: "field max_replicas invalid",
		},
		{
			name: "Max replicas equal to min replicas",
			modify: func(s *relServiceSpec) {
				s.Stateful, s.Volumes, s.Autoscale = false, nil, &relAutoscale{2, 2, 0, 50}
			},
		},
		{
			name: "Autoscale without target",
			modify: func(s *relServiceSpec) {
				s.Stateful, s.Volumes, s.Autoscale = false, nil, &relAutoscale{1, 2, 0, 0}
			},
			errContains: "group gr1",
		},
		{
			name:        "Http health check without path",
			modify:      func(s *relServiceSpec) { s.HealthCheck.Path = "" },
			errContains: "field path invalid",
		},
		{
			name:   "Tcp health check without path",
			modify: func(s *relServiceSpec) { s.HealthCheck = &relHealthCheck{Type: "tcp"} },
		},
		{
			name:        "Ingress port is not a service port",
			modify:      func(s *relServiceSpec) { s.Ingress.Port = 8080 },
			errContains: "value is not part of ../ports.port",
		},
		{
			name:        "Invalid env name",
			modify:      func(s *relServiceSpec) { s.Env["foo"] = "1" },
			errContains: "field env invalid",
		},
		{
			name:        "Invalid secret name",
			modify:      func(s *relServiceSpec) { s.Secrets = []string{"bar"} },
			errContains: "field secrets invalid",
		},
		{
			name:        "Secret collides with env",
			modify:      func(s *relServiceSpec) { s.Secrets = []string{"FOO"} },
			errContains: "value is part of env",
		},
		{
			name:        "Volumes on stateless service",
			modify:      func(s *relServiceSpec) { s.Stateful = false },
			errContains: "field volumes invalid",
		},
		{
			name:        "Duplicate volume",
			modify:      func(s *relServiceSpec) { s.Volumes[1].Name = "data" },
			errContains: "value data of name is not unique",
		},
		{
			name:        "Invalid mount path",
			modify:      func(s *relServiceSpec) { s.Volumes[1].MountPath = "logs" },
			errContains: "path must be absolute",
		},
		{
			name:        "Invalid file path",
			modify:      func(s *relServiceSpec) { s.Files[0].Path = "config" },
			errContains: "path must be absolute",
		},
	}

	r := newRelValidator()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			spec := validRelServiceSpec()
			test.modify(spec)
			err := r.Validate(spec)
			if test.errContains == "" {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), test.errContains)
			}
		})
	}
}

func TestValidateRelationsJson(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected relImage
		wantErr  bool
	}{
		{name: "Missing required ref", input: `{"image":{"type":"oci"}}`, wantErr: true},
		{name: "Empty required ref", input: `{"image":{"type":"oci","ref":""}}`, wantErr: true},
		{name: "Given ref", input: `{"image":{"type":"oci","ref":"nginx"}}`, expected: relImage{Type: "oci", Ref: "nginx"}},
		{name: "Missing ref that is not required", input: `{"image":{"type":"go","path":"./cmd"}}`, expected: relImage{Type: "go", Path: "./cmd"}},
		{name: "Missing required path", input: `{"image":{"type":"dockerfile"}}`, wantErr: true},
	}

	type request struct {
		Image relImage `json:"image" vld:"-"`
	}

	r := NewValidator()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonMap := map[string]any{}
			require.NoError(t, json.Unmarshal([]byte(test.input), &jsonMap))

			result := &request{}
			err := r.ValidateAndUpdate(jsonMap, result)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result.Image)
			}
		})
	}
}

func TestValidateRelationPaths(t *testing.T) {
	type Leaf struct {
		Up   int `json:"up" vld:"equ$../value"`
		Root int `json:"root" vld:"equ$/value"`
		Two  int `json:"two" vld:"equ$../../value"`
	}
	type Middle struct {
		Value int  `json:"value"`
		Leaf  Leaf `json:"leaf" vld:"-"`
	}
	type Top struct {
		Value  int    `json:"value"`
		Middle Middle `json:"middle" vld:"-"`
	}

	r := NewValidator()
	assert.NoError(t, r.Validate(&Top{Value: 1, Middle: Middle{Value: 2, Leaf: Leaf{Up: 2, Root: 1, Two: 1}}}))

	err := r.Validate(&Top{Value: 1, Middle: Middle{Value: 2, Leaf: Leaf{Up: 1, Root: 1, Two: 1}}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "field up invalid")

	err = r.Validate(&Top{Value: 1, Middle: Middle{Value: 2, Leaf: Leaf{Up: 2, Root: 2, Two: 1}}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "field root invalid")

	type AboveTop struct {
		Value int `json:"value" vld:"equ$../value"`
	}
	err = r.Validate(&AboveTop{Value: 1})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "goes above the top level")
}

func TestValidateReferences(t *testing.T) {
	t.Run("Missing referenced field", func(t *testing.T) {
		validations := []model.Validation{{Key: "max", Type: model.Int, Requirement: "min$min"}}
		_, err := NewValidator().ValidateWithValidation(map[string]any{"max": 3.0}, validations)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "referenced field min is not given")
	})

	t.Run("Reference to json number", func(t *testing.T) {
		validations := []model.Validation{
			{Key: "min", Type: model.Int, Requirement: "-"},
			{Key: "max", Type: model.Int, Requirement: "min$min"},
		}
		r := NewValidator()
		_, err := r.ValidateWithValidation(map[string]any{"min": 2.0, "max": 3.0}, validations)
		assert.NoError(t, err)
		_, err = r.ValidateWithValidation(map[string]any{"min": 4.0, "max": 3.0}, validations)
		assert.Error(t, err)
	})

	t.Run("Reference to a list needs frm or nfr", func(t *testing.T) {
		validations := []model.Validation{{Key: "value", Type: model.Int, Requirement: "equ$list"}}
		_, err := NewValidator().ValidateWithValidation(map[string]any{"value": 1.0, "list": []any{1.0}}, validations)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "has to be a single value")
	})

	t.Run("Relation to missing field is null", func(t *testing.T) {
		validations := []model.Validation{{Key: "value", Type: model.String, Requirement: "rel(other equyes) => min1"}}
		_, err := NewValidator().ValidateWithValidation(map[string]any{"value": ""}, validations)
		assert.NoError(t, err)
	})

	t.Run("Relation as plain condition", func(t *testing.T) {
		validations := []model.Validation{{Key: "value", Type: model.String, Requirement: "rel(other equyes)"}}
		r := NewValidator()
		_, err := r.ValidateWithValidation(map[string]any{"value": "", "other": "yes"}, validations)
		assert.NoError(t, err)
		_, err = r.ValidateWithValidation(map[string]any{"value": "", "other": "no"}, validations)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "related field other invalid")
	})

	t.Run("Without the other fields", func(t *testing.T) {
		r := NewValidator()
		err := r.ValidateValueWithParser("x", &model.Validation{Requirement: "rel(other equyes) => min1"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "can only be validated together with the other fields")

		err = r.ValidateValueWithParser("x", &model.Validation{Requirement: "equ$other || min1"})
		require.Error(t, err, "A missing scope is a configuration error, also inside an OR group")
	})
}

func TestValidateImplication(t *testing.T) {
	tests := []struct {
		requirement string
		value       any
		wantErr     bool
	}{
		{requirement: "equoci => min5", value: "oci", wantErr: true},
		{requirement: "min3 => max5", value: "ab"},
		{requirement: "min3 => max5", value: "abcd"},
		{requirement: "min3 => max5", value: "abcdef", wantErr: true},
		{requirement: "min1 && max3 => equab || equabc", value: "abc"},
		{requirement: "min1 && max3 => equab || equabc", value: "a", wantErr: true},
		{requirement: "min1 && max3 => equab || equabc", value: "abcde"},
		{requirement: "min1 => min2 => min3", value: "ab", wantErr: true},
		{requirement: "min1 => min2 => min3", value: "a"},
		{requirement: "(min3 => max5) && conb", value: "abcdef", wantErr: true},
		{requirement: "(min3 => max5) && conb", value: "ab"},
	}

	r := NewValidator()
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v with %v", test.requirement, test.value), func(t *testing.T) {
			err := r.ValidateValueWithParser(test.value, &model.Validation{Requirement: test.requirement})
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	t.Run("Error names the premise", func(t *testing.T) {
		err := r.ValidateValueWithParser("abcdef", &model.Validation{Requirement: "min3 => max5"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "required because of min'3'")
	})
}

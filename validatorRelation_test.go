package validator

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/siherrmann/validator/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

type injectRequest struct {
	Owner  string `json:"owner" vld:"rel(${role} equadmin) || equ${user_id}"`
	Role   string `json:"role,omitempty" vld:"rel(${role} equadmin) || max0"`
	Key    string `json:"key" vld:"equ${api_key}"`
	Amount int    `json:"amount" vld:"min1 && max${credit_limit}"`
	Region string `json:"region" vld:"frm${allowed_regions}"`
}

func injectContext(role string) context.Context {
	ctx := Inject(context.Background(), "role", role)
	ctx = Inject(ctx, "user_id", "u1")
	ctx = Inject(ctx, "api_key", "s3cr3t-key")
	ctx = Inject(ctx, "credit_limit", 100)
	ctx = Inject(ctx, "allowed_regions", []string{"region-eu", "region-us"})
	return ctx
}

func validInjectRequest() *injectRequest {
	return &injectRequest{Owner: "u1", Key: "s3cr3t-key", Amount: 50, Region: "region-eu"}
}

func TestRunFieldCondition(t *testing.T) {
	t.Run("Service spec", func(t *testing.T) {
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
	})

	t.Run("Json input with missing keys", func(t *testing.T) {
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
	})

	t.Run("Paths to other levels", func(t *testing.T) {
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
	})

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

func TestRunInjectedCondition(t *testing.T) {
	t.Run("Permissions, secret, limit and list", func(t *testing.T) {
		tests := []struct {
			name        string
			role        string
			modify      func(r *injectRequest)
			errContains string
		}{
			{name: "Valid request", role: "user", modify: func(r *injectRequest) {}},
			{name: "User changes other owner", role: "user", modify: func(r *injectRequest) { r.Owner = "u2" }, errContains: "field owner invalid"},
			{name: "Admin changes other owner", role: "admin", modify: func(r *injectRequest) { r.Owner = "u2" }},
			{name: "User sets role", role: "user", modify: func(r *injectRequest) { r.Role = "admin" }, errContains: "field role invalid"},
			{name: "Admin sets role", role: "admin", modify: func(r *injectRequest) { r.Role = "admin" }},
			{name: "Wrong key", role: "user", modify: func(r *injectRequest) { r.Key = "guess" }, errContains: "field key invalid"},
			{name: "Amount above limit", role: "user", modify: func(r *injectRequest) { r.Amount = 101 }, errContains: "field amount invalid"},
			{name: "Region not allowed", role: "user", modify: func(r *injectRequest) { r.Region = "asia" }, errContains: "field region invalid"},
		}

		r := NewValidator()
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				request := validInjectRequest()
				test.modify(request)
				err := r.ValidateContext(injectContext(test.role), request)
				if test.errContains == "" {
					assert.NoError(t, err)
				} else {
					require.Error(t, err)
					assert.Contains(t, err.Error(), test.errContains)
				}
			})
		}
	})

	t.Run("Errors do not contain values", func(t *testing.T) {
		request := &injectRequest{Owner: "u2", Role: "admin", Key: "guess", Amount: 101, Region: "asia"}
		err := NewValidator().ValidateContext(injectContext("user"), request)
		require.Error(t, err)

		for _, secret := range []string{"s3cr3t-key", "u1", "100", "region-eu", "region-us"} {
			assert.NotContains(t, err.Error(), secret, "Expected the injected value not to be in the error")
		}
		for _, name := range []string{"${api_key}", "${user_id}", "${credit_limit}", "${allowed_regions}", "${role}"} {
			assert.Contains(t, err.Error(), name, "Expected the name of the injected value in the error")
		}
	})

	r := NewValidator()

	t.Run("Missing value is a configuration error", func(t *testing.T) {
		_, err := r.ValidateWithValidationContext(context.Background(), map[string]any{"key": "x"},
			[]model.Validation{{Key: "key", Type: model.String, Requirement: "equ${api_key}"}})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "injected value ${api_key} is not given")
	})

	t.Run("Missing value is a configuration error also in an OR group", func(t *testing.T) {
		_, err := r.ValidateWithValidationContext(context.Background(), map[string]any{"key": "x"},
			[]model.Validation{{Key: "key", Type: model.String, Requirement: "equ${api_key} || min1"}})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "injected value ${api_key} is not given")
	})

	t.Run("Missing value is null in rel", func(t *testing.T) {
		_, err := r.ValidateWithValidationContext(context.Background(), map[string]any{"role": "admin"},
			[]model.Validation{{Key: "role", Type: model.String, Requirement: "rel(${role} equadmin) || max0"}})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "field role invalid")
	})

	t.Run("Without context", func(t *testing.T) {
		err := r.Validate(validInjectRequest())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "is not given, it has to be injected")
	})

	t.Run("Secret as bytes", func(t *testing.T) {
		type Request struct {
			Key string `json:"key" vld:"equ${api_key}"`
		}
		ctx := Inject(context.Background(), "api_key", []byte("s3cr3t"))
		r := NewValidator()
		assert.NoError(t, r.ValidateContext(ctx, &Request{Key: "s3cr3t"}))
		assert.Error(t, r.ValidateContext(ctx, &Request{Key: "s3cr3"}))
	})

	t.Run("Not equal", func(t *testing.T) {
		type Request struct {
			NewPassword string `json:"new_password" vld:"min8 && neq${current_password}"`
		}
		ctx := Inject(context.Background(), "current_password", "old-password")
		r := NewValidator()
		assert.NoError(t, r.ValidateContext(ctx, &Request{NewPassword: "new-password"}))

		err := r.ValidateContext(ctx, &Request{NewPassword: "old-password"})
		require.Error(t, err)
		assert.NotContains(t, err.Error(), "old-password")
	})

	t.Run("Concurrent requests", func(t *testing.T) {
		type Request struct {
			Owner string `json:"owner" vld:"equ${user_id}"`
		}

		r := NewValidator()
		var wg sync.WaitGroup
		for i := range 50 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				userID := fmt.Sprintf("u%d", i)
				ctx := Inject(context.Background(), "user_id", userID)
				for range 50 {
					assert.NoError(t, r.ValidateContext(ctx, &Request{Owner: userID}))
					assert.Error(t, r.ValidateContext(ctx, &Request{Owner: "other"}))
				}
			}()
		}
		wg.Wait()
	})
}

func TestInjectedError(t *testing.T) {
	v := &model.AstValue{ConditionType: model.EQUAL, Ref: &model.FieldRef{Raw: "${api_key}", Path: []string{"api_key"}, Injected: true}}
	assert.EqualError(t, injectedError(v), "value does not fulfill equ${api_key}")
}

func TestRunResolvedCondition(t *testing.T) {
	r := NewValidator()
	condition := func(conditionType model.ConditionType) *model.AstValue {
		return &model.AstValue{Type: model.CONDITION, ConditionType: conditionType, Ref: &model.FieldRef{Raw: "other", Path: []string{"other"}}}
	}

	tests := []struct {
		name    string
		input   any
		v       *model.AstValue
		value   any
		wantErr bool
	}{
		{name: "Single value", input: 5, v: condition(model.MIN_VALUE), value: 3.0},
		{name: "Single value not fulfilled", input: 2, v: condition(model.MIN_VALUE), value: 3.0, wantErr: true},
		{name: "List for frm", input: 2, v: condition(model.FROM), value: []any{1.0, 2.0}},
		{name: "List for nfr", input: 2, v: condition(model.NOT_FROM), value: []any{1.0, 2.0}, wantErr: true},
		{name: "List for other conditions", input: 2, v: condition(model.EQUAL), value: []any{2.0}, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fatal, err := r.runResolvedCondition(context.Background(), test.input, test.v, test.value, nil)
			assert.False(t, fatal)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestResolveFieldRef(t *testing.T) {
	top := &scope{current: map[string]any{
		"name":  "top",
		"ports": []any{map[string]any{"port": 80.0}, map[string]any{"port": 443.0}},
	}}
	middle := &scope{current: map[string]any{"name": "middle"}, parent: top}
	leaf := &scope{current: map[string]any{"name": "leaf", "inner": map[string]any{"value": 1.0}}, parent: middle}

	tests := []struct {
		name      string
		ref       string
		expected  any
		wantFound bool
		wantErr   bool
	}{
		{name: "Same level", ref: "name", expected: "leaf", wantFound: true},
		{name: "Inner field", ref: "inner.value", expected: 1.0, wantFound: true},
		{name: "One level up", ref: "../name", expected: "middle", wantFound: true},
		{name: "Two levels up", ref: "../../name", expected: "top", wantFound: true},
		{name: "Top level", ref: "/name", expected: "top", wantFound: true},
		{name: "Field of every item in a list", ref: "/ports.port", expected: []any{80.0, 443.0}, wantFound: true},
		{name: "Missing field", ref: "missing"},
		{name: "Above the top level", ref: "../../../name", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ref, err := model.ParseFieldRef(test.ref)
			require.NoError(t, err)

			value, found, err := resolveFieldRef(leaf, ref)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.wantFound, found)
			assert.Equal(t, test.expected, value)
		})
	}
}

func TestLookupField(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		field     string
		expected  any
		wantFound bool
	}{
		{name: "JsonMap", value: map[string]any{"a": 1}, field: "a", expected: 1, wantFound: true},
		{name: "Missing key", value: map[string]any{"a": 1}, field: "b"},
		{name: "Typed map", value: map[string]string{"a": "b"}, field: "a", expected: "b", wantFound: true},
		{name: "List of JsonMaps", value: []any{map[string]any{"a": 1}, map[string]any{"b": 2}}, field: "a", expected: []any{1}, wantFound: true},
		{name: "Typed list", value: []map[string]int{{"a": 1}, {"a": 2}}, field: "a", expected: []any{1, 2}, wantFound: true},
		{name: "Null", value: nil, field: "a"},
		{name: "Single value", value: "a", field: "a"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, found := lookupField(test.value, test.field)
			assert.Equal(t, test.wantFound, found)
			assert.Equal(t, test.expected, value)
		})
	}
}

func TestPremiseToString(t *testing.T) {
	condition := &model.AstValue{Type: model.CONDITION, ConditionType: model.MIN_VALUE, ConditionValue: "3", Operator: model.IMPLIES}
	assert.Equal(t, "min'3'", premiseToString(condition))

	group := &model.AstValue{Type: model.GROUP, Operator: model.IMPLIES, ConditionGroup: model.ConditionGroup{
		{Type: model.CONDITION, ConditionType: model.MIN_VALUE, ConditionValue: "1", Operator: model.AND},
		{Type: model.CONDITION, ConditionType: model.MAX_VALUE, ConditionValue: "5"},
	}}
	assert.Equal(t, "(min'1' && max'5')", premiseToString(group))
}

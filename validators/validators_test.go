package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEqual(t *testing.T) {
	type args struct {
		v       any
		compare any
	}
	tests := []struct {
		name  string
		args  args
		equal bool
	}{
		{
			name: "Valid equal",
			args: args{
				v:       5,
				compare: 5,
			},
			equal: true,
		},
		{
			name: "Invalid equal",
			args: args{
				v:       5,
				compare: 6,
			},
			equal: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Equal(test.args.v, test.args.compare)
			if test.equal {
				assert.True(t, result, "Expected equal but got not equal")
			} else {
				assert.False(t, result, "Expected not equal but got equal")
			}
		})
	}
}

func TestNotEqual(t *testing.T) {
	type args struct {
		v       any
		compare any
	}
	tests := []struct {
		name  string
		args  args
		equal bool
	}{
		{
			name: "Valid not equal",
			args: args{
				v:       5,
				compare: 6,
			},
			equal: true,
		},
		{
			name: "Invalid not equal",
			args: args{
				v:       5,
				compare: 5,
			},
			equal: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := NotEqual(test.args.v, test.args.compare)
			if test.equal {
				assert.True(t, result, "Expected not equal but got equal")
			} else {
				assert.False(t, result, "Expected equal but got not equal")
			}
		})
	}
}

func TestMin(t *testing.T) {
	type args struct {
		v   int
		min int
	}
	tests := []struct {
		name  string
		args  args
		valid bool
	}{
		{
			name: "Valid min",
			args: args{
				v:   10,
				min: 5,
			},
			valid: true,
		},
		{
			name: "Invalid min",
			args: args{
				v:   10,
				min: 15,
			},
			valid: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Min(test.args.v, test.args.min)
			assert.Equal(t, test.valid, result, "Expected validation result does not match")
		})
	}
}

func TestMax(t *testing.T) {
	type args struct {
		v   int
		max int
	}
	tests := []struct {
		name  string
		args  args
		valid bool
	}{
		{
			name: "Valid max",
			args: args{
				v:   10,
				max: 15,
			},
			valid: true,
		},
		{
			name: "Invalid max",
			args: args{
				v:   10,
				max: 5,
			},
			valid: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Max(test.args.v, test.args.max)
			assert.Equal(t, test.valid, result, "Expected validation result does not match")
		})
	}
}

func TestContains(t *testing.T) {
	type args struct {
		v       any
		element string
	}
	tests := []struct {
		name    string
		args    args
		valid   bool
		wantErr bool
	}{
		{
			name: "Valid contains string",
			args: args{
				v:       "banana",
				element: "ana",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid contains string",
			args: args{
				v:       "banana",
				element: "ane",
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Valid contains array string",
			args: args{
				v:       []string{"apple", "banana", "cherry"},
				element: "banana",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid contains array string",
			args: args{
				v:       []string{"apple", "banana", "cherry"},
				element: "orange",
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Valid contains array bool",
			args: args{
				v:       []bool{true, false, true},
				element: "false",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid contains array bool",
			args: args{
				v:       []bool{true, true, true},
				element: "false",
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid compare type bool",
			args: args{
				v:       []bool{true, true, true},
				element: "banana",
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid contains array int",
			args: args{
				v:       []int{1, 2, 3},
				element: "2",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid contains array int",
			args: args{
				v:       []int{1, 2, 3},
				element: "4",
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid compare type int",
			args: args{
				v:       []int{1, 2, 3},
				element: "banana",
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid contains array int8",
			args: args{
				v:       []int8{1, 2, 3},
				element: "2",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid contains array int8",
			args: args{
				v:       []int8{1, 2, 3},
				element: "4",
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid compare type int8",
			args: args{
				v:       []int8{1, 2, 3},
				element: "banana",
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid contains array int16",
			args: args{
				v:       []int16{1, 2, 3},
				element: "2",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid contains array int16",
			args: args{
				v:       []int16{1, 2, 3},
				element: "4",
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid compare type int16",
			args: args{
				v:       []int16{1, 2, 3},
				element: "banana",
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid contains array int32",
			args: args{
				v:       []int32{1, 2, 3},
				element: "2",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid contains array int32",
			args: args{
				v:       []int32{1, 2, 3},
				element: "4",
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid compare type int32",
			args: args{
				v:       []int32{1, 2, 3},
				element: "banana",
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid contains array int64",
			args: args{
				v:       []int64{1, 2, 3},
				element: "2",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid contains array int64",
			args: args{
				v:       []int64{1, 2, 3},
				element: "4",
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid compare type int64",
			args: args{
				v:       []int64{1, 2, 3},
				element: "banana",
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid contains array uint",
			args: args{
				v:       []uint{1, 2, 3},
				element: "2",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid contains array uint",
			args: args{
				v:       []uint{1, 2, 3},
				element: "4",
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid compare type uint",
			args: args{
				v:       []uint{1, 2, 3},
				element: "banana",
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid contains array uint8",
			args: args{
				v:       []uint8{1, 2, 3},
				element: "2",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid contains array uint8",
			args: args{
				v:       []uint8{1, 2, 3},
				element: "4",
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid compare type uint8",
			args: args{
				v:       []uint8{1, 2, 3},
				element: "banana",
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid contains array uint16",
			args: args{
				v:       []uint16{1, 2, 3},
				element: "2",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid contains array uint16",
			args: args{
				v:       []uint16{1, 2, 3},
				element: "4",
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid compare type uint16",
			args: args{
				v:       []uint16{1, 2, 3},
				element: "banana",
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid contains array uint32",
			args: args{
				v:       []uint32{1, 2, 3},
				element: "2",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid contains array uint32",
			args: args{
				v:       []uint32{1, 2, 3},
				element: "4",
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid compare type uint32",
			args: args{
				v:       []uint32{1, 2, 3},
				element: "banana",
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid contains array uint64",
			args: args{
				v:       []uint64{1, 2, 3},
				element: "2",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid contains array uint64",
			args: args{
				v:       []uint64{1, 2, 3},
				element: "4",
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid compare type uint64",
			args: args{
				v:       []uint64{1, 2, 3},
				element: "banana",
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid contains array float32",
			args: args{
				v:       []float32{1.1, 2.2, 3.3},
				element: "2.2",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid contains array float32",
			args: args{
				v:       []float32{1.1, 2.2, 3.3},
				element: "4.4",
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid compare type float32",
			args: args{
				v:       []float32{1, 2, 3},
				element: "banana",
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid contains array float64",
			args: args{
				v:       []float64{1.1, 2.2, 3.3},
				element: "2.2",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid contains array float64",
			args: args{
				v:       []float64{1.1, 2.2, 3.3},
				element: "4.4",
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid compare type float64",
			args: args{
				v:       []float64{1, 2, 3},
				element: "banana",
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid contains map with string values",
			args: args{
				v:       map[string]string{"a": "1", "b": "2", "c": "3"},
				element: "a",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid contains map with string values",
			args: args{
				v:       map[string]string{"a": "1", "b": "2", "c": "3"},
				element: "d",
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Valid contains map with bool values",
			args: args{
				v:       map[string]bool{"a": true, "b": false, "c": true},
				element: "a",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid contains map with bool values",
			args: args{
				v:       map[string]bool{"a": true, "b": false, "c": true},
				element: "d",
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Valid contains map with int values",
			args: args{
				v:       map[string]int{"a": 1, "b": 2, "c": 3},
				element: "a",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid contains map with int values",
			args: args{
				v:       map[string]int{"a": 1, "b": 2, "c": 3},
				element: "d",
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Valid contains map with int8 values",
			args: args{
				v:       map[string]int8{"a": 1, "b": 2, "c": 3},
				element: "a",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid contains map with int8 values",
			args: args{
				v:       map[string]int8{"a": 1, "b": 2, "c": 3},
				element: "d",
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Valid contains map with int16 values",
			args: args{
				v:       map[string]int16{"a": 1, "b": 2, "c": 3},
				element: "a",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid contains map with int16 values",
			args: args{
				v:       map[string]int16{"a": 1, "b": 2, "c": 3},
				element: "d",
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Valid contains map with int32 values",
			args: args{
				v:       map[string]int32{"a": 1, "b": 2, "c": 3},
				element: "a",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid contains map with int32 values",
			args: args{
				v:       map[string]int32{"a": 1, "b": 2, "c": 3},
				element: "d",
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Valid contains map with int64 values",
			args: args{
				v:       map[string]int64{"a": 1, "b": 2, "c": 3},
				element: "a",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid contains map with int64 values",
			args: args{
				v:       map[string]int64{"a": 1, "b": 2, "c": 3},
				element: "d",
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Valid contains map with uint values",
			args: args{
				v:       map[string]uint{"a": 1, "b": 2, "c": 3},
				element: "a",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid contains map with uint values",
			args: args{
				v:       map[string]uint{"a": 1, "b": 2, "c": 3},
				element: "d",
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Valid contains map with uint8 values",
			args: args{
				v:       map[string]uint8{"a": 1, "b": 2, "c": 3},
				element: "a",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid contains map with uint8 values",
			args: args{
				v:       map[string]uint8{"a": 1, "b": 2, "c": 3},
				element: "d",
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Valid contains map with uint16 values",
			args: args{
				v:       map[string]uint16{"a": 1, "b": 2, "c": 3},
				element: "a",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid contains map with uint16 values",
			args: args{
				v:       map[string]uint16{"a": 1, "b": 2, "c": 3},
				element: "d",
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Valid contains map with uint32 values",
			args: args{
				v:       map[string]uint32{"a": 1, "b": 2, "c": 3},
				element: "a",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid contains map with uint32 values",
			args: args{
				v:       map[string]uint32{"a": 1, "b": 2, "c": 3},
				element: "d",
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Valid contains map with uint64 values",
			args: args{
				v:       map[string]uint64{"a": 1, "b": 2, "c": 3},
				element: "a",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid contains map with uint64 values",
			args: args{
				v:       map[string]uint64{"a": 1, "b": 2, "c": 3},
				element: "d",
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Valid contains map with float32 values",
			args: args{
				v:       map[string]float32{"a": 1.1, "b": 2.2, "c": 3.3},
				element: "a",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid contains map with float32 values",
			args: args{
				v:       map[string]float32{"a": 1.1, "b": 2.2, "c": 3.3},
				element: "d",
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Valid contains map with float64 values",
			args: args{
				v:       map[string]float64{"a": 1.1, "b": 2.2, "c": 3.3},
				element: "a",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid contains map with float64 values",
			args: args{
				v:       map[string]float64{"a": 1.1, "b": 2.2, "c": 3.3},
				element: "d",
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Valid contains map with any values",
			args: args{
				v:       map[string]any{"a": 1.1, "b": 2.2, "c": 3.3},
				element: "a",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid contains map with any values",
			args: args{
				v:       map[string]any{"a": 1.1, "b": 2.2, "c": 3.3},
				element: "d",
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Valid contains map with float64 values",
			args: args{
				v:       map[string]float64{"a": 1.1, "b": 2.2, "c": 3.3},
				element: "a",
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Unsupported type",
			args: args{
				v:       args{},
				element: "d",
			},
			valid:   false,
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := Contains(test.args.v, test.args.element)
			if test.wantErr {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.Equal(t, test.valid, result, "Expected validation result does not match")
			}
		})
	}
}

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

func TestFrom(t *testing.T) {
	type args struct {
		v    any
		from string
		not  bool
	}
	tests := []struct {
		name    string
		args    args
		valid   bool
		wantErr bool
	}{
		{
			name: "Valid from string",
			args: args{
				v:    "apple",
				from: "apple,banana",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Valid not from string",
			args: args{
				v:    "cherry",
				from: "apple,banana",
				not:  true,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from string",
			args: args{
				v:    "cherry",
				from: "apple,banana",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition from string",
			args: args{
				v:    "cherry",
				from: "",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from bool",
			args: args{
				v:    true,
				from: "true",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from bool",
			args: args{
				v:    false,
				from: "true",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition type from bool",
			args: args{
				v:    true,
				from: "banana",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from int",
			args: args{
				v:    5,
				from: "5,6,7",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from int",
			args: args{
				v:    8,
				from: "5,6,7",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition type from int",
			args: args{
				v:    5,
				from: "banana",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from int8",
			args: args{
				v:    int8(5),
				from: "5,6,7",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from int8",
			args: args{
				v:    int8(8),
				from: "5,6,7",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition type from int8",
			args: args{
				v:    int8(5),
				from: "banana",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from int16",
			args: args{
				v:    int16(5),
				from: "5,6,7",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from int16",
			args: args{
				v:    int16(8),
				from: "5,6,7",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition type from int16",
			args: args{
				v:    int16(5),
				from: "banana",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from int32",
			args: args{
				v:    int32(5),
				from: "5,6,7",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from int32",
			args: args{
				v:    int32(8),
				from: "5,6,7",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition type from int32",
			args: args{
				v:    int32(5),
				from: "banana",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from int64",
			args: args{
				v:    int64(5),
				from: "5,6,7",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from int64",
			args: args{
				v:    int64(8),
				from: "5,6,7",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition type from int64",
			args: args{
				v:    int64(5),
				from: "banana",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from uint",
			args: args{
				v:    uint(5),
				from: "5,6,7",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from uint",
			args: args{
				v:    uint(8),
				from: "5,6,7",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition type from uint",
			args: args{
				v:    uint(5),
				from: "banana",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from uint8",
			args: args{
				v:    uint8(5),
				from: "5,6,7",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from uint8",
			args: args{
				v:    uint8(8),
				from: "5,6,7",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition type from uint8",
			args: args{
				v:    uint8(5),
				from: "banana",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from uint16",
			args: args{
				v:    uint16(5),
				from: "5,6,7",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from uint16",
			args: args{
				v:    uint16(8),
				from: "5,6,7",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition type from uint16",
			args: args{
				v:    uint16(5),
				from: "banana",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from uint32",
			args: args{
				v:    uint32(5),
				from: "5,6,7",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from uint32",
			args: args{
				v:    uint32(8),
				from: "5,6,7",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition type from uint32",
			args: args{
				v:    uint32(5),
				from: "banana",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from uint64",
			args: args{
				v:    uint64(5),
				from: "5,6,7",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from uint64",
			args: args{
				v:    uint64(8),
				from: "5,6,7",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition type from uint64",
			args: args{
				v:    uint64(5),
				from: "banana",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from float32",
			args: args{
				v:    float32(5.5),
				from: "5.5,6.6,7.7",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from float32",
			args: args{
				v:    float32(8.8),
				from: "5.5,6.6,7.7",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition type from float32",
			args: args{
				v:    float32(5.5),
				from: "banana",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from float64",
			args: args{
				v:    float64(5.5),
				from: "5.5,6.6,7.7",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from float64",
			args: args{
				v:    float64(8.8),
				from: "5.5,6.6,7.7",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition type from float64",
			args: args{
				v:    float64(5.5),
				from: "banana",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from []string",
			args: args{
				v:    []string{"apple", "banana"},
				from: "apple,banana,cherry",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from []string",
			args: args{
				v:    []string{"apple", "orange"},
				from: "apple,banana,cherry",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition from []string",
			args: args{
				v:    []string{"apple", "banana"},
				from: "",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from []bool",
			args: args{
				v:    []bool{true, false},
				from: "true,false",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from []bool",
			args: args{
				v:    []bool{true, false},
				from: "true",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition type from []bool",
			args: args{
				v:    []bool{true, false},
				from: "banana",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from []int",
			args: args{
				v:    []int{1, 2},
				from: "1,2,3",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from []int",
			args: args{
				v:    []int{1, 4},
				from: "1,2,3",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition type from []int",
			args: args{
				v:    []int{1, 2},
				from: "banana",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from []int8",
			args: args{
				v:    []int8{1, 2},
				from: "1,2,3",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from []int8",
			args: args{
				v:    []int8{1, 4},
				from: "1,2,3",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition type from []int8",
			args: args{
				v:    []int8{1, 2},
				from: "banana",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from []int16",
			args: args{
				v:    []int16{1, 2},
				from: "1,2,3",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from []int16",
			args: args{
				v:    []int16{1, 4},
				from: "1,2,3",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition type from []int16",
			args: args{
				v:    []int16{1, 2},
				from: "banana",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from []int32",
			args: args{
				v:    []int32{1, 2},
				from: "1,2,3",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from []int32",
			args: args{
				v:    []int32{1, 4},
				from: "1,2,3",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition type from []int32",
			args: args{
				v:    []int32{1, 2},
				from: "banana",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from []int64",
			args: args{
				v:    []int64{1, 2},
				from: "1,2,3",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from []int64",
			args: args{
				v:    []int64{1, 4},
				from: "1,2,3",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition type from []int64",
			args: args{
				v:    []int64{1, 2},
				from: "banana",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from []uint",
			args: args{
				v:    []uint{1, 2},
				from: "1,2,3",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from []uint",
			args: args{
				v:    []uint{1, 4},
				from: "1,2,3",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition type from []uint",
			args: args{
				v:    []uint{1, 2},
				from: "banana",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from []uint8",
			args: args{
				v:    []uint8{1, 2},
				from: "1,2,3",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from []uint8",
			args: args{
				v:    []uint8{1, 4},
				from: "1,2,3",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition type from []uint8",
			args: args{
				v:    []uint8{1, 2},
				from: "banana",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from []uint16",
			args: args{
				v:    []uint16{1, 2},
				from: "1,2,3",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from []uint16",
			args: args{
				v:    []uint16{1, 4},
				from: "1,2,3",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition type from []uint16",
			args: args{
				v:    []uint16{1, 2},
				from: "banana",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from []uint32",
			args: args{
				v:    []uint32{1, 2},
				from: "1,2,3",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from []uint32",
			args: args{
				v:    []uint32{1, 4},
				from: "1,2,3",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition type from []uint32",
			args: args{
				v:    []uint32{1, 2},
				from: "banana",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from []uint64",
			args: args{
				v:    []uint64{1, 2},
				from: "1,2,3",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from []uint64",
			args: args{
				v:    []uint64{1, 4},
				from: "1,2,3",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition type from []uint64",
			args: args{
				v:    []uint64{1, 2},
				from: "banana",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from []float32",
			args: args{
				v:    []float32{1.1, 2.2},
				from: "1.1,2.2,3.3",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from []float32",
			args: args{
				v:    []float32{1.1, 4.4},
				from: "1.1,2.2,3.3",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition type from []float32",
			args: args{
				v:    []float32{1.1, 2.2},
				from: "banana",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from []float64",
			args: args{
				v:    []float64{1.1, 2.2},
				from: "1.1,2.2,3.3",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from []float64",
			args: args{
				v:    []float64{1.1, 4.4},
				from: "1.1,2.2,3.3",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition type from []float64",
			args: args{
				v:    []float64{1.1, 2.2},
				from: "banana",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from map[string]string",
			args: args{
				v:    map[string]string{"a": "1", "b": "2"},
				from: "a,b,c",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from map[string]string",
			args: args{
				v:    map[string]string{"a": "1", "d": "4"},
				from: "a,b,c",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition from map[string]string",
			args: args{
				v:    map[string]string{"a": "1", "b": "2"},
				from: "",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from map[string]bool",
			args: args{
				v:    map[string]bool{"a": true, "b": false},
				from: "a,b,c",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from map[string]bool",
			args: args{
				v:    map[string]bool{"a": true, "d": false},
				from: "a,b,c",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition from map[string]bool",
			args: args{
				v:    map[string]bool{"a": true, "b": false},
				from: "",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from map[string]int",
			args: args{
				v:    map[string]int{"a": 1, "b": 2},
				from: "a,b,c",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from map[string]int",
			args: args{
				v:    map[string]int{"a": 1, "d": 4},
				from: "a,b,c",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition from map[string]int",
			args: args{
				v:    map[string]int{"a": 1, "b": 2},
				from: "",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from map[string]int8",
			args: args{
				v:    map[string]int8{"a": 1, "b": 2},
				from: "a,b,c",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from map[string]int8",
			args: args{
				v:    map[string]int8{"a": 1, "d": 4},
				from: "a,b,c",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition from map[string]int8",
			args: args{
				v:    map[string]int8{"a": 1, "b": 2},
				from: "",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from map[string]int16",
			args: args{
				v:    map[string]int16{"a": 1, "b": 2},
				from: "a,b,c",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from map[string]int16",
			args: args{
				v:    map[string]int16{"a": 1, "d": 4},
				from: "a,b,c",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition from map[string]int16",
			args: args{
				v:    map[string]int16{"a": 1, "b": 2},
				from: "",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from map[string]int32",
			args: args{
				v:    map[string]int32{"a": 1, "b": 2},
				from: "a,b,c",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from map[string]int32",
			args: args{
				v:    map[string]int32{"a": 1, "d": 4},
				from: "a,b,c",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition from map[string]int32",
			args: args{
				v:    map[string]int32{"a": 1, "b": 2},
				from: "",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from map[string]int64",
			args: args{
				v:    map[string]int64{"a": 1, "b": 2},
				from: "a,b,c",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from map[string]int64",
			args: args{
				v:    map[string]int64{"a": 1, "d": 4},
				from: "a,b,c",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition from map[string]int64",
			args: args{
				v:    map[string]int64{"a": 1, "b": 2},
				from: "",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from map[string]uint",
			args: args{
				v:    map[string]uint{"a": 1, "b": 2},
				from: "a,b,c",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from map[string]uint",
			args: args{
				v:    map[string]uint{"a": 1, "d": 4},
				from: "a,b,c",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition from map[string]uint",
			args: args{
				v:    map[string]uint{"a": 1, "b": 2},
				from: "",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from map[string]uint8",
			args: args{
				v:    map[string]uint8{"a": 1, "b": 2},
				from: "a,b,c",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from map[string]uint8",
			args: args{
				v:    map[string]uint8{"a": 1, "d": 4},
				from: "a,b,c",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition from map[string]uint8",
			args: args{
				v:    map[string]uint8{"a": 1, "b": 2},
				from: "",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from map[string]uint16",
			args: args{
				v:    map[string]uint16{"a": 1, "b": 2},
				from: "a,b,c",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from map[string]uint16",
			args: args{
				v:    map[string]uint16{"a": 1, "d": 4},
				from: "a,b,c",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition from map[string]uint16",
			args: args{
				v:    map[string]uint16{"a": 1, "b": 2},
				from: "",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from map[string]uint32",
			args: args{
				v:    map[string]uint32{"a": 1, "b": 2},
				from: "a,b,c",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from map[string]uint32",
			args: args{
				v:    map[string]uint32{"a": 1, "d": 4},
				from: "a,b,c",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition from map[string]uint32",
			args: args{
				v:    map[string]uint32{"a": 1, "b": 2},
				from: "",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from map[string]uint64",
			args: args{
				v:    map[string]uint64{"a": 1, "b": 2},
				from: "a,b,c",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from map[string]uint64",
			args: args{
				v:    map[string]uint64{"a": 1, "d": 4},
				from: "a,b,c",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition from map[string]uint64",
			args: args{
				v:    map[string]uint64{"a": 1, "b": 2},
				from: "",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from map[string]float32",
			args: args{
				v:    map[string]float32{"a": 1.1, "b": 2.2},
				from: "a,b,c",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from map[string]float32",
			args: args{
				v:    map[string]float32{"a": 1.1, "d": 4.4},
				from: "a,b,c",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition from map[string]float32",
			args: args{
				v:    map[string]float32{"a": 1.1, "b": 2.2},
				from: "",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from map[string]float64",
			args: args{
				v:    map[string]float64{"a": 1.1, "b": 2.2},
				from: "a,b,c",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from map[string]float64",
			args: args{
				v:    map[string]float64{"a": 1.1, "d": 4.4},
				from: "a,b,c",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition from map[string]float64",
			args: args{
				v:    map[string]float64{"a": 1.1, "b": 2.2},
				from: "",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Valid from map[string]any",
			args: args{
				v:    map[string]any{"a": 1.1, "b": "test"},
				from: "a,b,c",
				not:  false,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "Invalid from map[string]any",
			args: args{
				v:    map[string]any{"a": 1.1, "d": "test"},
				from: "a,b,c",
				not:  false,
			},
			valid:   false,
			wantErr: false,
		},
		{
			name: "Invalid condition from map[string]any",
			args: args{
				v:    map[string]any{"a": 1.1, "b": "test"},
				from: "",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
		{
			name: "Unsupported type",
			args: args{
				v:    struct{}{},
				from: "5,6,7",
				not:  false,
			},
			valid:   false,
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := From(test.args.v, test.args.from, test.args.not)
			if test.wantErr {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.Equal(t, test.valid, result, "Expected validation result does not match")
			}
		})
	}
}

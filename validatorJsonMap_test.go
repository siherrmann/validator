package validator

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"

	"github.com/siherrmann/validator/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestStruct struct {
	String             string             `json:"string"`
	Bool               bool               `json:"bool"`
	Int                int                `json:"int"`
	Int8               int8               `json:"int8"`
	Int16              int16              `json:"int16"`
	Int32              int32              `json:"int32"`
	Int64              int64              `json:"int64"`
	Uint               uint               `json:"uint"`
	Uint8              uint8              `json:"uint8"`
	Uint16             uint16             `json:"uint16"`
	Uint32             uint32             `json:"uint32"`
	Uint64             uint64             `json:"uint64"`
	Float32            float32            `json:"float32"`
	Float64            float64            `json:"float64"`
	BoolStringOn       bool               `json:"boolStringOn"`
	BoolStringOff      bool               `json:"boolStringOff"`
	IntString          int                `json:"intString"`
	Int8String         int8               `json:"int8String"`
	Int16String        int16              `json:"int16String"`
	Int32String        int32              `json:"int32String"`
	Int64String        int64              `json:"int64String"`
	UintString         uint               `json:"uintString"`
	Uint8String        uint8              `json:"uint8String"`
	Uint16String       uint16             `json:"uint16String"`
	Uint32String       uint32             `json:"uint32String"`
	Uint64String       uint64             `json:"uint64String"`
	Float32String      float32            `json:"float32String"`
	Float64String      float64            `json:"float64String"`
	BoolJson           bool               `json:"boolJson"`
	IntJson            int                `json:"intJson"`
	Int8Json           int8               `json:"int8Json"`
	Int16Json          int16              `json:"int16Json"`
	Int32Json          int32              `json:"int32Json"`
	Int64Json          int64              `json:"int64Json"`
	UintJson           uint               `json:"uintJson"`
	Uint8Json          uint8              `json:"uint8Json"`
	Uint16Json         uint16             `json:"uint16Json"`
	Uint32Json         uint32             `json:"uint32Json"`
	Uint64Json         uint64             `json:"uint64Json"`
	Float32Json        float32            `json:"float32Json"`
	Float64Json        float64            `json:"float64Json"`
	StringArray        []string           `json:"stringArray"`
	BoolArray          []bool             `json:"boolArray"`
	IntArray           []int              `json:"intArray"`
	Int8Array          []int8             `json:"int8Array"`
	Int16Array         []int16            `json:"int16Array"`
	Int32Array         []int32            `json:"int32Array"`
	Int64Array         []int64            `json:"int64Array"`
	UintArray          []uint             `json:"uintArray"`
	Uint8Array         []uint8            `json:"uint8Array"`
	Uint16Array        []uint16           `json:"uint16Array"`
	Uint32Array        []uint32           `json:"uint32Array"`
	Uint64Array        []uint64           `json:"uint64Array"`
	Float32Array       []float32          `json:"float32Array"`
	Float64Array       []float64          `json:"float64Array"`
	StringArrayJson    []string           `json:"stringArrayJson"`
	BoolArrayJson      []bool             `json:"boolArrayJson"`
	IntArrayJson       []int              `json:"intArrayJson"`
	Int8ArrayJson      []int8             `json:"int8ArrayJson"`
	Int16ArrayJson     []int16            `json:"int16ArrayJson"`
	Int32ArrayJson     []int32            `json:"int32ArrayJson"`
	Int64ArrayJson     []int64            `json:"int64ArrayJson"`
	UintArrayJson      []uint             `json:"uintArrayJson"`
	Uint8ArrayJson     []uint8            `json:"uint8ArrayJson"`
	Uint16ArrayJson    []uint16           `json:"uint16ArrayJson"`
	Uint32ArrayJson    []uint32           `json:"uint32ArrayJson"`
	Uint64ArrayJson    []uint64           `json:"uint64ArrayJson"`
	Float32ArrayJson   []float32          `json:"float32ArrayJson"`
	Float64ArrayJson   []float64          `json:"float64ArrayJson"`
	MapStringString    map[string]string  `json:"mapStringString"`
	MapStringBool      map[string]bool    `json:"mapStringBool"`
	MapStringInt       map[string]int     `json:"mapStringInt"`
	MapStringFloat64   map[string]float64 `json:"mapStringFloat64"`
	MapStringInterface map[string]any     `json:"mapStringInterface"`
	MapIntStringJson   map[int]string     `json:"mapIntStringJson"`
	MapInt8StringJson  map[int8]string    `json:"mapInt8StringJson"`
}

func TestMapJsonToStruct(t *testing.T) {
	testStruct := &TestStruct{}
	json := model.JsonMap{
		"string":                 "test",
		"bool":                   true,
		"int":                    int(1),
		"int8":                   int8(2),
		"int16":                  int16(3),
		"int32":                  int32(4),
		"int64":                  int64(5),
		"uint":                   uint(6),
		"uint8":                  uint8(7),
		"uint16":                 uint16(8),
		"uint32":                 uint32(9),
		"uint64":                 uint64(10),
		"float32":                float32(11.11),
		"float64":                float64(12.12),
		"boolStringOn":           "on",
		"boolStringOff":          "off",
		"intString":              "1",
		"int8String":             "2",
		"int16String":            "3",
		"int32String":            "4",
		"int64String":            "5",
		"uintString":             "6",
		"uint8String":            "7",
		"uint16String":           "8",
		"uint32String":           "9",
		"uint64String":           "10",
		"float32String":          "11.11",
		"float64String":          "12.12",
		"boolJson":               "true",
		"intJson":                float64(1),
		"int8Json":               float64(2),
		"int16Json":              float64(3),
		"int32Json":              float64(4),
		"int64Json":              float64(5),
		"uintJson":               float64(6),
		"uint8Json":              float64(7),
		"uint16Json":             float64(8),
		"uint32Json":             float64(9),
		"uint64Json":             float64(10),
		"float32Json":            float64(11.11),
		"float64Json":            float64(12.12),
		"stringArray":            []string{"a", "b", "c"},
		"boolArray":              []bool{true, false},
		"intArray":               []int{1, 2, 3},
		"int8Array":              []int8{4, 5, 6},
		"int16Array":             []int16{7, 8, 9},
		"int32Array":             []int32{10, 11, 12},
		"int64Array":             []int64{13, 14, 15},
		"uintArray":              []uint{6, 7, 8},
		"uint8Array":             []uint8{19, 20, 21},
		"uint16Array":            []uint16{16, 17, 18},
		"uint32Array":            []uint32{22, 23, 24},
		"uint64Array":            []uint64{25, 26, 27},
		"float32Array":           []float32{31.31, 32.32},
		"float64Array":           []float64{33.33, 34.34},
		"stringArrayJson":        []any{"a", "b", "c"},
		"boolArrayJson":          []any{"true", "false"},
		"intArrayJson":           []any{1.0, 2.0, 3.0},
		"int8ArrayJson":          []any{4.0, 5.0, 6.0},
		"int16ArrayJson":         []any{7.0, 8.0, 9.0},
		"int32ArrayJson":         []any{10.0, 11.0, 12.0},
		"int64ArrayJson":         []any{13.0, 14.0, 15.0},
		"uintArrayJson":          []any{6.0, 7.0, 8.0},
		"uint8ArrayJson":         []any{19.0, 20.0, 21.0},
		"uint16ArrayJson":        []any{16.0, 17.0, 18.0},
		"uint32ArrayJson":        []any{22.0, 23.0, 24.0},
		"uint64ArrayJson":        []any{25.0, 26.0, 27.0},
		"float32ArrayJson":       []any{31.31, 32.32},
		"float64ArrayJson":       []any{33.33, 34.34},
		"mapStringString":        map[string]string{"key1": "value1", "key2": "value2"},
		"mapStringBool":          map[string]bool{"key1": true, "key2": false},
		"mapStringInt":           map[string]int{"key1": 1, "key2": 2},
		"mapStringInt8":          map[string]int8{"key1": 3, "key2": 4},
		"mapStringUint8":         map[string]uint8{"key1": 5, "key2": 6},
		"mapStringFloat32":       map[string]float32{"key1": 7.7, "key2": 8.8},
		"mapStringFloat64":       map[string]float64{"key1": 1.1, "key2": 2.2},
		"mapStringInterface":     map[string]any{"key1": "value1", "key2": "value2"},
		"mapStringStringJson":    map[string]any{"key1": "value1", "key2": "value2"},
		"mapStringBoolJson":      map[string]any{"key1": "true", "key2": "false"},
		"mapStringIntJson":       map[string]any{"key1": 1.0, "key2": 2.0},
		"mapStringInterfaceJson": map[string]any{"key1": "value1", "key2": "value2"},
		"mapIntStringJson":       map[int]string{1: "value1", 2: "value2"},
		"mapInt8StringJson":      map[int8]string{3: "value3", 4: "value4"},
	}

	err := MapJsonMapToStruct(json, testStruct)
	assert.NoError(t, err, "Expected no error when mapping json to struct")
	assert.Equal(t, "test", testStruct.String, "Expected string to match")
	assert.Equal(t, true, testStruct.Bool, "Expected bool to match")
	assert.Equal(t, 1, testStruct.Int, "Expected int to match")
	assert.Equal(t, int8(2), testStruct.Int8, "Expected int8 to match")
	assert.Equal(t, int16(3), testStruct.Int16, "Expected int16 to match")
	assert.Equal(t, int32(4), testStruct.Int32, "Expected int32 to match")
	assert.Equal(t, int64(5), testStruct.Int64, "Expected int64 to match")
	assert.Equal(t, uint(6), testStruct.Uint, "Expected uint to match")
	assert.Equal(t, uint8(7), testStruct.Uint8, "Expected uint8 to match")
	assert.Equal(t, uint16(8), testStruct.Uint16, "Expected uint16 to match")
	assert.Equal(t, uint32(9), testStruct.Uint32, "Expected uint32 to match")
	assert.Equal(t, uint64(10), testStruct.Uint64, "Expected uint64 to match")
	assert.Equal(t, float32(11.11), testStruct.Float32, "Expected float32 to match")
	assert.Equal(t, float64(12.12), testStruct.Float64, "Expected float64 to match")
	assert.Equal(t, true, testStruct.BoolStringOn, "Expected boolStringOn to match")
	assert.Equal(t, false, testStruct.BoolStringOff, "Expected boolStringOff to match")
	assert.Equal(t, 1, testStruct.IntString, "Expected intString to match")
	assert.Equal(t, int8(2), testStruct.Int8String, "Expected int8String to match")
	assert.Equal(t, int16(3), testStruct.Int16String, "Expected int16String to match")
	assert.Equal(t, int32(4), testStruct.Int32String, "Expected int32String to match")
	assert.Equal(t, int64(5), testStruct.Int64String, "Expected int64String to match")
	assert.Equal(t, uint(6), testStruct.UintString, "Expected uintString to match")
	assert.Equal(t, uint8(7), testStruct.Uint8String, "Expected uint8String to match")
	assert.Equal(t, uint16(8), testStruct.Uint16String, "Expected uint16String to match")
	assert.Equal(t, uint32(9), testStruct.Uint32String, "Expected uint32String to match")
	assert.Equal(t, uint64(10), testStruct.Uint64String, "Expected uint64String to match")
	assert.Equal(t, float32(11.11), testStruct.Float32String, "Expected float32String to match")
	assert.Equal(t, float64(12.12), testStruct.Float64String, "Expected float64String to match")
	assert.Equal(t, true, testStruct.BoolJson, "Expected boolJson to match")
	assert.Equal(t, 1, testStruct.IntJson, "Expected intJson to match")
	assert.Equal(t, int8(2), testStruct.Int8Json, "Expected int8Json to match")
	assert.Equal(t, int16(3), testStruct.Int16Json, "Expected int16Json to match")
	assert.Equal(t, int32(4), testStruct.Int32Json, "Expected int32Json to match")
	assert.Equal(t, int64(5), testStruct.Int64Json, "Expected int64Json to match")
	assert.Equal(t, uint(6), testStruct.UintJson, "Expected uintJson to match")
	assert.Equal(t, uint8(7), testStruct.Uint8Json, "Expected uint8Json to match")
	assert.Equal(t, uint16(8), testStruct.Uint16Json, "Expected uint16Json to match")
	assert.Equal(t, uint32(9), testStruct.Uint32Json, "Expected uint32Json to match")
	assert.Equal(t, uint64(10), testStruct.Uint64Json, "Expected uint64Json to match")
	assert.Equal(t, float32(11.11), testStruct.Float32Json, "Expected float32Json to match")
	assert.Equal(t, float64(12.12), testStruct.Float64Json, "Expected float64Json to match")
	assert.Equal(t, []string{"a", "b", "c"}, testStruct.StringArray, "Expected stringArray to match")
	assert.Equal(t, []bool{true, false}, testStruct.BoolArray, "Expected boolArray to match")
	assert.Equal(t, []int{1, 2, 3}, testStruct.IntArray, "Expected intArray to match")
	assert.Equal(t, []int8{4, 5, 6}, testStruct.Int8Array, "Expected int8Array to match")
	assert.Equal(t, []int16{7, 8, 9}, testStruct.Int16Array, "Expected int16Array to match")
	assert.Equal(t, []int32{10, 11, 12}, testStruct.Int32Array, "Expected int32Array to match")
	assert.Equal(t, []int64{13, 14, 15}, testStruct.Int64Array, "Expected int64Array to match")
	assert.Equal(t, []uint{6, 7, 8}, testStruct.UintArray, "Expected uintArray to match")
	assert.Equal(t, []uint8{19, 20, 21}, testStruct.Uint8Array, "Expected uint8Array to match")
	assert.Equal(t, []uint16{16, 17, 18}, testStruct.Uint16Array, "Expected uint16Array to match")
	assert.Equal(t, []uint32{22, 23, 24}, testStruct.Uint32Array, "Expected uint32Array to match")
	assert.Equal(t, []uint64{25, 26, 27}, testStruct.Uint64Array, "Expected uint64Array to match")
	assert.Equal(t, []float32{31.31, 32.32}, testStruct.Float32Array, "Expected float32Array to match")
	assert.Equal(t, []float64{33.33, 34.34}, testStruct.Float64Array, "Expected float64Array to match")
	assert.Equal(t, []string{"a", "b", "c"}, testStruct.StringArrayJson, "Expected stringArrayJson to match")
	assert.Equal(t, []bool{true, false}, testStruct.BoolArrayJson, "Expected boolArrayJson to match")
	assert.Equal(t, []int{1, 2, 3}, testStruct.IntArrayJson, "Expected intArrayJson to match")
	assert.Equal(t, []int8{4, 5, 6}, testStruct.Int8ArrayJson, "Expected int8ArrayJson to match")
	assert.Equal(t, []int16{7, 8, 9}, testStruct.Int16ArrayJson, "Expected int16ArrayJson to match")
	assert.Equal(t, []int32{10, 11, 12}, testStruct.Int32ArrayJson, "Expected int32ArrayJson to match")
	assert.Equal(t, []int64{13, 14, 15}, testStruct.Int64ArrayJson, "Expected int64ArrayJson to match")
	assert.Equal(t, []uint{6, 7, 8}, testStruct.UintArrayJson, "Expected uintArrayJson to match")
	assert.Equal(t, []uint8{19, 20, 21}, testStruct.Uint8ArrayJson, "Expected uint8ArrayJson to match")
	assert.Equal(t, []uint16{16, 17, 18}, testStruct.Uint16ArrayJson, "Expected uint16ArrayJson to match")
	assert.Equal(t, []uint32{22, 23, 24}, testStruct.Uint32ArrayJson, "Expected uint32ArrayJson to match")
	assert.Equal(t, []uint64{25, 26, 27}, testStruct.Uint64ArrayJson, "Expected uint64ArrayJson to match")
	assert.Equal(t, []float32{31.31, 32.32}, testStruct.Float32ArrayJson, "Expected float32ArrayJson to match")
	assert.Equal(t, []float64{33.33, 34.34}, testStruct.Float64ArrayJson, "Expected float64ArrayJson to match")
	assert.Equal(t, map[string]string{"key1": "value1", "key2": "value2"}, testStruct.MapStringString, "Expected mapStringString to match")
	assert.Equal(t, map[string]bool{"key1": true, "key2": false}, testStruct.MapStringBool, "Expected mapStringBool to match")
	assert.Equal(t, map[string]int{"key1": 1, "key2": 2}, testStruct.MapStringInt, "Expected mapStringInt to match")
	assert.Equal(t, map[string]float64{"key1": 1.1, "key2": 2.2}, testStruct.MapStringFloat64, "Expected mapStringFloat64 to match")
	assert.Equal(t, map[string]any{"key1": "value1", "key2": "value2"}, testStruct.MapStringInterface, "Expected mapStringInterface to match")
	assert.Equal(t, map[int]string{1: "value1", 2: "value2"}, testStruct.MapIntStringJson, "Expected mapIntStringJson to match")
	assert.Equal(t, map[int8]string{3: "value3", 4: "value4"}, testStruct.MapInt8StringJson, "Expected mapInt8StringJson to match")
}

func TestUnmarshalRequestToJsonMap(t *testing.T) {
	t.Run("Valid JSON request", func(t *testing.T) {
		data := []byte(`{"name":"apple","age":2}`)
		req, err := http.NewRequest("POST", "/", bytes.NewBuffer(data))
		require.NoError(t, err, "Expected no error creating request")
		req.Header.Set("Content-Type", "application/json")

		mapOut, err := UnmarshalRequestToJsonMap(req)
		assert.NoError(t, err, "Expected no error unmarshaling request to JsonMap")
		assert.Equal(t, "apple", mapOut["name"], "Expected name to be 'apple'")
		assert.Equal(t, float64(2), mapOut["age"], "Expected age to be 2")
	})

	t.Run("Invalid request", func(t *testing.T) {
		_, err := UnmarshalRequestToJsonMap(nil)
		assert.Error(t, err, "Expected error unmarshaling invalid JSON request")
		assert.Contains(t, err.Error(), "request is nil", "Expected error to contain nil error")
	})
}

func TestUnmapRequestToJsonMap(t *testing.T) {
	t.Run("Valid form request", func(t *testing.T) {
		form := url.Values{}
		form.Set("name", "apple")
		form.Set("age", "2")
		req, err := http.NewRequest("POST", "/", bytes.NewBufferString(form.Encode()))
		require.NoError(t, err, "Expected no error creating request")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		mapOut, err := UnmapRequestToJsonMap(req)
		assert.NoError(t, err, "Expected no error unmapping request to JsonMap")
		assert.Equal(t, "apple", mapOut["name"], "Expected name to be 'apple'")
		assert.Equal(t, float64(2), mapOut["age"], "Expected age to be 2")
	})

	t.Run("Invalid request", func(t *testing.T) {
		_, err := UnmapRequestToJsonMap(nil)
		assert.Error(t, err, "Expected error unmapping invalid request")
		assert.Contains(t, err.Error(), "request is nil", "Expected error to contain nil error")
	})
}

func TestUnmarshalJsonToJsonMap(t *testing.T) {
	t.Run("Valid JSON", func(t *testing.T) {
		jsonData := []byte(`{"key1": "value1", "key2": 2, "key3": true}`)

		mapOut, err := UnmarshalJsonToJsonMap(jsonData)
		assert.NoError(t, err, "Expected no error when unmarshaling JSON to JsonMap")
		assert.Equal(t, "value1", mapOut["key1"], "Expected key1 to match")
		assert.Equal(t, float64(2), mapOut["key2"], "Expected key2 to match")
		assert.Equal(t, true, mapOut["key3"], "Expected key3 to match")
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		jsonData := []byte(`<html><body>Invalid JSON</body></html>`)

		_, err := UnmarshalJsonToJsonMap(jsonData)
		assert.Error(t, err, "Expected error when unmarshaling invalid JSON")
		assert.Contains(t, err.Error(), "error unmarshaling:", "Expected error to contain JSON parsing error")
	})
}

func TestUnmapUrlValuesToJsonMap(t *testing.T) {
	t.Run("Valid URL values", func(t *testing.T) {
		values := url.Values{}
		values.Set("name", "apple")
		values.Set("age", "2")
		values["array"] = []string{"fruit", "food"}
		values["map"] = []string{"{\"fruit\": \"apple\", \"food\": \"banana\"}"}

		mapOut, err := UnmapUrlValuesToJsonMap(values)
		assert.NoError(t, err, "Expected no error unmapping URL values to JsonMap")
		assert.Equal(t, "apple", mapOut["name"], "Expected name to be 'apple'")
		assert.Equal(t, float64(2), mapOut["age"], "Expected age to be 2")
		assert.ElementsMatch(t, []string{"fruit", "food"}, mapOut["array"], "Expected array to match")
		assert.Equal(t, map[string]any{"food": "banana", "fruit": "apple"}, mapOut["map"], "Expected map to match")
	})

	t.Run("Empty URL values", func(t *testing.T) {
		values := url.Values{}

		mapOut, err := UnmapUrlValuesToJsonMap(values)
		assert.NoError(t, err, "Expected no error unmapping empty URL values to JsonMap")
		assert.Empty(t, mapOut, "Expected empty JsonMap for empty URL values")
	})
}

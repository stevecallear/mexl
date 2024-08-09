package types_test

import (
	"reflect"
	"testing"

	"github.com/stevecallear/mexl/types"
)

var (
	nativeMap = map[string]any{
		"i": int64(1),
		"b": true,
		"s": "s",
		"m": map[string]any{"f": float64(1.1)},
		"a": []any{int64(2)},
	}

	objectMap = types.Map{
		"i": &types.Integer{Value: 1},
		"b": &types.Boolean{Value: true},
		"s": &types.String{Value: "s"},
		"m": types.Map{
			"f": &types.Float{Value: 1.1},
		},
		"a": types.Array{&types.Integer{Value: 2}},
	}

	nativeArray = []any{
		int64(1),
		true,
		"s",
		map[string]any{"f": float64(1.1)},
		[]any{int64(2)},
	}

	objectArray = types.Array{
		&types.Integer{Value: 1},
		&types.Boolean{Value: true},
		&types.String{Value: "s"},
		types.Map{
			"f": &types.Float{Value: 1.1},
		},
		types.Array{&types.Integer{Value: 2}},
	}
)

func TestToMap(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]any
		exp   types.Map
		err   bool
	}{
		{
			name:  "valid",
			input: nativeMap,
			exp:   objectMap,
		},
		{
			name:  "invalid",
			input: map[string]any{"x": struct{}{}},
			err:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act, err := types.ToMap(tt.input)
			if err != nil && !tt.err {
				t.Fatalf("got %v, expected nil", err)
			}
			if err == nil && tt.err {
				t.Fatal("got nil, expected error")
			}
			if !reflect.DeepEqual(act, tt.exp) {
				t.Errorf("got %v, expected %v", act, tt.exp)
			}
		})
	}
}

func TestToObject(t *testing.T) {
	tests := []struct {
		name  string
		input any
		exp   types.Object
		err   bool
	}{
		{
			name:  "int",
			input: int(1),
			exp:   &types.Integer{Value: 1},
		},
		{
			name:  "int8",
			input: int8(1),
			exp:   &types.Integer{Value: 1},
		},
		{
			name:  "int16",
			input: int16(1),
			exp:   &types.Integer{Value: 1},
		},
		{
			name:  "int32",
			input: int32(1),
			exp:   &types.Integer{Value: 1},
		},
		{
			name:  "int64",
			input: int64(1),
			exp:   &types.Integer{Value: 1},
		},
		{
			name:  "uint",
			input: uint(1),
			exp:   &types.Integer{Value: 1},
		},
		{
			name:  "uint8",
			input: uint8(1),
			exp:   &types.Integer{Value: 1},
		},
		{
			name:  "uint16",
			input: uint16(1),
			exp:   &types.Integer{Value: 1},
		},
		{
			name:  "uint32",
			input: uint32(1),
			exp:   &types.Integer{Value: 1},
		},
		{
			name:  "uint64",
			input: uint64(1),
			exp:   &types.Integer{Value: 1},
		},
		{
			name:  "float32",
			input: float32(1.0), // todo: this is a hack to avoid rounding errors
			exp:   &types.Float{Value: 1.0},
		},
		{
			name:  "float64",
			input: float64(1.1),
			exp:   &types.Float{Value: 1.1},
		},
		{
			name:  "bool",
			input: true,
			exp:   &types.Boolean{Value: true},
		},
		{
			name:  "string",
			input: "abc",
			exp:   &types.String{Value: "abc"},
		},
		{
			name:  "map",
			input: nativeMap,
			exp:   objectMap,
		},
		{
			name:  "array",
			input: nativeArray,
			exp:   objectArray,
		},
		{
			name:  "invalid",
			input: struct{}{},
			err:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act, err := types.ToObject(tt.input)
			if err != nil && !tt.err {
				t.Fatalf("got %v, expected nil", err)
			}
			if err == nil && tt.err {
				t.Fatal("got nil, expected error")
			}
			if !reflect.DeepEqual(act, tt.exp) {
				t.Errorf("got %v, expected %v", act, tt.exp)
			}
		})
	}
}

func TestToNative(t *testing.T) {
	tests := []struct {
		name  string
		input types.Object
		exp   any
		err   bool
	}{
		{
			name:  "null",
			input: &types.Null{},
			exp:   nil,
		},
		{
			name:  "integer",
			input: &types.Integer{Value: 1},
			exp:   int64(1),
		},
		{
			name:  "float",
			input: &types.Float{Value: 1.1},
			exp:   float64(1.1),
		},
		{
			name:  "string",
			input: &types.String{Value: "a"},
			exp:   "a",
		},
		{
			name:  "bool",
			input: &types.Boolean{Value: true},
			exp:   true,
		},
		{
			name:  "map",
			input: objectMap,
			exp:   nativeMap,
		},
		{
			name:  "array",
			input: objectArray,
			exp:   nativeArray,
		},
		{
			name:  "invalid",
			input: types.Func(nil),
			err:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act, err := types.ToNative(tt.input)
			if err != nil && !tt.err {
				t.Fatalf("got %v, expected nil", err)
			}
			if err == nil && tt.err {
				t.Fatal("got nil, expected error")
			}
			if !reflect.DeepEqual(act, tt.exp) {
				t.Errorf("got %v, expected %v", act, tt.exp)
			}
		})
	}
}

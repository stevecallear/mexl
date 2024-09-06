package types_test

import (
	"strings"
	"testing"

	"github.com/stevecallear/mexl/types"
)

func TestInteger_Type(t *testing.T) {
	exp := types.TypeInteger
	act := new(types.Integer).Type()

	if act != exp {
		t.Errorf("got %s, expected %s", act, exp)
	}
}

func TestInteger_Equal(t *testing.T) {
	sut := &types.Integer{Value: 1}

	tests := []struct {
		name string
		cmp  types.Object
		exp  bool
	}{
		{
			name: "not equal (type)",
			cmp:  &types.Null{},
			exp:  false,
		},
		{
			name: "not equal (int)",
			cmp:  &types.Integer{Value: 2},
			exp:  false,
		},
		{
			name: "not equal (float)",
			cmp:  &types.Float{Value: 1.0},
			exp:  false,
		},
		{
			name: "equal (pointer)",
			cmp:  sut,
			exp:  true,
		},
		{
			name: "equal (int)",
			cmp:  &types.Integer{Value: 1},
			exp:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act := sut.Equal(tt.cmp)
			if act != tt.exp {
				t.Errorf("got %v, expected %v", act, tt.exp)
			}
		})
	}
}

func TestInteger_Inspect(t *testing.T) {
	sut := &types.Integer{Value: 1}
	exp := "1"
	act := sut.Inspect()

	if act != exp {
		t.Errorf("got %s, expected %s", act, exp)
	}
}

func TestFloat_Type(t *testing.T) {
	exp := types.TypeFloat
	act := new(types.Float).Type()

	if act != exp {
		t.Errorf("got %s, expected %s", act, exp)
	}
}

func TestFloat_Equal(t *testing.T) {
	sut := &types.Float{Value: 1.0}

	tests := []struct {
		name string
		cmp  types.Object
		exp  bool
	}{
		{
			name: "not equal (type)",
			cmp:  &types.Null{},
			exp:  false,
		},
		{
			name: "not equal (int)",
			cmp:  &types.Integer{Value: 1},
			exp:  false,
		},
		{
			name: "not equal (float)",
			cmp:  &types.Float{Value: 2.2},
			exp:  false,
		},
		{
			name: "equal (pointer)",
			cmp:  sut,
			exp:  true,
		},
		{
			name: "equal (float)",
			cmp:  &types.Float{Value: 1.0},
			exp:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act := sut.Equal(tt.cmp)
			if act != tt.exp {
				t.Errorf("got %v, expected %v", act, tt.exp)
			}
		})
	}
}

func TestFloat_Inspect(t *testing.T) {
	sut := &types.Float{Value: 1.1}
	exp := "1.1"
	act := sut.Inspect()

	if act != exp {
		t.Errorf("got %s, expected %s", act, exp)
	}
}

func TestBoolean_Type(t *testing.T) {
	exp := types.TypeBoolean
	act := new(types.Boolean).Type()

	if act != exp {
		t.Errorf("got %s, expected %s", act, exp)
	}
}

func TestBoolean_Equal(t *testing.T) {
	sut := &types.Boolean{Value: true}

	tests := []struct {
		name string
		cmp  types.Object
		exp  bool
	}{
		{
			name: "not equal (type)",
			cmp:  &types.Null{},
			exp:  false,
		},
		{
			name: "not equal (value)",
			cmp:  &types.Boolean{Value: false},
			exp:  false,
		},
		{
			name: "equal (pointer)",
			cmp:  sut,
			exp:  true,
		},
		{
			name: "equal (value)",
			cmp:  &types.Boolean{Value: true},
			exp:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act := sut.Equal(tt.cmp)
			if act != tt.exp {
				t.Errorf("got %v, expected %v", act, tt.exp)
			}
		})
	}
}

func TestBoolean_Inspect(t *testing.T) {
	sut := &types.Boolean{Value: true}
	exp := "true"
	act := sut.Inspect()

	if act != exp {
		t.Errorf("got %s, expected %s", act, exp)
	}
}

func TestString_Type(t *testing.T) {
	exp := types.TypeString
	act := new(types.String).Type()

	if act != exp {
		t.Errorf("got %s, expected %s", act, exp)
	}
}

func TestString_Equal(t *testing.T) {
	sut := &types.String{Value: "a"}

	tests := []struct {
		name string
		cmp  types.Object
		exp  bool
	}{
		{
			name: "not equal (type)",
			cmp:  &types.Null{},
			exp:  false,
		},
		{
			name: "not equal (case)",
			cmp:  &types.String{Value: "A"},
			exp:  false,
		},
		{
			name: "not equal (value)",
			cmp:  &types.String{Value: "b"},
			exp:  false,
		},
		{
			name: "equal (pointer)",
			cmp:  sut,
			exp:  true,
		},
		{
			name: "equal (value)",
			cmp:  &types.String{Value: "a"},
			exp:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act := sut.Equal(tt.cmp)
			if act != tt.exp {
				t.Errorf("got %v, expected %v", act, tt.exp)
			}
		})
	}
}

func TestString_Inspect(t *testing.T) {
	sut := &types.String{Value: "a"}
	exp := `"a"`
	act := sut.Inspect()

	if act != exp {
		t.Errorf("got %s, expected %s", act, exp)
	}
}

func TestNull_Type(t *testing.T) {
	exp := types.TypeNull
	act := new(types.Null).Type()

	if act != exp {
		t.Errorf("got %s, expected %s", act, exp)
	}
}

func TestNull_Equal(t *testing.T) {
	sut := &types.Null{}

	tests := []struct {
		name string
		cmp  types.Object
		exp  bool
	}{
		{
			name: "not equal (type)",
			cmp:  &types.Integer{},
			exp:  false,
		},
		{
			name: "equal (pointer)",
			cmp:  sut,
			exp:  true,
		},
		{
			name: "equal (type)",
			cmp:  &types.Null{},
			exp:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act := sut.Equal(tt.cmp)
			if act != tt.exp {
				t.Errorf("got %v, expected %v", act, tt.exp)
			}
		})
	}
}

func TestNull_Inspect(t *testing.T) {
	sut := &types.Null{}
	exp := "null"
	act := sut.Inspect()

	if act != exp {
		t.Errorf("got %s, expected %s", act, exp)
	}
}

func TestArray_Type(t *testing.T) {
	exp := types.TypeArray
	act := types.Array{}.Type()

	if act != exp {
		t.Errorf("got %s, expected %s", act, exp)
	}
}

func TestArray_Equal(t *testing.T) {
	tests := []struct {
		name string
		sut  types.Array
		cmp  types.Object
		exp  bool
	}{
		{
			name: "not equal (type)",
			sut: types.Array{
				&types.String{Value: "a"},
			},
			cmp: &types.Null{},
			exp: false,
		},
		{
			name: "not equal (length)",
			sut: types.Array{
				&types.String{Value: "a"},
			},
			cmp: types.Array{},
			exp: false,
		},
		{
			name: "not equal (elements)",
			sut: types.Array{
				&types.String{Value: "a"},
			},
			cmp: types.Array{
				&types.Integer{Value: 1},
			},
			exp: false,
		},
		{
			name: "equal",
			sut: types.Array{
				&types.String{Value: "a"},
				types.Array{
					&types.Integer{Value: 1},
				},
			},
			cmp: types.Array{
				&types.String{Value: "a"},
				types.Array{
					&types.Integer{Value: 1},
				},
			},
			exp: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act := tt.sut.Equal(tt.cmp)
			if act != tt.exp {
				t.Errorf("got %v, expected %v", act, tt.exp)
			}
		})
	}
}

func TestArray_Inspect(t *testing.T) {
	sut := types.Array{
		&types.Integer{Value: 1},
		&types.String{Value: "a"},
		types.Array{
			&types.Float{Value: 2.2},
		},
	}

	exp := `[1, "a", [2.2]]`
	act := sut.Inspect()

	if act != exp {
		t.Errorf("got %s, expected %s", act, exp)
	}
}

func TestMap_Type(t *testing.T) {
	exp := types.TypeMap
	act := types.Map{}.Type()

	if act != exp {
		t.Errorf("got %s, expected %s", act, exp)
	}
}

func TestMap_Equal(t *testing.T) {
	tests := []struct {
		name string
		sut  types.Map
		cmp  types.Object
		exp  bool
	}{
		{
			name: "not equal (type)",
			sut: types.Map{
				"s": &types.String{Value: "a"},
			},
			cmp: &types.Null{},
			exp: false,
		},
		{
			name: "not equal (length)",
			sut: types.Map{
				"s": &types.String{Value: "a"},
			},
			cmp: types.Map{},
			exp: false,
		},
		{
			name: "not equal (elements)",
			sut: types.Map{
				"s": &types.String{Value: "a"},
			},
			cmp: types.Map{
				"i": &types.Integer{Value: 1},
			},
			exp: false,
		},
		{
			name: "equal",
			sut: types.Map{
				"i": &types.Integer{Value: 1},
				"s": &types.String{Value: "a"},
				"m": types.Map{
					"b": &types.Boolean{Value: true},
				},
				"a": types.Array{
					&types.Float{Value: 2.2},
				},
			},
			cmp: types.Map{
				"i": &types.Integer{Value: 1},
				"s": &types.String{Value: "a"},
				"m": types.Map{
					"b": &types.Boolean{Value: true},
				},
				"a": types.Array{
					&types.Float{Value: 2.2},
				},
			},
			exp: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act := tt.sut.Equal(tt.cmp)
			if act != tt.exp {
				t.Errorf("got %v, expected %v", act, tt.exp)
			}
		})
	}
}

func TestMap_Inspect(t *testing.T) {
	sut := types.Map{
		"i": &types.Integer{Value: 1},
		"s": &types.String{Value: "a"},
		"m": types.Map{
			"b": &types.Boolean{Value: true},
		},
		"a": types.Array{
			&types.Float{Value: 2.2},
		},
	}

	act := sut.Inspect()

	// map 'random' ordering means that a direct comparison would be flakey
	if !strings.HasPrefix(act, "{") || !strings.HasSuffix(act, "}") {
		t.Errorf("got %s, expected map string", act)
	}
}

func TestFunc_Type(t *testing.T) {
	exp := types.TypeFunc
	act := types.Func(nil).Type()

	if act != exp {
		t.Errorf("got %s, expected %s", act, exp)
	}
}

func TestFunc_Equal(t *testing.T) {
	sut := types.Func(nil)

	tests := []struct {
		name string
		cmp  types.Object
		exp  bool
	}{
		{
			name: "not equal (different type)",
			cmp:  &types.Null{},
			exp:  false,
		},
		{
			name: "not equal (same type)",
			cmp:  sut,
			exp:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act := sut.Equal(tt.cmp)
			if act != tt.exp {
				t.Errorf("got %v, expected %v", act, tt.exp)
			}
		})
	}
}

func TestFunc_Inspect(t *testing.T) {
	sut := types.Func(nil)
	exp := "func"
	act := sut.Inspect()

	if act != exp {
		t.Errorf("got %s, expected %s", act, exp)
	}
}

package types_test

import (
	"reflect"
	"testing"

	"github.com/stevecallear/mexl/types"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name string
		obj  types.Object
		typ  types.Type
		exp  types.Object
		ok   bool
	}{
		{
			name: "same type",
			obj:  &types.Integer{Value: 1},
			typ:  types.TypeInteger,
			exp:  &types.Integer{Value: 1},
			ok:   true,
		},
		{
			name: "int to float",
			obj:  &types.Integer{Value: 1},
			typ:  types.TypeFloat,
			exp:  &types.Float{Value: 1.0},
			ok:   true,
		},
		{
			name: "null to int",
			obj:  new(types.Null),
			typ:  types.TypeInteger,
			exp:  new(types.Integer),
			ok:   true,
		},
		{
			name: "null to float",
			obj:  new(types.Null),
			typ:  types.TypeFloat,
			exp:  new(types.Float),
			ok:   true,
		},
		{
			name: "null to string",
			obj:  new(types.Null),
			typ:  types.TypeString,
			exp:  new(types.String),
			ok:   true,
		},
		{
			name: "null to bool",
			obj:  new(types.Null),
			typ:  types.TypeBoolean,
			exp:  new(types.Boolean),
			ok:   true,
		},
		{
			name: "null to array",
			obj:  new(types.Null),
			typ:  types.TypeArray,
			exp:  types.Array{},
			ok:   true,
		},
		{
			name: "null to map",
			obj:  new(types.Null),
			typ:  types.TypeMap,
			exp:  types.Map{},
			ok:   true,
		},
		{
			name: "invalid",
			obj:  &types.String{Value: "1"},
			typ:  types.TypeInteger,
			exp:  &types.String{Value: "1"},
			ok:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act, ok := types.Convert(tt.obj, tt.typ)
			if ok != tt.ok {
				t.Errorf("got %v, expected %v", ok, tt.ok)
			}
			if !reflect.DeepEqual(act, tt.exp) {
				t.Errorf("got %v, expected %v", act, tt.exp)
			}
		})
	}
}

func TestCoerce(t *testing.T) {
	type pair struct {
		left  types.Object
		right types.Object
	}

	tests := []struct {
		name  string
		input pair
		exp   pair
	}{
		{
			name: "null null ignore",
			input: pair{
				left:  new(types.Null),
				right: new(types.Null),
			},
			exp: pair{
				left:  new(types.Null),
				right: new(types.Null),
			},
		},
		{
			name: "null int coerce",
			input: pair{
				left:  new(types.Null),
				right: new(types.Integer),
			},
			exp: pair{
				left:  new(types.Integer),
				right: new(types.Integer),
			},
		},
		{
			name: "float null coerce",
			input: pair{
				left:  new(types.Float),
				right: new(types.Null),
			},
			exp: pair{
				left:  new(types.Float),
				right: new(types.Float),
			},
		},
		{
			name: "int float coerce",
			input: pair{
				left:  &types.Integer{Value: 1},
				right: &types.Float{Value: 1.1},
			},
			exp: pair{
				left:  &types.Float{Value: 1.0},
				right: &types.Float{Value: 1.1},
			},
		},
		{
			name: "float int coerce",
			input: pair{
				left:  &types.Float{Value: 1.1},
				right: &types.Integer{Value: 1},
			},
			exp: pair{
				left:  &types.Float{Value: 1.1},
				right: &types.Float{Value: 1.0},
			},
		},
		{
			name: "invalid ignore",
			input: pair{
				left:  &types.String{Value: "1"},
				right: &types.Integer{Value: 1},
			},
			exp: pair{
				left:  &types.String{Value: "1"},
				right: &types.Integer{Value: 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var act pair
			act.left, act.right = types.Coerce(tt.input.left, tt.input.right)

			if !reflect.DeepEqual(act, tt.exp) {
				t.Errorf("got %v, expected %v", act, tt.exp)
			}
		})
	}
}

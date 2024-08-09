package types

import (
	"fmt"
)

func ToMap(m map[string]any) (Map, error) {
	o, err := ToObject(m)
	if err != nil {
		return nil, err
	}

	return o.(Map), nil
}

func ToObject(a any) (Object, error) {
	var err error
	switch v := a.(type) {
	case int:
		return &Integer{Value: int64(v)}, nil

	case int8:
		return &Integer{Value: int64(v)}, nil

	case int16:
		return &Integer{Value: int64(v)}, nil

	case int32:
		return &Integer{Value: int64(v)}, nil

	case int64:
		return &Integer{Value: v}, nil

	case uint:
		return &Integer{Value: int64(v)}, nil

	case uint8:
		return &Integer{Value: int64(v)}, nil

	case uint16:
		return &Integer{Value: int64(v)}, nil

	case uint32:
		return &Integer{Value: int64(v)}, nil

	case uint64:
		return &Integer{Value: int64(v)}, nil

	case float32:
		return &Float{Value: float64(v)}, nil

	case float64:
		return &Float{Value: v}, nil

	case string:
		return &String{Value: v}, nil

	case bool:
		return &Boolean{Value: v}, nil

	case Func:
		return v, nil

	case func(...Object) (Object, error):
		return Func(v), nil

	case []any:
		a := make(Array, len(v))
		for i, e := range v {
			a[i], err = ToObject(e)
			if err != nil {
				return nil, err
			}
		}
		return a, nil

	case map[string]any:
		m := make(Map, len(v))
		for k, vv := range v {
			m[k], err = ToObject(vv)
			if err != nil {
				return nil, err
			}
		}
		return m, nil

	default:
		return nil, fmt.Errorf("invalid type: %+T", a)
	}
}

func ToNative(o Object) (_ any, err error) {
	switch t := o.(type) {
	case *Null:
		return nil, nil

	case *Boolean:
		return t.Value, nil

	case *Integer:
		return t.Value, nil

	case *Float:
		return t.Value, nil

	case *String:
		return t.Value, nil

	case Array:
		v := make([]any, len(t))
		for i, e := range t {
			v[i], err = ToNative(e)
			if err != nil {
				return nil, err
			}
		}
		return v, nil

	case Map:
		m := make(map[string]any, len(t))
		for k, v := range t {
			m[k], err = ToNative(v)
			if err != nil {
				return nil, err
			}
		}
		return m, nil

	default:
		return nil, fmt.Errorf("failed to convert to native type: %s", o.Type())
	}
}

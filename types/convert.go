package types

import "fmt"

var defaults = map[Type]Object{
	TypeInteger: new(Integer),
	TypeFloat:   new(Float),
	TypeString:  new(String),
	TypeBoolean: new(Boolean),
	TypeArray:   Array{},
	TypeMap:     Map{},
}

func Convert(o Object, t Type) (Object, bool) {
	ot := o.Type()
	switch {
	case ot == t:
		return o, true

	case ot == TypeInteger && t == TypeFloat:
		return &Float{Value: float64(o.(*Integer).Value)}, true

	case ot == TypeNull:
		return getDefault(t), true

	default:
		return o, false
	}
}

func Coerce(left, right Object) (Object, Object) {
	lt, rt := left.Type(), right.Type()

	switch {
	case lt == TypeNull && rt != TypeNull:
		left = getDefault(rt)

	case lt != TypeNull && rt == TypeNull:
		right = getDefault(lt)

	case lt == TypeInteger && rt == TypeFloat:
		left, _ = Convert(left, rt)

	case lt == TypeFloat && rt == TypeInteger:
		right, _ = Convert(right, lt)
	}

	return left, right
}

func getDefault(t Type) Object {
	o, ok := defaults[t]
	if !ok {
		panic(fmt.Errorf("no default for type: %s", t))
	}
	return o
}

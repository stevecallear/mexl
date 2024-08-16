# mexl

[![build](https://github.com/stevecallear/mexl/actions/workflows/build.yml/badge.svg)](https://github.com/stevecallear/mexl/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/stevecallear/mexl/graph/badge.svg?token=3JBUN06BOD)](https://codecov.io/gh/stevecallear/mexl)
[![Go Report Card](https://goreportcard.com/badge/github.com/stevecallear/mexl)](https://goreportcard.com/report/github.com/stevecallear/mexl)

`mexl` is a basic expression language. It is heavily informed by the [Writing An Intepreter/Compiler In Go](https://interpreterbook.com/) books by Thorsten Ball and was primarily written as a learning exercise.

It is intended to offer the simplicity of [`rules`](https://github.com/nikunjy/rules) with some of the performance and flexibility of [`expr`](https://github.com/expr-lang/expr), while simplifying expressions with partial type coercion. The API and syntax are not currently stable so I would not recommend use in production. Generally I would recommend `expr` in that scenario unless there are specific syntactic or functional requirements.


## Getting Started
```
go get github.com/stevecallear/mexl@latest
```

```
const input = `lower(user.email) ew "@email.com" or "beta" in user.roles`

program, err := mexl.Compile(input)
if err != nil {
	log.Fatal(err)
}

env := map[string]any{
	"user": map[string]any{
		"email": "Test@Email.com",
		"roles": []any{"admin", "beta"},
	},
}

out, err := mexl.Run(program, env)
if err != nil {
	log.Fatal(err)
}

fmt.Println(out)
// Output: true
```

## Types
`mexl` is generally statically typed, but uses type coercion where appropriate to avoid casts and excessive null checking.

|Type   |Go Type	     |
|---    |---             |
|integer|`int64`         |
|float  |`float64`       |
|string |`string`        |
|boolean|`bool`          |
|array  |`[]any`         |
|map    |`map[string]any`|
|null   |`nil`           |

### Null Coalescing
Nulls are coalesced by default. The following expression would evaluate to null.
```
out, err := mexl.Eval("x.y.z", nil) // x is not defined, out is nil
```

### Type Coercion
To avoid casts and null checking, partial type coercion is applied at runtime.

Numeric expressions containing float and integer values will result in a float value:
```
out, err := mexl.Eval("1 + 1.1", nil) // out is 2.1
```

Division of numeric values will result in a float value if the result cannot be represented by an integer:
```
out, err := mexl.Eval("3 / 2", nil) // out is 1.5
```

If a binary expression contains a null value it will be coerced to the default value of the appropriate type:
```
out, err := mexl.Eval("null + 1", nil) // out is 1
```

Null coercion is intended to avoid constant null checks or runtime errors when the input values and there types cannot be guaranteed. For example:
```
out, err := mexl.Eval(`lower(x.y) ew "abc"`, nil) // x is not defined, out is false
```

Null checks can be performed using the `eq` or `ne` operators if required:
```
out, err := mexl.Eval("x.y eq null", nil) // out is true
```

## Operations
The following operations are supported:

|Operator|Alternative|Result               |
|---     |---        |---                  |
|`eq`    |`==`       |equal                |
|`ne`    |`!=`     	 |not equal            |
|`lt`    |`<`        |less than            |
|`gt`    |`>`        |greater than         |
|`le`    |`<=`       |less than or equal   |
|`ge`    |`>=`       |greater than or equal|
|`sw`    |           |starts with          |
|`ew`    |           |ends with            |
|`in`    |           |in an array or string|
|`and`   |`&&`       |and                  |
|`or`    |`\|\|`     |or                   |
|`not`   |`!`        |not                  |

## Functions
The following built in functions are available:

|Function|Go Equivalent    |Result                                    |
|---	 |---		       |---                                       |
|`len`   |`len`            |the length of the string, array or map    |		
|`lower` |`strings.ToLower`|the lowercase representation of the string|
|`upper` |`strings.ToUpper`|the uppercase representation of the string|

### Custom Functions
Custom functions can be defined as part of the environment. `mexl` does not use reflection internally so the function must conform to the `types.Func` definition.

```
env := map[string]any {
	"word": "abc",
	"reverse": func(args ...types.Object) (types.Object, error) {
		if len(args) != 1 {
			return nil, errors.New("invalid args length")
		}
		if args[0].Type() != types.TypeString {
			return nil, errors.New("invalid arg type)
		}

		runes := []rune(args[0].(*types.String).Value)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		
		return &types.String{Value: string(runes)}, nil
	},
}

const input = `reverse(word)`

out, err := mexl.Eval(input, env)
if err != nil {
	log.Fatal(err)
}

fmt.Println(out)
// Output: cba
```


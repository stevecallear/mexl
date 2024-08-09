# mexl

[![build](https://github.com/stevecallear/mexl/actions/workflows/build.yml/badge.svg)](https://github.com/stevecallear/mexl/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/stevecallear/mexl/graph/badge.svg?token=3JBUN06BOD)](https://codecov.io/gh/stevecallear/mexl)
[![Go Report Card](https://goreportcard.com/badge/github.com/stevecallear/mexl)](https://goreportcard.com/report/github.com/stevecallear/mexl)

`mexl` is a basic expression language. It is heavily informed by the [Writing An Intepreter/Compiler In Go](https://interpreterbook.com/) books by Thorsten Ball and was primarily written as a learning exercise.

It is intended to offer the simplicity of [`rules`](https://github.com/nikunjy/rules) with some of the performance and flexibility of [`expr`](https://github.com/expr-lang/expr). With that said, if I was deploying an expression language to production it would definitely be the latter unless there was a specific syntactic or functional need.


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
`mexl` is generally statically typed, but with in-built conversion to remove the need for casts. For example, any numeric infix expression that involves a float value will result in a float output. The same is true of integer division where, for simplicity, dividing two integers will result in a float output.

|Type   |Go Type	     |
|---    |---             |
|integer|`int64`         |
|float  |`float64`       |
|string |`string`        |
|boolean|`bool`          |
|array  |`[]any`         |
|map    |`map[string]any`|
|null   |`nil`           |

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


# validate

[![Go Reference](https://pkg.go.dev/badge/github.com/protogodev/validate/vulndb.svg)][1]

Input validation made easy for Go interface methods.


## Installation

Make a custom build of [protogo](https://github.com/protogodev/protogo):

```bash
$ protogo build --plugin=github.com/protogodev/validate
```

Or build from a local fork:

```bash
$ protogo build --plugin=github.com/protogodev/validate=../my-fork
```

<details open>
  <summary> Usage </summary>

```bash
$ protogo validate -h
Usage: protogo validate <source-file> <interface-name>

Arguments:
  <source-file>       source file
  <interface-name>    interface name

Flags:
  -h, --help       Show context-sensitive help.

      --out="."    output directory
      --fmt        whether to make the generated code formatted
```
</details>


## Validation Syntax


| Operator / Validator | Validating Equivalent(s)                                                                                                                                    | Example              |
|----------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------|
| `!`                  | [Not](https://pkg.go.dev/github.com/RussellLuo/validating/v3#Not)                                                                                           | `!lt(0)`             |
| `&&`                 | [All / And](https://pkg.go.dev/github.com/RussellLuo/validating/v3#All)                                                                                     | `gt(0) && lt(100)`   | 
| `\|\|`               | [Any / Or](https://pkg.go.dev/github.com/RussellLuo/validating/v3#Or)                                                                                       | `eq(0) \|\| eq(1)`   |
| `nonzero`            | [Nonzero](https://pkg.go.dev/github.com/RussellLuo/validating/v3#Nonzero)                                                                                   | `nonzero`            | 
| `zero`               | [Zero](https://pkg.go.dev/github.com/RussellLuo/validating/v3#Zero)                                                                                         | `zero`               | 
| `len`                | [LenString](https://pkg.go.dev/github.com/RussellLuo/validating/v3#LenString) / [LenSlice](https://pkg.go.dev/github.com/RussellLuo/validating/v3#LenSlice) | `len(0, 100)`        | 
| `runecount`          | [RuneCount](https://pkg.go.dev/github.com/RussellLuo/validating/v3#RuneCount)                                                                               | `runecount(0, 100)`  |
| `eq`                 | [Eq](https://pkg.go.dev/github.com/RussellLuo/validating/v3#Eq)                                                                                             | `eq(10)`             |
| `ne`                 | [Ne](https://pkg.go.dev/github.com/RussellLuo/validating/v3#Ne)                                                                                             | `ne(-1)`             |
| `gt`                 | [Gt](https://pkg.go.dev/github.com/RussellLuo/validating/v3#Gt)                                                                                             | `gt(0)`              |
| `gte`                | [Gte](https://pkg.go.dev/github.com/RussellLuo/validating/v3#Gte)                                                                                           | `gte(0)`             |
| `lt`                 | [Lt](https://pkg.go.dev/github.com/RussellLuo/validating/v3#Lt)                                                                                             | `lt(100)`            |
| `lte`                | [Lte](https://pkg.go.dev/github.com/RussellLuo/validating/v3#Lte)                                                                                           | `lte(100)`           |
| `xrange`             | [Range](https://pkg.go.dev/github.com/RussellLuo/validating/v3#Range)                                                                                       | `xrange(0, 100)`     |
| `in`                 | [In](https://pkg.go.dev/github.com/RussellLuo/validating/v3#In)                                                                                             | `in("yes", "no")`    |
| `nin`                | [Nin](https://pkg.go.dev/github.com/RussellLuo/validating/v3#Nin)                                                                                           | `nin("yes", "no")`   |
| `match`              | [Match](https://pkg.go.dev/github.com/RussellLuo/validating/v3#Match)                                                                                       | ``match(`^\w+$`)``   |
| `_`                  | A special validator that means to use the nested `Schema()` of the struct argument                                                                          | `_`                  |


## Examples

See [examples](examples).


## Documentation

Check out the [documentation][1].


## License

[MIT](LICENSE)


[1]: https://pkg.go.dev/github.com/protogodev/validate

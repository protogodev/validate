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
  -h, --help             Show context-sensitive help.

      --out="."          output directory
      --fmt              whether to make the generated code formatted
      --custom=STRING    the declaration file of custom validators
```
</details>


## Quick Start

**NOTE**: The following code is located in [helloworld](examples/helloworld).

1. Define the interface

    ```go
    type Service interface {
        SayHello(ctx context.Context, name string) (message string, err error)
    }
    ```

2. Implement the service

    ```go
    type Greeter struct{}

    func (g *Greeter) SayHello(ctx context.Context, name string) (string, error) {
        return "Hello " + name, nil
    }
    ```

3. Add the validation annotations

    ```go
    type Service interface {
        // @schema:
        //   name: len(0, 10) && match(`^\w+$`)
        SayHello(ctx context.Context, name string) (message string, err error)
    }
    ```

4. Generate the validation middleware

    ```bash
    $ cd examples/helloworld
    $ protogo validate ./service.go Service
    ```

5. Use the middleware for input validation

    ```go
    func main() {
        var svc helloworld.Service = &helloworld.Greeter{}
        svc = helloworld.ValidateMiddleware(nil)(svc)

        message, err = svc.SayHello(context.Background(), "!Tracey")
        fmt.Printf("message: %q, err: %v\n", message, err)

        // Output:
        // message: "", err: name: INVALID(invalid format)
    }
    ```


## Validation Syntax


| Operator / Validator | Validating / Vext Equivalent(s)                                                                                                                             | Example                             |
|----------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------------|
| `!`                  | [Not](https://pkg.go.dev/github.com/RussellLuo/validating/v3#Not)                                                                                           | `!lt(0)`                            |
| `&&`                 | [All / And](https://pkg.go.dev/github.com/RussellLuo/validating/v3#All)                                                                                     | `gt(0) && lt(10)`                   |
| `\|\|`               | [Any / Or](https://pkg.go.dev/github.com/RussellLuo/validating/v3#Any)                                                                                      | `eq(0) \|\| eq(1)`                  |
| `nonzero`            | [Nonzero](https://pkg.go.dev/github.com/RussellLuo/validating/v3#Nonzero)                                                                                   | `nonzero`                           |
| `zero`               | [Zero](https://pkg.go.dev/github.com/RussellLuo/validating/v3#Zero)                                                                                         | `zero`                              |
| `len`                | [LenString](https://pkg.go.dev/github.com/RussellLuo/validating/v3#LenString) / [LenSlice](https://pkg.go.dev/github.com/RussellLuo/validating/v3#LenSlice) | `len(0, 10)`                        |
| `runecnt`            | [RuneCount](https://pkg.go.dev/github.com/RussellLuo/validating/v3#RuneCount)                                                                               | `runecnt(0, 10)`                    |
| `eq`                 | [Eq](https://pkg.go.dev/github.com/RussellLuo/validating/v3#Eq)                                                                                             | `eq(1)`                             |
| `ne`                 | [Ne](https://pkg.go.dev/github.com/RussellLuo/validating/v3#Ne)                                                                                             | `ne(2)`                             |
| `gt`                 | [Gt](https://pkg.go.dev/github.com/RussellLuo/validating/v3#Gt)                                                                                             | `gt(0)`                             |
| `gte`                | [Gte](https://pkg.go.dev/github.com/RussellLuo/validating/v3#Gte)                                                                                           | `gte(0)`                            |
| `lt`                 | [Lt](https://pkg.go.dev/github.com/RussellLuo/validating/v3#Lt)                                                                                             | `lt(10)`                            |
| `lte`                | [Lte](https://pkg.go.dev/github.com/RussellLuo/validating/v3#Lte)                                                                                           | `lte(10)`                           |
| `xrange`             | [Range](https://pkg.go.dev/github.com/RussellLuo/validating/v3#Range)                                                                                       | `xrange(0, 10)`                     |
| `in`                 | [In](https://pkg.go.dev/github.com/RussellLuo/validating/v3#In)                                                                                             | `in(0, 1)`                          |
| `nin`                | [Nin](https://pkg.go.dev/github.com/RussellLuo/validating/v3#Nin)                                                                                           | `nin("Y", "N")`                     |
| `match`              | [Match](https://pkg.go.dev/github.com/RussellLuo/validating/v3#Match)                                                                                       | ``match(`^\w+$`)``                  |
| `email`              | [Email](https://pkg.go.dev/github.com/RussellLuo/vext#Email)                                                                                                | `email`                             |
| `ip`                 | [IP](https://pkg.go.dev/github.com/RussellLuo/vext#IP)                                                                                                      | `ip`                                |
| `time`               | [Time](https://pkg.go.dev/github.com/RussellLuo/vext#Time)                                                                                                  | `time("2006-01-02T15:04:05Z07:00")` |
| `_`                  | A special validator that means to use the nested `Schema()` of the struct argument.                                                                         | `_`                                 |


## Examples

See [examples](examples).


## Documentation

Check out the [documentation][1].


## License

[MIT](LICENSE)


[1]: https://pkg.go.dev/github.com/protogodev/validate

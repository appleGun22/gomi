# gomi
Primitive macros in go, micros.

All gomi files should end with `.gomi`

gomi introduces 2 features,
simple text replacement micro `#mi`
and the `shout` keyword, a micro for errors

Declaration of **ALL** micros should take place at the top of the file, between `package` and `import` keywords

```go
package main

#mi PI 3.14159
#mi obj_type Object.content.message.GetType()

import (
  "fmt"
  "os"
)
```
### shout

`shout` is a micro for errors, example usage:

```go
shout err
// or
shout e := obj.Err()
```

will get converted to

```go
if err != nil {
  panic(err)
}
// or
if e := obj.Err(); e != nil {
  panic(e)
}
```

The default shout error handler is set to `panic(V)`, but you can change it by using a micro
```go
#shout log.Fatal(V)
```

## How to use

#### setup
1) install `parser.go`
2) build it by running `go build -o gomi.exe parser.go`
3) add the file path to the gomi executable into your `PATH`

#### usage
cd into `.gomi` file location and run `gomi gen sample.gomi` to generate `.go` file,

or, you can use it just like the go compiler , just change `go run ....` to `gomi run ....`

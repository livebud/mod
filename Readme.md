# Mod

[![Go Reference](https://pkg.go.dev/badge/github.com/livebud/mod.svg)](https://pkg.go.dev/github.com/livebud/mod)

A package for finding and manipulating `go.mod` files.

## Features

- Recursively traverses up the filesystem looking for `go.mod`
- Supports binary builds with `-trimpath`
- Extracted from [Bud](github.com/livebud/bud)

## Install

```sh
go get github.com/livebud/mod
```

## Example

```go
func main() {
  module, err := mod.Find(".")
  if err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
  fmt.Println(module.Dir())
  fmt.Println(module.Import())
}
```

```sh
$ go run .
$GOPATH/src/github.com/livebud/mod
github.com/livebud/mod
```

### Binary with `-trimpath`

```sh
$ go build -trimpath -ldflags " -X github.com/livebud/mod.path=$(go list -m) -X github.com/livebud/mod.dir=$(go list -m -f {{.Dir}})" .
$ ./main
$GOPATH/src/github.com/livebud/mod
github.com/livebud/mod
```

## Contributors

- Matt Mueller ([@mattmueller](https://twitter.com/mattmueller))

## License

MIT

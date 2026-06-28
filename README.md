# go-json

A dynamic, JSON-compatible value model for Go. `value.Value` is an alias for `any` constrained to the JSON value space (`nil | bool | int64 | float64 | string | []Value | map[string]Value`), with typed, constant-error accessors (`AsObject`, `AsList`, `AsString`, `AsInt`, `AsFloat`, `AsBool`), a `KindOf` discriminator, and the coercion rules a small expression language needs (`Truthy`, `Equal`, `Compare`, `Add`). A decoded JSON document is already a `Value`, so it interoperates directly with `encoding/json`. Depends only on the standard library.

## Install

```sh
go get github.com/gomatic/go-json
```

## Usage

```go
package main

import (
	"encoding/json"
	"fmt"

	value "github.com/gomatic/go-json"
)

func main() {
	var v value.Value
	_ = json.Unmarshal([]byte(`{"name":"ada","age":36}`), &v)

	obj, err := value.AsObject(v)
	if err != nil {
		panic(err)
	}
	name, _ := value.AsString(obj["name"])
	fmt.Println(name, value.KindOf(obj["age"])) // ada KindFloat

	sum, _ := value.Add("hello, ", name) // string concatenation
	fmt.Println(sum)                       // hello, ada
}
```

## Errors

Every typed accessor returns a sentinel matchable with `errors.Is`: `ErrNotObject`, `ErrNotList`, `ErrNotString`, `ErrNotNumber`, `ErrNotBool`, and `ErrIncomparable` (from `Compare`).

## Build & test

The `Makefile`, `.golangci.yaml`, `.editorconfig`, `.gitignore`, and `.github/` are the canonical gomatic Go toolchain, owned and distributed by [`nicerobot/tools.repository`](https://github.com/nicerobot/tools.repository) — do not edit them in-tree; per-repo changes belong in a `Makefile.local`. Run the full gate (lint, staticcheck, govulncheck, 100% coverage) with `make check`.

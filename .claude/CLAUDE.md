# go-json

A dynamic, JSON-compatible value model: `value.Value` (alias for `any` over the JSON value space), typed constant-error accessors (`AsObject`/`AsList`/`AsString`/`AsInt`/`AsFloat`/`AsBool`), `KindOf`, and coercion (`Truthy`/`Equal`/`Compare`/`Add`). Extracted from `gomatic/cirql`'s `value` package.

- Package `value`, stdlib-only (no test deps beyond the standard `testing`). Every accessor returns a sentinel `Error` const matchable with `errors.Is`.
- Gate: gofumpt, vet, staticcheck, govulncheck, gocognit ≤ 7, 100% coverage. Shared config (`Makefile`, `.golangci.yaml`, `.github/`, …) is owned by `nicerobot/tools.repository` — never edit in-tree; use `Makefile.local`.

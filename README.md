# cond

[![Go Reference](https://pkg.go.dev/badge/github.com/bold-minds/cond.svg)](https://pkg.go.dev/github.com/bold-minds/cond)
[![Build](https://img.shields.io/github/actions/workflow/status/bold-minds/cond/test.yaml?branch=main&label=tests)](https://github.com/bold-minds/cond/actions/workflows/test.yaml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/bold-minds/cond)](go.mod)

**The two conditionals Go deliberately left out — a typed ternary and a presence check.**

Go's language omits the ternary operator and has no single-call way to ask "is this value blank?" — both are intentional language choices. `cond` gives you both in two functions, and only two, with zero external dependencies.

```go
// Typed ternary — no three-line if/else for a single assignment
label := cond.If[string](user.IsActive, "active", "inactive")

// Presence check — nil, zero, empty, whitespace, or empty collection
if cond.IsEmpty(req.Name) {
    return errors.New("name required")
}
```

## ✨ Why cond?

- 🎯 **Two functions, crystal clear** — `If[T]` and `IsEmpty`. Nothing else, ever.
- ⚡ **Zero reflection** — pure type switches; safe for hot paths.
- 🦥 **Lazy evaluation opt-in** — pass `func() T` as a branch and only the chosen side runs. Eager by default.
- 🪶 **Zero dependencies** — pure Go stdlib.
- 💥 **Fail loud on programming errors** — `If[T]` panics on a mistyped branch argument. Silent zero values hide bugs.

## 📦 Installation

```bash
go get github.com/bold-minds/cond
```

Requires Go 1.21 or later.

## 🎯 Quick Start

```go
package main

import (
    "fmt"

    "github.com/bold-minds/cond"
)

func main() {
    // Typed ternary — choose an int without an if/else block
    retries := cond.If[int](isProduction(), 5, 1)

    // Different type, same call
    msg := cond.If[string](retries > 3, "high-retry mode", "low-retry mode")

    // Lazy branches — expensive() only runs when the condition selects it
    result := cond.If[string](cacheMiss, func() string { return expensive() }, "cached")

    // Presence check — handles nil, zero, empty string, whitespace, empty slice, empty map
    if cond.IsEmpty(msg) {
        fmt.Println("blank")
    }

    _ = result
}

func isProduction() bool { return true }
func expensive() string  { return "…" }
```

## 📚 API

### `If[T any](condition bool, trueVal, falseVal any) T`

Returns `trueVal` if `condition` is true, `falseVal` otherwise. The branch arguments are inspected at runtime:

- **Direct value (eager):** `cond.If[string](ok, "yes", "no")` — both values already exist.
- **`func() T` (lazy):** `cond.If[string](ok, func() string { return expensive() }, "fallback")` — only the chosen branch runs.

You can mix: one branch can be a direct value, the other a `func() T`. Selection happens per-branch.

**Panic contract:** if a branch argument is neither a `T` nor a `func() T`, `If` panics with a message identifying which branch failed (`cond: type assertion failed for trueVal` / `cond: type assertion failed for falseVal`). This is intentional — a mistyped branch is a bug, not a runtime condition to paper over.

### `IsEmpty(value any) bool`

Reports whether a value is blank in the Ruby sense. Returns `true` for:

| Kind | "Empty" when |
|---|---|
| `nil` | always |
| any numeric type | value is `0` |
| `string` | `strings.TrimSpace(s) == ""` |
| `bool` | `false` |
| `[]any`, `[]int`, `[]string` | `len == 0` |
| `map[string]any`, `map[string]int`, `map[any]any` | `len == 0` |
| `*string` | nil pointer, or points to an empty string |
| `*int`, `*float64`, `*bool`, `*any` | nil pointer |

Types not explicitly handled return `false` (treated as present). `IsEmpty` deliberately does not use reflection; if you need presence checks on custom struct types, perform them yourself before calling.

## 🧭 When to use `cond.If` vs plain `if`

Use `cond.If` when you're assigning a single value based on a boolean and the alternative is a multi-line if/else that interrupts the flow of the surrounding code:

```go
// A natural cond.If use case
retries := cond.If[int](isProduction, 5, 1)
```

Use plain `if` when you're doing anything more — side effects, multiple statements, or branches that differ in structure. `cond.If` is for expressions, not statements.

## 🧭 When to use `cond.IsEmpty`

Use `cond.IsEmpty` when you need a single check that covers "nil, zero, empty, or whitespace-only" across mixed input types. Common cases: HTTP request field validation, config defaulting, CLI flag sanity checks.

Don't use it for business-logic predicates. "Is this user's age valid?" is not a presence check — write the real predicate.

## 🔗 Related bold-minds libraries

- [`bold-minds/to`](https://github.com/bold-minds/to) — safe type conversion with fallbacks. Pair with `cond` when you want `cond.If` to produce a different output type than your inputs.
- [`bold-minds/dig`](https://github.com/bold-minds/dig) — nested data navigation. Pair with `cond` when `IsEmpty` needs to check a deeply nested value.
- [`bold-minds/txt`](https://github.com/bold-minds/txt), [`bold-minds/each`](https://github.com/bold-minds/each), [`bold-minds/list`](https://github.com/bold-minds/list) — other small, opinionated Go utility libraries in the same family.

## 🚫 Non-goals

- **No `Must*` variants.** Either a function panics on programming errors (like `If`'s branch mistype) or it doesn't. No parallel panicking twin of a safe function.
- **No reflection-based `IsEmpty`.** If your type needs it, write your own check.
- **No additional conditionals.** `cond` will always be two functions. If you want `Unless`, `IfElse`, `Case`, `Switch`, or similar, use Go's native constructs.

## 📄 License

MIT — see [LICENSE](LICENSE).

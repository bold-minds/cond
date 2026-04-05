# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] — 2026-04-05

### Changed
- **Breaking:** `IsEmpty` pointer semantics are now symmetric across all
  pointer element types. `IsEmpty(p)` is true exactly when `p` is nil,
  for every pointer type — it never dereferences. Previously `*string`
  alone also returned true for a non-nil pointer to an empty string,
  which was inconsistent with `*int`, `*bool`, and `*float64`, and
  undocumented. Callers who want to test the pointee should dereference
  first: `IsEmpty(*p)`.

### Fixed
- `If[T]` no longer panics with a raw Go nil-func call when given a
  typed-nil `func() T` branch argument. A typed nil in an interface is
  not the nil interface, so it survives the `val == nil` guard; the
  lazy-thunk branch now checks `fn == nil` before calling and treats a
  nil thunk identically to an untyped nil branch — zero value when `T`
  is nilable, descriptive panic otherwise.
- `IsEmpty` now correctly reports typed nil pointers of user-defined and
  unlisted types as empty. Previously, `IsEmpty((*MyStruct)(nil))` fell
  through the type switch and returned `false` because of Go's
  interface-nil trap — the exact bug the function was supposed to
  prevent. A reflect-based fallback in the default branch now handles
  any nilable pointer, interface, chan, or func.
- `IsEmpty` now correctly reports empty slices, maps, arrays, and
  channels of element/key types not enumerated by the fast-path type
  switch (`[]byte`, `[]float64`, `map[string]string`, `map[int]int`, and
  similar). Previously these returned `false`.
- `If[T]` now accepts the untyped `nil` literal as a branch argument
  when `T` is a nilable kind (pointer, slice, map, chan, func,
  interface), returning the zero value of `T`. Previously this panicked
  because `nil` did not satisfy either the direct-value or lazy-thunk
  type assertions. Passing `nil` to a non-nilable `T` still panics, with
  a message identifying the problem.
- `If[T]` panic messages now include both the actual and expected types,
  e.g. `cond: type assertion failed for trueVal: got int, want string
  or func() string`, instead of a bare `type assertion failed for
  trueVal`.

### Added
- `ExampleIf`, `ExampleIf_lazy`, `ExampleIf_nilBranch`, `ExampleIsEmpty`,
  `ExampleIsEmpty_collections`, `ExampleIsEmpty_nilPointer` — godoc
  examples rendered on pkg.go.dev and verified by `go test`.
- Test coverage for typed nil pointers, nil branches in `If`, `NaN`,
  arrays via the reflect fallback, nil channels, and nil funcs. Package
  coverage is now 100%.

### Changed
- `.github/workflows/test.yaml` now runs `golangci-lint` on every push
  and PR, tests against both Go 1.21 and the current stable release,
  enforces the 80% coverage floor inline, and uploads the coverage
  profile as an artifact.
- `.golangci.yml` depguard config now allows the package self-import in
  test files (previously would have blocked them, but lint was not
  wired up in CI so the error was invisible).
- `scripts/validate.sh` now clears only the test cache (`go clean
  -testcache`) instead of the machine-wide build cache, and resolves
  its mode from `$CI` when no positional argument is passed, so
  `CI=true ./scripts/validate.sh` implies CI mode.
- `README.md`: "zero reflection" bullet replaced with a more accurate
  "type-switch fast path, reflect only as a fallback" description.
  Documented the nil-branch behavior of `If`, the `func() T` ambiguity
  when `T` is itself a function type, the expanded `IsEmpty` type
  coverage, and the `math.NaN()` non-empty semantics.

## [0.1.0] — Initial release

### Added
- `If[T any](condition bool, trueVal, falseVal any) T` — typed ternary with automatic eager/lazy evaluation based on whether the branch arguments are direct values or `func() T`.
- `IsEmpty(value any) bool` — presence check for nil, zero numerics, whitespace-only strings, empty collections, and nil pointers. Type-switch based (no reflection).
- Panic contract on `If[T]` when a branch argument is neither a `T` nor a `func() T`, with a message identifying which branch failed.
- 100% coverage of exported functions, benchmarks for every path, zero external dependencies.

### Requires
- Go 1.21 or later

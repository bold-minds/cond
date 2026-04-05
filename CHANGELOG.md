# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] — Initial release

### Added
- `If[T any](condition bool, trueVal, falseVal any) T` — typed ternary with automatic eager/lazy evaluation based on whether the branch arguments are direct values or `func() T`.
- `IsEmpty(value any) bool` — presence check for nil, zero numerics, whitespace-only strings, empty collections, and nil pointers. Type-switch based (no reflection).
- Panic contract on `If[T]` when a branch argument is neither a `T` nor a `func() T`, with a message identifying which branch failed.
- 100% coverage of exported functions, benchmarks for every path, zero external dependencies.

### Requires
- Go 1.21 or later

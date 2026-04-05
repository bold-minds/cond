// Package cond provides conditional primitives that Go's language
// intentionally omits: a typed ternary and a zero-value/blank check.
//
// Two functions, one file, zero dependencies. This package pairs with
// the other bold-minds utility libraries (txt, to, dig, each, list,
// num, kv, dt) but has no dependency on any of them.
package cond

import (
	"strings"
)

// If returns trueVal when condition is true, falseVal otherwise. It
// supports both eager and lazy evaluation, chosen automatically by the
// runtime type of the branch arguments:
//
//   - Direct values (eager):
//     If[string](ok, "yes", "no")
//   - func() T (lazy — only the chosen branch is evaluated):
//     If[string](ok, func() string { return expensive() }, "fallback")
//
// When a branch argument is neither a T nor a func() T, If panics with
// a message identifying which branch failed its type assertion. The
// panic is intentional: a mistyped branch is a programming error, and
// silently returning the zero value would hide bugs.
func If[T any](condition bool, trueVal, falseVal any) T {
	if condition {
		if fn, ok := trueVal.(func() T); ok {
			return fn()
		}
		val, ok := trueVal.(T)
		if !ok {
			panic("cond: type assertion failed for trueVal")
		}
		return val
	}
	if fn, ok := falseVal.(func() T); ok {
		return fn()
	}
	val, ok := falseVal.(T)
	if !ok {
		panic("cond: type assertion failed for falseVal")
	}
	return val
}

// IsEmpty reports whether a value is "blank" in the Ruby sense: nil,
// the zero value of a numeric/boolean/string type, a whitespace-only
// string, or an empty collection. It also handles nil pointers and
// pointers to empty strings.
//
// The check is type-switch based (no reflection). Types not explicitly
// handled return false — they are treated as present by default, since
// IsEmpty cannot meaningfully introspect arbitrary user types without
// reflection.
func IsEmpty(value any) bool {
	if value == nil {
		return true
	}

	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v) == ""
	case int:
		return v == 0
	case int8:
		return v == 0
	case int16:
		return v == 0
	case int32:
		return v == 0
	case int64:
		return v == 0
	case uint:
		return v == 0
	case uint8:
		return v == 0
	case uint16:
		return v == 0
	case uint32:
		return v == 0
	case uint64:
		return v == 0
	case float32:
		return v == 0
	case float64:
		return v == 0
	case bool:
		return !v
	case []any:
		return len(v) == 0
	case []int:
		return len(v) == 0
	case []string:
		return len(v) == 0
	case map[string]any:
		return len(v) == 0
	case map[string]int:
		return len(v) == 0
	case map[any]any:
		return len(v) == 0
	case *string:
		return v == nil || *v == ""
	case *int:
		return v == nil
	case *float64:
		return v == nil
	case *bool:
		return v == nil
	default:
		if ptr, ok := value.(*any); ok {
			return ptr == nil
		}
		return false
	}
}

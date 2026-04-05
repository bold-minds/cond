// Package cond provides conditional primitives that Go's language
// intentionally omits: a typed ternary and a zero-value/blank check.
//
// Two functions, one file, zero external dependencies. This package
// pairs with the other bold-minds utility libraries but has no
// dependency on any of them.
//
// # Design notes
//
// The hot path uses a type switch with no reflection. Reflection is
// used only as a fallback for kinds the switch cannot enumerate —
// pointers to arbitrary user types and collections with arbitrary
// element/key types — and for the nil-branch safety check in If.
// Code paths that hit the reflect fallback are explicitly documented
// on the functions below.
package cond

import (
	"fmt"
	"reflect"
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
// # Nil branches
//
// A branch argument that is the untyped nil literal is treated as the
// zero value of T, provided T is a nilable kind (pointer, slice, map,
// channel, func, interface). This makes expressions like
// If[*string](ok, nil, &s) work as written. If T is not nilable (e.g.
// int, string, struct), a nil branch is a programming error and If
// panics with a descriptive message.
//
// # func() T ambiguity when T is itself a function type
//
// If T is itself a function type (e.g. func() string) the lazy-branch
// detection treats a func() T argument as a thunk and calls it,
// expecting a T back. To pass a plain function value as a direct
// branch of a function-typed T, wrap it so the outer type is not
// func() T — or use an explicit cast — otherwise the runtime
// assertion will fail and If will panic.
//
// # Panic contract
//
// When a branch argument is neither a T, a func() T, nor a valid
// untyped nil, If panics with a message naming the offending branch
// and including both the actual and expected types. The panic is
// intentional: a mistyped branch is a programming error, and silently
// returning the zero value would hide bugs.
func If[T any](condition bool, trueVal, falseVal any) T {
	if condition {
		return resolveBranch[T](trueVal, "trueVal")
	}
	return resolveBranch[T](falseVal, "falseVal")
}

// resolveBranch extracts a T from a branch argument, handling the
// eager value, lazy thunk, and untyped-nil cases. It is internal to
// the If implementation.
func resolveBranch[T any](val any, which string) T {
	var zero T

	// Untyped nil branch — valid only when T is a nilable kind.
	if val == nil {
		if isNilableKind[T]() {
			return zero
		}
		panic(fmt.Sprintf(
			"cond: %s is nil, but type %T is not nilable",
			which, zero,
		))
	}

	// Lazy thunk path.
	if fn, ok := val.(func() T); ok {
		if fn == nil {
			// A typed-nil func() T survives the val == nil check above
			// (typed nil in an interface is not the nil interface).
			// Treat it exactly like an untyped nil branch: zero value
			// when T is nilable, panic otherwise.
			if isNilableKind[T]() {
				return zero
			}
			panic(fmt.Sprintf(
				"cond: %s is a nil func() %T, and type %T is not nilable",
				which, zero, zero,
			))
		}
		return fn()
	}

	// Eager value path.
	if v, ok := val.(T); ok {
		return v
	}

	panic(fmt.Sprintf(
		"cond: type assertion failed for %s: got %T, want %T or func() %T",
		which, val, zero, zero,
	))
}

// isNilableKind reports whether T's zero value may legitimately be
// expressed as the untyped nil literal. Used by resolveBranch to
// decide whether a nil branch argument is acceptable.
func isNilableKind[T any]() bool {
	var zero T
	t := reflect.TypeOf(&zero).Elem()
	switch t.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map,
		reflect.Chan, reflect.Func, reflect.Interface:
		return true
	}
	return false
}

// IsEmpty reports whether a value is "blank" in the Ruby sense: nil,
// the zero value of a numeric/boolean type, a whitespace-only string,
// an empty collection, or a nil pointer.
//
// # Type coverage
//
// The common cases are handled by a type switch with no reflection.
// For types the switch does not enumerate — pointers to arbitrary
// user types, slices/maps/arrays/channels of arbitrary element types
// — IsEmpty falls through to a reflect-based check that handles any
// nilable pointer/interface and any Len-bearing collection. This
// fallback is intentional: without it, IsEmpty(p) on a typed nil
// pointer of an unhandled type would hit Go's interface-nil trap and
// return false, which is the opposite of the function's intent.
//
// User-defined struct types (other than the reflect-fallback nilable
// and collection kinds) return false — IsEmpty cannot meaningfully
// introspect struct fields, and doing so via reflection would invite
// surprises.
//
// # Pointer semantics
//
// For any pointer type, IsEmpty(p) is true exactly when p is nil. It
// never dereferences. In particular, IsEmpty(&"") and IsEmpty(&0)
// both return false — the pointer is present. This is deliberately
// symmetric across all pointer element types; callers who want to
// test the pointee should dereference first: IsEmpty(*p).
//
// # Floats: NaN, infinities, and signed zero
//
// IsEmpty(math.NaN()) returns false — a NaN is a present value, just
// not equal to zero. Check math.IsNaN first if you want NaN treated
// as blank. IsEmpty(math.Inf(+1)) and IsEmpty(math.Inf(-1)) also
// return false for the same reason. IsEmpty(math.Copysign(0, -1))
// returns true, because IEEE 754 has -0.0 == 0.0.
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
		// Pointer semantics: nil check only. A non-nil pointer to an
		// empty string is a present value. Callers who want to test
		// the pointee should dereference first.
		return v == nil
	case *int:
		return v == nil
	case *float64:
		return v == nil
	case *bool:
		return v == nil
	case *any:
		return v == nil
	default:
		return isEmptyReflect(value)
	}
}

// isEmptyReflect handles the fall-through case for IsEmpty: any type
// the type switch does not enumerate. It covers nilable kinds
// (pointer, interface, chan, func, map, slice) and Len-bearing
// collections. Everything else returns false.
//
// This is the only reflect use on the IsEmpty code path and is only
// reached for types the explicit type switch cannot match.
func isEmptyReflect(value any) bool {
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface,
		reflect.Chan, reflect.Func:
		return rv.IsNil()
	case reflect.Slice, reflect.Map:
		// Nil slice/map is zero-length, so IsNil is subsumed by Len.
		return rv.Len() == 0
	case reflect.Array:
		return rv.Len() == 0
	}
	return false
}

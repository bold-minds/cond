package cond_test

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/bold-minds/cond"
)

func TestIf(t *testing.T) {
	// true condition
	if got := cond.If[string](true, "yes", "no"); got != "yes" {
		t.Errorf("Expected 'yes' for true condition, got %v", got)
	}
	// false condition
	if got := cond.If[string](false, "yes", "no"); got != "no" {
		t.Errorf("Expected 'no' for false condition, got %v", got)
	}
	// different types
	if got := cond.If[int](true, 42, 0); got != 42 {
		t.Errorf("Expected 42 for true condition, got %v", got)
	}
	if got := cond.If[int](false, 42, 0); got != 0 {
		t.Errorf("Expected 0 for false condition, got %v", got)
	}
	// compose with IsEmpty
	if got := cond.If[string](cond.IsEmpty(""), "empty", "not empty"); got != "empty" {
		t.Errorf("Expected 'empty', got %v", got)
	}
	if got := cond.If[string](cond.IsEmpty("hello"), "empty", "not empty"); got != "not empty" {
		t.Errorf("Expected 'not empty', got %v", got)
	}
	// nil pointer
	var nilPtr *string
	if got := cond.If[string](nilPtr == nil, "is nil", "not nil"); got != "is nil" {
		t.Errorf("Expected 'is nil', got %v", got)
	}
	// slice branches
	emptySlice := []int{}
	nonEmptySlice := []int{1, 2, 3}
	if got := cond.If[[]int](len(emptySlice) == 0, nonEmptySlice, emptySlice); len(got) != 3 {
		t.Errorf("Expected non-empty slice, got %v", got)
	}
	if got := cond.If[[]int](len(nonEmptySlice) == 0, nonEmptySlice, emptySlice); len(got) != 0 {
		t.Errorf("Expected empty slice, got %v", got)
	}

	t.Run("lazy_evaluation", func(t *testing.T) {
		trueFn := func() string { return "lazy true" }
		if got := cond.If[string](true, trueFn, "eager false"); got != "lazy true" {
			t.Errorf("Expected 'lazy true', got %v", got)
		}
		falseFn := func() string { return "lazy false" }
		if got := cond.If[string](false, "eager true", falseFn); got != "lazy false" {
			t.Errorf("Expected 'lazy false', got %v", got)
		}

		trueFn2 := func() int { return 42 }
		falseFn2 := func() int { return 0 }
		if got := cond.If[int](true, trueFn2, falseFn2); got != 42 {
			t.Errorf("Expected 42, got %v", got)
		}
		if got := cond.If[int](false, trueFn2, falseFn2); got != 0 {
			t.Errorf("Expected 0, got %v", got)
		}
	})

	t.Run("nil_branch_nilable_T", func(t *testing.T) {
		// Pointer T: nil literal on either branch is the zero value (nil).
		s := "hello"
		if got := cond.If[*string](true, nil, &s); got != nil {
			t.Errorf("Expected nil *string for true branch, got %v", got)
		}
		if got := cond.If[*string](false, &s, nil); got != nil {
			t.Errorf("Expected nil *string for false branch, got %v", got)
		}
		if got := cond.If[*string](true, &s, nil); got == nil || *got != "hello" {
			t.Errorf("Expected &\"hello\", got %v", got)
		}

		// Slice T: nil branch.
		if got := cond.If[[]int](true, nil, []int{1}); got != nil {
			t.Errorf("Expected nil []int, got %v", got)
		}

		// Map T: nil branch.
		if got := cond.If[map[string]int](false, map[string]int{"a": 1}, nil); got != nil {
			t.Errorf("Expected nil map, got %v", got)
		}

		// Interface T: nil branch.
		if got := cond.If[any](true, nil, "x"); got != nil {
			t.Errorf("Expected nil any, got %v", got)
		}
	})

	t.Run("nil_branch_non_nilable_T_panics", func(t *testing.T) {
		defer func() {
			r := recover()
			if r == nil {
				t.Error("Expected panic when passing nil branch for non-nilable type int")
				return
			}
			msg := fmt.Sprint(r)
			if !strings.Contains(msg, "trueVal is nil") {
				t.Errorf("Expected panic identifying trueVal nil, got: %v", r)
			}
			if !strings.Contains(msg, "not nilable") {
				t.Errorf("Expected panic mentioning 'not nilable', got: %v", r)
			}
		}()
		_ = cond.If[int](true, nil, 0)
	})

	t.Run("type_assertion_failures", func(t *testing.T) {
		t.Run("trueVal_panic", func(t *testing.T) {
			defer func() {
				r := recover()
				if r == nil {
					t.Error("Expected panic for trueVal type assertion failure")
					return
				}
				msg := fmt.Sprint(r)
				if !strings.Contains(msg, "type assertion failed for trueVal") {
					t.Errorf("Expected specific panic message, got: %v", r)
				}
				// Panic should include actual and expected type info.
				if !strings.Contains(msg, "got int") {
					t.Errorf("Expected panic to name the actual type, got: %v", r)
				}
				if !strings.Contains(msg, "want string") {
					t.Errorf("Expected panic to name the expected type, got: %v", r)
				}
			}()
			var trueVal any = 42
			var falseVal any = "false"
			_ = cond.If[string](true, trueVal, falseVal)
		})

		t.Run("falseVal_panic", func(t *testing.T) {
			defer func() {
				r := recover()
				if r == nil {
					t.Error("Expected panic for falseVal type assertion failure")
					return
				}
				msg := fmt.Sprint(r)
				if !strings.Contains(msg, "type assertion failed for falseVal") {
					t.Errorf("Expected specific panic message, got: %v", r)
				}
				if !strings.Contains(msg, "got int") {
					t.Errorf("Expected panic to name the actual type, got: %v", r)
				}
				if !strings.Contains(msg, "want string") {
					t.Errorf("Expected panic to name the expected type, got: %v", r)
				}
			}()
			var trueVal any = "true"
			var falseVal any = 42
			_ = cond.If[string](false, trueVal, falseVal)
		})
	})
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		value    any
		expected bool
	}{
		// nil
		{nil, true},

		// strings
		{"", true},
		{"  ", true},
		{"\t\n", true},
		{"hello", false},
		{" hello ", false},

		// integers
		{int(0), true},
		{int(42), false},
		{int8(0), true},
		{int8(42), false},
		{int16(0), true},
		{int16(42), false},
		{int32(0), true},
		{int32(42), false},
		{int64(0), true},
		{int64(42), false},

		// unsigned integers
		{uint(0), true},
		{uint(42), false},
		{uint8(0), true},
		{uint8(42), false},
		{uint16(0), true},
		{uint16(42), false},
		{uint32(0), true},
		{uint32(42), false},
		{uint64(0), true},
		{uint64(42), false},

		// floats
		{float32(0), true},
		{float32(42.5), false},
		{float64(0), true},
		{float64(42.5), false},

		// NaN is present, not empty.
		{math.NaN(), false},

		// booleans
		{false, true},
		{true, false},

		// slices
		{[]any{}, true},
		{[]any{1}, false},
		{[]int{}, true},
		{[]int{1}, false},
		{[]string{}, true},
		{[]string{"a"}, false},

		// maps
		{map[string]any{}, true},
		{map[string]any{"a": 1}, false},
		{map[string]int{}, true},
		{map[string]int{"a": 1}, false},
		{map[any]any{}, true},
		{map[any]any{"a": 1}, false},

		// Reflect-fallback slice/map types that the switch does not list.
		{[]byte{}, true},
		{[]byte("x"), false},
		{[]float64{}, true},
		{[]float64{1.5}, false},
		{map[string]string{}, true},
		{map[string]string{"a": "b"}, false},
		{map[int]int{}, true},
		{map[int]int{1: 1}, false},
	}

	for _, test := range tests {
		if got := cond.IsEmpty(test.value); got != test.expected {
			t.Errorf("IsEmpty(%v) = %v, expected %v", test.value, got, test.expected)
		}
	}

	t.Run("pointers", func(t *testing.T) {
		var nilStrPtr *string
		if !cond.IsEmpty(nilStrPtr) {
			t.Error("Expected nil string pointer to be empty")
		}
		emptyStr := ""
		if !cond.IsEmpty(&emptyStr) {
			t.Error("Expected pointer to empty string to be empty")
		}
		nonEmptyStr := "hello"
		if cond.IsEmpty(&nonEmptyStr) {
			t.Error("Expected pointer to non-empty string to not be empty")
		}
		var nilIntPtr *int
		if !cond.IsEmpty(nilIntPtr) {
			t.Error("Expected nil int pointer to be empty")
		}
		var nilFloatPtr *float64
		if !cond.IsEmpty(nilFloatPtr) {
			t.Error("Expected nil float64 pointer to be empty")
		}
		var nilBoolPtr *bool
		if !cond.IsEmpty(nilBoolPtr) {
			t.Error("Expected nil bool pointer to be empty")
		}
		var nilAnyPtr *any
		if !cond.IsEmpty(nilAnyPtr) {
			t.Error("Expected nil *any pointer to be empty")
		}
	})

	t.Run("typed_nil_pointers_via_reflect_fallback", func(t *testing.T) {
		// These pointer types are not in the explicit type switch.
		// Without the reflect fallback, IsEmpty would return false for
		// them (Go's interface-nil trap) — the exact bug the fallback
		// exists to prevent.
		var nilStructPtr *struct{ X int }
		if !cond.IsEmpty(nilStructPtr) {
			t.Error("Expected nil *struct to be empty via reflect fallback")
		}
		var nilInt64Ptr *int64
		if !cond.IsEmpty(nilInt64Ptr) {
			t.Error("Expected nil *int64 to be empty via reflect fallback")
		}
		var nilFloat32Ptr *float32
		if !cond.IsEmpty(nilFloat32Ptr) {
			t.Error("Expected nil *float32 to be empty via reflect fallback")
		}
		var nilByteSlicePtr *[]byte
		if !cond.IsEmpty(nilByteSlicePtr) {
			t.Error("Expected nil *[]byte to be empty via reflect fallback")
		}

		// Non-nil pointer to a struct is present.
		s := struct{ X int }{X: 1}
		if cond.IsEmpty(&s) {
			t.Error("Expected non-nil *struct to not be empty")
		}
	})

	t.Run("arrays_via_reflect_fallback", func(t *testing.T) {
		// Zero-length arrays are empty; populated arrays are not. The
		// reflect fallback is the only path that handles arrays since
		// the type switch cannot enumerate every possible array type.
		var emptyArr [0]int
		if !cond.IsEmpty(emptyArr) {
			t.Error("Expected [0]int to be empty")
		}
		var populatedArr [3]int
		if cond.IsEmpty(populatedArr) {
			t.Error("Expected [3]int to not be empty")
		}
	})

	t.Run("nil_channels_and_funcs", func(t *testing.T) {
		var nilCh chan int
		if !cond.IsEmpty(nilCh) {
			t.Error("Expected nil chan to be empty")
		}
		var nilFn func()
		if !cond.IsEmpty(nilFn) {
			t.Error("Expected nil func to be empty")
		}
	})

	t.Run("unknown_types", func(t *testing.T) {
		type CustomStruct struct {
			Field string
		}
		if cond.IsEmpty(CustomStruct{Field: "test"}) {
			t.Error("Expected unknown struct type to not be empty (default case)")
		}
		// Zero-value struct is still not empty — IsEmpty does not
		// introspect user struct fields.
		if cond.IsEmpty(CustomStruct{}) {
			t.Error("Expected zero-value struct to not be empty (default case)")
		}
	})
}

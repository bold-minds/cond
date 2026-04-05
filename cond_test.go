package cond_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bold-minds/cond"
)

func Test_If(t *testing.T) {
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

	t.Run("type_assertion_failures", func(t *testing.T) {
		t.Run("trueVal_panic", func(t *testing.T) {
			defer func() {
				r := recover()
				if r == nil {
					t.Error("Expected panic for trueVal type assertion failure")
					return
				}
				if !strings.Contains(fmt.Sprint(r), "type assertion failed for trueVal") {
					t.Errorf("Expected specific panic message, got: %v", r)
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
				if !strings.Contains(fmt.Sprint(r), "type assertion failed for falseVal") {
					t.Errorf("Expected specific panic message, got: %v", r)
				}
			}()
			var trueVal any = "true"
			var falseVal any = 42
			_ = cond.If[string](false, trueVal, falseVal)
		})
	})
}

func Test_IsEmpty(t *testing.T) {
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

	t.Run("unknown_types", func(t *testing.T) {
		type CustomStruct struct {
			Field string
		}
		if cond.IsEmpty(CustomStruct{Field: "test"}) {
			t.Error("Expected unknown type to not be empty (default case)")
		}
	})
}

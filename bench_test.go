package cond_test

import (
	"testing"

	"github.com/bold-minds/cond"
)

// =============================================================================
// If benchmarks
// =============================================================================

// The eager-evaluation benchmarks deliberately use loop-carried
// values for the branch arguments. With static literals the Go
// compiler can hoist the call out of the loop and the numbers become
// meaningless (and suspiciously alloc-free). Using values derived
// from the loop variable also surfaces the boxing-to-any cost that
// real call sites with non-constant branches incur.

func Benchmark_If_EagerEvaluation_True(b *testing.B) {
	trueVals := []string{"a", "b", "c", "d"}
	falseVals := []string{"w", "x", "y", "z"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cond.If[string](true, trueVals[i&3], falseVals[i&3])
	}
}

func Benchmark_If_EagerEvaluation_False(b *testing.B) {
	trueVals := []string{"a", "b", "c", "d"}
	falseVals := []string{"w", "x", "y", "z"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cond.If[string](false, trueVals[i&3], falseVals[i&3])
	}
}

func Benchmark_If_LazyEvaluation_True(b *testing.B) {
	trueFunc := func() string { return "true_value" }
	falseFunc := func() string { return "false_value" }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cond.If[string](true, trueFunc, falseFunc)
	}
}

func Benchmark_If_LazyEvaluation_False(b *testing.B) {
	trueFunc := func() string { return "true_value" }
	falseFunc := func() string { return "false_value" }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cond.If[string](false, trueFunc, falseFunc)
	}
}

func Benchmark_If_IntegerType(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = cond.If[int](i%2 == 0, 42, 0)
	}
}

func Benchmark_If_BooleanType(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = cond.If[bool](i%2 == 0, true, false)
	}
}

func Benchmark_If_NestedCalls(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = cond.If[string](i%2 == 0,
			cond.If[string](i%4 == 0, "divisible_by_4", "divisible_by_2"),
			"odd")
	}
}

// =============================================================================
// IsEmpty benchmarks
// =============================================================================

func Benchmark_IsEmpty_String_Empty(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = cond.IsEmpty("")
	}
}

func Benchmark_IsEmpty_String_NonEmpty(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = cond.IsEmpty("hello")
	}
}

func Benchmark_IsEmpty_String_Whitespace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = cond.IsEmpty("   ")
	}
}

func Benchmark_IsEmpty_Int_Zero(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = cond.IsEmpty(0)
	}
}

func Benchmark_IsEmpty_Int_NonZero(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = cond.IsEmpty(42)
	}
}

func Benchmark_IsEmpty_Float_Zero(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = cond.IsEmpty(0.0)
	}
}

func Benchmark_IsEmpty_Float_NonZero(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = cond.IsEmpty(3.14)
	}
}

func Benchmark_IsEmpty_Bool_False(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = cond.IsEmpty(false)
	}
}

func Benchmark_IsEmpty_Bool_True(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = cond.IsEmpty(true)
	}
}

func Benchmark_IsEmpty_Slice_Empty(b *testing.B) {
	slice := []int{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cond.IsEmpty(slice)
	}
}

func Benchmark_IsEmpty_Slice_NonEmpty(b *testing.B) {
	slice := []int{1, 2, 3}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cond.IsEmpty(slice)
	}
}

func Benchmark_IsEmpty_Map_Empty(b *testing.B) {
	m := map[string]any{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cond.IsEmpty(m)
	}
}

func Benchmark_IsEmpty_Map_NonEmpty(b *testing.B) {
	m := map[string]any{"key": "value"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cond.IsEmpty(m)
	}
}

func Benchmark_IsEmpty_Nil(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = cond.IsEmpty(nil)
	}
}

func Benchmark_IsEmpty_Pointer_Nil(b *testing.B) {
	var ptr *string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cond.IsEmpty(ptr)
	}
}

func Benchmark_IsEmpty_Pointer_NonNil(b *testing.B) {
	str := "hello"
	ptr := &str
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cond.IsEmpty(ptr)
	}
}

func Benchmark_IsEmpty_MixedTypes(b *testing.B) {
	values := []any{
		"", "hello", 0, 42, false, true,
		[]int{}, []int{1, 2, 3},
		map[string]any{}, map[string]any{"key": "value"},
		nil,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, val := range values {
			_ = cond.IsEmpty(val)
		}
	}
}

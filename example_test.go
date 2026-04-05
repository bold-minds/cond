package cond_test

import (
	"fmt"

	"github.com/bold-minds/cond"
)

func ExampleIf() {
	isProd := true
	retries := cond.If[int](isProd, 5, 1)
	fmt.Println(retries)
	// Output: 5
}

func ExampleIf_lazy() {
	// Only the chosen branch is evaluated. The expensive side runs
	// zero times when the condition selects the cheap side.
	msg := cond.If[string](false,
		func() string { return "expensive" },
		"cheap",
	)
	fmt.Println(msg)
	// Output: cheap
}

func ExampleIf_nilBranch() {
	// A nil literal is a valid branch for pointer, slice, map, chan,
	// func, and interface types. Here the true branch returns a typed
	// nil *string without an extra variable.
	var s = "hello"
	got := cond.If[*string](false, &s, nil)
	fmt.Println(got == nil)
	// Output: true
}

func ExampleIsEmpty() {
	fmt.Println(cond.IsEmpty(""))
	fmt.Println(cond.IsEmpty("   "))
	fmt.Println(cond.IsEmpty("hello"))
	fmt.Println(cond.IsEmpty(0))
	fmt.Println(cond.IsEmpty(42))
	// Output:
	// true
	// true
	// false
	// true
	// false
}

func ExampleIsEmpty_collections() {
	fmt.Println(cond.IsEmpty([]int{}))
	fmt.Println(cond.IsEmpty([]int{1, 2, 3}))
	fmt.Println(cond.IsEmpty(map[string]int{}))
	fmt.Println(cond.IsEmpty(map[string]int{"a": 1}))
	// Output:
	// true
	// false
	// true
	// false
}

func ExampleIsEmpty_nilPointer() {
	// Nil pointers of any type are reported as empty, including
	// pointers to user-defined structs.
	type User struct{ Name string }
	var u *User
	fmt.Println(cond.IsEmpty(u))
	// Output: true
}

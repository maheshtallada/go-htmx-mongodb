package main

import "testing"

// Test files in Go usually end with _test.go.
// A Go test is just a regular function whose name begins with Test.
// The testing package gives us the *testing.T type, which lets us report failures.
//
// Why this file matters:
// - It checks whether our functions in mathx.go behave correctly.
// - If a test fails, the Go tool tells us which function failed and why.
// - This is the basic pattern used in real Go projects.
//
// Common test syntax explained:
// - func TestName(t *testing.T) { ... }
//   Every test function must accept a pointer to testing.T.

// - t.Errorf("...", actual, expected)
//   Reports a failure and shows the formatted message.
// - %d is an integer format placeholder.
// - %q is a quoted string format placeholder.

//
// How Go finds tests:
// - The function name must start with Test.
// - Example: TestAdd, TestDivMod, TestGreet
// - If it doesn't begin with Test, Go will ignore it.

func TestAdd(t *testing.T) {
	actual := Add(3, 5)
	expected := 8

	if actual != expected {
		t.Errorf("Add(3, 5) = %d; want %d", actual, expected)
	}
}

func TestDivMod(t *testing.T) {
	q, r := DivMod(10, 3)
	expectedQ, expectedR := 3, 1

	if q != expectedQ || r != expectedR {
		t.Errorf("DivMod(10, 3) = (%d, %d); want (%d, %d)", q, r, expectedQ, expectedR)
	}
}

func TestGreet(t *testing.T) {
	actual := Greet("Mahesh")
	expected := "Hello, Mahesh!"

	if actual != expected {
		t.Errorf("Greet(\"Mahesh\") = %q; want %q", actual, expected)
	}

	actualEmpty := Greet("")
	expectedEmpty := "Hello, world!"

	if actualEmpty != expectedEmpty {
		t.Errorf("Greet(\"\") = %q; want %q", actualEmpty, expectedEmpty)
	}
}

func TestSum(t *testing.T) {
	actual := Sum(1, 2, 3, 4)
	expected := 10

	if actual != expected {
		t.Errorf("Sum(1, 2, 3, 4) = %d; want %d", actual, expected)
	}
}

// How to run tests:
// 1. Open a terminal in the project folder.
// 2. Run:
//
//    go test ./...
//
// This command tells Go to run all tests in the current module.
// If everything passes, you will see output like:
//
//    ok   <module-name>
//
// Example in this project:
//
//    cd day1/go
//    go test ./...
//
// You can also run just a specific test:
//
//    go test ./... -run TestAdd
//
// The -run option filters tests to only those whose names match the pattern.

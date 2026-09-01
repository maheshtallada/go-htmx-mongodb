package main

import (
    "strconv"
    "strings"
)

// Engineer models a team member.
// At a high level, a struct is a way to group related data together.
// Here, one Engineer contains a name, role, years of experience, and a list of skills.
// This makes it easy to work with one employee as one object instead of separate variables.
//
// Example idea:
//   engineer := Engineer{Name: "Mahesh", Role: "Lead", years: 8, skills: []string{"Go", "MongoDB"}}
//
// Exported fields (capitalized) can be used from other packages.
// Unexported fields (lowercase) are only visible inside the same package.
type Engineer struct {
    Name   string
    Role   string
    Years  int
    Skills []string
}

// NewEngineer builds an Engineer with the given fields.
// This is a constructor pattern: instead of manually creating the struct in many places,
// we centralize the creation logic in one function.
//
// The variadic parameter "skills ...string" means the caller can pass any number of skills,
// such as "Go", "MongoDB", or "Python".
func NewEngineer(name string, role string, years int, skills ...string) Engineer {
    return Engineer{
        Name:   name,
        Role:   role,
        Years:  years,
        Skills: skills,
    }
}

// High-level difference:
// - A function is standalone and called as Function(value).
// - A method is associated with a type and called as value.Method().
// In Go, methods are just functions with a receiver, but they are written differently.

// Summary returns "<Name> (<Role>), <Years>y".
// Best practice:
// - Use strconv.Itoa when converting an int to its decimal text.
// - Do not use string(e.Years), because that converts the integer to a single rune/code point instead of the digit characters.
// - For a few concatenations, + is fine and easy to read.
// - For production-grade assembly of many strings, prefer strings.Builder or fmt.Sprintf depending on the use case.
//
// Why `string(e.Years)` failed:
// In Go, `string(n)` converts the integer to a single Unicode code point (a rune), not to the string "5".
// Example: `string(5)` is not "5"; it is a rune whose value is 5, which is not a decimal digit string.
//
// What is a rune?
// A rune is Go's alias for a single Unicode code point, typically used to represent one character like 'A', 'β', or '中'.
// It is not the same as a decimal string like "42".
func (e Engineer) Summary() string {
    return e.Name + " (" + e.Role + "), " + strconv.Itoa(e.Years) + "y"
}

// Promote increases Years by 1 and updates Role. It must MUTATE the engineer.
func (e *Engineer) Promote(newRole string) {
    e.Years++
    e.Role = newRole
}

// Describe returns a sentence for the engineer.
// Best practice for readable output:
// - Keep the method focused on a single formatted sentence.
// - Use string concatenation here because the expression is short and clear.
// - For larger, dynamic output, prefer strings.Builder or fmt.Sprintf/strconv for controlled formatting.
func (e Engineer) Describe() string {
    return "I am " + e.Name + ", a " + e.Role + "."
}

// PrintAll returns each item's description joined by newlines.
// Production-grade approach:
// - We use a strings.Builder because we are appending multiple strings in a loop.
// - This avoids repeated reallocation that happens when using repeated `+` concatenation.
// - The final result is joined with "\n" exactly once between items.
func PrintAll(items []Describer) string {
    var b strings.Builder
    for i, item := range items {
        if i > 0 {
            b.WriteString("\n")
        }
        b.WriteString(item.Describe())
    }
    return b.String()
}

// After the method, here's the flow of the loop more clearly:
// - `range items` is handled by Go itself.
// - First iteration: `i = 0`, `item = items[0]`.
// - Second iteration: `i = 1`, `item = items[1]`.
// - Third iteration: `i = 2`, `item = items[2]`.
// - This continues until the end of the slice.
// - There is no manual `i++` in our code because the runtime updates `i` automatically.
// - The `if i > 0` check adds a newline before every item after the first one.
// - So the output becomes: first item, newline, second item, newline, third item, ...

// Stretch A4: embedding promotes methods to the outer type.
// A Manager can use Engineer.Describe() without writing a new method.
type Manager struct {
    Engineer
    Reports int
}

// This line is a compile-time check: Manager satisfies Describer automatically.
var _ Describer = Manager{}

func main() {

}
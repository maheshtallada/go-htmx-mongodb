# Day 2 Notes

## A1 Micro-questions

### 1. What is the zero value of an Engineer struct?

The zero value of an Engineer is:

```go
Engineer{
    Name:  "",
    Role:  "",
    years: 0,
    skills: nil,
}
```

Explanation:
- strings default to empty string ""
- int default to 0
- slice default to nil (no elements)

So a newly declared zero-value Engineer has no meaningful data yet.

### 2. What does `%v` print vs `%q`?

`%v` prints the default value formatting.
- For strings, it prints the raw string without quotes.
- For structs, it prints their fields in a default Go format, usually as `{field1:value1 field2:value2}`.

Example:
```go
fmt.Printf("%v\n", "Go")   // Go
fmt.Printf("%q\n", "Go")   // "Go"
```

Example with a struct:
```go
engineer := Engineer{Name: "Mahesh", Role: "Lead", years: 8}
fmt.Printf("%v\n", engineer)
// Output: {Mahesh Lead 8 []}
```

Here, `%v` prints the struct using Go's default formatting:
- field values are shown without extra labels
- the exact layout is the default Go representation
- for slices, it shows the list inside brackets

`%q` is best when you want to show a string with quotes, especially when debugging values or comparing output.

Example:
```go
fmt.Printf("%q\n", engineer.Name) // "Mahesh"
```

`%v` is general-purpose, while `%q` is specifically useful for strings and quoted output.


A2:
trap check: (Trap check: If Promote had a value receiver (e Engineer), would the mutation persist? Write down why.)

No. A value receiver gets a copy of the struct. Any change inside the method changes the copy, not the original variable outside the method.

```go
func (e Engineer) Promote(newRole string) {
    e.Role = newRole
    e.Years++
}
```

This does not persist the update to the original `Engineer` value. The original remains unchanged.

To mutate the original, use a pointer receiver:

```go
func (e *Engineer) Promote(newRole string) {
    e.Role = newRole
    e.Years++
}
```

Another valid pattern is a value receiver that returns a new value:

```go
func (e Engineer) Promote(newRole string) Engineer {
    e.Role = newRole
    e.Years++
    return e
}
```

Then the caller must assign the returned value back:

```go
engineer = engineer.Promote("SDE III")
```

## A3 Micro-questions

### 1. How does Go decide a type satisfies an interface? (compare to Java implements)

Go does not use an explicit `implements` keyword.

A type satisfies an interface automatically when it has all the methods required by that interface.

```go
type Describer interface {
    Describe() string
}

type Engineer struct {
    Name string
}

func (e Engineer) Describe() string {
    return "I am " + e.Name
}
```

Here, `Engineer` satisfies `Describer` because it has `Describe() string`.

This is different from Java, where a class must declare `implements SomeInterface`.

In Go, satisfaction is structural and compile-time checked automatically.

### 2. What is an empty interface `interface{}`? What can it hold? Use with type assertions/switches.

`interface{}` means: “this can hold any type.”

It is the empty interface because it has no methods.

```go
var x interface{} = 42
x = "hello"
x = true
x = Engine{Name: "Asha"}
```

Any value can be assigned to an `interface{}`.

Type assertion:

```go
v, ok := x.(string)
if ok {
    fmt.Println("string:", v)
}
```

Type switch:

```go
switch v := x.(type) {
case string:
    fmt.Println("string:", v)
case int:
    fmt.Println("int:", v)
default:
    fmt.Println("other type:", v)
}
```

This is useful when you want to accept values of many different types and then inspect them at runtime.

### 3. Why `interface{}` holds any value; use with type assertions/switches?

Because Go interfaces are satisfied by behavior, not by declaration. The empty interface has no methods, so every concrete type satisfies it.

This makes it good for generic containers, runtime values, and APIs that can accept many types.

But it should be used carefully, because once you store a value in `interface{}`, you must recover the original concrete type using a type assertion or type switch before using it as that type.

Example:

```go
var value interface{} = 7
n := value.(int)
fmt.Println(n + 1) // 8
```

A safe version is:

```go
if n, ok := value.(int); ok {
    fmt.Println(n + 1)
}
```

Type switches are cleaner when handling multiple possible types.

### 4. When embedding: `type Manager struct { Engineer; Reports int }` — does a `Manager` value also use `Describe()`? And can an Engineer's method set change? Give the syntax.

Yes. Embedding a type promotes its methods to the outer type.

```go
type Manager struct {
    Engineer
    Reports int
}
```

A `Manager` value can call methods of `Engineer` directly:

```go
m := Manager{
    Engineer: Engineer{Name: "Asha", Role: "Manager"},
    Reports:  5,
}

fmt.Println(m.Describe())
```

This works because the embedded `Engineer` field's methods are promoted to `Manager`.

Can an Engineer's method set change?

Yes, but not by embedding alone. The method set of `Engineer` is defined by the methods declared on `Engineer` itself. It can change if you add or remove methods on the `Engineer` type.

Example:

```go
type Engineer struct {
    Name string
}

func (e Engineer) Describe() string {
    return "I am " + e.Name
}

func (e Engineer) Work() string {
    return "coding"
}
```

Now `Engineer` has the method set `{Describe, Work}`.

If you later add:

```go
func (e Engineer) Promote() {}
```

then the method set expands to include `Promote`.

Embedded fields promote those methods to the outer type, but they do not change the original type's method set unless you modify the original type definition itself.



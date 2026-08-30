package main

import "fmt"

// Printing in Go:
// - fmt.Println(...) prints values as plain text. It does NOT interpret format verbs like %T.
// - fmt.Printf(format, value...) interprets format verbs and is used when we want control over output.
//
// Common format verbs:
// %T -> prints the type of the value
// %v -> prints the value in default format
// %d -> integer
// %f -> float
// %s -> string
// %t -> bool
// %q -> string in double quotes
//
// Example:
// fmt.Printf("value=%v type=%T\n", name, name)
// This prints: value=Go type=string

func typesDemo() {
    // Variable declarations
    var language string = "Go"
    var version float64 = 1.20
    var isCompiled bool = true
    var year int = 2009

    // Short declaration (:=) and var assignment
    shortName := "C"
    var designedBy = "Google"

    // Constant declaration
    const TOPIC = "Programming Languages"

    // fmt.Printf uses the format string. %T tells Go to print the variable type.
    fmt.Printf("language = %v | type = %T\n", language, language)
    fmt.Printf("version = %v | type = %T\n", version, version)
    fmt.Printf("isCompiled = %v | type = %T\n", isCompiled, isCompiled)
    fmt.Printf("year = %v | type = %T\n", year, year)
    fmt.Printf("shortName = %v | type = %T\n", shortName, shortName)
    fmt.Printf("designedBy = %v | type = %T\n", designedBy, designedBy)
    fmt.Printf("TOPIC = %v | type = %T\n", TOPIC, TOPIC)

    var name string
    var number int
    var percentile float64
    var allGood bool

    fmt.Println("|", name, number, percentile, allGood) // prints zero values for each type
}

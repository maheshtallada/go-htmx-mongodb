package main

import "fmt"

func main() {
    fmt.Println("Day 1 ready")

//     typesDemo()

    fmt.Println(Add(3, 5))

    q, r := DivMod(10, 3)
    fmt.Println(q, r)

    fmt.Println(Greet("Mahesh"))
}

// <go run .>  is used to compile and run the Go program in the current directory.
// When you execute this command, it will build the Go code in the `main.go` file and then run the resulting executable.
// In this case, it will output "Day 1 ready" to the console.

// Q: what happens if you run `go run .` in a directory that contains multiple Go files with a `main` package?
// Ans: If you run `go run .` in a directory that contains multiple Go files with a `main` package,
// the Go compiler will compile all the Go files in that directory that belong to the `main` package and then execute the resulting program.

// Q: Should the file name be `main.go` for the `go run .` command to work correctly?
// Ans: No, the file name does not have to be `main.go` for the `go run .` command to work correctly.
// The `go run .` command will compile and run all Go files in the current directory that belong to the `main` package,
//  regardless of their individual file names. However, it is a common convention to name the main file `main.go` for clarity and organization.

// Q: what if there are multiple .go files in the directory with package main in various files ?
// Ans: If there are multiple `.go` files in the directory with the `package main`, the `go run .` command will compile all of those files together as part of the same program. The Go compiler will treat them as a single package, and it will look for the `main` function across all those files to determine where to start execution.
//
// Q: what if there are multiple `main` functions in different files within the same package?
// Ans: If there are multiple `main` functions in different files within the same package, the Go compiler will produce an error when you try to run `go run .`. The Go language specification requires that there be exactly one `main` function in the `main` package, as it serves as the entry point for the program. Having multiple `main` functions will lead to a compilation error indicating that there are duplicate definitions of the `main` function.

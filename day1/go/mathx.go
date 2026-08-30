package main

// Add returns the sum of a and b.
func Add(a, b int) int {
    return a + b
}

// DivMod returns quotient and remainder (multiple return values).
func DivMod(a, b int) (int, int) {
    return a / b, a % b
}

// Greet returns "Hello, <name>!" or "Hello, world!" if name is empty.
func Greet(name string) string {
    if name == "" {
        return "Hello, world!"
    }
    return "Hello, " + name + "!"
}

func Sum(nums ...int) int {
    total := 0

    for _, num := range nums {
        total += num
    }
    return total
}
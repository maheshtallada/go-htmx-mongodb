package main

import (
    "errors"
    "fmt"
)

// ErrInsufficientFunds is a sentinel error callers can check for.
// Sentinel means it is a predefined error value used for comparison.
var ErrInsufficientFunds = errors.New("insufficient funds")

// Withdraw returns the new balance, or an error if amount > balance.
func Withdraw(balance, amount int) (int, error) {
    if balance >= amount {
        return balance - amount, nil
    }

    // Returning 0 here is acceptable because the caller should check the error
    // before relying on the balance value. In Go, the convention is:
    // if err != nil { ignore/handle the returned value }
    //
    // We return the sentinel directly because it is a standard package-level error.
    // Callers can match it with errors.Is(err, ErrInsufficientFunds).
    //
    // If we want extra context, we can wrap it instead:
    // return 0, fmt.Errorf("withdraw %d from %d: %w", amount, balance, ErrInsufficientFunds)
    return 0, ErrInsufficientFunds
}

// Pay withdraws `amount` for `item`. On failure it wraps the underlying error
// with context: "pay <item>: <cause>".
func Pay(balance, amount int, item string) (int, error) {
    bal, err := Withdraw(balance, amount)
    if err != nil {
        return 0, fmt.Errorf("pay %s: %w", item, err)
    }
    return bal, nil
}

// ValidationError carries which field failed.
type ValidationError struct {
    Field string
    Msg   string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Msg)
     // Sprintf stands for Sprintf stands for:
//                            S = String
//                            printf = formatted print
//                            In Go, fmt.Sprintf creates and returns a formatted string instead of printing it to the screen.
}

// Validate returns a *ValidationError if name is empty.
func Validate(name string) error {
    if name == "" {
        return &ValidationError{Field: "name", Msg: "is required"}
    }
    return nil
}
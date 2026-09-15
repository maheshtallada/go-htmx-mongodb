package main

import (
    "errors"
    "testing"
)

func TestWithdraw(t *testing.T) {
    bal, err := Withdraw(100, 30)
    if err != nil || bal != 70 {
        t.Fatalf("got (%d, %v); want (70, nil)", bal, err)
    }

    // _ is the blank identifier: it means "ignore this returned value".
    // Here we only care about the error, so we discard the balance.
    _, err = Withdraw(50, 80)

    // errors.Is checks whether err matches ErrInsufficientFunds.
    // This is the preferred way to compare errors, especially if they are wrapped.
    if !errors.Is(err, ErrInsufficientFunds) {
        t.Fatalf("want ErrInsufficientFunds, got %v", err)
    }
}


func TestPayWrap(t *testing.T) {
    _, err := Pay(20, 50, "coffee")
    if !errors.Is(err, ErrInsufficientFunds) {
        t.Fatalf("wrapped error should still match sentinel; got %v", err)
    }
    if got := err.Error(); got != "pay coffee: insufficient funds" {
        t.Fatalf("message = %q", got)
    }
}


func TestErrorsAs(t *testing.T) {
    err := Validate("")
    var ve *ValidationError
    if !errors.As(err, &ve) {
        t.Fatalf("expected a *ValidationError, got %v", err)
    }
    if ve.Field != "name" {
        t.Fatalf("field = %q; want name", ve.Field)
    }
}
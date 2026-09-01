package main

// Describer is anything that can describe itself.
type Describer interface {
    Describe() string
}
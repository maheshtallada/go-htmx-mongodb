package main

import (
    "reflect"
    "testing"
)

func TestFilter(t *testing.T) {
    got := Filter([]int{1, 2, 3, 4, 5, 6})
    want := []int{2, 4, 6}
    if !reflect.DeepEqual(got, want) {
        t.Fatalf("got %v; want %v", got, want)
    }
}

func TestSumRange(t *testing.T) {
    if SumRange([]int{1, 2, 3, 4}) != 10 {
        t.Fatal("sum wrong")
    }
}
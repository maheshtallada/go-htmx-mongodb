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

func TestWordCount(t *testing.T) {
    got := WordCount([]string{"go", "htmx", "go", "mongo", "go"})
    if got["go"] != 3 || got["htmx"] != 1 {
        t.Fatalf("counts = %v", got)
    }
    if _, ok := got["missing"]; ok {
        t.Fatal("missing key should not exist")
    }
}

func TestMapSlice(t *testing.T) {
    lengths := MapSlice([]string{"go", "htmx"}, func(s string) int { return len(s) })
    if lengths[0] != 2 || lengths[1] != 4 {
        t.Fatalf("lengths = %v", lengths)
    }
}

func TestKeys(t *testing.T) {
    input := map[string]int{"go": 1, "htmx": 2}
    got := Keys(input)

    // Map iteration order is unspecified, so check membership instead of order.
    gotSet := make(map[string]bool, len(got))
    for _, key := range got {
        gotSet[key] = true
    }

    if len(gotSet) != len(input) || !gotSet["go"] || !gotSet["htmx"] {
        t.Fatalf("Keys(%v) = %v", input, got)
    }
}
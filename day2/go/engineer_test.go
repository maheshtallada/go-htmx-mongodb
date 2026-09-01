package main

import "testing"

func TestNewEngineer(t *testing.T) {
    e := NewEngineer("Asha", "SDE II", 5, "Go", "MongoDB")
    if e.Name != "Asha" || e.Years != 5 {
        t.Fatalf("bad fields %+v", e)
    }
    if len(e.Skills) != 2 {
        t.Fatalf("skills = %v; want 2", e.Skills)
    }
}

func TestMethods(t *testing.T) {
    e := NewEngineer("Asha", "SDE II", 5, "Go")
    if e.Summary() != "Asha (SDE II), 5y" {
        t.Fatalf("summary = %q", e.Summary())
    }
    e.Promote("SDE III")
    if e.Role != "SDE III" || e.Years != 6 {
        t.Fatalf("after promote: %+v", e)
    }
}

func TestInterfaces(t *testing.T) {
    items := []Describer{
        NewEngineer("Ravi", "SDE III", 8, "Go"),
        NewEngineer("Asha", "SDE II", 5, "Go"),
    }
    got := PrintAll(items)
    want := "I am Ravi, a SDE III.\nI am Asha, a SDE II."
    if got != want {
        t.Fatalf("got:\n%s\nwant:\n%s", got, want)
    }
}
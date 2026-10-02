package main

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"
)

// A1
func TestSquareAll(t *testing.T) {
	got := SquareAll([]int{1, 2, 3, 4})
	want := []int{1, 4, 9, 16}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v; want %v", got, want)
	}
}

// go test -race ./... shows no data race --> to
//  execute this we need a cgo or gcc c compiler

// A2
func TestChannels(t *testing.T) {
	if SumChan(Produce(5)) != 15 {
		t.Fatal("sum wrong")
	}
}

// A3
func TestContextCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	_, err := SlowDouble(ctx, 21)
	if err == nil {
		t.Fatal("expected context deadline error")
	}
}

// A4
func TestCounterConcurrentIncrement(t *testing.T) {
	var counter Counter
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()
	if got := counter.Value(); got != 100 {
		t.Fatalf("counter = %d; want 100", got)
	}
}

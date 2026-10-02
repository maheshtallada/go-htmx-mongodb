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

// A4 asks us to let 100 goroutines increment one shared counter and verify the
// result is exactly 100 without a data race; this test coordinates and checks it.
func TestCounterConcurrentIncrement(t *testing.T) {
	var counter Counter    // One shared counter instance is used by all workers.
	var wg sync.WaitGroup  // Separate from the mutex: waits for all workers to finish.

	for i := 0; i < 100; i++ {
		wg.Add(1) // Count this worker before starting it.
		go func() {
			defer wg.Done()       // Tell the test this worker is finished.
			counter.Increment()   // The Counter method's mutex protects the shared update.
		}()
	}

	wg.Wait() // Do not check the total until all 100 increments have completed.
	if got := counter.Value(); got != 100 { // Read through the synchronized method and assert A4's expected total.
		t.Fatalf("counter = %d; want 100", got)
	}
}

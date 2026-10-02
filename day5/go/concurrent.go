package main // Keep these concurrency exercises in the same package as their tests.

import (
	"context" // Provides cancellation signals and their associated errors.
	"sync"    // Provides the WaitGroup and Mutex synchronization primitives.
	"time"    // Provides the delay and timer used by SlowDouble.
)

// SquareAll computes the square of every number concurrently and keeps input order.
func SquareAll(nums []int) []int {
	out := make([]int, len(nums)) // Reserve one output slot per input; each goroutine owns one slot.
	var wg sync.WaitGroup         // Track all workers so this function waits for their writes.

	for i, n := range nums { // Read each input index and number before starting its worker.
		wg.Add(1)                    // Register this worker before it starts, so Wait cannot finish too early.
		go func(index, number int) { // Pass loop values as arguments so each worker gets its own copy.
			defer wg.Done()              // Mark this worker complete even when its function returns.
			out[index] = number * number // Write only this worker's distinct output index.
		}(i, n) // Start the worker with the current index and number.
	}

	wg.Wait()  // Wait until every worker has written its result.
	return out // Return the completed slice in the original input order.
}

// Produce sends the integers from 1 through n, then closes the returned channel.
func Produce(n int) <-chan int {
	ch := make(chan int) // Create an unbuffered channel for values sent by the producer.

	go func() { // Run sending separately so the caller can receive while production happens.
		defer close(ch)                       // Close the channel after all sends, including when n is less than 1.
		for value := 1; value <= n; value++ { // Generate every integer in the inclusive range 1..n.
			ch <- value // Send this value; it waits until a receiver is ready.
		}
	}()

	return ch // Expose receive-only access so callers can read but cannot send or close it.
}

// SumChan reads all values from ch until closed and returns the sum.
func SumChan(ch <-chan int) int {
	sum := 0                // Start at zero, the additive identity for integer addition.
	for value := range ch { // Receive each value and stop automatically when ch is closed.
		sum += value // Add the received value to the running total.
	}
	return sum // Return the total after the channel has closed and drained.
}

// SlowDouble returns n multiplied by 2 after 100ms, unless ctx is cancelled first.
// When cancellation wins, it returns zero and the context's cancellation error.
func SlowDouble(ctx context.Context, n int) (int, error) {
	timer := time.NewTimer(100 * time.Millisecond) // Start a timer that represents the normal delay.
	defer timer.Stop()                             // Release timer resources if cancellation happens first.

	select { // Wait for either the delay to finish or the caller to cancel.
	case <-timer.C: // The timer fired, so the requested delay elapsed.
		return n * 2, nil // Return the doubled number and no error.
	case <-ctx.Done(): // The context was cancelled or its deadline expired.
		return 0, ctx.Err() // Return the context's reason for stopping.
	}
}

// Counter stores an integer that can be safely shared by concurrent goroutines.
type Counter struct {
	mu    sync.Mutex // Protect value so only one goroutine accesses it at a time.
	value int        // Hold the current count; its zero value starts the count at zero.
}

// Increment adds one to the counter while holding its mutex.
func (c *Counter) Increment() {
	c.mu.Lock()         // Acquire exclusive access before reading or changing value.
	defer c.mu.Unlock() // Release exclusive access when this method returns.
	c.value++           // Update the protected count without a data race.
}

// Value returns the counter's current value while holding its mutex.
func (c *Counter) Value() int {
	c.mu.Lock()         // Acquire the same lock used by Increment for a consistent read.
	defer c.mu.Unlock() // Release the lock whether the method returns normally or not.
	return c.value      // Copy the protected value to the caller.
}

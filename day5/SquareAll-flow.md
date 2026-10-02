# `SquareAll` execution flow

Example input: `nums := []int{1, 2, 3, 4}`

```mermaid
flowchart TD
    A["SquareAll starts"] --> B["Create out with 4 slots: [0, 0, 0, 0]"]
    B --> C["Create WaitGroup: counter = 0"]
    C --> D["Loop through nums in order"]

    D --> E1["i=0, n=1"]
    D --> E2["i=1, n=2"]
    D --> E3["i=2, n=3"]
    D --> E4["i=3, n=4"]

    E1 --> F1["wg.Add(1), start worker"]
    E2 --> F2["wg.Add(1), start worker"]
    E3 --> F3["wg.Add(1), start worker"]
    E4 --> F4["wg.Add(1), start worker"]

    F1 --> G1["Worker: out[0] = 1*1"]
    F2 --> G2["Worker: out[1] = 2*2"]
    F3 --> G3["Worker: out[2] = 3*3"]
    F4 --> G4["Worker: out[3] = 4*4"]

    G1 --> H1["defer wg.Done(): counter - 1"]
    G2 --> H2["defer wg.Done(): counter - 1"]
    G3 --> H3["defer wg.Done(): counter - 1"]
    G4 --> H4["defer wg.Done(): counter - 1"]

    H1 --> I["wg.Wait() returns when counter = 0"]
    H2 --> I
    H3 --> I
    H4 --> I

    I --> J["Return out: [1, 4, 9, 16]"]
```

## What happens in order?

1. **Sequential setup:** `SquareAll` allocates the output slice and starts its loop. The loop visits each input in order.
2. **Launch workers:** For each number, `wg.Add(1)` increments the WaitGroup counter, then `go` starts a worker with that item's index and value.
3. **Concurrent work:** The workers calculate and write results. Their completion order can differ from the input order; each writes to its own output index.
4. **Wait at the barrier:** Once the loop has started all workers, `wg.Wait()` blocks until all workers have called `Done`.
5. **Return:** Only after the counter reaches zero does the function return the completed output slice.

> **Important:** “Concurrent” means the workers can make progress independently and their execution may overlap; it does not guarantee that they all run at the exact same instant. The output order is preserved by writing each answer to `out[index]`, not by controlling which worker finishes first.

## Simplified timeline

```text
SquareAll (main goroutine): allocate out
  loop: Add(1) -> launch worker 0
       Add(1) -> launch worker 1
       Add(1) -> launch worker 2
       Add(1) -> launch worker 3
  Wait() -----------------------------------------------+
                                                       |
Worker 0: calculate 1*1 -> write out[0] -> Done() -----|
Worker 1: calculate 2*2 -> write out[1] -> Done() -----|-- counter is 0
Worker 2: calculate 3*3 -> write out[2] -> Done() -----|   Wait returns
Worker 3: calculate 4*4 -> write out[3] -> Done() -----+   return [1,4,9,16]
```

The worker lines are shown separately to illustrate parallel work. The loop
launches workers one by one, but after each launch the worker can run while
the loop continues launching the others.

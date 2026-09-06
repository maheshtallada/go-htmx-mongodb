# Day 5 Assignment — Test-Driven Learning (Go · HTMX · MongoDB · English)

> **How this works:** Same as before. Make each **failing test or checkpoint**
> pass. Answer keys at the bottom — try first.
>
> **Time-box:** ~3 hours + 30–45m English.
>
> **Day 5 scope**
> - **Go:** goroutines, channels, `sync` (WaitGroup/Mutex), `context`
> - **HTMX:** small interactive UI patterns (search + auto-refresh/polling)
> - **MongoDB:** schema design — embedding vs referencing
> - **English:** give a clear opinion, modal verbs, 5 vocab words
>
> **Warm-up:** finish any red items from Day 4 first.

---

## 0. Setup (5 min)

```bash
docker start mongo-day1 || docker run -d --name mongo-day5 -p 27017:27017 mongo:7
mkdir -p day5/go && cd day5/go
go mod init day5
```

---

## PART A — Go: Concurrency (~60 min)

### A1. Goroutines + WaitGroup
**Task:** In `concurrent.go`:

```go
// SquareAll computes squares concurrently and returns them in input order.
func SquareAll(nums []int) []int
```

Run one goroutine per number, collect results in order (index-based write is
safe if each goroutine writes its own slot).

**Checkpoint A1** — `concurrent_test.go`:

```go
package main

import (
	"reflect"
	"testing"
)

func TestSquareAll(t *testing.T) {
	got := SquareAll([]int{1, 2, 3, 4})
	want := []int{1, 4, 9, 16}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v; want %v", got, want)
	}
}
```

- [ ] Uses `sync.WaitGroup`
- [ ] Order preserved
- [ ] `go test -race ./...` shows **no data race**

> **Micro-questions:**
>
> 1. What does `sync.WaitGroup` do (`Add`/`Done`/`Wait`)?
> 2. Why is writing to distinct slice indices race-free but `append` from
>    goroutines is not?

### A2. Channels
**Task:** Add:

```go
// Produce sends 1..n on the returned channel, then closes it.
func Produce(n int) <-chan int

// SumChan reads all values from ch until closed and returns the sum.
func SumChan(ch <-chan int) int
```

**Checkpoint A2**:

```go
func TestChannels(t *testing.T) {
	if SumChan(Produce(5)) != 15 {
		t.Fatal("sum wrong")
	}
}
```

- [ ] `Produce` closes the channel when done
- [ ] `SumChan` uses `for v := range ch`

> **Trap check:** What happens if you read from a channel that is never closed?
> What about sending on a closed channel?

### A3. `context` for cancellation (the tested core)
**Task:** Add:

```go
// SlowDouble returns n*2 after 100ms, but respects ctx cancellation.
// If ctx is done first, it returns ctx.Err().
func SlowDouble(ctx context.Context, n int) (int, error)
```

**Checkpoint A3**:

```go
func TestContextCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	_, err := SlowDouble(ctx, 21)
	if err == nil {
		t.Fatal("expected context deadline error")
	}
}
```

- [ ] Uses `select` on `ctx.Done()` vs a timer
- [ ] Returns `ctx.Err()` on cancellation

**Checkpoint A (all tests, race detector on)**

```bash
go test -race ./...
# expected: PASS (ok day5)
```

> **Micro-questions:**
>
> 1. What is a **deadlock** and how does Go detect some at runtime?
> 2. When do you use a `Mutex` vs a channel?
> 3. Why pass `context.Context` as the **first** parameter by convention?

### A4. Stretch (optional)
Use a `sync.Mutex` to build a concurrency-safe counter incremented by 100
goroutines; assert the total is exactly 100.
- [ ] No race, correct total

---

---

# Day 5 Assignment

## Part B — HTMX Polling & Search

```go
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, page)
	})
	http.HandleFunc("/now", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<span>%s</span>", time.Now().Format(time.Kitchen))
	})
	http.ListenAndServe(":8080", nil)
}

const page = `<!doctype html>
<html><head><script src="https://unpkg.com/htmx.org@2"></script></head>
<body>
	<div hx-get="/now" hx-trigger="every 2s" hx-swap="innerHTML">loading...</div>
</body></html>`
```

**Checkpoint B1**
- [ ] Time updates every 2s without reload
- [ ] Uses `hx-trigger="every 2s"`

### B2. Search + refresh combo

**Task:** Add a filterable list (reuse Day 4 search) **and** a manual "Refresh" button using `hx-get` on the same target.

**Checkpoint B2**
- [ ] Typing filters live
- [ ] Refresh button re-fetches the list
- [ ] Both target the same `#results` container

> **Micro-questions:**
>
> 1. When is polling (`every 2s`) the right choice vs. server-sent events / websockets?
> 2. How do you **stop** a poll from the server side? (hint: `HX-Trigger` / return no `hx-trigger`)
> 3. What's the risk of a too-frequent poll?

### B3. Stretch (optional)

Make the poll interval back off (2s → 5s) after N refreshes by returning a fragment with a new `hx-trigger`.
- [ ] Interval changes dynamically

---

## PART C — MongoDB Schema Design: Embed vs Reference (~30 min)

### C1. Embedding

**Task:** Model a blog **post** with comments embedded:

```js
use day5
db.posts.insertOne({
    title: "Learning Go",
    author: "Mahesh",
    comments: [
        { user: "Asha", text: "Great!" },
        { user: "Ravi", text: "Helpful" }
    ]
})
```

Query: fetch the post and read `comments.0.user`.
- [ ] One document holds the post + its comments
- [ ] You can read a nested comment

### C2. Referencing

**Task:** Model **users and orders as separate collections** linked by `userId`:

```js
db.customers.insertOne({ _id: 1, name: "Neha" })
db.invoices.insertMany([
    { userId: 1, amount: 100 },
    { userId: 1, amount: 250 }
])
// "join" the two:
db.invoices.aggregate([
    { $match: { userId: 1 } },
    { $lookup: { from: "customers", localField: "userId", foreignField: "_id", as: "customer" } }
])
```

- [ ] Two collections linked by `userId`
- [ ] `$lookup` joins them

### C3. Decision (the tested core – written)

For each scenario, choose **embed** or **reference** and justify:

1. Product + its 3 fixed spec fields.
2. User + their (potentially thousands of) activity events.
3. Order + its line items (read together, rarely change).
4. Blog post + author profile shared across many posts.

- [ ] All 4 decided with a one-line reason

### C4. Understanding (written)

1. What's the **16MB document limit** and how does it affect embedding?
2. Rule of thumb: "embed for **____ together**, reference for **____ / unbounded growth**."
3. One downside of referencing (hint: extra queries / `$lookup` cost).

- [ ] All 3 answered

---

## PART D — English (~30–45 min)

### D1. Vocabulary

One sentence each: `important`, `useful`, `effective`, `practical`, `appropriate`.

- [ ] 5 original sentences

### D2. Grammar — modal verbs

Choose the best modal (`can`, `should`, `must`, `may`, `might`):

1. You ______ always validate input on the server. *(strong obligation)*
2. We ______ use MongoDB here – it fits the flexible schema. *(recommendation)*
3. This design ______ scale to 10x, but we should load-test. *(possibility)*
4. Developers ______ read the runbook before an on-call shift. *(requirement)*
5. You ______ deploy on Friday if the tests pass. *(permission)*

- [ ] Attempt all 5

### D3. Speaking – give a clear opinion (record yourself)

Speak **60–90 seconds** answering *"Was choosing HTMX over a SPA a good decision? Why?"*

Structure: **Opinion → 2 reasons → example → conclusion.**

**Checkpoint D3**
- [ ] Recorded once, listened back
- [ ] Clear stance in the first sentence
- [ ] At least one modal verb used correctly

---

## Day 5 Done — Definition of Complete

- [ ] Go: `go test -race ./...` green (A1–A3, optional A4)
- [ ] HTMX: polling clock (B1) + search/refresh (B2)
- [ ] MongoDB: embedding + referencing + `$lookup` + decisions (C1–C3)
- [ ] English: vocab + modal grammar + opinion recording
- [ ] `answers/day5-notes.md` complete

Score yourself: **_/5 sections green.** 🎉 **Week 1 complete!**

---

## 🔑 Answer Key (try first, then check)

<details>
<summary>Go answers</summary>

```go
// concurrent.go
package main

import (
	"context"
	"sync"
	"time"
)

func SquareAll(nums []int) []int {
	out := make([]int, len(nums))
	var wg sync.WaitGroup
	for i, n := range nums {
		wg.Add(1)
		go func(i, n int) {
			defer wg.Done()
			out[i] = n * n
		}(i, n)
	}
	wg.Wait()
	return out
}

func Produce(n int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 1; i <= n; i++ {
			ch <- i
		}
	}()
	return ch
}

func SumChan(ch <-chan int) int {
	sum := 0
	for v := range ch {
		sum += v
	}
	return sum
}

func SlowDouble(ctx context.Context, n int) (int, error) {
	select {
	case <-time.After(100 * time.Millisecond):
		return n * 2, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}
```

- **A1.1** `WaitGroup`: `Add(n)` sets the counter, `Done()` decrements, `Wait()` blocks until zero.
- **A1.2** Distinct indices are separate memory; `append` can reallocate/race on shared length.
- **A2 trap** Reading a never-closed channel blocks forever (possible deadlock); sending on a closed channel **panics**.
- **A3.2** Use a `Mutex` to protect shared state; use a channel to **pass ownership/communicate**. “Don’t communicate by sharing memory; share memory by communicating.”

</details>

<details>
<summary>HTMX answers</summary>

- **B.1** Polling suits simple, low-frequency updates; use SSE/websockets for high-frequency or push-driven data.
- **B.2** Return a response without the repeating `hx-trigger`, or send an `HX-Trigger` event that a listener uses to cancel; simplest is to swap the polling element with one that no longer polls.
- **B.3** Too-frequent polls waste server/network resources and can amplify load under many clients.

</details>

<details>
<summary>MongoDB answers</summary>

- **C3** 1. **Embed** (small, fixed, read together). 2. **Reference** (unbounded growth → would blow the doc size). 3. **Embed** (read together, stable). 4. **Reference** (shared entity, avoid duplication).
- **C4.1** A single BSON document can't exceed **16MB**; embedding unbounded arrays risks hitting it.
- **C4.2** Embed for **read together**, reference for **shared data / unbounded growth**.
- **C4.3** Referencing needs extra round-trips or `$lookup`, which costs more than reading one embedded doc.

</details>

<details>
<summary>English answers</summary>

**D2 grammar:** 1. **must** · 2. **should** · 3. **might/may** · 4. **must** · 5. **can/may**

**Meanings:**
- `must` = strong obligation/requirement.
- `should` = advice/recommendation.
- `can` / `may` = ability/permission.
- `might` / `may` = possibility.

</details>

---

## Tomorrow (Day 6 preview)

`net/http` routing with Chi + middleware · HTMX with real routes and partial responses · MongoDB with the Go driver (context-aware CRUD repository) · English: polished 90-second “tell me about yourself” + sentence structure.

**Week 2 begins.**

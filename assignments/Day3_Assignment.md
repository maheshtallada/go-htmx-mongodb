# Day 3 Assignment — Test-Driven Learning (Go · HTMX · MongoDB · English)

> **How this works:** Same as before. Each task tells you **what to
> build/answer**, gives you a **failing test or checkpoint**, and you make it
> pass by learning just enough. Answer keys are at the bottom — no peeking until
> you try.

> **Time-box:** ~3 hours + 30–45m English. Aim for **green checkmarks**, not
> perfection.

> **Day 3 scope**
> - **Go:** errors as values, wrapping (`fmt.Errorf` + `%w`), `errors.Is` /
>   `errors.As`, no exceptions
> - **HTMX:** forms, server-side validation, redirects (`HX-Redirect`)
> - **MongoDB:** query filters, projections, sorting, limits
> - **English:** talk about a challenge, past simple vs past continuous, 5 vocab
>   words

> **Warm-up:** finish any red items from Day 2 first (5–10 min).

---

## 0. Setup (5 min)

```bash
go version
docker start mongo-day1 || docker run -d --name mongo-day3 -p 27017:27017 mongo:7
mkdir -p day3/go && cd day3/go
go mod init day3
```

- [ ] Go builds and runs
- [ ] MongoDB reachable on `localhost:27017`

---

## PART A — Go: Errors as Values (~60 min)

### A1. Returning errors (not throwing)

**Task:** In `bank.go`, create:

```go
package main

import "errors"

// ErrInsufficientFunds is a sentinel error callers can check for.
var ErrInsufficientFunds = errors.New("insufficient funds")

// Withdraw returns the new balance, or an error if amount > balance.
func Withdraw(balance, amount int) (int, error)
```

**Checkpoint A1:** write `bank_test.go` and make it pass:

```go
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

    _, err = Withdraw(50, 80)
    if !errors.Is(err, ErrInsufficientFunds) {
        t.Fatalf("want ErrInsufficientFunds, got %v", err)
    }
}
```

- [ ] Happy path returns new balance, `nil` error
- [ ] Overdraw returns `ErrInsufficientFunds`

> **Micro-questions:**
> 1. Why does Go return errors instead of throwing exceptions?
> 2. What is a **sentinel error**?

### A2. Wrapping errors with `%w`

**Task:** Add a higher-level function that wraps the low-level error with
context:

```go
// Pay withdraws `amount` for `item`. On failure it wraps the underlying error
// with context: "pay <item>: <cause>".
func Pay(balance, amount int, item string) (int, error)
```

**Checkpoint A2:**

```go
func TestPayWrap(t *testing.T) {
    _, err := Pay(20, 50, "coffee")
    if !errors.Is(err, ErrInsufficientFunds) {
        t.Fatalf("wrapped error should still match sentinel; got %v", err)
    }
    if got := err.Error(); got != "pay coffee: insufficient funds" {
        t.Fatalf("message = %q", got)
    }
}
```

- [ ] `Pay` wraps with `fmt.Errorf("pay %s: %w", item, err)`
- [ ] `errors.Is` still finds the sentinel **through** the wrap

### A3. `errors.As` (the tested core)

**Task:** Define a custom error type carrying data, then extract it:

```go
// ValidationError carries which field failed.
type ValidationError struct {
    Field string
    Msg   string
}

func (e *ValidationError) Error() string

// Validate returns a *ValidationError if name is empty.
func Validate(name string) error
```

**Checkpoint A3:**

```go
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
```

- [ ] `errors.As` extracts the concrete type
- [ ] You understand `Is` (identity) vs `As` (type + data)

**Checkpoint A (all Go tests)**

```bash
go test ./...
# expected: PASS (ok day3)
```

> **Micro-questions:**
> 1. Difference between `errors.Is` and `errors.As`?
> 2. What does `%w` do that `%v` does not?
> 3. When should you wrap vs. return the error unchanged?

### A4. Stretch (optional)

Chain two wraps (`Pay` → `Withdraw`) and prove `errors.Is` still finds the root
cause through **both** layers.

- [ ] Multi-level unwrap works

---

## PART B — HTMX Forms, Validation & Redirects (~45 min)

### B1. A form that validates server-side

**Task:** In `day3/htmx/server.go`, serve a signup form (`GET /`) and handle
`POST /signup`. Validate that `email` contains `@`. On error, return the form
fragment **with an error message** and HTTP `422`.

Minimal starter:

```go
package main

import (
    "fmt"
    "net/http"
    "strings"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, page)
    })
    http.HandleFunc("/signup", handleSignup)
    http.ListenAndServe(":8080", nil)
}

func handleSignup(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    email := r.FormValue("email")
    if !strings.Contains(email, "@") {
        w.WriteHeader(http.StatusUnprocessableEntity) // 422
        fmt.Fprintf(w, `<p class="error" id="msg">Invalid email</p>`)
        return
    }
    // TODO B2: on success, redirect via HX-Redirect
}

const page = `<!doctype html>
<html><head><script src="https://unpkg.com/htmx.org@2"></script></head>
<body>
    <form hx-post="/signup" hx-target="#msg" hx-swap="outerHTML">
        <input name="email" placeholder="you@example.com">
        <button type="submit">Sign up</button>
    </form>
    <p id="msg"></p>
</body></html>`
```

**Checkpoint B1** (manual): submit `bad-email`.

- [ ] Error message appears inline, no reload
- [ ] Response status is **422**

### B2. Redirect on success (`HX-Redirect`)

**Task:** When the email is valid, respond with the `HX-Redirect` header so the
browser navigates to `/welcome`.

```go
w.Header().Set("HX-Redirect", "/welcome")
w.WriteHeader(http.StatusOK)
```

Add a `GET /welcome` that prints `Welcome!`.

**Checkpoint B2:** submit `me@site.com`.

- [ ] Browser navigates to `/welcome`
- [ ] You used the `HX-Redirect` response header

> **Micro-questions:**
> 1. Why return `422` instead of `200` for a validation error?
> 2. Difference between `HX-Redirect` and `HX-Location`?
> 3. Why do server-side validation even if you also validate in the browser?

### B3. Stretch (optional)

Add a required `name` field; show **two** independent error slots (name + email)
using out-of-band swaps.

- [ ] Both fields validate independently

---

## PART C — MongoDB Queries: Filters, Projections, Sort, Limit (~30 min)

### C1. Seed

```js
use day3
db.orders.insertMany([
  { customer: "Asha", total: 250, status: "paid", items: 3 },
  { customer: "Ravi", total: 90, status: "pending", items: 1 },
  { customer: "Mahesh", total: 500, status: "paid", items: 5 },
  { customer: "Neha", total: 150, status: "paid", items: 2 },
  { customer: "Kiran", total: 40, status: "pending", items: 1 }
])
```

- [ ] 5 documents inserted

### C2. Filters & operators

1. Orders with `total >= 150`.
2. `{status: "paid" and items >= 3}`.
3. Orders where `status` is `"paid"` **or** `total < 100` (Use `$or`).

Starter:

```js
db.orders.find({ total: { $gte: 150 } })
```

- [ ] Q1 returns Asha, Mahesh, Neha
- [ ] Q2 returns Asha, Mahesh
- [ ] Q3 uses `$or`

### C3. Projection + sort + limit (the tested core)

**Task:** Return the **top 2** orders by `total` (descending), showing only
`customer` and `total` (hide `_id`).

```js
db.orders.find({}, { _id: 0, customer: 1, total: 1 }).sort({ total: -1 }).limit(2)
```

**Checkpoint C3**

- [ ] Returns Mahesh(500) then Asha(250)
- [ ] Only `customer` + `total` shown
- [ ] `.sort()` then `.limit()` chained correctly

### C4. Understanding (written)

1. What does `1` vs `0` mean in a projection? Can you mix them?
2. What does `.sort({ total: -1 })` mean? What is `1`?
3. Why use `.limit()` for pagination performance?

- [ ] All 3 answered

---

## PART D — English (~30–45 min)

### D1. Vocabulary

Write **one sentence for each** about your real work:
`challenge`, `solution`, `impact`, `result`, `improvement`.

- [ ] 5 original sentences

### D2. Grammar — past simple vs past continuous

Fill the correct form, then check the key:

1. While I _____ (debug) the service, the alert _____ (fire).
2. Yesterday I _____ (fix) the production bug in 20 minutes.
3. The team _____ (deploy) when the database _____ (go) down.
4. I _____ (review) code all morning.
5. When the incident started, we _____ (run) the load test.

- [ ] Attempt all 5

### D3. Speaking — describe a challenge (record yourself)

Speak **90 seconds** in **STAR** format answering *"Tell me about a technical
challenge you solved."*

- **Situation → Task → Action → Result**

**Checkpoint D3**

- [ ] Recorded once, listened back
- [ ] Followed STAR order
- [ ] Ended with a measurable **result**

---

## ✅ Day 3 Done — Definition of Complete

- [ ] Go: `go test ./...` green (A1–A3, optional A4)
- [ ] HTMX: inline validation error (B1) + `HX-Redirect` success (B2)
- [ ] MongoDB: filters, projection, sort, limit all correct (C2–C3)
- [ ] English: vocab + grammar attempt + STAR challenge recording
- [ ] `answers/day3-notes.md` has all micro-question answers

Score yourself: **__/5 sections green.**

---

## 🔑 Answer Key (try first, then check)

<details>
<summary>Go answers</summary>

```go
// bank.go
package main

import (
    "errors"
    "fmt"
)

var ErrInsufficientFunds = errors.New("insufficient funds")

func Withdraw(balance, amount int) (int, error) {
    if amount > balance {
        return balance, ErrInsufficientFunds
    }
    return balance - amount, nil
}

func Pay(balance, amount int, item string) (int, error) {
    bal, err := Withdraw(balance, amount)
    if err != nil {
        return balance, fmt.Errorf("pay %s: %w", item, err)
    }
    return bal, nil
}

type ValidationError struct {
    Field string
    Msg   string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Msg)
}

func Validate(name string) error {
    if name == "" {
        return &ValidationError{Field: "name", Msg: "is required"}
    }
    return nil
}
```

- **A1.1** Errors as values = explicit handling, no hidden control flow; the
  compiler forces you to deal with them.
- **A1.2** A **sentinel** is a predefined error value you compare against with
  `errors.Is`.
- **A3** `errors.Is` checks **identity** (same sentinel, through wraps);
  `errors.As` checks **type** and lets you read the concrete error's fields.
- **A.%** `%w` wraps so the chain is unwrappable by `Is`/`As`; `%v` just
  formats text and loses the chain.

</details>

<details>
<summary>HTMX answers</summary>

```go
// success branch of handleSignup
w.Header().Set("HX-Redirect", "/welcome")
w.WriteHeader(http.StatusOK)
```

- **B.1** `422 Unprocessable Entity` tells the client the request was
  well-formed but semantically invalid — correct for validation.
- **B.2** `HX-Redirect` does a client-side full navigation; `HX-Location` does
  an HTMX-style AJAX navigation (no full reload).
- **B.3** Never trust the client — the browser can be bypassed; the server is
  the source of truth.

</details>

<details>
<summary>MongoDB answers</summary>

```js
// C2.2
db.orders.find({ status: "paid", items: { $gte: 3 } })
// C2.3
db.orders.find({ $or: [ { status: "paid" }, { total: { $lt: 100 } } ] })
// C3
db.orders.find({}, { _id: 0, customer: 1, total: 1 }).sort({ total: -1 }).limit(2)
```

- **C4.1** `1` = include, `0` = exclude. You can't mix includes and excludes
  (except excluding `_id`).
- **C4.2** `-1` = descending, `1` = ascending.
- **C4.3** `.limit()` caps documents scanned/returned — combined with an index
  it avoids full collection scans.

</details>

<details>
<summary>English answers</summary>

**D2 grammar:**

1. **was debugging** / **fired** (background action + interrupting event)
2. **fixed** (past simple — finished, "yesterday")
3. **was deploying** / **went** (in-progress + interruption)
4. **was reviewing** (past continuous — ongoing duration "all morning")
5. **were running** (past continuous — in progress when something happened)

**Rule of thumb:**

- **Past continuous** → an action *in progress* in the past (was/were + -ing).
- **Past simple** → a *completed* action; often the shorter event that
  interrupts the continuous one.

</details>

---

## Tomorrow (Day 4 preview)

Slices, maps, `range`, generics intro · HTMX `hx-trigger` + loading states +
swap patterns · MongoDB indexes, compound indexes, `explain()` · English:
strengths/weaknesses + comparatives/superlatives. Bring your Day 3 red items.

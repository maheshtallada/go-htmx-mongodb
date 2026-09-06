# Day 4 Assignment — Test-Driven Learning (Go · HTMX · MongoDB · English)

> **How this works:** Same as before. Build/answer to make each **failing test or checkpoint** pass. Answer keys at the bottom — try first.

> **Time-box:** ~3 hours + 30–45m English.

> **Day 4 scope**
> - **Go:** slices, maps, `range`, generics intro
> - **HTMX:** `hx-trigger`, loading states / indicators, swap patterns
> - **MongoDB:** indexes, compound indexes, `explain()` basics
> - **English:** strengths & weaknesses, comparatives & superlatives, 5 vocab words
>
> **Warm-up:** finish any red items from Day 3 first.

---

## 0. Setup (5 min)

```bash
docker start mongo-day1 || docker run -d --name mongo-day4 -p 27017:27017 mongo:7
mkdir -p day4/go && cd day4/go
go mod init day4
```

---

## PART A — Go: Slices, Maps, Range, Generics (~60 min)

### A1. Slices & `range`

**Task:** In `collections.go`:

```go
// Filter returns only the even numbers, preserving order.
func Filter(nums []int) []int

// SumRange returns the sum of nums using range.
func SumRange(nums []int) int
```

**Checkpoint A1** — write `collections_test.go`:

```go
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
```

- [ ] `Filter` returns a new slice
- [ ] You used `range` in both

> **Micro-questions:**
> 1. Difference between a slice's **length** and **capacity**?
> 2. What happens to the underlying array when `append` exceeds capacity?

### A2. Maps

**Task:** Add:

```go
// WordCount returns how many times each word appears.
func WordCount(words []string) map[string]int
```

**Checkpoint A2:**

```go
func TestWordCount(t *testing.T) {
    got := WordCount([]string{"go", "htmx", "go", "mongo", "go"})
    if got["go"] != 3 || got["htmx"] != 1 {
        t.Fatalf("counts = %v", got)
    }
    if _, ok := got["missing"]; ok {
        t.Fatal("missing key should not exist")
    }
}
```

- [ ] Map counts correctly
- [ ] You used the `value, ok := m[key]` comma-ok idiom in the test

> **Trap check:** Is iteration order of a Go map guaranteed? Write down the answer.

### A3. Generics (the tested core)

**Task:** Write a generic `Map` (transform) function:

```go
// MapSlice applies fn to each element and returns a new slice.
func MapSlice[T any, U any](in []T, fn func(T) U) []U
```

**Checkpoint A3:**

```go
func TestMapSlice(t *testing.T) {
    lengths := MapSlice([]string{"go", "htmx"}, func(s string) int { return len(s) })
    if lengths[0] != 2 || lengths[1] != 4 {
        t.Fatalf("lengths = %v", lengths)
    }
}
```

- [ ] Generic type parameters `[T any, U any]` used
- [ ] Works across different in/out types

**Checkpoint A (all tests)**

```bash
go test ./...
# expected: PASS (ok day4)
```

> **Micro-questions:**
> 1. What is a **type constraint** (e.g. `comparable`, `any`)?
> 2. When are generics worth it vs. just using `any` + assertions?
> 3. What's the zero value of a `map` vs a `slice`, and which is safe to read from?

### A4. Stretch (optional)

Add `Keys[K comparable, V any](m map[K]V) []K` returning the map's keys.

- [ ] Generic over both key and value types

---

## PART B — HTMX Triggers, Indicators & Swap Patterns (~45 min)

### B1. `hx-trigger` on keyup (live search)

**Task:** In `day4/htmx/server.go`, serve a search box that filters **fruit** `/search?q=` on `keyup changed delay:300ms` and swaps results into `#results`.

Starter:

```go
package main

import (
    "fmt"
    "net/http"
    "strings"
)

var fruits = []string{"apple", "apricot", "banana", "cherry", "grape", "mango"}

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, page)
    })

    http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
        q := strings.ToLower(r.URL.Query().Get("q"))
        for _, f := range fruits {
            if q == "" || strings.Contains(f, q) {
                fmt.Fprintf(w, "<li>%s</li>", f)
            }
        }
    })

    http.ListenAndServe(":8080", nil)
}
```

```html
const page = `<!doctype html>
<html><head><script src="https://unpkg.com/htmx.org@2"></script></head>
<body>
    <input type="search" name="q" placeholder="Search fruit..."
           hx-get="/search" hx-trigger="keyup changed delay:300ms"
           hx-target="#results">
    <span class="htmx-indicator">Searching...</span>
    <ul id="results"></ul>
</body></html>`
```

**Checkpoint B1:** type `"ap"`.

- [ ] Results filter live (apple, apricot appear)
- [ ] Fires on `keyup changed delay:300ms` (debounced)

### B2. Loading indicator

**Task:** Add a `class="htmx-indicator"` element (done above) and confirm it appears during the request. Add an artificial `time.Sleep(400 * time.Millisecond)` in the handler to see it.

**Checkpoint B2:**

- [ ] Indicator shows while the request is in flight, hides after
- [ ] You understand `htmx-indicator` + `hx-indicator`

> **Micro-questions:**
> 1. What events can `hx-trigger` listen to besides `keyup`?
> 2. What do the modifiers `changed`, `delay:300ms`, `once` do?
> 3. Name three `hx-swap` values and what each does.

### B3. Stretch (optional)

Add `hx-trigger="load"` on an element so it fetches content **as soon as the page loads**.

- [ ] Content auto-loads on page render

---

## PART C — MongoDB Indexes & `explain()` (~30 min)

### C1. Seed a bigger collection

```js
use day4
for (let i = 0; i < 1000; i++) {
  db.users.insertOne({
    name: "user" + i,
    age: 18 + (i % 50),
    city: ["NYC","LA","SF","CHI"][i % 4],
    active: i % 2 === 0
  })
}
```

```js
db.users.countDocuments()   // expected: 1000
```

- [ ] 1000 docs inserted

### C2. Baseline `explain()` (no index)

```js
db.users.find({ city: "SF", age: { $gte: 30 } }).explain("executionStats")
```

Note the `totalDocsExamined` and `stage` (`COLLSCAN`).

- [ ] You recorded `COLLSCAN` and docs examined ≈ 1000

### C3. Create indexes (the tested core)

1. Single-field index on `city`.
2. **Compound** index on `{ city: 1, age: 1 }`.

```js
db.users.createIndex({ city: 1 })
db.users.createIndex({ city: 1, age: 1 })
db.users.getIndexes()
```

Re-run the `explain()` from C2.

**Checkpoint C3**

- [ ] Stage is now `IXSCAN` (not `COLLSCAN`)
- [ ] `totalDocsExamined` dropped dramatically
- [ ] `getIndexes()` shows your two new indexes

### C4. Understanding (written)

1. What is the **ESR rule** (Equality, Sort, Range) for compound index field order?
2. Why does `{ city: 1, age: 1 }` help `find({city, age})` but `{ age: 1, city: 1 }` might not for a `city`-only query?
3. What's the cost of adding indexes (writes, storage)?

- [ ] All 3 answered

---

## PART D — English (~30–45 min)

### D1. Vocabulary

One sentence each: `adaptable`, `organized`, `analytical`, `proactive`, `reliable`.

- [ ] 5 original sentences

### D2. Grammar — comparatives & superlatives

Fill in:

1. Go compiles _____ (fast) than Java. → *faster*
2. This is the _____ (efficient) service we run.
3. MongoDB is _____ (flexible) than a rigid SQL schema.
4. This is the _____ (bad) bug I have seen this year.
5. HTMX is _____ (simple) than a full SPA framework.

- [ ] Attempt all 5

### D3. Speaking — strengths & weaknesses (record yourself)

Speak **60–90 seconds**: give **2 strengths** (with examples) and **1 weakness** + how you're improving it.

**Checkpoint D3**

- [ ] Recorded once, listened back
- [ ] Each strength backed by an example
- [ ] Weakness includes a concrete improvement plan (no fake weakness)

---

## ✅ Day 4 Done — Definition of Complete

- [ ] Go: `go test ./...` green (A1–A3, optional A4)
- [ ] HTMX: live search (B1) + loading indicator (B2)
- [ ] MongoDB: `explain()` shows `IXSCAN` after indexing (C2–C3)
- [ ] English: vocab + grammar attempt + strengths/weaknesses recording
- [ ] `answers/day4-notes.md` complete

Score yourself: **__/5 sections green.**

---

## 🔑 Answer Key (try first, then check)

<details>
<summary>Go answers</summary>

```go
// collections.go
package main

func Filter(nums []int) []int {
    out := make([]int, 0, len(nums))
    for _, n := range nums {
        if n%2 == 0 {
            out = append(out, n)
        }
    }
    return out
}

func SumRange(nums []int) int {
    sum := 0
    for _, n := range nums {
        sum += n
    }
    return sum
}

func WordCount(words []string) map[string]int {
    m := make(map[string]int)
    for _, w := range words {
        m[w]++
    }
    return m
}

func MapSlice[T any, U any](in []T, fn func(T) U) []U {
    out := make([]U, len(in))
    for i, v := range in {
        out[i] = fn(v)
    }
    return out
}

func Keys[K comparable, V any](m map[K]V) []K {
    out := make([]K, 0, len(m))
    for k := range m {
        out = append(out, k)
    }
    return out
}
```

- **A1.1** Length = current elements; capacity = allocated slots before regrowth.
- **A1.2** `append` allocates a **new** larger array and copies — old references don't see new elements.
- **A2 trap** Map iteration order is **randomized** by design; never rely on it.
- **A3.1** A constraint limits what types are allowed (`comparable` = usable with `==`; `any` = anything).
- **A3.3** Zero value of a map is `nil` (writing panics, reading is safe/returns zero); nil slice is safe to `append` and `range`.

</details>

<details>
<summary>HTMX answers</summary>

- **B.1** `hx-trigger` can listen to `keyup`, `change`, `submit`, `load`, `mouseenter`, custom events, `every 2s` (polling), etc.
- **B.2** `changed` = only fire if value changed; `delay:300ms` = debounce; `once` = fire a single time.
- **B.3** `innerHTML` (replace content), `outerHTML` (replace element), `beforeend`/`afterbegin` (append/prepend), `delete`, `none`.

</details>

<details>
<summary>MongoDB answers</summary>

- **C4.1** **ESR**: put **Equality** fields first, then **Sort** fields, then **Range** fields in a compound index.
- **C4.2** A compound index is usable **left-to-right** (prefix rule); `{city, age}` covers `city`-only and `city+age`, but `{age, city}` can't serve a `city`-only query efficiently.
- **C4.3** Indexes speed reads but **slow writes** (each write updates every relevant index) and consume storage/RAM.

</details>

<details>
<summary>English answers</summary>

**D2 grammar:** 1. **faster** · 2. **most efficient** · 3. **more flexible** · 4. **worst** · 5. **simpler**

**Rules:**
- Short words: `-er` / `-est` (fast → faster → fastest).
- Long words: `more` / `most` (efficient → more → most efficient).
- Irregular: good → better → best, bad → worse → worst.

</details>

---

## Tomorrow (Day 5 preview)

Goroutines, channels, `sync`, `context` · HTMX interactive patterns (search + refresh) · MongoDB schema design (embedding vs referencing) · English: giving a clear opinion + modal verbs. Bring your Day 4 red items.

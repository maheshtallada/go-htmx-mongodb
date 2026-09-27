# Day 8 Assignment — Go, HTMX & MongoDB

## PART A — Go Pagination + MongoDB skip/limit

### A1. Offset pagination (`skip/limit`)

**Task:** In `paginate.go`, a pure helper you can unit-test without a DB!

```go
// Page returns the slice for a 1-based page and pageSize.
// Out-of-range pages return an empty slice (never panic).
func Page[T any](items []T, page, pageSize int) []T
```

**Checkpoint A1** — `paginate_test.go`:

```go
package main

import (
    "reflect"
    "testing"
)

func TestPage(t *testing.T) {
    items := []int{1, 2, 3, 4, 5, 6, 7}

    if got := Page(items, 1, 3); !reflect.DeepEqual(got, []int{1, 2, 3}) {
        t.Fatalf("page1 = %v", got)
    }

    if got := Page(items, 3, 3); !reflect.DeepEqual(got, []int{7}) {
        t.Fatalf("page3 = %v", got)
    }

    if got := Page(items, 99, 3); len(got) != 0 {
        t.Fatalf("out-of-range should be empty, got %v", got)
    }
}
```

- [ ] Correct slice per page
- [ ] Out-of-range returns empty, no panic
- [ ] Generic over `T`

> **Micro-questions:**
> 1. How do `skip` / `limit` map to `page` / `pageSize`?
> 2. Why does deep offset pagination (`skip` large) get slow in MongoDB?

### A2. Metadata (the tested core)

**Task:** Add:

```go
type PageInfo struct {
    Page       int
    PageSize   int
    TotalItems int
    TotalPages int
    HasNext    bool
}

// Paginate computes PageInfo for the given total.
func Paginate(page, pageSize, total int) PageInfo
```

**Checkpoint A2:**

```go
func TestPaginate(t *testing.T) {
    p := Paginate(2, 10, 25)
    if p.TotalPages != 3 || !p.HasNext {
        t.Fatalf("got %+v", p)
    }

    last := Paginate(3, 10, 25)
    if last.HasNext {
        t.Fatalf("last page should have no next")
    }
}
```

- [ ] `TotalPages` uses ceiling division
- [ ] `HasNext` correct on middle and last pages

### A3. Driver query with `skip/limit` (integration)

**Task:** Using the Mongo driver, fetch page 2 (size 5) of a `products` collection sorted by `name`, via:

```go
options.Find().SetSkip(5).SetLimit(5).SetSort(...)
```

- [ ] Returns at most 5 docs, correct offset
- [ ] Sort applied server-side

**Checkpoint A (unit tests)**

```bash
go test ./...
```

> **Micro-questions:**
> 1. What's the difference between **offset** and **cursor (range) pagination**?
> 2. Why is cursor pagination better for large datasets / infinite scroll?
> 3. What index supports `sort({name:1})` efficiently?

---

## PART B — HTMX Inline Field Validation (~45 min)

### B1. Field-level validation on blur

**Task:** In `day8/htmx/server.go`, a signup form where the **username** field validates on `blur` via `hx-post="/validate/username"` into a per-field error slot.

Starter idea:

```html
<input name="username" hx-post="/validate/username"
       hx-trigger="blur" hx-target="#username-error" hx-swap="innerHTML">
<span id="username-error"></span>
```

**Server rule:** username must be ≥ 3 chars and alphanumeric.

**Checkpoint B1:**

- [ ] Typing a bad username shows an error when you tab away
- [ ] Error is scoped to that field's slot

### B2. Whole-form submit with errors preserved

**Task:** On `POST /signup`, if any field is invalid, re-render the form fragment with **all** error messages and a `422`. On success, `HX-Redirect` to `/done`.

**Checkpoint B2:**

- [ ] Invalid submit shows all field errors, keeps entered values
- [ ] Valid submit redirects

> **Micro-questions:**
> 1. Why validate both on **blur** (per field) and on submit (whole form)?
> 2. How do you keep the user's typed values when re-rendering after an error?
> 3. What status code signals validation failure?

### B3. Stretch (optional)

Disable the submit button until all fields are valid using an out-of-band swap.

- [ ] Submit enables only when valid

---

## PART C — MongoDB Validation, Upserts, Write Concern (~30 min)

### C1. Schema validation rules

**Task:** Create a collection with a JSON-schema validator:

```js
use day8

db.createCollection("members", {
  validator: { $jsonSchema: {
    bsonType: "object",
    required: ["email", "age"],
    properties: {
      email: { bsonType: "string", pattern: "@" },
      age: { bsonType: "int", minimum: 18 }
    }
  }}
})
```

Try an invalid insert (missing email / age < 18) and confirm it's rejected.

**Checkpoint C1**

- [ ] Valid doc inserts
- [ ] Invalid doc is **rejected** by the validator

### C2. Upsert (the tested core)

**Task:** Upsert a member by email — Insert if absent, update if present!

```js
db.members.updateOne(
  { email: "a@x.com" },
  { $set: { age: 30 } },
  { upsert: true }
)
```

Run it twice; confirm only **one** doc exists.

**Checkpoint C2**

- [ ] First run inserts, second run updates (no duplicate)
- [ ] `matchedCount` / `upsertedId` behave as expected

### C3. Write concern (written)

1. What does write concern `w: "majority"` guarantee?
2. What does `w: 1` vs `w: 0` mean?
3. Trade-off: stronger write concern vs latency?

- [ ] All 3 answered

---

## PART D — English (~30–45 min)

### D1. Vocabulary

One sentence each: `disagreement`, `alignment`, `compromise`, `resolution`, `support`.

- [ ] 5 original sentences

### D2. Grammar — reported speech

Convert to reported speech:

1. He said, "I will fix the bug today."
2. She said, "We are deploying now."
3. They said, "We have finished testing."
4. My manager said, "You did a great job."
5. The client said, "We need it by Friday."

- [ ] Attempt all 5

### D3. Speaking — conflict & resolution (record yourself)

Speak **90 seconds** answering: **"Tell me about a disagreement with a teammate and how you resolved it."** Use STAR + at least one reported-speech sentence.

**Checkpoint D3**

- [ ] Recorded, listened back
- [ ] Stayed professional (focus on resolution, not blame)
- [ ] Used at least one reported-speech sentence

---

## ✅ Day 8 Done — Definition of Complete

- [ ] Go: `go test ./...` green (A1–A2) + driver skip/limit (A3)
- [ ] HTMX: per-field validation + form-level errors + redirect (B2)
- [ ] MongoDB: validator (C1) + upsert (C2)
- [ ] English: vocab + reported speech + conflict recording
- [ ] `answers/day8-notes.md` complete

Score yourself: **__/5 sections green.**

---

# 🔑 Answer Key (try first, then check)

<details>
<summary>Go answers</summary>

```go
// paginate.go
package main

func Page[T any](items []T, page, pageSize int) []T {
    if page < 1 || pageSize < 1 {
        return []T{}
    }

    start := (page - 1) * pageSize
    if start >= len(items) {
        return []T{}
    }

    end := start + pageSize
    if end > len(items) {
        end = len(items)
    }

    return items[start:end]
}

func Paginate(page, pageSize, total int) PageInfo {
    totalPages := 0
    if pageSize > 0 {
        totalPages = (total + pageSize - 1) / pageSize // ceil
    }

    return PageInfo{
        Page:       page,
        PageSize:   pageSize,
        TotalItems: total,
        TotalPages: totalPages,
        HasNext:    page < totalPages,
    }
}
```

- **A1.1** `skip = (page-1)*pageSize`, `limit = pageSize`.
- **A1.2** Large `skip` still scans/discards skipped docs — O(skip) work.
- **A3.1** Offset uses `skip/limit`; cursor uses a `_id`/field > `lastSeen` filter.
- **A3.2** Cursor pagination is O(pageSize), stable under inserts, and ideal for infinite scroll.

</details>

<details>
<summary>HTMX answers</summary>

- **B1.1** Per-field `blur` gives fast feedback; submit validation is the authoritative gate (never trust client).
- **B2.1** Re-render the form with the submitted values bound back into the inputs (`value="{{.Username}}"`).
- **B3.1** `422 Unprocessable Entity`.

</details>

<details>
<summary>MongoDB answers</summary>

- **C3.1** `w: "majority"` acknowledges only after a majority of replica-set members persist the write — survives failover.
- **C3.2** `w:1` = acknowledged by primary only; `w:0` = fire-and-forget (no ack).
- **C3.3** Stronger concern = more durability but higher latency (waits for more nodes).

</details>

<details>
<summary>English answers</summary>

**D2 reported speech:**

1. He said (that) he **would** fix the bug **that day**.
2. She said (that) they **were deploying** then.
3. They said (that) they **had finished** testing.
4. My manager said (that) I **had done** a great job.
5. The client said (that) they **needed** it by Friday.

**Rule:** shift tense back one step (`will`→`would`, present→past, present perfect→past perfect) and adjust time words (`today`→`that day`).

</details>

---

## Tomorrow (Day 9 preview)

Testing with `testing` + `httptest`, testing HTMX endpoints & HTML fragments, testing MongoDB code with fixtures/mocks → cleanup. English: Leadership & ownership + active vs passive voice.

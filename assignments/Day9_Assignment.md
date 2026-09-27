# Day 9 Assignment — Test-Driven Learning (Go · HTMX · MongoDB · English)

> **How this works:** Make each **failing test or checkpoint** at the bottom — try first. Answer keys are below.
>
> **Time-box:** ~3 hours + 30–45m English.

## Day 9 scope

- **Go:** testing (`testing`, table tests, `httptest`), coverage
- **HTMX:** testing HTMX endpoints and HTML fragments
- **MongoDB:** testing DB code with fixtures/mocks + cleanup
- **English:** leadership & ownership, active vs passive voice, 5 vocab words

> **Warm-up:** finish any red items from Day 8 first.

---

## 0. Setup (5 min)

```bash
docker start mongo-day1 || docker run -d --name mongo-day9 -p 27017:27017 mongo:7
mkdir -p day9 && cd day9
go mod init day9
```

---

# PART A — Go: Testing Mastery (~60 min)

### A1. Table-driven tests

**Task:** In `slugify.go`:

```go
// Slugify lowercases and replaces spaces with hyphens; trims surrounding spaces.
func Slugify(s string) string
```

**Checkpoint A1** — write a **table-driven** test:

```go
package main

import "testing"

func TestSlugify(t *testing.T) {
    cases := []struct{ in, want string }{
        {"Hello World", "hello-world"},
        {"  Go HTMX  ", "go-htmx"},
        {"Already-Slug", "already-slug"},
        {"", ""},
    }

    for _, c := range cases {
        if got := Slugify(c.in); got != c.want {
            t.Errorf("Slugify(%q) = %q; want %q", c.in, got, c.want)
        }
    }
}
```

- [ ] All rows pass
- [ ] You used `t.Errorf` (continue) not `t.Fatalf` (stop inside the loop)

> **Micro-questions:**
> 1. Why are table-driven tests idiomatic in Go?
> 2. `t.Run(name, ...)` — what do subtests give you?

### A2. `httptest` for a handler (the tested core)

**Task:** A handler `TaskHandler` returning JSON `{"id":1,"title":"Ship"}`. Test it with `httptest`.

**Checkpoint A2:**

```go
func TestTaskHandler(t *testing.T) {
    req := httptest.NewRequest("GET", "/task", nil)
    rec := httptest.NewRecorder()
    TaskHandler(rec, req)

    if rec.Code != 200 {
        t.Fatalf("code = %d", rec.Code)
    }

    if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
        t.Fatalf("content-type = %q", ct)
    }
}
```

- [ ] Status + content-type asserted
- [ ] Body is valid JSON

### A3. Test helpers + coverage

**Task:** Add a `t.Helper()`-based assertion helper and run coverage.

```bash
go test -cover ./...
go test -coverprofile=cov.out ./... && go tool cover -func=cov.out | tail -1
```

- [ ] Coverage reported
- [ ] Helper marked with `t.Helper()`

> **Micro-questions:**
> 1. What does `t.Helper()` change in failure output?
> 2. What's the difference between `go test`, `-run`, and `-v`?
> 3. What is a **golden file** test?

---

# PART B — Testing HTMX Endpoints (~45 min)

### B1. Assert a fragment response

**Task:** Given a `FragmentHandler` returning `<li id="eng-1">Mahesh</li>`, write an `httptest` test asserting the body **contains** the expected element and the correct id.

**Checkpoint B1:**
- [ ] Test checks status 200
- [ ] Test asserts the fragment HTML (e.g. `strings.Contains`)

### B2. Assert HX behavior via headers

**Task:** A handler that sets `HX-Redirect: /done` on success. Test that the header is present.

**Checkpoint B2:**
- [ ] Test asserts `HX-Redirect` header value
- [ ] You understand how to test HTMX responses without a browser

> **Micro-questions:**
> 1. How do you test HTMX behavior server-side without a real browser?
> 2. What should a fragment test assert (structure vs exact string)?
> 3. When would you reach for a browser-based test (Playwright) instead?

---

# PART C — Testing MongoDB Code (~30 min)

### C1. Fixture setup + teardown

**Task:** Write a Go test that, against a **test database** (`day9_test`), inserts fixtures in setup and **drops the collection** in cleanup with `t.Cleanup(...)`.

```go
func setupColl(t *testing.T) *mongo.Collection {
    // connect, pick day9_test db, unique collection name
    col := client.Database("day9_test").Collection("tasks_" + t.Name())
    t.Cleanup(func() { col.Drop(context.Background()) })
    return col
}
```

**Checkpoint C1:**
- [ ] Each test gets an isolated collection
- [ ] `t.Cleanup` drops it afterward (no leftover data)

### C2. Test a repository method (the tested core)

**Task:** Insert 3 tasks, call `FindAll`, assert count == 3.

- [ ] Repo method verified against real Mongo
- [ ] Test is repeatable (clean each run)

### C3. Mock vs real (written)

1. When is an **in-memory/mocked** repo better than hitting real Mongo?
2. What's the risk of mocking the driver (false confidence)?
3. What is a **testcontainers** approach and why is it popular?

- [ ] All 3 answered

---

# PART D — English (~30–45 min)

### D1. Vocabulary

One sentence each: `accountability`, `initiative`, `guidance`, `coordination`, `influence`

- [ ] 5 original sentences

### D2. Grammar — active vs passive voice

Convert active → passive (or note when passive is better):

1. The team deployed the service. → *(passive)*
2. We fixed the bug. → *(passive)*
3. A junior engineer wrote this module. → *(passive)*
4. Mongo stores the documents. → *(passive)*
5. Rewrite one passive back to a **stronger active** sentence.

- [ ] Attempt all 5

### D3. Speaking — leadership & ownership (record yourself)

Speak **90 seconds** answering **"Tell me about a time you led or took ownership."** STAR format. Emphasize **your** actions (active voice).

**Checkpoint D3**
- [ ] Recorded, listened back
- [ ] Used mostly **active voice** ("I led", "I decided")
- [ ] Clear ownership + measurable result

---

## ✅ Day 9 Done — Definition of Complete

- [ ] Go: table tests + `httptest` green, coverage reported (A1–A3)
- [ ] HTMX: fragment + header tests (B1–B2)
- [ ] MongoDB: fixtures + cleanup + repo test (C1–C2)
- [ ] English: vocab + active/passive + leadership recording
- [ ] `answers/day9-notes.md` complete

**Score yourself: __/5 sections green.** 🎉 **Week 2 complete!**

---

# 🔑 Answer Key (try first, then check)

<details>
<summary>Go answers</summary>

```go
// slugify.go
package main

import "strings"

func Slugify(s string) string {
    s = strings.TrimSpace(s)
    s = strings.ToLower(s)
    return strings.ReplaceAll(s, " ", "-")
}
```

```go
// task handler
func TaskHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(`{"id":1,"title":"Ship"}`))
}
```

- **A1.1** Table tests add cases with zero boilerplate and give one clear failure per case.
- **A1.2** `t.Run` creates named subtests you can filter with `-run` and report independently.
- **A3.1** `t.Helper()` makes failures point at the **caller's** line, not the helper's.
- **A3.3** A golden file stores expected output; the test compares actual output to it (great for HTML/JSON).

</details>

<details>
<summary>HTMX answers</summary>

```go
func TestFragment(t *testing.T) {
    rec := httptest.NewRecorder()
    FragmentHandler(rec, httptest.NewRequest("GET", "/frag", nil))
    if !strings.Contains(rec.Body.String(), `id="eng-1"`) {
        t.Fatalf("body = %q", rec.Body.String())
    }
}
```

- **B1.1** Call the handler with `httptest` and assert on the recorded body/headers — no browser needed.
- **B2.2** Assert on **structure** (ids, presence of elements) rather than brittle exact strings.
- **B.3** Use a browser test when you need real JS execution / actual swap behavior end-to-end.

</details>

<details>
<summary>MongoDB answers</summary>

- **C3.1** Mock/in-memory when you want fast, isolated unit tests of logic that doesn't depend on real query behavior.
- **C3.2** Mocking the driver can hide real query/serialization bugs — green tests, broken production.
- **C3.3** **testcontainers** spins up a real MongoDB in Docker per test run — realistic + isolated + CI-friendly.

</details>

<details>
<summary>English answers</summary>

**D2 passive:**

1. The service **was deployed** by the team.
2. The bug **was fixed**.
3. This module **was written** by a junior engineer.
4. The documents **are stored** by Mongo.
5. Stronger active example: **"I fixed the bug."**

**Rule:** Prefer **active voice** in interviews to show ownership; use passive only when the actor is unknown or unimportant.

</details>

---

## Tomorrow (Day 10 preview)

Build, Docker, env config, graceful shutdown · HTMX page shell/layout + reusable fragments · MongoDB Atlas/local Docker, connection strings, backups · English: failure & learning + conditionals. **Week 3 begins next.**

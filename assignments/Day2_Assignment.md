# Day 2 Assignment — Test-Driven Learning (Go · HTMX · MongoDB · English)

## How This Works

Same as Day 1. Each task tells you **what to build/answer**, gives you a failing test or checkpoint, and you make it pass by learning just enough. Look things up as you go. Answer keys are at the bottom (no peeking until you try).

## Time-box: ~3 hours + 30-45m English

Aim for green checkmarks, not perfection.

---

## Day 2 Scope

- **Go**: structs, methods, interfaces (implicit satisfaction), pointers vs values
- **HTMX**: server-rendered HTML with `html/template`, partial updates, out-of-band swaps
- **MongoDB**: CRUD from the shell — insert, find, update, delete
- **English**: explain a project end-to-end, past simple vs present perfect, 5 vocab words

## Warm-up: Finish any red items from Day 1 first (5–10 min).

---

## 0. Setup (5 min) — Prove Day 1 Still Works

### Checkpoint 0.1 — Reuse yesterday's project or make a fresh one

Paste outputs into `answers/day2-notes.md`.

```bash
go version          # still there
docker start mongo-day1    # restart yesterday's container (or run a new one)
# docker run -d --name mongo-day2 -p 27017:27017 mongo:7
```

**Checkboxes:**
- [ ] Go builds and runs
- [ ] MongoDB is reachable on `localhost:27017`
- [ ] Compass connected

**Create today's project:**

```bash
mkdir -p day2/go && cd day2/go
go mod init day2
```

---

## PART A — Go: Structs, Methods, Interfaces (~60 min)

### A1. Structs & constructors

**Task:** In `engineer.go`, define a struct and a constructor:

```go
// Engineer models a team member.
type Engineer struct {
    Name   string
    Role   string
    Years  int
    Skills []string
}

// NewEngineer builds an Engineer with the given fields.
func NewEngineer(name, role string, years int, skills ...string) Engineer
```

**Checkpoint A1** — Write `engineer_test.go` and make it pass:

```go
package main

import "testing"

func TestNewEngineer(t *testing.T) {
    e := NewEngineer("Asha", "SDE II", 5, "Go", "MongoDB")
    if e.Name != "Asha" || e.Years != 5 {
        t.Fatalf("bad fields %q", e)
    }
    if len(e.Skills) != 2 {
        t.Fatalf("skills = %v; want 2", e.Skills)
    }
}
```

**Checkboxes:**
- [ ] Struct compiles and test passes
- [ ] You used a variadic `skills ...string`

**Micro-questions:**
1. What is the zero value of an `Engineer` struct?
2. What does `%v` print vs `%q`?

---

### A2. Methods (value vs pointer receivers)

**Task:** Add two methods to `Engineer`:

```go
// Summary returns "<Name> (<Role>), <Years>y".
func (e Engineer) Summary() string

// Promote increases Years by 1 and updates Role. It must MUTATE the engineer.
func (e *Engineer) Promote(newRole string)
```

**Checkpoint A2** — Extend the test:

```go
func TestMethods(t *testing.T) {
    e := NewEngineer("Asha", "SDE II", 5, "Go")
    if e.Summary() != "Asha (SDE II), 5y" {
        t.Fatalf("summary = %q", e.Summary())
    }
    e.Promote("SDE III")
    if e.Role != "SDE III" || e.Years != 6 {
        t.Fatalf("after promote: %v", e)
    }
}
```

**Checkboxes:**
- [ ] `Summary()` uses a value receiver
- [ ] `Promote()` uses a pointer receiver and the change sticks

**Trap check:** If `Promote` had a value receiver (e `Engineer`), would the mutation persist? Write down why.

---

### A3. Interfaces (implicit satisfaction — the tested core)

**Task:** Define an interface and make `Engineer` satisfy it without declaring `implements`:

```go
// Describer is anything that can describe itself.
type Describer interface {
    Describe() string
}
```

Add a `Describe()` method to `Engineer` that returns `"I am <Name>, a <Role>."` Then write a free function:

```go
// PrintAll returns each item's description joined by newlines.
func PrintAll(items []Describer) string
```

**Checkpoint A3:**

```go
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
```

**Checkboxes:**
- [ ] `Engineer` satisfies `Describer` with no explicit declaration
- [ ] `PrintAll` works over the interface, not the concrete type

**Checkpoint A (all Go tests):**

```bash
go test ./...
# expected: PASS (ok day2)
```

**Micro-questions:**
1. How does Go decide a type satisfies an interface? (compare to Java `implements`)
2. What is an empty interface `interface{}`? What can it hold? Use with type assertions/switches.
3. Why / `interface{}` holds any value; use with type assertions/switches?
4. When embedding: type `Manager struct { Engineer; Reports int }` — does a `Manager` value also use `Describe()`? And can an `Engineer`'s method set change? Give the syntax.

---

### A4. Stretch (optional, 10 min)

Add a `Manager` struct that embeds `Engineer` and add `Reports int`. Prove that a `Manager` value can also be used as a `Describer`.

**Checkbox:**
- [ ] Embedding gives `Manager` the `Describe()` method for free

---

## PART B — HTMX with Server-Rendered Templates (~45 min)

Yesterday you returned raw strings. Today use Go's `html/template` and do partial updates.

### B1. Template + list rendering

**Task:** In `day2/htmx/`, create `server.go`. Keep an in-memory slice of engineers and render them with a template.

**Minimal starter:**

```go
package main

import (
    "html/template"
    "net/http"
)

type Engineer struct {
    ID    int
    Name  string
    Role  string
}

var engineers = []Engineer{
    {1, "Mahesh", "Senior Lead"},
    {2, "Asha", "SDE II"},
}
var nextID = 3

var tmpl = template.Must(template.New("page").Parse(pageHTML))

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        tmpl.ExecuteTemplate(w, "page", engineers)
    })
    // TODO B2: POST /engineers → append + return the new <li> fragment
    // TODO B3: DELETE /engineers/{id} → remove + return empty (200 OK + outerHTML swap removes the row)
    http.ListenAndServe(":8080", nil)
}

const pageHTML = `<!DOCTYPE html>
<html>
<head><script src="https://unpkg.com/htmx.org@2"></script></head>
<body>
<h1>Engineers</h1>
<ul id="list">
{{range .}}
  <li id="eng-{{.ID}}">
    {{.Name}} - {{.Role}}
    <!-- TODO B1: add form that requests in new engineer -->
  </li>
{{end}}
</ul>
</body>
</html>`
```

**Checkpoint B1 (manual):** Open http://localhost:8080 — you see the needed list rendered by the template.

**Checkboxes:**
- [ ] List renders from the Go slice via `{{range .}}`

---

### B2. hx-post that appends a row (partial update)

**Task:** Add a form (`name`, `role`) that `hx-post` s to `/engineers`. The handler appends to the slice and returns only the new `<li>` fragment. Use `hx-target="#list"` and `hx-swap="beforeend"` so the new row added without a reload.

**Checkpoint B2:** Submit a new engineer.

**Checkboxes:**
- [ ] A new `<li>` appears at the bottom of `#list`, no full page reload
- [ ] The server returns just the fragment, not the whole page
- [ ] You used `hx-swap="beforeend"`

---

### B3. hx-delete that removes a row

**Task:** Give each `<li>` a delete button that `hx-delete` s to `/engineers/{id}`. The handler removes it and returns empty (200 OK + `outerHTML` swap removes the row itself — great for counters, flash messages, silos).

**Checkpoint B3:** Delete a row.

**Checkboxes:**
- [ ] The row disappears without a reload
- [ ] You used `hx-target="closest li"` (or `hx-swap="outerHTML"` — return an empty body).

**Micro-questions:**
1. Why is `html/template` safer than `fmt.Fprintf` for HTML?
2. What HTTP verbs can HTMX issue? PUT / PATCH, DELETE (plain `<form>` cannot?).
3. What is an out-of-band swap ( `hx-swap-oob` ) and when is it useful?

---

### B4. Stretch (optional)

Return an out-of-band fragment that also updates a header count `Total: N` whenever you add/delete.

**Checkbox:**
- [ ] Count updates via `hx-swap-oob="true"`

---

## PART C — MongoDB CRUD from the Shell (~30 min)

Use mongosh (and watch changes live in Compass).

### C1. Create (insert)

**Task:** In database `day2`, collection `tasks`, insert these 3 documents:

```javascript
use day2
db.tasks.insertMany([
  { title: "Write API", status: "todo", points: 5, tags: ["Go", "http"] },
  { title: "Design UI", status: "doing", points: 3, tags: ["UI"] },
  { title: "Model data", status: "todo", points: 8, tags: ["mongo", "schema"] }
])
```

**Checkpoint C1:**

```javascript
db.tasks.countDocuments()  // expected: 3
```

**Checkboxes:**
- [ ] 3 documents inserted, visible in Compass

---

### C2. Read (find + projection)

**Task:** Write queries to answer:

1. All status `"todo"` tasks, show only `title` and `points` (hide `_id`).
2. One task by title `"Design UI"` using `findOne`.

**Starter for Q1:**

```javascript
db.tasks.find({ status: "todo" }, { _id: 0, title: 1, points: 1 })
```

**Checkboxes:**
- [ ] Q1 returns "Write API" and "Model data" only, with just `title + points`
- [ ] `findOne` returns exactly one document

---

### C3. Update (the tested core)

1. Set `"Design UI"` to status: `"done"` with `updateOne` + `$set`.
2. Increment `points` by 1 for all `todo` tasks with `updateMany` + `$inc`.
3. Add tag `"urgent"` to `"Write API"` using `$push` (or `$addToSet`).

**Starter:**

```javascript
db.tasks.updateOne({ title: "Design UI" }, { $set: { status: "done" } })
```

**Checkpoint C3:**

```javascript
db.tasks.find({ status: "done" }).count()     // expected: 1
db.tasks.findOne({ title: "Write API" }).tags // includes "urgent"
```

**Checkboxes:**
- [ ] `updateOne` changed exactly one doc
- [ ] `$inc` bumped points on all matches
- [ ] `$set`, `$inc`, `$push` — describe what each does in one line each

---

### C4. Delete

1. Delete the `"Model data"` task with `deleteOne`.
2. Confirm the count is now 2.

**Checkboxes:**
- [ ] `deleteOne` removed one document
- [ ] `countDocuments()` returns 2

---

### C5. Understanding (written)

Answer in notes:

1. Difference between `updateOne` and `updateMany`.
2. What does an upsert do? Write the option flag.
3. Difference between `$set`, `$inc`, `$push`, `$addToSet` in one line each.

**Checkbox:**
- [ ] All 3 answered in your own words

---

## PART D — English (~30-45 min)

### D1. Vocabulary (use, don't memorize)

Write one sentence for each word, about your real work: `architecture`, `scalable`, `deadline`, `stakeholder`, `deliver`.

**Checkpoint:**
- [ ] 5 original sentences, each using the word correctly

---

### D2. Grammar — past simple vs present perfect

Fill the correct form, then check against the key:

1. Last year, I _____ (migrate) our monolith to microservices.
2. I _____ (work) with MongoDB since 2021.
3. Yesterday, the team _____ (deploy) the new release.
4. I _____ (never / use) Kafka before this project.
5. In 2020, we _____ (build) our first Go service.

**Checkpoint:**
- [ ] Attempt all 5 before checking answers

---

### D3. Speaking — explain a project end-to-end (record yourself)

Speak 90–120 seconds answering **"Walk me through a project you led."** Cover:

- the problem / context (past simple)
- what you built + your role (past simple + present perfect)
- the outcome / impact (numbers if possible)
- what you would do differently (would + verb)

**Checkpoint D3:**

- [ ] Recorded once, listened back
- [ ] Clear beginning → end (no rambling)
- [ ] Used at least 2 past-simple and 1 present-perfect sentence
- [ ] Under 2 minutes

---

## ✅ Day 2 Done — Definition of Complete

**Completion checklist:**

- [ ] Go: `go test ./...` is green (A1–A3, optional A4)
- [ ] HTMX: list renders (B1), add appends a row (B2), delete removes a row (B3)
- [ ] MongoDB: CRUD all four work (C1–C4)
- [ ] English: vocab + grammar attempt + recorded project walkthrough
- [ ] `answers/day2-notes.md` has all micro-question answers

**Score yourself:** _/5 sections green. Anything red = tomorrow's warm-up.

---

## 🔑 Answer Key (try first, then check)

<details>
<summary>Go answers</summary>

### A1 Micro-questions:
1. The zero value of an `Engineer` struct is `Engineer{0, "", 0, nil}` — all fields at their zero value. (`"" ` for strings, `0` for ints, `nil` for slices).
2. `%v` prints any value in a default format. `%q` quotes strings and escapes special chars.

### A2 Trap check:
If `Promote` had a value receiver, the mutation would NOT persist because Go copies the struct when you call a value method. The method modifies the copy, not the original. You must use a pointer receiver (`*Engineer`) to mutate.

### A3 Micro-questions:
1. Go uses **structural subtyping** — if a type has all the methods of an interface, it satisfies it. No need to declare `implements`. Java requires explicit `implements`.
2. An empty interface `interface{}` holds any value. You use **type assertions** (`x.(string)`) or **type switches** (`switch x.(type)`) to extract the real type.
3. (duplicate of 2 above)
4. Embedding: `type Manager struct { Engineer; Reports int }` — a `Manager` value automatically inherits all methods of `Engineer`, including `Describe()`. The method set is promoted. No, an `Engineer`'s method set is fixed at declaration. You can't add methods to a type after it's defined. To use embedding: `func (m Manager) Describe() string { return m.Engineer.Describe() }` or just rely on promotion.

### A4 (Stretch):
Embedding `Engineer` in `Manager` means `Manager` gets `Describe()` for free without writing it again. A `Manager` value satisfies `Describer` just like `Engineer` does.

</details>

<details>
<summary>HTMX answers</summary>

### B1:
The template uses `{{range .}}` to iterate over the slice and render each engineer as a `<li>`.

### B2:
```html
<form hx-post="/engineers" hx-target="#list" hx-swap="beforeend" hx-on-after-request="this.reset()">
  <input name="name" placeholder="Name" required>
  <input name="role" placeholder="Role" required>
  <button type="submit">Add Engineer</button>
</form>
```

The handler parses the form, appends to the slice, and returns just the `<li>` fragment.

### B3:
```html
<button hx-delete="/engineers/{{.ID}}" hx-target="closest li" hx-swap="outerHTML swap:1s">Delete</button>
```

The handler removes from the slice and returns an empty body (200 OK). `outerHTML` swap removes the `<li>` itself.

### B4 Micro-questions:
1. `html/template` **auto-escapes** dangerous characters (`<`, `>`, `&`) to prevent XSS. `fmt.Fprintf` does no escaping.
2. HTMX can issue `GET`, `POST`, `PUT`, `PATCH`, `DELETE`. Plain `<form>` can only do `GET` and `POST`.
3. An **out-of-band swap** (`hx-swap-oob="true"`) lets the server return a fragment that targets a different element on the page. Example: add a task → server returns the new `<li>` AND a header `<span>Total: 4</span>` that replaces the old count in one response.

</details>

<details>
<summary>MongoDB answers</summary>

### C2 Queries:

1. All `"todo"` tasks, showing only `title` and `points`:
   ```javascript
   db.tasks.find({ status: "todo" }, { _id: 0, title: 1, points: 1 })
   ```
   Result: "Write API" (5), "Model data" (8)

2. One task by title:
   ```javascript
   db.tasks.findOne({ title: "Design UI" })
   ```

### C3 Updates:

1. Set "Design UI" to status "done":
   ```javascript
   db.tasks.updateOne({ title: "Design UI" }, { $set: { status: "done" } })
   ```

2. Increment points by 1 for all `"todo"` tasks:
   ```javascript
   db.tasks.updateMany({ status: "todo" }, { $inc: { points: 1 } })
   ```

3. Add tag "urgent" to "Write API":
   ```javascript
   db.tasks.updateOne({ title: "Write API" }, { $push: { tags: "urgent" } })
   ```
   (or `$addToSet` to avoid duplicates)

### C5 Understanding:

1. **Difference between `updateOne` and `updateMany`:** `updateOne` modifies the first matching document; `updateMany` modifies all matches.
2. **What is an upsert?** An upsert updates a document if it exists, or inserts it if it doesn't. Use `{ upsert: true }` as an option.
3. **Operators in one line each:**
   - `$set`: replace the field value.
   - `$inc`: add a number to the field.
   - `$push`: append a value to an array (allows duplicates).
   - `$addToSet`: append only if the value doesn't already exist in the array.

</details>

<details>
<summary>English answers</summary>

### D2 Grammar:

1. Last year, I **migrated** (past simple — finished action at a specific past time).
2. I **have worked** (present perfect — started in past, continues now).
3. Yesterday, the team **deployed** (past simple — yesterday is a specific point).
4. I **have never used** (present perfect — never in my life up to now).
5. In 2020, we **built** (past simple — 2020 is a specific past year).

**Rule of thumb:**
- **Past simple** — finished action at a specific past time (yesterday, last year, in 2020).
- **Present perfect** — unfinished time / experience up to now (since, for, ever, never, yet).

</details>

---

## 📚 Tomorrow (Day 3 preview)

**Errors as values (mapping, errors, nil)** · **HTMX form validation + nitwits** · **MongoDB query filters, projections, sorting, limits** · **English: handling "tell me about a challenge" answers**. Bring your Day 2 red items.

---

**Great work! 🚀**

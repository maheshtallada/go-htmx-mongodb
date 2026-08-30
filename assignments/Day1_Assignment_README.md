# Day 1 Assignment — Test-Driven Learning (Go · HTMX · MongoDB · English)

## How This Works

You haven't studied yet — that's the point. Each task tells you **what to build/answer**, gives you a failing test or checkpoint, and you make it pass by learning just enough. Answer keys are at the bottom (no peeking until you try).

## Time-box: ~3 hours + 30-45m English

Don't aim for perfect — aim for **green checkmarks**.

---

## Day 1 Scope

- **Go**: syntax, types, functions, packages, `go mod`
- **HTMX**: `hx-get`, `hx-post`, `hx-target`, `hx-swap`
- **MongoDB**: documents vs collections, BSON, Compass, mongosh
- **English**: self-introduction, present simple vs present perfect, 5 vocab words

---

## 0. Setup (10 min) — Prove Your Tools Work

### Checkpoint 0.1 — Install and verify versions

Paste outputs into `answers/day1-notes.md`.

```bash
go version          # expect go1.2x.x
mongod --version    # or use MongoDB Atlas (cloud) / Docker
```

### Run MongoDB via Docker (if you don't want a local install)

```bash
docker run -d --name mongo-day1 -p 27017:27017 mongo:7
```

**Checkboxes:**
- [ ] Go prints a version
- [ ] MongoDB is reachable on `localhost:27017` (or Atlas URI ready)
- [ ] MongoDB Compass installed and connected

---

## PART A — Go (~60 min)

### Create the project

```bash
mkdir -p day1/go && cd day1/go
go mod init day1
```

---

### A1. Module + package basics

**Task:** Create `main.go` in `package main` with a `main()` that prints `"Day 1 ready"`.

**Checkpoint A1** — Running the program prints exactly "Day 1 ready":

```bash
go run .
# expected output:
# Day 1 ready
```

**Checkboxes:**
- [ ] `go.mod` exists with module name `day1`
- [ ] `go run .` prints exactly "Day 1 ready"

**Micro-questions (answer in notes):**
1. What does `go mod init` create and why?
2. Why must an executable program use `package main`?
3. What is the difference between `go run .` and `go build`?

---

### A2. Types & variables

**Task:** In a new file `types.go`, declare one of each and print them with their types using `fmt.Printf("%T", x);`

- an `int`, a `float64`, a `string`, a `bool`
- one variable using `:=` and one using `var`
- one constant

**Checkpoint A2** — Running the program prints 6 lines, each showing a value and its type:

```bash
go run .
# expected: 6 lines, each with a type like:
# 42
# 3.14
# hello
# true
# ...
```

**Checkboxes:**
- [ ] You used both `:=` and `var`
- [ ] You used `const`
- [ ] You printed types with `%T`

**Trap check:** What is the zero value of `int`, `string`, and `bool`? Write them down.

---

### A3. Functions (the tested core)

**Task:** Create `mathx.go` with these three functions:

```go
// Add returns the sum of a and b.
func Add(a, b int) int

// DivMod returns quotient and remainder (multiple return values).
func DivMod(a, b int) (int, int)

// Greet returns "Hello, <name>!" or "Hello, world!" if name is empty.
func Greet(name string) string
```

**Now write the test yourself in `mathx_test.go`, then make it pass:**

```go
package main

import "testing"

func TestAdd(t *testing.T) {
    if Add(2, 3) != 5 {
        t.Fatalf("Add(2,3) = %d; want 5", Add(2, 3))
    }
}

func TestDivMod(t *testing.T) {
    q, r := DivMod(17, 5)
    if q != 3 || r != 2 {
        t.Fatalf("DivMod(17,5) = (%d,%d); want (3,2)", q, r)
    }
}

func TestGreet(t *testing.T) {
    if Greet("Mahesh") != "Hello, Mahesh!" {
        t.Fatalf("got %q", Greet("Mahesh"))
    }
    if Greet("") != "Hello, world!" {
        t.Fatalf("empty name got %q", Greet(""))
    }
}
```

**Checkpoint A3** — Run tests:

```bash
go test ./...
# expected: PASS (ok day1)
```

**Checkboxes:**
- [ ] All three tests pass
- [ ] `Greet("")` correctly returns the default

**Micro-questions:**
1. Why does Go use capitalized function names for exported functions?
2. How are Go's multiple return values different from Java?
3. What does `t.Fatalf` do vs `t.Error`?

---

### A4. Stretch (optional, 10 min)

Add `func Sum(nums ...int) int` (variadic) and a test proving `Sum(1,2,3,4) = 10`.

**Checkbox:**
- [ ] Variadic function works

---

## PART B — HTMX (~45 min)

You'll serve a tiny page from Go and wire 4 HTMX attributes.

### B1. Minimal server

**Task:** In `day1/htmx/`, create `server.go` serving:

- **GET /** → returns an HTML page (with the htmx script tag from a CDN)
- **GET /time** → returns the current time as an HTML fragment (just a `<span>`)
- **POST /echo** → reads a form field `msg` and returns `<p>You said: {msg}</p>`

**Minimal starter:**

```go
package main

import (
    "fmt"
    "net/http"
    "time"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, page)
    })
    http.HandleFunc("/time", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "<span>%s</span>", time.Now().Format(time.Kitchen))
    })
    http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
        r.ParseForm()
        fmt.Fprintf(w, "<p>You said: %s</p>", r.FormValue("msg"))
    })
    http.ListenAndServe(":8080", nil)
}

const page = `<!DOCTYPE html>
<html>
<head><script src="https://unpkg.com/htmx.org@2"></script></head>
<body>
<!-- TODO B2: add hx-get button that swaps the time -->
<!-- TODO B3: add hx-post form that echoes a message -->
</body>
</html>`
```

---

### B2. hx-get + hx-target + hx-swap

**Task:** Add a button that, when clicked, fetches `/time` and puts the result into a `<div id="clock">`.

**Required attributes to use:** `hx-get`, `hx-target`, `hx-swap`.

**Checkpoint B2 (manual):** Open http://localhost:8080, click the button.

**Checkboxes:**
- [ ] Clicking the button shows the current time inside `#clock`
- [ ] You used `hx-target="#clock"`
- [ ] You chose a `hx-swap` value and can explain it

---

### B3. hx-post

**Task:** Add a form with a text input named `msg` and a submit button that uses `hx-post` to `/echo` and swaps the response below the form.

**Checkpoint B3 (manual):** Type "hello", submit.

**Checkbox:**
- [ ] Page shows "You said: hello" without a full page reload

**Micro-questions:**
1. What is the default value of `hx-swap` if you don't set it?
2. What's the difference between `innerHTML` and `outerHTML` swaps?
3. Why does HTMX return HTML fragments instead of JSON here?

---

## PART C — MongoDB (~30 min)

Use Compass (GUI) and the mongosh shell so you feel both.

### C1. Create data (documents vs collections)

**Task:** In database `day1`, collection `engineers`, insert these 3 documents:

```javascript
use day1
db.engineers.insertMany([
  { name: "Mahesh", role: "Senior Lead", years: 12, skills: ["Java", "Spring"] },
  { name: "Asha", role: "SDE II", years: 5, skills: ["Go", "MongoDB"] },
  { name: "Ravi", role: "SDE III", years: 8, skills: ["Go", "HTMX", "MongoDB"] }
])
```

**Checkpoint C1:**

```javascript
db.engineers.countDocuments()  // expected: 3
```

**Checkboxes:**
- [ ] 3 documents inserted
- [ ] You can see the collection in Compass

---

### C2. Query (filters + projection)

**Task:** Write queries to answer:

1. All engineers with `years >= 8` — return only `name` and `years` (hide `_id`).
2. All engineers who know "Go".
3. Sort everyone by `years` descending.

**Starter for Q1:**

```javascript
db.engineers.find({ years: { $gte: 8 } }, { _id: 0, name: 1, years: 1 })
```

**Checkpoint C2**

**Checkboxes:**
- [ ] Q1 returns Mahesh(12) and Ravi(8) only, with just `name + years`
- [ ] Q2 returns Asha and Ravi
- [ ] Q3 orders Mahesh → Ravi → Asha

---

### C3. BSON understanding (written)

**Answer in notes:**

1. What BSON type is `years`? What about `skills`?
2. What is `_id` and what type does MongoDB give it by default?
3. One difference between a document and a row, and a collection and a table.

**Checkbox:**
- [ ] All 3 answered in your own words

---

## PART D — English (~30-45 min)

### D1. Vocabulary (use, don't memorize)

Write one sentence for each word, about your real work:

`professional`, `confident`, `collaboration`, `responsibility`, `experience`

**Checkpoint:**
- [ ] 5 original sentences, each using the word correctly

---

### D2. Grammar — present simple vs present perfect

Fill the correct form, then check against the key:

1. I _____ (work) as a software engineer for eight years. (still true → ?)
2. I _____ (lead) three teams so far in my career.
3. Every day, I _____ (review) pull requests.
4. I _____ (just/finish) a migration project.
5. We _____ (use) MongoDB in our current system.

**Checkpoint:**
- [ ] Attempt all 5 before checking answers

---

### D3. Speaking — self-introduction (record yourself)

Speak 60-90 seconds answering **"Tell me about yourself."** Cover:

- who you are + current role (present simple)
- key experience/achievements (present perfect)
- what you're learning now (present continuous)

**Checkpoint D3**

**Checkboxes:**
- [ ] Recorded once, listened back
- [ ] Used at least 2 present-perfect sentences ("I have led...", "I have worked...")
- [ ] Under 90 seconds

---

## ✅ Day 1 Done — Definition of Complete

**Completion checklist:**

- [ ] Go: `go test ./...` is green (A3 + optional A4)
- [ ] HTMX: button swaps time (B2) and form echoes message (B3)
- [ ] MongoDB: 3 queries return correct results (C2)
- [ ] English: vocab sentences + grammar attempt + recorded intro
- [ ] `answers/day1-notes.md` has all micro-question answers

**Score yourself:** _/5 sections green. **Anything red = tomorrow's warm-up.**

---

## 🔑 Answer Key (try first, then check)

### Go answers

**A1 Micro-questions:**
1. `go mod init` creates `go.mod` — it declares your module path, allows Go to manage dependencies, and enables other packages to import your code.
2. `package main` is required because Go looks for a `main()` function to execute. Without it, you just have a library.
3. `go run .` compiles and executes in one step (temporary). `go build` compiles to a binary file on disk for distribution.

**A2 Trap check:**
- `int` zero value: `0`
- `string` zero value: `""` (empty string)
- `bool` zero value: `false`

**A3 Micro-questions:**
1. Capitalization signals export. Capitalized names are exported (visible outside the package); lowercase names are private to the package.
2. Go allows multiple return values natively (tuples). Java requires a wrapper object or an array.
3. `t.Fatalf` stops the test immediately and marks it failed. `t.Error` logs the error but continues the test.

---

### HTMX answers

**B2 Micro-questions:**
1. The default value of `hx-swap` is `innerHTML` — it replaces the inner HTML of the target element.
2. `innerHTML` replaces content inside the element; `outerHTML` replaces the element itself (including its tags).
3. HTMX returns HTML fragments because the server sends back ready-to-insert HTML. JSON requires the client to render (extra work); fragments go straight to the DOM.

---

### MongoDB answers

**C2 Queries:**

1. All engineers with `years >= 8`, returning only `name` and `years`:
   ```javascript
   db.engineers.find({ years: { $gte: 8 } }, { _id: 0, name: 1, years: 1 })
   ```
   Result: Mahesh (12), Ravi (8)

2. All engineers who know "Go":
   ```javascript
   db.engineers.find({ skills: "Go" })
   ```
   Result: Asha, Ravi

3. Sort everyone by `years` descending:
   ```javascript
   db.engineers.find().sort({ years: -1 })
   ```
   Result: Mahesh (12) → Ravi (8) → Asha (5)

**C3 BSON understanding:**
1. `years` is an Int32 (or Int64). `skills` is an Array of Strings.
2. `_id` is a unique identifier. MongoDB creates it as an ObjectId by default (12-byte value with timestamp, machine ID, process ID, and counter).
3. A document is like a row, but flexible (different fields per document). A collection is like a table, but schema-less (no fixed schema).

---

### English answers

**D1 Vocabulary example sentences:**
- Professional: "I always maintain professional communication with my clients and colleagues."
- Confident: "I feel confident solving backend problems after working with distributed systems for two years."
- Collaboration: "Collaboration with the frontend team helped us ship the API on time."
- Responsibility: "I take responsibility for code quality and test coverage in my team."
- Experience: "My experience with cloud infrastructure gives me an edge in system design."

**D2 Grammar:**
1. I **have worked** (present perfect — started in the past, still true)
2. I **have led** (present perfect — multiple times, still relevant)
3. Every day, I **review** (present simple — habitual action)
4. I **have just finished** (present perfect — recent)
5. We **use** (present simple — ongoing state)

---

## 📚 Tomorrow (Day 2 preview)

**Structs/methods/interfaces** · server-rendered partial updates · **MongoDB CRUD from code** · "explain a project end-to-end" in English.

**Bring your Day 1 red items** — they're tomorrow's warm-up.

---

**Good luck! 🚀**

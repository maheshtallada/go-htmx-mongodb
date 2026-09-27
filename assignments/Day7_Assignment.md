# Day 7 Assignment — Test-Driven Learning (Go · HTMX · MongoDB · English)

> **How this works:** Make each **failing test or checkpoint** pass. Answer keys are at the bottom — try first.
>
> **Time-box:** ~3 hours + 30–45m English.

> **Day 7 scope**
> - **Go:** JSON encode/decode, config from env, structured logging with `slog`
> - **HTMX:** templates + partials for **list + detail** views
> - **MongoDB:** aggregation pipeline basics + BSON↔struct mapping
> - **English:** explain a project end-to-end + sequence words (`first`, `then`, `finally`)

> **Warm-up:** finish any red items from Day 6 first.

---

## 0. Setup (5 min)

```bash
docker start mongo-day7 || docker run -d --name mongo-day7 -p 27017:27017 mongo:7
mkdir -p day7 && cd day7
go mod init day7
```

---

# PART A — Go: JSON, Config, `slog` (~60 min)

## A1. JSON marshal/unmarshal

**Task:** In `model.go`:

```go
type Task struct {
    ID    int    `json:"id"`
    Title string `json:"title"`
    Done  bool   `json:"done"`
}
```

```go
// ToJSON serializes a Task.
func ToJSON(t Task) (string, error)

// FromJSON parses JSON into a Task.
func FromJSON(s string) (Task, error)
```

### Checkpoint A1 — `model_test.go`

```go
package main

import "testing"

func TestJSONRoundTrip(t *testing.T) {
    in := Task{ID: 1, Title: "Ship", Done: true}
    s, err := ToJSON(in)
    if err != nil {
        t.Fatal(err)
    }
    out, err := FromJSON(s)
    if err != nil || out != in {
        t.Fatalf("round-trip failed: %v %v", out, err)
    }
}
```

- [ ] Round-trip returns the same struct
- [ ] JSON keys are lowercase via tags

> **Micro-questions:**
> 1. What controls the JSON field name?
> 2. What does an unexported (lowercase) field do to JSON output?

## A2. Config from env (the tested core)

**Task:** Add:

```go
// Port reads PORT from env, defaulting to "8080" when unset/empty.
func Port(getenv func(string) string) string
```

(Injecting `getenv` makes it testable.)

### Checkpoint A2

```go
func TestPort(t *testing.T) {
    if Port(func(string) string { return "" }) != "8080" {
        t.Fatal("default failed")
    }
    if Port(func(string) string { return "9090" }) != "9090" {
        t.Fatal("override failed")
    }
}
```

- [ ] Default returned when empty
- [ ] Env value overrides default

## A3. Structured logging with `slog`

**Task:** Use `log/slog` to log a JSON line with attributes:

```go
func LogStartup(logger *slog.Logger, port string) {
    logger.Info("server starting", "port", port, "env", "dev")
}
```

### Checkpoint A3 (manual)

Run:

```bash
go test ./...
```

and confirm JSON output like:

```json
{"time":"...","level":"INFO","msg":"server starting","port":"8080","env":"dev"}
```

- [ ] Uses `slog.NewJSONHandler`
- [ ] Log line includes structured attributes

### Checkpoint A (all tests)

```bash
go test ./...
```

> **Micro-questions:**
> 1. Why structured (JSON) logs over `fmt.Println`?
> 2. What are `slog` **levels** and how do you set the minimum?
> 3. What is a **logger with context** (`With`) good for?

---

# PART B — HTMX List + Detail Views (~45 min)

## B1. Template files (List)

**Task:** In `day7/htmx/`, use `html/template` with `{{define}}` blocks: a `list` template rendering task rows, each linking to a detail view via `hx-get="/tasks/{{id}}"`.

### Checkpoint B1

- [ ] List renders from a slice
- [ ] Each row can request its detail fragment

## B2. Detail fragment swapped into a panel

**Task:** `GET /tasks/{id}` returns a **detail fragment**; clicking a row swaps it into `#detail` (`hx-target="#detail"`, `hx-swap="innerHTML"`).

### Checkpoint B2

- [ ] Clicking a task shows its detail in the panel, no reload
- [ ] Detail route returns only the fragment

> **Micro-questions:**
> 1. What does `template.Must` do and when does it panic?
> 2. Difference between `{{define}}` / `{{template}}` / `{{block}}`?
> 3. How do you pass data into a nested template?

## B3. Stretch (optional)

Add an “active row” highlight using an out-of-band swap that updates a `#selected` label.

- [ ] Selected task name appears via OOB

---

# PART C — MongoDB Aggregation Basics (~30 min)

## C1. Seed

```js
use day7

db.sales.insertMany([
  { region: "West", product: "A", qty: 5, price: 10 },
  { region: "West", product: "B", qty: 2, price: 20 },
  { region: "East", product: "A", qty: 7, price: 10 },
  { region: "East", product: "C", qty: 1, price: 50 },
  { region: "West", product: "A", qty: 3, price: 10 }
])
```

- [ ] 5 docs inserted

## C2. `$group` + `$sum` (the tested core)

**Task:** Total revenue (`qty * price`) **per region**.

```js
db.sales.aggregate([
  { $project: { region: 1, revenue: { $multiply: ["$qty", "$price"] } } },
  { $group: { _id: "$region", total: { $sum: "$revenue" } } },
  { $sort: { total: -1 } }
])
```

### Checkpoint C2

- [ ] West and East each get a total
- [ ] Sorted by total descending
- [ ] You used `$project`, `$group`, `$sum`, `$sort`

## C3. `$match` + `$group`

**Task:** For product `"A"` only, total `qty` across all regions.

- [ ] Returns a single grouped total for product A

## C4. BSON ↔ struct (written)

1. Map this pipeline result to a Go struct — what `bson` tags would `_id`, `total` need?
2. What is the difference between `bson.M`, `bson.D`, and `bson.A`?
3. Why does aggregation order (stage sequence) matter for performance?

- [ ] All 3 answered

---

# PART D — English (~30–45 min)

## D1. Vocabulary

One sentence each: `requirement`, `design`, `implementation`, `testing`, `deployment`.

- [ ] 5 original sentences

## D2. Grammar — sequence words

Rewrite this into 5 ordered steps using `first`, `then`, `after that`, `next`, `finally`:

> "We gather requirements, design the API, build it, test it, and deploy."

- [ ] 5 steps, each starting with a sequence word

## D3. Speaking — explain a project end-to-end (record yourself)

Speak **90–120 seconds** walking through **one project's lifecycle**:

`requirement → design → implementation → testing → deployment → outcome`

Use sequence words.

### Checkpoint D3

- [ ] Recorded, listened back
- [ ] Clear ordered flow (no jumping around)
- [ ] Ended with the outcome/impact

---

## ✅ Day 7 Done — Definition of Complete

- [ ] Go: `go test ./...` green (A1–A2) + `slog` JSON output (A3)
- [ ] HTMX: list (B1) + detail fragment (B2)
- [ ] MongoDB: `$group` / `$sum` revenue per region (C2–C3)
- [ ] English: vocab + sequence-words steps + project walkthrough recording
- [ ] `answers/day7-notes.md` complete

Score yourself: **__/5 sections green.**

---

# 🔑 Answer Key (try first, then check)

<details>
<summary>Go answers</summary>

```go
// model.go
package main

import (
    "encoding/json"
    "log/slog"
    "os"
)

type Task struct {
    ID    int    `json:"id"`
    Title string `json:"title"`
    Done  bool   `json:"done"`
}

func ToJSON(t Task) (string, error) {
    b, err := json.Marshal(t)
    return string(b), err
}

func FromJSON(s string) (Task, error) {
    var t Task
    err := json.Unmarshal([]byte(s), &t)
    return t, err
}

func Port(getenv func(string) string) string {
    if p := getenv("PORT"); p != "" {
        return p
    }
    return "8080"
}

func LogStartup(logger *slog.Logger, port string) {
    logger.Info("server starting", "port", port, "env", "dev")
}

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    LogStartup(logger, Port(os.Getenv))
}
```

- **A1.1:** The `json:"name"` struct tag controls the JSON key.
- **A1.2:** Unexported fields are not serialized by `encoding/json`.
- **A3.1:** Structured logs are machine-parseable, searchable, and filterable in log tools.
- **A3.2:** `slog` levels include DEBUG, INFO, WARN, and ERROR; the minimum is controlled by the handler's `slog.Level` configuration.
- **A3.3:** `logger.With("requestID", id)` returns a child logger that stamps that attribute on every line.

</details>

<details>
<summary>HTMX answers</summary>

- **B1.1:** `template.Must` wraps a `(*Template, error)` and **panics** on parse error — useful at startup to fail fast.
- **B2.2:** `{{define "x"}}` declares a named template; `{{template "x" .}}` renders it; `{{block "x" .}}` defines and renders with a default (overridable).
- **B3.1:** Pass a value as the pipeline argument: `{{template "row" .Task}}`.

</details>

<details>
<summary>MongoDB answers</summary>

```js
// C3
db.sales.aggregate([
  { $match: { product: "A" } },
  { $group: { _id: "$product", totalQty: { $sum: "$qty" } } }
])
```

- **C4.1:** `bson:"_id"` and `bson:"total"` on the struct fields.
- **C4.2:** `bson.M` = unordered map; `bson.D` = ordered slice of key/value pairs (order matters, e.g. for commands); `bson.A` = array.
- **C4.3:** Putting `$match` early filters documents before expensive stages (and can use indexes), reducing work downstream.

</details>

<details>
<summary>English answers</summary>

**D2 sequence:**
1. **First**, we gather requirements.
2. **Then**, we design the API.
3. **After that**, we build it.
4. **Next**, we test it.
5. **Finally**, we deploy.

</details>

---

## Tomorrow (Day 8 preview)

DB access with the driver + pagination patterns; HTMX form submit + inline validation errors; MongoDB validation rules, upserts, write concern; English: conflict & resolution + reported speech.

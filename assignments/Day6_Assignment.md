# Day 6 Assignment --- Test-Driven Learning (Go · HTMX · MongoDB · English)

> 💡 **How this works:** Make each **failing test or checkpoint** pass.
> Answer keys at the bottom --- try first.

> **Time-box:** \~3 hours + 30--45m English.

> **Day 6 scope** --- **Week 2 begins: APIs + UI wiring** - **Go:**
> net/http, routing with **Chi**, middleware - **HTMX:** wiring HTMX to
> real routes returning partial responses - **MongoDB:** connect with
> the official Go driver, context-aware CRUD repository - **English:**
> polished 90-second "tell me about yourself", sentence structure & word
> order

> **Warm-up:** Finish any red items from Day 5 first.

------------------------------------------------------------------------

## 0. Setup (5 min)

``` bash
docker start mongo-day1 || docker run -d --name mongo-day6 -p 27017:27017 mongo:7
mkdir -p day6 && cd day6
go mod init day6
go get github.com/go-chi/chi/v5
go get go.mongodb.org/mongo-driver/mongo
```

------------------------------------------------------------------------

## PART A --- Go: net/http + Chi + Middleware (x60 min)

### A1. A Chi router with two routes

**Task:** In `main.go`, build a Chi router: - GET `/health` → `200` and
body `"ok"` - GET `/greet/{name}` → `"Hello, <name>!"`

Starter:

``` go
package main

import (
    "net/http"

    "github.com/go-chi/chi/v5"
)

func Router() http.Handler {
    r := chi.NewRouter()
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("ok"))
    })
    // TODO A1: GET /greet/{name}
    return r
}

func main() {
    http.ListenAndServe(":8080", Router())
}
```

**Checkpoint A1** --- write `main_test.go` using `httptest`:

``` go
package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHealth(t *testing.T) {
    req := httptest.NewRequest("GET", "/health", nil)
    rec := httptest.NewRecorder()
    Router().ServeHTTP(rec, req)
    if rec.Code != 200 || rec.Body.String() != "ok" {
        t.Fatalf("got %d %q", rec.Code, rec.Body.String())
    }
}

func TestGreet(t *testing.T) {
    req := httptest.NewRequest("GET", "/greet/Mahesh", nil)
    rec := httptest.NewRecorder()
    Router().ServeHTTP(rec, req)
    if rec.Body.String() != "Hello, Mahesh!" {
        t.Fatalf("got %q", rec.Body.String())
    }
}
```

-   [ ] Both routes pass
-   [ ] You used `chi.URLParam(r, "name")`

> **Micro-questions:** 1. What does `http.Handler` vs `http.HandlerFunc`
> mean? 2. How does Chi extract URL path params?

### A2. Middleware (the tested core)

**Task:** Write a middleware that adds a response header
`X-Request-ID: test-123` (hard-coded is fine for the test) and wrap the
router with it.

``` go
func RequestID(next http.Handler) http.Handler
```

**Checkpoint A2**:

``` go
func TestMiddleware(t *testing.T) {
    req := httptest.NewRequest("GET", "/health", nil)
    rec := httptest.NewRecorder()
    RequestID(Router()).ServeHTTP(rec, req)
    if rec.Header().Get("X-Request-ID") == "" {
        t.Fatal("missing X-Request-ID header")
    }
}
```

-   [ ] Middleware is `func(http.Handler) http.Handler`
-   [ ] Header is present on the response

**Checkpoint A (all tests)**

``` bash
go test ./...
```

> **Micro-questions:** 1. What is the middleware pattern (wrapping
> handlers)? 2. Name 3 things middleware is commonly used for. 3. Where
> does Chi's built-in `middleware.Logger` fit?

### A3. Stretch (optional)

Add `middleware.Recoverer` (Chi built-in) and a route that panics; prove
the server returns `500` instead of crashing. - \[ \] Panic is recovered
→ `500`

------------------------------------------------------------------------

## PART B --- HTMX Wired to Real Routes (x45 min)

### B1. Serve page + fragment from the same Chi router

**Task:** Extend the router: - GET `/` → full HTML page (with htmx
script) - GET `/fragment` → a `<div>` fragment - Button on the page does
`hx-get="/fragment"` into `#slot`

**Checkpoint B1** (manual + test): - \[ \] Button swaps the fragment
into `#slot` - \[ \] GET `/fragment` returns **only** the fragment
(assert via `httptest`)

### B2. Detect HTMX requests

**Task:** In a handler, branch on the `HX-Request` header: if present,
return the fragment; otherwise return the full page. This lets one route
serve both.

``` go
if r.Header.Get("HX-Request") == "true" {
    // return fragment
} else {
    // return full page
}
```

**Checkpoint B2** - \[ \] With `HX-Request: true` → fragment - \[ \]
Without it → full page

> **Micro-questions:** 1. What headers does HTMX add to requests
> (`HX-Request`, `HX-Target`, `HX-Trigger`)? 2. Why is "one URL, two
> representations (page vs fragment)" useful? 3. What's the benefit of
> returning HTML fragments over JSON here?

------------------------------------------------------------------------

## PART C --- MongoDB with the Go Driver (x30 min)

### C1. Connect + ping

**Task:** In `repo.go`, connect to Mongo with a `context` and ping it.

``` go
package main

import (
    "context"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(ctx context.Context, uri string) (*mongo.Client, error) {
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        return nil, err
    }
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    return client, client.Ping(ctx, nil)
}
```

**Checkpoint C1** (needs Mongo running): - \[ \]
`Connect(ctx, "mongodb://localhost:27017")` returns no error - \[ \]
Ping succeeds

### C2. Context-aware CRUD repository (the tested core)

**Task:** Add a `Task` struct + repository methods that all take
`context.Context`:

``` go
type Task struct {
    ID    string `bson:"_id,omitempty"`
    Title string `bson:"title"`
    Done  bool   `bson:"done"`
}

type TaskRepo struct{ col *mongo.Collection }

func (r *TaskRepo) Insert(ctx context.Context, t Task) (string, error)
func (r *TaskRepo) FindAll(ctx context.Context) ([]Task, error)
```

**Checkpoint C2** (integration-style): - \[ \] `Insert` adds a doc,
returns its id - \[ \] `FindAll` returns it back - \[ \] Every method
takes `ctx` as the first arg

### C3. Understanding (written)

1.  Why does every driver call take a `context.Context`?
2.  What do the `bson:"title"` struct tags do?
3.  What does `omitempty` on `_id` achieve on insert?

-   [ ] All 3 answered

------------------------------------------------------------------------

## PART D --- English (x30--45 min)

### D1. Vocabulary

One sentence each: `background`, `exposure`, `responsibility`,
`delivery`, `achievement`.

-   [ ] 5 original sentences

### D2. Grammar --- sentence structure & word order

Reorder into a correct sentence:

1.  (a senior engineer / for eight years / I / have been)
2.  (microservices / we / in production / run / dozens of)
3.  (yesterday / the team / a critical bug / fixed)
4.  (learning / currently / Go and HTMX / I am)
5.  (reliable / to build / systems / my goal is)

-   [ ] Attempt all 5

### D3. Speaking --- polished self-introduction (record yourself)

Record a final 90-second "tell me about yourself": - who you are +
current role (present simple) - 2 key achievements (present perfect) -
what you're learning now (present continuous) - one sentence on what
you're looking for

**Checkpoint D3** - \[ \] Recorded, listened back, under 90s - \[ \]
Clear structure, no rambling - \[ \] At least 2 present-perfect
sentences

------------------------------------------------------------------------

## 🟩 Day 6 Done --- Definition of Complete

-   [ ] Go: `go test ./...` green (A1--A2, optional A3)
-   [ ] HTMX: fragment route (B1) + `HX-Request` branching (B2)
-   [ ] MongoDB: driver connect (C1) + context-aware repo (C2)
-   [ ] English: vocab + word-order + polished intro recording
-   [ ] `answers/day6-notes.md` complete

Score yourself: \*\*\_\_/5 sections green.\*\*

------------------------------------------------------------------------

## 📝 Answer Key (try first, then check)

```{=html}
<details>
```
```{=html}
<summary>
```
Go answers
```{=html}
</summary>
```
``` go
// A1: greet route
r.Get("/greet/{name}", func(w http.ResponseWriter, r *http.Request) {
    name := chi.URLParam(r, "name")
    fmt.Fprintf(w, "Hello, %s!", name)
})

// A2: middleware
func RequestID(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-Request-ID", "test-123")
        next.ServeHTTP(w, r)
    })
}
```

-   **A1.1** `http.Handler` is an interface with `ServeHTTP`;
    `http.HandlerFunc` adapts a plain func to that interface.
-   **A1.2** Chi matches `{name}` patterns and exposes them via
    `chi.URLParam`.
-   **A2** Middleware wraps a handler, running code before/after
    `next.ServeHTTP` --- used for logging, auth, request IDs, recovery,
    CORS, timing.

```{=html}
</details>
```
```{=html}
<details>
```
```{=html}
<summary>
```
HTMX answers
```{=html}
</summary>
```
``` go
func handleIndex(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("HX-Request") == "true" {
        fmt.Fprint(w, `<div id="slot">fragment</div>`)
        return
    }
    fmt.Fprint(w, fullPage)
}
```

-   **B.1** HTMX adds `HX-Request: true`, plus `HX-Target`,
    `HX-Trigger`, `HX-Current-URL`, etc.
-   **B.2** One canonical URL can serve a full page (direct visit /
    refresh) and a fragment (HTMX swap) --- good for progressive
    enhancement and shareable URLs.
-   **B.3** Fragments let the browser swap HTML directly --- no
    client-side JSON parsing or templating (hypermedia approach).

```{=html}
</details>
```
```{=html}
<details>
```
```{=html}
<summary>
```
MongoDB answers
```{=html}
</summary>
```
``` go
func (r *TaskRepo) Insert(ctx context.Context, t Task) (string, error) {
    res, err := r.col.InsertOne(ctx, t)
    if err != nil {
        return "", err
    }
    return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *TaskRepo) FindAll(ctx context.Context) ([]Task, error) {
    cur, err := r.col.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    defer cur.Close(ctx)
    var out []Task
    if err := cur.All(ctx, &out); err != nil {
        return nil, err
    }
    return out, nil
}
```

-   **C3.1** `context` enables timeouts/cancellation/deadlines to
    propagate to the DB call.
-   **C3.2** `bson:"title"` maps the Go field to the BSON field name.
-   **C3.3** `omitempty` on `_id` lets Mongo generate the ObjectId on
    insert instead of sending an empty value.

```{=html}
</details>
```
```{=html}
<details>
```
```{=html}
<summary>
```
English answers
```{=html}
</summary>
```
**D2 word order:**

1.  I have been a senior engineer for eight years.
2.  We run dozens of microservices in production.
3.  Yesterday the team fixed a critical bug.
4.  I am currently learning Go and HTMX.
5.  My goal is to build reliable systems.

**Rule:** English default order is **Subject → Verb → Object**, with
time/place at the start or end.
```{=html}
</details>
```

------------------------------------------------------------------------

## Tomorrow (Day 7 preview)

JSON, config, structured logging with `slog` + HTMX templates/partials
for list + detail views + MongoDB aggregation pipeline basics +
BSON↔struct mapping + English: explain a project end-to-end + sequence
words.

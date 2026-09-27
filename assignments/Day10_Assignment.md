# Day 10 Assignment --- Test-Driven Learning (Go · HTMX · MongoDB · English)

> **How this works:** Make each **failing test or checkpoint** pass.
> Answer keys at the bottom --- try first.
>
> **Time-box:** \~3 hours + 30--45m English.
>
> **Day 10 scope** - **Go:** build, Docker, env config, **graceful
> shutdown** - **HTMX:** page shell / layout + reusable fragments -
> **MongoDB:** Atlas vs local Docker, connection strings, backup/restore
> basics - **English:** failure & learning, conditionals
> (`if`/`when`/`unless`), 5 vocab words
>
> **Warm-up:** finish any red items from Day 9 first.

------------------------------------------------------------------------

## 0. Setup (5 min)

``` bash
docker start mongo-day10 || docker run -d --name mongo-day10 -p 27017:27017 mongo:7
mkdir -p day10 && cd day10
go mod init day10
```

------------------------------------------------------------------------

## PART A --- Go: Build, Docker, Graceful Shutdown (\~60 min)

### A1. Graceful shutdown (the tested core)

**Task:** In `main.go`, run an `http.Server` and shut it down cleanly on
`SIGINT`/`SIGTERM` using `signal.NotifyContext` +
`server.Shutdown(ctx)`.

Starter:

``` go
package main

import (
    "context"
    "errors"
    "log"
    "net/http"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    srv := &http.Server{Addr: ":8080", Handler: http.DefaultServeMux}
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("ok"))
    })

    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    go func() {
        if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
            log.Fatal(err)
        }
    }()
    log.Println("listening on :8080")

    <-ctx.Done() // wait for signal
    log.Println("shutting down...")

    shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(shutdownCtx); err != nil {
        log.Printf("shutdown error: %v", err)
    }
    log.Println("stopped cleanly")
}
```

**Checkpoint A1** (manual): run, hit `/health`, then `Ctrl+C`. - \[ \]
Server logs `"shutting down"` then `"stopped cleanly"` - \[ \] In-flight
requests are allowed to finish (test with a slow handler) - \[ \] Uses
`signal.NotifyContext` + `srv.Shutdown`

> **Micro-questions:** 1. Why is graceful shutdown important in
> Kubernetes (SIGTERM → grace period)? 2. What does
> `http.ErrServerClosed` signal?

### A2. Dockerize (multi-stage build)

**Task:** Write a multi-stage `Dockerfile`:

``` dockerfile
# build stage
FROM golang:1.22 AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app ./...

# run stage (tiny)
FROM gcr.io/distroless/static-debian12
COPY --from=build /app /app
EXPOSE 8080
ENTRYPOINT ["/app"]
```

``` bash
docker build -t day10 .
docker run -p 8080:8080 day10
```

**Checkpoint A2:** - \[ \] Image builds - \[ \] Container serves
`/health` - \[ \] Final image is small (distroless/scratch)

### A3. Env config

**Task:** Read `PORT` and `MONGO_URI` from env with sane defaults; log
them at startup. - \[ \] Configurable via `-e PORT=9090`

> **Micro-questions:** 1. Why multi-stage builds (small final image, no
> build tools shipped)? 2. Why `CGO_ENABLED=0` for a static binary? 3.
> Why read config from env (12-factor)?

------------------------------------------------------------------------

## PART B --- HTMX Layout & Reusable Fragments (\~45 min)

### B1. Base layout + content block

**Task:** In `day10/htmx/`, create a base layout template with a
`{{block "content" .}}` and pages that fill it.

**Checkpoint B1:** - \[ \] Two pages share one layout (header/footer) -
\[ \] Each provides its own `content`

### B2. Reusable row/card fragment

**Task:** Extract a `card` fragment used by both a 1st page and an HTMX
partial response (append on add).

**Checkpoint B2:** - \[ \] The same `card` template renders in full-page
and in fragment responses - \[ \] No duplicated HTML between page and
fragment

> **Micro-questions:** 1. Why keep fragments as separate named
> templates? 2. How does a layout avoid repeating `<head>`/scripts on
> every page? 3. What's the DRY benefit of one fragment for page +
> partial?

------------------------------------------------------------------------

## PART C --- MongoDB Ops: Connections & Backups (\~30 min)

### C1. Connection strings (written + try)

1.  Local: `mongodb://localhost:27017`
2.  Atlas SRV: `mongodb+srv://user:pass@cluster0.xxxx.mongodb.net/db`
3.  With options: `?retryWrites=true&w=majority`

-   [ ] You can explain each part (scheme, host, auth, options)

### C2. Backup & restore (the tested core)

**Task:** Dump and restore a database using the Mongo tools:

``` bash
# seed something first, then:
docker exec mongo-day10 mongodump --db day10 --out /tmp/backup
docker exec mongo-day10 mongorestore --drop --db day10 /tmp/backup/day10
```

**Checkpoint C2:** - \[ \] `mongodump` creates a backup - \[ \]
`mongorestore --drop` restores it cleanly

### C3. Understanding (written)

1.  Difference between `mongodump`/`mongorestore` and
    `mongoexport`/`mongoimport`?
2.  What does Atlas give you that self-hosted doesn't (managed backups,
    scaling)?
3.  Why never commit connection strings with credentials?

-   [ ] ALL 3 answered

------------------------------------------------------------------------

## PART D --- English (\~30--45 min)

### D1. Vocabulary

One sentence each: `mistake`, `lesson`, `feedback`, `growth`,
`resilience`

-   [ ] 5 original sentences

### D2. Grammar --- conditionals

Complete with the right conditional (`if`/`when`/`unless`):

1.  \_\_\_\_\_\_ the tests pass, we deploy. *(condition)*
2.  We roll back \_\_\_\_\_\_ the error rate spikes. *(negative
    condition = if not stable)*
3.  \_\_\_\_\_\_ I get an alert, I check the dashboard first. *(habit →
    zero conditional)*
4.  If we \_\_\_\_\_\_ (add) an index, the query \_\_\_\_\_\_ (be)
    faster. *(first conditional)*
5.  We won't ship \_\_\_\_\_\_ QA signs off. *(unless)*

-   [ ] Attempt all 5

### D3. Speaking --- failure & learning (record yourself)

Speak **90 seconds** answering *"Tell me about a failure and what you
learned."* Be honest; end on the **lesson** and how you applied it.

**Checkpoint D3** - \[ \] Recorded, listened back - \[ \] Real failure
(not humble-brag) - \[ \] Clear lesson + how you changed afterward

------------------------------------------------------------------------

## ✅ Day 10 Done --- Definition of Complete

-   [ ] Go: graceful shutdown (A1) + Docker image (A2) + env config (A3)
-   [ ] HTMX: shared layout (B1) + reusable fragment (B2)
-   [ ] MongoDB: connection strings + dump/restore (C1--C2)
-   [ ] English: vocab + conditionals + failure recording
-   [ ] `answers/day10-notes.md` complete

Score yourself: \*\*\_\_/5 sections green.\*\*

------------------------------------------------------------------------

## 🔑 Answer Key (try first, then check)

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
-   **A1.1** K8s sends `SIGTERM`, waits `terminationGracePeriodSeconds`,
    then `SIGKILL`; graceful shutdown drains in-flight requests within
    that window.
-   **A1.2** `http.ErrServerClosed` is returned by `ListenAndServe`
    after a normal `Shutdown` --- treat it as expected, not an error.
-   **A3.1** Multi-stage keeps compilers out of the final image →
    smaller, fewer CVEs.
-   **A3.2** `CGO_ENABLED=0` builds a fully static binary that runs on
    `scratch`/distroless.
-   **A3.3** Env config makes the same image portable across
    environments (12-factor).

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
-   **B.1** Separate named templates are reusable in both full pages and
    partial responses.
-   **B.2** A base layout defines `<head>`, scripts, header/footer once;
    pages only supply the `content` block.
-   **B.3** One fragment for page + partial = no drift between initial
    render and HTMX updates.

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
-   **C3.1** `mongodump`/`restore` use binary BSON (full fidelity,
    backups); `mongoexport`/`import` use JSON/CSV (human-readable, data
    exchange).
-   **C3.2** Atlas provides managed backups, monitoring, scaling,
    patching, and easy HA.
-   **C3.3** Committed credentials leak secrets; use env vars / secret
    managers instead.

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
**D2 conditionals:** 1. **If/When** 2. **if/when** (roll back if it
spikes) 3. **When** (routine) 4. **add / will be** 5. **unless**

**Rules:** - Zero conditional (`when`/`if` + present) = facts/habits. -
First conditional (`if` + present, `will` + verb) = real future
possibility. - `unless` = "if not".

```{=html}
</details>
```

------------------------------------------------------------------------

## Tomorrow (Day 11 preview)

Project layout, clean architecture, dependency injection · HTMX
progressive-enhancement mindset · MongoDB repository boundaries + schema
versioning · English: stakeholder updates + concise business English.
**Week 3 begins.**

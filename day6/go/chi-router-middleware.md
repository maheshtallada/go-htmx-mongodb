# Day 6 Go: Chi routes and middleware

## What the program builds

`Router` returns an `http.Handler`, which is the standard-library interface
implemented by anything that can serve an HTTP request. Chi's router is one
such handler. `main` gives the fully wrapped handler to `ListenAndServe`.

The three endpoints demonstrate routing:

- `GET /health` responds with `ok`.
- `GET /greet/Mahesh` responds with `Hello, Mahesh!`.
- `GET /panic` deliberately panics so Chi's `Recoverer` middleware can turn
  the panic into an HTTP 500 response rather than terminate the server.

The greet route contains `{name}` as a named path segment. Chi stores the
matched value on the request, and `chi.URLParam(req, "name")` retrieves it.

## Middleware is a wrapper

An `http.Handler` has a `ServeHTTP(ResponseWriter, *Request)` method.
`http.HandlerFunc` adapts a function with that signature into an `http.Handler`.
`RequestID` takes the next handler and returns a new handler. The returned
handler adds a header, then calls `next.ServeHTTP` to continue the request.
This is the same general wrapping/decorator idea used for logging,
authentication, tracing, and request limits.

The test calls `RequestID(Router())` to compose the wrapper around the router.
The production `main` does the same. Chi's built-in `Recoverer` is registered
on the router with `r.Use`, so it protects route handlers from panics.

## How the tests work

`httptest.NewRequest` creates an in-memory request and
`httptest.NewRecorder` captures the response. Calling `ServeHTTP` runs the
same handler chain without opening a network socket. Tests verify route
status/body, the middleware header, and panic recovery.

Run from this directory with `go test ./...`.

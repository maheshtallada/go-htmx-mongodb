
package main // Keep the router and its handlers in the executable's main package.

import (
	"net/http" // Defines HTTP handlers and request/response types.

	"github.com/go-chi/chi/v5"            // Provides the Chi router and path parameter lookup.
	"github.com/go-chi/chi/v5/middleware" // Provides reusable HTTP middleware, including panic recovery.
)

// Router builds the HTTP routes and returns them behind the standard handler interface.
func Router() http.Handler {
	r := chi.NewRouter()        // Create the Chi router that matches requests to route handlers.
	r.Use(middleware.Recoverer) // Recover handler panics and turn them into HTTP 500 responses.

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) { // Register the GET /health endpoint.
		_, _ = w.Write([]byte("ok"))
	})
	r.Get("/greet/{name}", func(w http.ResponseWriter, req *http.Request) { // {name} captures one path segment.
		_, _ = w.Write([]byte("Hello, " + chi.URLParam(req, "name") + "!"))
	})
	r.Get("/panic", func(http.ResponseWriter, *http.Request) { // Test-only route demonstrating Recoverer behavior.
		panic("intentional recovery exercise")
	})

	return r // Chi's router implements http.Handler, so it can be passed to net/http.
}

// RequestID is middleware: it wraps a handler to add a response header to its requests.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) { // Adapt this function to the http.Handler interface.
		w.Header().Set("X-Request-ID", "test-123") // Set the header before the wrapped handler writes the response.
		next.ServeHTTP(w, req)                     // Continue through the chain, eventually reaching the router.
	})
}

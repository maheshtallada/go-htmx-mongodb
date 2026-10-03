package main // Build the Day 6 HTTP server executable.

import (
	"log"      // Reports server startup errors.
	"net/http" // Provides the HTTP server.
)

func main() {
	handler := RequestID(Router())                   // Wrap the Chi router so each response receives the request ID header.
	log.Fatal(http.ListenAndServe(":8080", handler)) // Listen on port 8080 and serve until an error stops the server.
}

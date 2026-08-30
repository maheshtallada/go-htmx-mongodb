package main

import (
    "fmt"
    "net/http"
    "os"
    "time"
)

// main starts the HTTP server.
// This is the entry point of our Go program.
func main() {
    // Each URL path is connected to a handler function.
    // A handler decides what response to send for that route.
    http.HandleFunc("/", homePage)
    http.HandleFunc("/time", showTime)
    http.HandleFunc("/echo", echoMessage)

    fmt.Println("Server running at http://localhost:8080")

    // Start listening for HTTP requests on port 8080.
    // If anything goes wrong, the program will stop and print the error.
    if err := http.ListenAndServe(":8080", nil); err != nil {
        panic(err)
    }
}

// Note on writing responses:
// - fmt.Fprintf(w, ...) writes formatted text directly to the http.ResponseWriter (the HTTP response body).
//   This is how handlers "send" HTML fragments or pages back to the browser.
// - fmt.Sprintf(...) builds and returns a string (useful if you need the string before writing it anywhere).
// - fmt.Println(...) prints to standard output (the server console), not the HTTP response.
//
// Important note about HTTP methods and handlers:
// - Handlers are registered by URL path (e.g., "/time" or "/echo"), not by HTTP method.
//   That means any HTTP method (GET, POST, PUT, etc.) will reach the same handler unless
//   the handler checks r.Method and rejects unsupported methods.
// - Earlier, if a handler didn't check r.Method, calling POST on a path intended for GET would still reach it
//   and often "appear" to work. For example, a GET-only response like a time span would still be returned on POST.
// - For clarity and safety, the handler should check r.Method and return 405 Method Not Allowed for unsupported methods.
//
// What's ideal in production:
// - Enforce HTTP methods (check r.Method) or use a router that supports method-based routes.
// - Return correct status codes (405 for wrong method, 400 for bad input, 500 for server error).
// - Set appropriate headers (e.g., Content-Type: text/html; charset=utf-8).
// - Parse and validate user input carefully and escape or sanitize output to avoid XSS.
// - Use middleware for logging, authentication, CORS, CSRF protection, rate limiting and timeouts.
// - For performance, parse templates once (text/template/html/template) and reuse them instead of reading files every request.
// - Prefer net/http servers with read/write timeouts and graceful shutdown handling.
//
// These practices make handlers predictable, secure, and production-ready.

// homePage reads the HTML file and sends it to the browser.
// This keeps the page markup separate from Go logic.
func homePage(w http.ResponseWriter, r *http.Request) {
    // Only serve the homepage for "/".
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }

    // Read the full HTML file from disk.
    htmlBytes, err := os.ReadFile("index.html")
    if err != nil {
        http.Error(w, "could not load page", http.StatusInternalServerError)
        return
    }

    // Write the HTML content to the response.
    fmt.Fprint(w, string(htmlBytes))
}

// showTime returns just an HTML fragment, not a full page.
// This is useful for HTMX because it can swap this fragment into an existing page.
func showTime(w http.ResponseWriter, r *http.Request) {
    // Differentiate HTTP methods inside a handler using r.Method.
    // Handlers are registered by path (e.g., "/time"), but the same path can be called with
    // GET, POST, etc. You can check r.Method to allow or reject methods.
    if r.Method != http.MethodGet {
        // If a non-GET method calls /time, respond with 405 Method Not Allowed.
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // time.Kitchen gives a readable time like 3:05PM.
    // Write an HTML fragment to the response; HTMX will insert it into the page.
    fmt.Fprintf(w, "<span>%s</span>", time.Now().Format(time.Kitchen))
}

// echoMessage reads the form field named "msg" from a POST request.
// Then it returns a small HTML fragment showing the message.
func echoMessage(w http.ResponseWriter, r *http.Request) {
    // Only allow POST for this endpoint, because it's intended to receive form data.
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // ParseForm reads form data from the request body.
    // For example: msg=hello
    if err := r.ParseForm(); err != nil {
        http.Error(w, "bad request", http.StatusBadRequest)
        return
    }

    // FormValue gets a single form field by name.
    msg := r.FormValue("msg")

    // HTMX expects HTML fragments, so we return a small snippet.
    fmt.Fprintf(w, "<p>You said: %s</p>", msg)
}

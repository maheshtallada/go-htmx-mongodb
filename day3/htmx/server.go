package main

import (
    "fmt"
    "net/http"
    "strings"
)

// This file implements the HTMX assignment parts B1, B2 and B3.
// B1: server-side validation returns 422 and HTML fragments for errors.
// B2: on success, respond with HX-Redirect header to /welcome.
// B3: supports independent error slots (name + email) via out-of-band swaps.

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, page)
    })
    http.HandleFunc("/signup", handleSignup)
    http.HandleFunc("/welcome", welcome)
    http.ListenAndServe(":8080", nil)
}

func handleSignup(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, "bad request", http.StatusBadRequest)
        return
    }

    name := strings.TrimSpace(r.FormValue("name"))
    email := strings.TrimSpace(r.FormValue("email"))

    var nameErr, emailErr string
    if name == "" {
        nameErr = `<p class="error" id="name-msg" hx-swap-oob="true">Name is required</p>`
    }
    if !strings.Contains(email, "@") {
        emailErr = `<p class="error" id="email-msg" hx-swap-oob="true">Invalid email</p>`
    }

    if nameErr != "" || emailErr != "" {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        w.WriteHeader(http.StatusUnprocessableEntity) // 422
        // Return both fragments; HTMX will perform the out-of-band swaps into
        // the elements with matching ids on the page.
        fmt.Fprint(w, nameErr)
        fmt.Fprint(w, emailErr)
        return
    }

    // Success: instruct HTMX to redirect the browser to /welcome
    w.Header().Set("HX-Redirect", "/welcome")
    w.WriteHeader(http.StatusOK)
}

func welcome(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fmt.Fprint(w, "<h1>Welcome!</h1>")
}

const page = `<!doctype html>
<html>
  <head>
    <meta charset="utf-8">
    <script src="https://unpkg.com/htmx.org@2"></script>
    <style>
      .error { color: red; }
      label { display:block; margin-bottom:8px; }
    </style>
  </head>
  <body>
    <h2>Sign up</h2>
    <!-- Form: uses out-of-band swaps for independent error slots -->
    <form hx-post="/signup" hx-swap="none">
      <label>
        Name: <input name="name" placeholder="Your name">
      </label>
      <div id="name-msg"></div>

      <label>
        Email: <input name="email" placeholder="you@example.com">
      </label>
      <div id="email-msg"></div>

      <button type="submit">Sign up</button>
    </form>

    <!-- HTMX will replace #name-msg and #email-msg when server returns
         elements with hx-swap-oob="true" -->
  </body>
</html>`
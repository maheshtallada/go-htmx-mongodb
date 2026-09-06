package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
    "strconv"
    "strings"
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
    http.HandleFunc("/", homePage)
    http.HandleFunc("/engineers", addEngineer)
    http.HandleFunc("/engineers/", deleteEngineer)

    log.Println("Server running at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func homePage(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }

    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    if err := tmpl.ExecuteTemplate(w, "page", engineers); err != nil {
        http.Error(w, "template error", http.StatusInternalServerError)
        log.Println("template execute error:", err)
    }
}

func addEngineer(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    if err := r.ParseForm(); err != nil {
        http.Error(w, "bad request", http.StatusBadRequest)
        return
    }

    name := strings.TrimSpace(r.FormValue("name"))
    role := strings.TrimSpace(r.FormValue("role"))
    if name == "" || role == "" {
        http.Error(w, "name and role are required", http.StatusBadRequest)
        return
    }

    e := Engineer{ID: nextID, Name: name, Role: role}
    nextID++
    engineers = append(engineers, e)

    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fmt.Fprintf(w, "<li id=\"eng-%d\">%s - %s <button type=\"button\" hx-delete=\"/engineers/%d\" hx-target=\"closest li\" hx-swap=\"outerHTML\">Delete</button></li>",
     e.ID,
     template.HTMLEscapeString(e.Name),
     template.HTMLEscapeString(e.Role),
     e.ID)
}

func deleteEngineer(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    idStr := strings.TrimPrefix(r.URL.Path, "/engineers/")
    if idStr == "" || idStr == "/engineers/" {
        http.Error(w, "missing id", http.StatusBadRequest)
        return
    }

    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "invalid id", http.StatusBadRequest)
        return
    }

    for i, e := range engineers {
        if e.ID == id {
            engineers = append(engineers[:i], engineers[i+1:]...)
            break
        }
    }

    w.WriteHeader(http.StatusOK)
}

const pageHTML = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <script src="https://unpkg.com/htmx.org@2"></script>
</head>
<body>
<h1>Engineers</h1>

<form hx-post="/engineers" hx-target="#list" hx-swap="beforeend">
  <input type="text" name="name" placeholder="Name" required>
  <input type="text" name="role" placeholder="Role" required>
  <button type="submit">Add</button>
</form>

<ul id="list">
{{range .}}
  <li id="eng-{{.ID}}">
    {{.Name}} - {{.Role}}
    <button type="button" hx-delete="/engineers/{{.ID}}" hx-target="closest li" hx-swap="outerHTML">Delete</button>
  </li>
{{end}}
</ul>
</body>
</html>`

//
// func postNewEngineer {
//     return fmt.sPrintF("`<li>`")
// }
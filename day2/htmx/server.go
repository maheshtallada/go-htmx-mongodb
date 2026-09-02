package main

import (
    "html/template"
    "net/http"
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
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        tmpl.ExecuteTemplate(w, "page", engineers)
    })
    // TODO B2: POST /engineers → append + return the new <li> fragment
    // TODO B3: DELETE /engineers/{id} → remove + return empty (200 OK + outerHTML swap removes the row)
    http.ListenAndServe(":8080", nil)
}

const pageHTML = `<!DOCTYPE html>
<html>
<head><script src="https://unpkg.com/htmx.org@2"></script></head>
<body>
<h1>Engineers</h1>
<ul id="list">
{{range .}}
  <li id="eng-{{.ID}}">
    {{.Name}} - {{.Role}}
    <!-- TODO B1: add form that requests in new engineer -->
  </li>
{{end}}
</ul>
</body>
</html>`
package controllers

import (
	"fmt"
	"html/template"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	tmpl := template.Must(template.ParseFiles(
		"views/index.html", "views/footer.html",
	))

	auth := false
	title := ""
	stmt := ""
	username := ""

	// login
	if session.Values["auth"] == true {
		if session.Values["username"] == "admin" {
			auth = true
			title = "Admin Page"
			stmt = "Welcome back, "
			username = session.Values["username"].(string)
		} else {
			auth = true
			title = "My Page"
			stmt = "Hello, "
			username = session.Values["username"].(string)
		}
	} else {
		title = "Top Page"
		stmt = "Hello, Guest"
	}

	templateData := TemplateDataIndex{
		auth,
		title,
		stmt,
		username,
	}
	if err := tmpl.ExecuteTemplate(w, "index", templateData); err != nil {
		fmt.Println(err)
	}
}

type TemplateDataIndex struct {
	Auth      bool
	Title     string
	Statement string
	Username  string
}

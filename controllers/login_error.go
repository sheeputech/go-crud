package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

func LoginError(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	if session.Values["auth"] == true {
		http.Redirect(w, r, "/", 303)
	} else {
		tmpl := template.Must(template.ParseFiles(
			"views/login_error.html", "views/footer.html",
		))

		title := "You aren't authorized."
		datetime := fmt.Sprint(time.Now().Format("2006-01-02 15:04:05"))

		templateData := TemplateDataLinErr{
			title,
			datetime,
		}
		if err := tmpl.ExecuteTemplate(w, "login_error", templateData); err != nil {
			fmt.Println(err)
		}
	}
}

type TemplateDataLinErr struct {
	Title    string
	Datetime string
}

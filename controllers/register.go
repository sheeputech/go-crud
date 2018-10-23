package controllers

import (
	"../models"
	"fmt"
	"html/template"
	"net/http"
	"sync"
	"time"
)

func Register(w http.ResponseWriter, r *http.Request) {
	s := time.Now()
	wg := new(sync.WaitGroup)
	session, _ := store.Get(r, "session-name")
	if session.Values["auth"] == true {
		tmpl := template.Must(template.ParseFiles(
			"views/register.html", "views/footer.html"))

		title := "Register"
		datetime := fmt.Sprint(time.Now().Format("2006-01-02 15:04:05"))
		username := session.Values["username"].(string)
		userId := session.Values["user_id"].(int)
		charstr := r.FormValue("charstr")

		var stmterr string

		if r.FormValue("reg") == "Save" {
			if len(charstr) > 0 && len(charstr) <= 200 {
				wg.Add(1)
				go func() {
					models.CString(userId, charstr)
					wg.Done()
				}()
				stmterr = "Registration Succeeded."
			} else if len(charstr) > 200 {
				stmterr = "That String is too long."
			} else {
				stmterr = "You can't register a blank text."
			}
		}

		templateData := TemplateDataRegister{
			title,
			datetime,
			username,
			stmterr,
		}

		if err := tmpl.ExecuteTemplate(w, "register", templateData); err != nil {
			fmt.Println(err)
		}
	} else {
		http.Redirect(w, r, "/login_error", 401)
	}
	wg.Wait()
	e := time.Now()
	fmt.Println(e.Sub(s).Seconds())
}

type TemplateDataRegister struct {
	Title    string
	Datetime string
	Username string
	StmtErr  string
}

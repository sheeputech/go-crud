package controllers

import (
	"../models"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

func List(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")

	if session.Values["auth"] == true {
		tmpl := template.Must(template.ParseFiles(
			"views/list.html", "views/footer.html",
		))

		title := "List"
		datetime := fmt.Sprint(time.Now().Format("2006-01-02 15:04:05"))
		username := session.Values["username"].(string)

		userId := session.Values["user_id"].(int)
		charstr := r.FormValue("updName" + r.FormValue("update"))

		var stmterr string
		var values map[int]string

		if updId, err := strconv.Atoi(r.FormValue("update")); err == nil {
			if len(charstr) > 0 && len(charstr) <= 200 {
				models.UString(updId, charstr)
				stmterr = "Registration Succeeded."
			} else if len(charstr) > 200 {
				stmterr = "That String is too long."
			} else {
				stmterr = "You can't register a blank text."
			}
		} else if delId, err := strconv.Atoi(r.FormValue("del")); err == nil {
			models.DString(delId)
			stmterr = "Delete Succeeded."
		} else {
			values, stmterr = models.ResearchString(userId)
		}
		values, _ = models.ResearchString(userId)

		templateData := TemplateDataList{
			title,
			datetime,
			values,
			username,
			stmterr,
		}
		if err := tmpl.ExecuteTemplate(w, "list", templateData); err != nil {
			fmt.Println(err)
		}
	} else {
		http.Redirect(w, r, "/login_error", 401)
	}
}

type TemplateDataList struct {
	Title    string
	Datetime string
	Values   map[int]string
	Username string
	StmtErr  string
}

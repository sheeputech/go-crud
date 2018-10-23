package controllers

import (
	"../models"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("views/signup.html", "views/footer.html"))

	title := "Sign Up"
	datetime := fmt.Sprint(time.Now().Format("2006-01-02 15:04:05"))
	doSignup, _ := strconv.ParseBool(r.FormValue("signup"))
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	result, statements := models.SignUp(false, doSignup, username, email, password)
	if result == false {
		username = ""
	}

	templateData := TemplateDataSignup{
		title,
		datetime,
		username,
		statements["stmterr"],
		statements["tooLongUser"],
		statements["tooLongEmail"],
		statements["tooLongPass"],
	}
	if err := tmpl.ExecuteTemplate(w, "signup", templateData); err != nil {
		fmt.Println(err)
	}
}

type TemplateDataSignup struct {
	Title         string
	Datetime      string
	Username      string
	StmtErr       string
	StmtLongUser  string
	StmtLongEmail string
	StmtLongPass  string
}

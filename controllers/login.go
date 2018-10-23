package controllers

import (
	"../models"
	"fmt"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var store = sessions.NewCookieStore([]byte("secret-key"))

func Login(w http.ResponseWriter, r *http.Request) {
	s := time.Now()
	tmpl := template.Must(template.ParseFiles("views/login.html", "views/footer.html"))

	doLogin, _ := strconv.ParseBool(r.FormValue("login"))
	username := r.FormValue("username")
	password := r.FormValue("password")

	stmterr := ""
	result := false
	var userId int

	wg := new(sync.WaitGroup)

	if doLogin == true {

		count, lastFailedTime := models.GetLoginFailCount(username, password)

		lo := "2006-01-02 15:04:05"
		loc, _ := time.LoadLocation("Asia/Tokyo")
		lft, _ := time.ParseInLocation(lo, lastFailedTime, loc)
		dur := time.Now().Sub(lft).Seconds()

		if count < 4 || (count >= 4 && count <= 10 && dur > 60) {
			result, userId = models.Login(username, password)
			if result == true {

				wg.Add(1)
				go func() {
					models.RefreshFailCntTemp(username, password)
					wg.Done()
				}()
				//models.RefreshFailCntTemp(username, password)

				session, err := store.Get(r, "session-name")
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				hashPass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
				session.Values["auth"] = true
				session.Values["user_id"] = userId
				session.Values["username"] = username
				session.Values["password"] = hashPass
				session.Options = &sessions.Options{
					MaxAge:   10800,
					HttpOnly: true,
				}
				session.Save(r, w)

				http.Redirect(w, r, "/", 303)
			} else {
				wg.Add(1)
				go func() {
					models.AddLoginFailCount(username, password)
					wg.Done()
				}()
				//models.AddLoginFailCount(username, password)
				stmterr = "You failed to login."
			}
		} else if count >= 4 && count < 10 && dur <= 60 {
			stmterr = "You failed to login more than predetermined number of times. Please wait for a while and try again."
		} else {
			stmterr = "Your account is currently frozen. We have sent you an URL for unfreezing to your email address. Please check it."
		}
	}

	templateData := TemplateDataLogin{
		"Login",
		fmt.Sprint(time.Now().Format("2006-01-02 15:04:05")),
		username,
		stmterr,
		result,
	}
	if err := tmpl.ExecuteTemplate(w, "login", templateData); err != nil {
		fmt.Println(err)
	}
	wg.Wait()
	e := time.Now()
	fmt.Println(e.Sub(s).Seconds())
}

type TemplateDataLogin struct {
	Title    string
	Datetime string
	Username string
	StmtErr  string
	Auth     bool
}

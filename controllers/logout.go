package controllers

import (
	"github.com/gorilla/sessions"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	session.Values["auth"] = false
	session.Options = &sessions.Options{
		MaxAge:   10800,
		HttpOnly: true,
	}
	session.Save(r, w)
	http.Redirect(w, r, "/", 303)
}

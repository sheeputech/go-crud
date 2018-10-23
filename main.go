package main

import (
	"./controllers"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
)

func main() {
	l, err := net.Listen("tcp", ":9000")
	if err != nil {
		return
	}
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.HandleFunc("/", controllers.Index)
	http.HandleFunc("/register", controllers.Register)
	http.HandleFunc("/list", controllers.List)
	http.HandleFunc("/signup", controllers.Signup)
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/logout", controllers.Logout)
	http.HandleFunc("/login_error", controllers.LoginError)
	log.Printf("Start Go HTTP Server")
	fcgi.Serve(l, nil)
}

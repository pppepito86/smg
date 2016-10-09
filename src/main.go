package main

import (
	"db"
	"fmt"
	"log"
	"net/http"
	"request"
	"submissions"
	"text/template"
)

func main() {
	db.OpenConnection()
	defer db.Close()

	go submissions.Checker()
	fmt.Println("server started")
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("../assets/"))))

	http.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe(":80", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	user := request.GetUser(*r)
	fmt.Println(r.URL, user.Id, r.RemoteAddr)

	if r.URL.Path == "/error.html" {
		errorHtml(w, r)
	} else if user.RoleName == "admin" {
		request.HandleAdmin(w, r, user)
	} else if user.RoleName == "teacher" {
		request.HandleTeacher(w, r, user)
	} else if user.RoleName == "user" {
		request.HandleUser(w, r, user)
	} else {
		request.HandleGuest(w, r)
	}
}

func errorHtml(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query()["error"]
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../error.html")
	t.Execute(w, msg)
}

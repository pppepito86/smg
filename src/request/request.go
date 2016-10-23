package request

import (
	"fmt"
	"net/http"
	"request/contest"
	"strings"
	"text/template"
)

func Route(w http.ResponseWriter, r *http.Request) {
	user := GetUser(*r)
	fmt.Println(r.URL, user.Id, r.RemoteAddr)

	if strings.Index(r.URL.Path, "/contest/") == 0 {
		contest.Route(w, r, user)
		return
	}

	if r.URL.Path == "/error.html" {
		errorHtml(w, r)
	} else if user.RoleName == "admin" {
		HandleAdmin(w, r, user)
	} else if user.RoleName == "teacher" {
		HandleTeacher(w, r, user)
	} else if user.RoleName == "user" {
		HandleUser(w, r, user)
	} else {
		HandleGuest(w, r)
	}
}

func errorHtml(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query()["error"]
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../error.html")
	t.Execute(w, msg)
}

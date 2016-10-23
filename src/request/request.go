package request

import (
	"fmt"
	"net/http"
	"text/template"
)

type InputValidator interface {
	Validate() error
}

type RequestHandler interface {
	Execute() error
}

func Route(w http.ResponseWriter, r *http.Request) {
	user := GetUser(*r)
	fmt.Println(r.URL, user.Id, r.RemoteAddr)

	/*
		if strings.Contains(r.URL.Path, "contest") && strings.Contains(r.URL.Path, "problem") {
			cid, _ := strconv.ParseInt(strings.Split(r.URL.Path, "/")[2], 10, 64)
			handler := contest.ProblemsHandler{
				R:    r,
				W:    w,
				User: user,
				Cid:  cid,
			}
			handler.Execute()
			return
		}
	*/

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

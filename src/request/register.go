package request

import (
	"db"
	"net/http"
	"strings"
)

func register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	errMsg := ""

	checkFieldsPresent([]string{"username", "firstname", "lastname", "email", "password", "passwordagain"}, r, &errMsg)
	if errMsg == "" {
		checkFieldsLength([]string{"username", "firstname", "lastname", "email", "password", "passwordagain"}, r, &errMsg)
	}
	if errMsg == "" {
		checkPasswordMatch("password", "passwordagain", r, &errMsg)
	}
	if errMsg != "" {
		http.Redirect(w, r, "/error.html?error="+errMsg, http.StatusFound)
	}

	user := db.User{
		-1,
		3,
		r.Form["username"][0],
		r.Form["firstname"][0],
		r.Form["lastname"][0],
		r.Form["email"][0],
		r.Form["password"][0],
		r.Form["password"][0],
		false,
		"",
	}
	user, err := db.CreateUser(user)
	if err != nil {
		http.Redirect(w, r, "/error.html?error="+err.Error(), http.StatusFound)
	}

	http.Redirect(w, r, "/index.html", http.StatusFound)
}

func checkFieldsPresent(fields []string, r *http.Request, err *string) {
	for _, f := range fields {
		if len(r.Form[f]) != 1 {
			*err = "There are missing registration fields"
			return
		}
	}
}

func checkFieldsLength(fields []string, r *http.Request, err *string) {
	for _, f := range fields {
		field := strings.TrimSpace(r.Form[f][0])
		if len(field) == 0 {
			*err = "You have not entered all registration fields"
			return
		}
	}
}

func checkPasswordMatch(pass1, pass2 string, r *http.Request, err *string) {
	p1 := r.Form[pass1][0]
	p2 := r.Form[pass1][0]
	if p1 != p2 {
		*err = "Passwords does not match"
	}
}

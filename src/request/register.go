package request

import (
	"db"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os/exec"
	"strings"

	"golang.org/x/crypto/bcrypt"
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
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.Form["password"][0]), bcrypt.DefaultCost)
	if err != nil {
		http.Redirect(w, r, "/error.html?error="+err.Error(), http.StatusFound)
		return
	}

	validationCodeByte, _ := bcrypt.GenerateFromPassword([]byte(r.Form["username"][0]+"-"+r.Form["email"][0]), bcrypt.DefaultCost)
	validationCode := url.QueryEscape(string(validationCodeByte))
	user := db.User{
		-1,
		3,
		r.Form["username"][0],
		r.Form["firstname"][0],
		r.Form["lastname"][0],
		r.Form["email"][0],
		string(hashedPassword),
		string(hashedPassword),
		false,
		"",
		validationCode,
	}
	user, err = db.CreateUser(user)
	if err != nil {
		http.Redirect(w, r, "/error.html?error="+err.Error(), http.StatusFound)
		return
	}
	db.CreateUserGroup(user.Id, 1)

	sendEmail(user.Email, user.ValidationCode)

	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../registered.html")
	t.Execute(w, nil)
}

func sendEmail(receiver, code string) {
	subject := "$(echo -e \"Account Validation at pesho.org\nContent-Type: text/html\")"

	sender := "registration@pesho.org"
	body := fmt.Sprintf("To validate your account at pesho.org, please click <a target=\"_blank\" href=\\\"http://www.pesho.org/emailvalidation?code=%s\\\">here</a>"+
		"<br>or just go to http://www.pesho.org/emailvalidation?code=%s", code, code)
	message := fmt.Sprintf("echo \"%s\" | mail -s \"%s\" -aFrom:\\<%s\\> %s", body, subject, sender, receiver)
	fmt.Println("message is: " + message)
	err := exec.Command("/bin/bash", "-c", message).Run()
	if err != nil {
		println(err.Error())
	}
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

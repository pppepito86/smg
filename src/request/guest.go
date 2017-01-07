package request

import (
	"crypto/rand"
	"db"
	"encoding/base64"
	"fmt"
	"net/http"
	"session"

	"golang.org/x/crypto/bcrypt"
)

func HandleGuest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/login" && r.Method == "POST" {
		login(w, r)
	} else if path == "/register" && r.Method == "POST" {
		register(w, r)
	} else if path == "/login.html" {
		http.ServeFile(w, r, "../login.html")
	} else {
		http.Redirect(w, r, "/login.html", http.StatusFound)
	}
}

func getSessionIdCookie(r http.Request) *http.Cookie {
	for _, cookie := range r.Cookies() {
		if cookie.Name == "session.id" {
			return cookie
		}
	}
	return new(http.Cookie)
}

func login(w http.ResponseWriter, r *http.Request) {
	if postLogin(w, r) {
		http.Redirect(w, r, "/myassignments.html", http.StatusFound)
	} else {
		http.Redirect(w, r, "/error.html?error=Login Failed", http.StatusFound)
	}
}

func postLogin(w http.ResponseWriter, r *http.Request) bool {
	r.ParseForm()
	user := r.Form["username"]
	pass := r.Form["password"]
	if len(user) != 1 || len(pass) != 1 {
		return false
	}
	return authenticate(w, *r, user[0], pass[0])
}

func authenticate(w http.ResponseWriter, r http.Request, username string, password string) bool {
	user, _ := db.GetUser(username)
	if user.UserName == "" || user.ValidationCode != "" {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordSalt), []byte(password))
	if err != nil {
		return false
	}
	val := fmt.Sprintf("%s-%s", user.UserName, RandomString())
	cookie := http.Cookie{Name: "session.id", Value: val}
	http.SetCookie(w, &cookie)
	session.SetAttribute(val, username)
	return true
}

func logout(w http.ResponseWriter, r *http.Request) {
	cookie := getSessionIdCookie(*r)
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
	session.RemoveAttribute(cookie.Value)
	http.Redirect(w, r, "/login.html", http.StatusFound)
}

func RandomString() string {
	size := 32
	rb := make([]byte, size)
	rand.Read(rb)
	rs := base64.URLEncoding.EncodeToString(rb)
	return rs
}

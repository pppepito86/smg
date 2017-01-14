package request

import (
	"crypto/rand"
	"db"
	"encoding/base64"
	"errors"
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
	if err := postLogin(w, r); err == nil {
		http.Redirect(w, r, "/myassignments.html", http.StatusFound)
	} else {
		http.Redirect(w, r, "/error.html?error="+err.Error(), http.StatusFound)
	}
}

func postLogin(w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()
	user := r.Form["username"]
	pass := r.Form["password"]
	if len(user) != 1 || len(pass) != 1 {
		return errors.New("Credentials not supplied")
	}
	return authenticate(w, *r, user[0], pass[0])
}

func authenticate(w http.ResponseWriter, r http.Request, username string, password string) error {
	user, _ := db.GetUser(username)
	if user.UserName == "" {
		return errors.New("Incorrect username or password")
	}
	if user.ValidationCode != "" {
		return errors.New("Username not activated! Plese check your email to activate.")
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordSalt), []byte(password))
	if err != nil {
		return errors.New("Incorrect username or password")
	}
	val := fmt.Sprintf("%s-%s", user.UserName, RandomString())
	cookie := http.Cookie{Name: "session.id", Value: val}
	http.SetCookie(w, &cookie)
	session.SetAttribute(val, username)
	return nil
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

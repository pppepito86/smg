package request

import (
	"db"
	"net/http"
	"session"
	"strings"
)

func GetUser(r http.Request) db.User {
	cookie := getSessionIdCookie(r)
	if session.ContainsKey(cookie.Value) {
		username := strings.Split(cookie.Value, "-")[0]
		user, _ := db.GetUser(username)
		return user
	} else {
		return db.User{}
	}
}

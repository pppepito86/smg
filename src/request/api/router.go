package api

import (
	"db"
	"net/http"
	"request/util"
	"strings"
)

func Route(w http.ResponseWriter, r *http.Request, user db.User) {
	path := r.URL.Path

	split := strings.Split(strings.TrimPrefix(path, "/api/v1/"), "/")
	if len(split) < 2 {
		//assignmentsAdminHtml(w, r)
		//TODO return error page
		return
	}

	page := split[0]

	var handler util.RequestHandler
	if page == "users" {
		UsersHandler(w, r, user, split[1:])
	}

	if handler == nil {
		return
	}

}

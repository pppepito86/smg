package api

import (
	"db"
	"encoding/json"
	"net/http"
	"strconv"
)

func UsersHandler(w http.ResponseWriter, r *http.Request, user db.User, path []string) {

	if len(path) == 0 {
		//assignmentsAdminHtml(w, r)
		//TODO return error page
		return
	}

	userId, err := strconv.ParseInt(path[0], 10, 64)
	if err != nil {
		if path[0] == "me" {
			userId = user.Id
		} else {
			// TODO return error
			return
		}
	}

	method := path[1]

	w.Header().Set("Content-Type", "application/json")

	if user.RoleName == "admin" || user.RoleName == "teacher" || user.Id == userId {

		u, err := db.GetUserById(userId)
		if err != nil {
			return
		}

		if method == "assignments" {
			var assignments []db.Assignment
			if u.RoleName == "admin" || u.RoleName == "teacher" {
				assignments, _ = db.ListAssignmentsForAuthor(u.Id)
			} else {
				assignments, _ = db.ListAssignmentsForUser(u)
			}
			js, _ := json.Marshal(assignments)
			w.Write(js)

		} else if method == "role" {
			js, _ := json.Marshal(u.RoleName)
			w.Write(js)
		}
	}

}

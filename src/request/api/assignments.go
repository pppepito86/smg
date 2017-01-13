package api

import (
	"db"
	"encoding/json"
	"net/http"
)

func AssignmentsHandler(w http.ResponseWriter, r *http.Request, user db.User, path []string) {

	if user.RoleName == "admin" || user.RoleName == "teacher" {
		w.Header().Set("Content-Type", "application/json")

		assignments, _ := db.ListAssignments()

		js, _ := json.Marshal(assignments)

		w.Write(js)
	}

}

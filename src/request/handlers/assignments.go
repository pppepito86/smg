package handlers

import (
	"db"
	"request/util"
)

type AssignmentsHandler struct {
	util.NoInputValidator
	util.RequestInfo
}

func (h *AssignmentsHandler) Execute() error {
	var assignments []db.Assignment

	if h.User.RoleName == "admin" || h.User.RoleName == "teacher" {
		assignments, _ = db.ListAssignments()
	} else {
		assignments, _ = db.ListAssignmentsForUser(h.User)
	}
	util.ServeHtml(h.W, h.User, "assignments.html", assignments)

	return nil
}

package handlers

import (
	"db"
	"request/util"
)

type AllAssignmentsHandler struct {
	util.NoInputValidator
	util.RequestInfo
}

func (h *AllAssignmentsHandler) Execute() error {
	var assignments []db.Assignment

	if h.User.RoleName != "admin" && h.User.RoleName != "teacher" {
		return nil
	}

	assignments, _ = db.ListAssignments()
	util.ServeHtml(h.W, h.User, "assignments.html", assignments)

	return nil
}

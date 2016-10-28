package handlers

import (
	"db"
	"request/util"
)

type MyProblemsHandler struct {
	util.NoInputValidator
	util.RequestInfo
}

func (h *MyProblemsHandler) Execute() error {
	if h.User.RoleName != "admin" && h.User.RoleName != "teacher" {
		return nil
	}

	problems, _ := db.ListProblemsWithAuthor(h.User.Id)
	util.ServeHtml(h.W, h.User, "problems.html", problems)

	return nil
}

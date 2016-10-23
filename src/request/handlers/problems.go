package handlers

import (
	"db"
	"request/util"
)

type ProblemsHandler struct {
	util.NoInputValidator
	util.RequestInfo
}

func (h *ProblemsHandler) Execute() error {
	if h.User.RoleName != "admin" {
		return nil
	}

	problems, _ := db.ListProblems()
	util.ServeHtml(h.W, h.User, "problems.html", problems)

	return nil
}

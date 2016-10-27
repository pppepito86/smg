package handlers

import (
	"db"
	"request/util"
	"strconv"
)

type ProblemHandler struct {
	util.NoInputValidator
	util.RequestInfo
}

func (h *ProblemHandler) Execute() error {
	if h.User.RoleName != "admin" && h.User.RoleName != "teacher" {
		return nil
	}

	id, _ := strconv.ParseInt(h.R.URL.Query()["id"][0], 10, 64)
	problem, _ := db.GetProblem(id)
	problem.LangLimits = util.LimitsFromString(problem.Languages)
	isAuthor := problem.AuthorId == h.User.Id
	util.ServeHtmlWithAuthor(h.W, h.User, "problem.html", problem, isAuthor)

	return nil
}

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
	if h.User.RoleName != "admin" {
		return nil
	}

	id, _ := strconv.ParseInt(h.R.URL.Query()["id"][0], 10, 64)
	problem, _ := db.GetProblem(id)
	problem.LangLimits = util.LimitsFromString(problem.Languages)
	util.ServeHtml(h.W, h.User, "problem.html", problem)

	return nil
}

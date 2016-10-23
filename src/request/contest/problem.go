package contest

import (
	"db"
	"net/http"
	"request/util"
	"strconv"
)

type ProblemHandler struct {
	util.NoInputValidator

	R    *http.Request
	W    http.ResponseWriter
	User db.User
	Cid  int64
	Args []string
}

func (h *ProblemHandler) Execute() error {
	apId, _ := strconv.ParseInt(h.Args[0], 10, 64)
	ap, _ := db.GetAssignmentProblem(apId)
	problem, _ := db.GetProblem(ap.ProblemId)
	problem.LangLimits = util.LimitsFromString(problem.Languages)
	response := util.Response{h.Cid, problem, ""}
	util.ServeContestHtml(h.W, h.R, h.User, "problem.html", response)

	return nil
}

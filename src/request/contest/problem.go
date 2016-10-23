package contest

import (
	"db"
	"request/util"
	"strconv"
)

type ProblemHandler struct {
	util.NoInputValidator
	ContestRequestInfo
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

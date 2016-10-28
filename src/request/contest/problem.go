package contest

import (
	"db"
	"net/http"
	"request/util"
	"strconv"
)

type ProblemHandler struct {
	util.NoInputValidator
	ContestRequestInfo
}

func (h *ProblemHandler) Execute() error {
	if !util.IsUserAssignedToContest(h.User, h.Assignment) {
		http.Redirect(h.W, h.R, "/error.html?error=\"You are not allowed to access this assignment\"", http.StatusFound)
		return nil
	}

	apId, _ := strconv.ParseInt(h.Args[0], 10, 64)
	ap, _ := db.GetAssignmentProblem(apId)
	problem, _ := db.GetProblem(ap.ProblemId)
	problem.LangLimits = util.LimitsFromString(problem.Languages)
	ServeContestHtml(h.ContestRequestInfo, "problem.html", problem)

	return nil
}

package contest

import (
	"db"
	"request/util"
	"strconv"
)

type EditHandler struct {
	util.NoInputValidator
	ContestRequestInfo
}

func (h *EditHandler) Execute() error {
	if h.User.RoleName != "admin" {
		return nil
	}

	assignment, _ := db.ListAssignment(h.Cid)
	aps, _ := db.ListAssignmentProblems(h.Cid)
	problems := ""
	for _, ap := range aps {
		problems += "," + strconv.FormatInt(ap.ProblemId, 10)
	}
	if len(problems) > 1 {
		problems = problems[1:len(problems)]
	}
	assignment.Problems = problems
	response := util.Response{h.Cid, assignment, ""}
	util.ServeContestHtml(h.W, h.R, h.User, "edit.html", response)

	return nil
}

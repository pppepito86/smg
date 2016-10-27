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
	if h.User.RoleName != "admin" && h.User.RoleName != "teacher" {
		return nil
	}

	assignment, _ := db.ListAssignment(h.Cid)
	if h.User.RoleName != "admin" && assignment.AuthorId != h.User.Id {
		return nil
	}

	aps, _ := db.ListAssignmentProblems(h.Cid)
	problems := ""
	for _, ap := range aps {
		problems += "," + strconv.FormatInt(ap.ProblemId, 10)
	}
	if len(problems) > 1 {
		problems = problems[1:len(problems)]
	}
	assignment.Problems = problems
	ServeContestHtml(h.ContestRequestInfo, "edit.html", assignment)

	return nil
}

package contest

import (
	"db"
	"request/util"
)

type SubmissionsHandler struct {
	util.NoInputValidator
	ContestRequestInfo
}

func (h *SubmissionsHandler) Execute() error {
	if h.User.RoleName != "admin" && h.User.RoleName != "teacher" {
		return nil
	}

	submissions, _ := db.ListSubmissionsForAssignment(h.Cid)
	ServeContestHtml(h.ContestRequestInfo, "submissions.html", submissions)

	return nil
}

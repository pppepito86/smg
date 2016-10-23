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
	if h.User.RoleName != "admin" {
		return nil
	}

	submissions, _ := db.ListSubmissionsForAssignment(h.Cid)
	response := util.Response{h.Cid, submissions, ""}
	util.ServeContestHtml(h.W, h.R, h.User, "submissions.html", response)

	return nil
}

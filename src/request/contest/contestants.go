package contest

import (
	"db"
	"request/util"
)

type ContestantsHandler struct {
	util.NoInputValidator
	ContestRequestInfo
}

func (h *ContestantsHandler) Execute() error {
	if h.User.RoleName != "admin" && h.User.RoleName != "teacher" {
		return nil
	}

	users, _ := db.ListUsersForAssignment(h.Cid)
	ServeContestHtml(h.ContestRequestInfo, "contestants.html", users)

	return nil
}

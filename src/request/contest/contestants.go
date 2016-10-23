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
	if h.User.RoleName != "admin" {
		return nil
	}

	users, _ := db.ListUsersForAssignment(h.Cid)
	response := util.Response{h.Cid, users, ""}
	util.ServeContestHtml(h.W, h.R, h.User, "contestants.html", response)

	return nil
}

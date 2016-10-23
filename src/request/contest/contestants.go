package contest

import (
	"db"
	"net/http"
	"request/util"
)

type ContestantsHandler struct {
	util.NoInputValidator
	ContestRequestInfo
}

func (h *ContestantsHandler) Execute() error {
	if !util.IsUserAssignedToContest(h.User, h.Cid) {
		http.Redirect(h.W, h.R, "/error.html?error=\"You are not allowed to access this assignment\"", http.StatusFound)
		return nil
	}

	users, _ := db.ListUsersForAssignment(h.Cid)
	response := util.Response{h.Cid, users, ""}
	util.ServeContestHtml(h.W, h.R, h.User, "contestants.html", response)

	return nil
}

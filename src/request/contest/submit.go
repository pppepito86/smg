package contest

import (
	"db"
	"net/http"
	"request/util"
)

type SubmitHandler struct {
	util.NoInputValidator
	ContestRequestInfo
}

func (h *SubmitHandler) Execute() error {
	if !util.IsUserAssignedToContest(h.User, h.Cid) {
		http.Redirect(h.W, h.R, "/error.html?error=\"You are not allowed to access this assignment\"", http.StatusFound)
		return nil
	}

	aps, _ := db.ListAssignmentProblems(h.Cid)
	response := util.Response{h.Cid, aps, ""}
	util.ServeContestHtml(h.W, h.R, h.User, "submitcode.html", response)

	return nil
}

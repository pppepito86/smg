package contest

import (
	"db"
	"net/http"
	"request/util"
)

type SubmissionsHandler struct {
	util.NoInputValidator
	ContestRequestInfo
}

func (h *SubmissionsHandler) Execute() error {
	if !util.IsUserAssignedToContest(h.User, h.Cid) {
		http.Redirect(h.W, h.R, "/error.html?error=\"You are not allowed to access this assignment\"", http.StatusFound)
		return nil
	}
	submissions, _ := db.ListSubmissionsForAssignment(h.Cid)
	response := util.Response{h.Cid, submissions, ""}
	util.ServeContestHtml(h.W, h.R, h.User, "submissions.html", response)

	return nil
}

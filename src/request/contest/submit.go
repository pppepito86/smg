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
	if !util.IsUserAssignedToContest(h.User, h.Assignment) {
		http.Redirect(h.W, h.R, "/error.html?error=\"You are not allowed to access this assignment\"", http.StatusFound)
		return nil
	}
	if h.User.RoleName != "admin" && h.User.RoleName != "teacher" && !h.Assignment.IsActive {
		http.Redirect(h.W, h.R, "/error.html?error=\"You are not allowed to submit in this assignment\"", http.StatusFound)
		return nil
	}

	aps, _ := db.ListAssignmentProblems(h.Cid)
	ServeContestHtml(h.ContestRequestInfo, "submitcode.html", aps)

	return nil
}

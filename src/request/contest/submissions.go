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
	if !util.IsUserAssignedToContest(h.User, h.Assignment) {
		http.Redirect(h.W, h.R, "/error.html?error=\"You are not allowed to access this assignment\"", http.StatusFound)
		return nil
	}
	if h.User.RoleName != "admin" && h.User.RoleName != "teacher" && !h.Assignment.HasFinished {
		http.Redirect(h.W, h.R, "/error.html?error=\"You are not allowed to submit in this assignment\"", http.StatusFound)
		return nil
	}
	var submissions []db.Submission
	if h.User.RoleName != "admin" && h.User.RoleName != "teacher" {
		submissions, _ = db.ListAcceptedSubmissionsForAssignment(h.Cid)
	} else {
		submissions, _ = db.ListSubmissionsForAssignment(h.Cid)
	}
	ServeContestHtml(h.ContestRequestInfo, "submissions.html", submissions)

	return nil
}

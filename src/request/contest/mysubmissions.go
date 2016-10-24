package contest

import (
	"db"
	"net/http"
	"request/util"
)

type MySubmissionsHandler struct {
	util.NoInputValidator
	ContestRequestInfo
}

func (h *MySubmissionsHandler) Execute() error {
	if !util.IsUserAssignedToContest(h.User, h.Cid) {
		http.Redirect(h.W, h.R, "/error.html?error=\"You are not allowed to access this assignment\"", http.StatusFound)
		return nil
	}

	mySubmissions, _ := db.ListMySubmissions(h.User.Id, h.Cid)
	ServeContestHtml(h.ContestRequestInfo, "mysubmissions.html", mySubmissions)

	return nil
}

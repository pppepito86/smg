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
	response := util.Response{h.Cid, mySubmissions, ""}
	util.ServeContestHtml(h.W, h.R, h.User, "mysubmissions.html", response)

	return nil
}

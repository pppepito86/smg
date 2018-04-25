package handlers

import (
	"db"
	"net/http"
	"request/util"
	"strconv"
)

type BanSubmissionHandler struct {
	util.NoInputValidator
	util.RequestInfo
}

func (h *BanSubmissionHandler) Execute() error {
	if h.User.RoleName != "admin" && h.User.RoleName != "teacher" {
		return nil
	}

	h.R.ParseForm()

	idStr := h.R.Form["id"][0]
	id, _ := strconv.ParseInt(idStr, 10, 64)
	submission, _ := db.ListSubmission(id)

	contestIdStr := strconv.FormatInt(submission.AssignmentId, 10)

	db.BanSubmission(id)

	http.Redirect(h.W, h.R, "/contest/" + contestIdStr + "/submission/" + idStr, http.StatusFound)
	return nil
}

package handlers

import (
	"db"
	"request/util"
	"strconv"
)

type AllSubmissionsHandler struct {
	util.NoInputValidator
	util.RequestInfo
}

func (h *AllSubmissionsHandler) Execute() error {
	if h.User.RoleName != "admin" {
		return nil
	}

	var MAX_LIMIT int = 5000

	h.R.ParseForm()
	limit, _ := strconv.Atoi(h.R.Form["limit"][0])
	if limit > MAX_LIMIT {
		limit = MAX_LIMIT
	}
	submissions, _ := db.ListSubmissions(limit);
	util.ServeHtml(h.W, h.User, "all-submissions.html", submissions)

	return nil
}
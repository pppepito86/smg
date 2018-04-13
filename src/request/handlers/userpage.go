package handlers

import (
	"request/util"
	"strconv"
	"db"
)

type UserPageHandler struct {
	util.NoInputValidator
	util.RequestInfo
}

func (h *UserPageHandler) Execute() error {
	if h.User.RoleName != "admin" && h.User.RoleName != "teacher" {
		return nil
	}

	h.R.ParseForm()
	userId, err := strconv.ParseInt(h.R.Form["id"][0], 10, 64)
	if err != nil {
		return nil
	}

	user, err := db.GetUserById(userId)
	if err != nil {
		return nil
	}

	submissions, err := db.ListMyAllSubmissions(userId)
	if err != nil {
		return nil
	}

	data := struct {
		User db.User
		Submissions []db.Submission
	} {
		user,
		submissions,
	}
	util.ServeHtml(h.W, h.User, "user-page.html", data)

	return nil
}

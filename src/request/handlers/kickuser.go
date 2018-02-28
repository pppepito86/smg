package handlers

import (
	"db"
	"net/http"
	"request/util"
	"strconv"
)

type KickUserHandler struct {
	util.NoInputValidator
	util.RequestInfo
}

func (h *KickUserHandler) Execute() error {
	if h.User.RoleName != "admin" && h.User.RoleName != "teacher" {
		return nil
	}

	h.R.ParseForm()
	user := h.R.Form["userid"]
	group := h.R.Form["groupid"]
	if len(user) != 1 || len(group) != 1 {
		h.W.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	userId, _ := strconv.ParseInt(user[0], 10, 64)
	groupId, _ := strconv.ParseInt(group[0], 10, 64)

	db.RemoveUserGroup(userId, groupId)

	http.Redirect(h.W, h.R, "/groups.html", http.StatusFound)
	return nil
}

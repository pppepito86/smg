package handlers

import (
	"db"
	"net/http"
	"request/util"
	"strconv"
)

type ChangeUserRoleHandler struct {
	util.NoInputValidator
	util.RequestInfo
}

func (h *ChangeUserRoleHandler) Execute() error {
	if h.User.RoleName != "admin" {
		return nil
	}

	h.R.ParseForm()
	user := h.R.Form["userid"]
	role := h.R.Form["roleid"]
	if len(user) != 1 || len(role) != 1 {
		h.W.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	userId, _ := strconv.ParseInt(user[0], 10, 64)
	roleId, _ := strconv.ParseInt(role[0], 10, 64)
	db.UpdateUserRole(userId, roleId)

	return nil
}

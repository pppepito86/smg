package handlers

import (
	"db"
	"request/util"
)

type GroupsHandler struct {
	util.NoInputValidator
	util.RequestInfo
}

func (h *GroupsHandler) Execute() error {
	var groups []db.Group

	if h.User.RoleName == "admin" || h.User.RoleName == "teacher" {
		groups = db.ListGroups()
	} else {
		groups, _ = db.ListGroupsForUser(h.User.Id)
	}

	util.ServeHtml(h.W, h.User, "groups.html", groups)

	return nil
}

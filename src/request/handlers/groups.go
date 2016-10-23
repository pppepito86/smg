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

	if h.User.RoleName == "admin" {
		groups = db.ListGroups()
	} else {
		groups, _ = db.ListGroupsForUser(h.User.Id)
	}

	util.ServeHtml(h.W, h.User, "groups.html", groups)

	return nil
}

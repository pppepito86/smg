package handlers

import (
	"db"
	"net/http"
	"request/util"
)

type JoinGroupHandler struct {
	util.NoInputValidator
	util.RequestInfo
}

func (h *JoinGroupHandler) Execute() error {
	if h.User.RoleName != "user" {
		return nil
	}

	if h.R.Method == "POST" {
		h.executePost()
	} else {
		h.executeGet()
	}

	return nil
}

func (h *JoinGroupHandler) executeGet() {
	util.ServeHtml(h.W, h.User, "joingroup.html", nil)
}

func (h *JoinGroupHandler) executePost() {
	h.R.ParseForm()
	groupName := h.R.Form["groupname"]
	group, _ := db.FindGroupByName(groupName[0])
	db.CreateUserGroup(h.User.Id, group.Id)
	http.Redirect(h.W, h.R, "/groups.html", http.StatusFound)
}

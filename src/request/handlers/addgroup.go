package handlers

import (
	"db"
	"net/http"
	"request/util"
)

type AddGroupHandler struct {
	util.NoInputValidator
	util.RequestInfo
}

func (h *AddGroupHandler) Execute() error {
	if h.User.RoleName != "admin" {
		return nil
	}

	if h.R.Method == "POST" {
		h.executePost()
	} else {
		h.executeGet()
	}

	return nil
}

func (h *AddGroupHandler) executeGet() {
	users, _ := db.ListUsers()
	util.ServeHtml(h.W, h.User, "addgroup.html", users)
}

func (h *AddGroupHandler) executePost() {
	h.R.ParseForm()
	name := h.R.Form["groupname"]
	description := h.R.Form["description"]
	if len(name) != 1 {
		h.W.WriteHeader(http.StatusInternalServerError)
		return
	}
	g := db.Group{
		GroupName:   name[0],
		Description: description[0],
		CreatorId:   h.User.Id,
	}
	g, _ = db.CreateGroup(g)
	http.Redirect(h.W, h.R, "/groups.html", http.StatusFound)

}

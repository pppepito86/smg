package handlers

import (
	"db"
	"request/util"
)

type UsersHandler struct {
	util.NoInputValidator
	util.RequestInfo
}

func (h *UsersHandler) Execute() error {
	if h.User.RoleName != "admin" {
		return nil
	}

	users, _ := db.ListUsers()
	util.ServeHtml(h.W, h.User, "users.html", users)

	return nil
}

type MyUsersHandler struct {
	util.NoInputValidator
	util.RequestInfo
}

func (h *MyUsersHandler) Execute() error {
	if h.User.RoleName != "admin" && h.User.RoleName != "teacher" {
		return nil
	}

	users, _ := db.ListMyUsers(h.User.Id)
	util.ServeHtml(h.W, h.User, "users.html", users)

	return nil
}

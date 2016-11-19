package handlers

import (
	"db"
	"request/util"
    "strconv"
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


type UsersInGroupHandler struct {
	util.NoInputValidator
	util.RequestInfo
}

func (h *UsersInGroupHandler) Execute() error {
	if h.User.RoleName != "admin" && h.User.RoleName != "teacher" {
		return nil
	}
    
    h.R.ParseForm()
	var groupId int64
	var err error
	if len(h.R.Form["id"]) == 1 {
		groupId, err = strconv.ParseInt(h.R.Form["id"][0], 10, 64)
		if err != nil {
			return nil
		}
	} else {
		return nil
	}
    
	users, _ := db.ListUsersInGroup(groupId)
	util.ServeHtml(h.W, h.User, "users.html", users)

	return nil
}
package handlers

import "request/util"

type GroupProgressHandler struct {
	util.NoInputValidator
	util.RequestInfo
}

func (h *GroupProgressHandler) Execute() error {
	if h.User.RoleName != "admin" && h.User.RoleName != "teacher" {
		return nil
	}

	util.ServeHtml(h.W, h.User, "groupprogress.html", nil)

	return nil
}

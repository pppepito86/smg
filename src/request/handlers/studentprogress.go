package handlers

import "request/util"

type StudentProgressHandler struct {
	util.NoInputValidator
	util.RequestInfo
}

func (h *StudentProgressHandler) Execute() error {
	if h.User.RoleName != "admin" {
		return nil
	}

	util.ServeHtml(h.W, h.User, "studentprogress.html", nil)

	return nil
}

package handlers

import (
	"net/http"
	"request/util"
	"session"
)

type LogoutHandler struct {
	util.NoInputValidator
	util.RequestInfo
}

func (h *LogoutHandler) Execute() error {
	cookie := util.GetSessionIdCookie(*h.R)
	cookie.MaxAge = -1
	http.SetCookie(h.W, cookie)
	session.RemoveAttribute(cookie.Value)
	http.Redirect(h.W, h.R, "/login.html", http.StatusFound)

	return nil
}

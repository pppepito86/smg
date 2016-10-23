package handlers

import "request/util"

type DashboardHandler struct {
	util.NoInputValidator
	util.RequestInfo
}

func (h *DashboardHandler) Execute() error {
	util.ServeHtml(h.W, h.User, "dashboard.html", nil)

	return nil
}

package handlers

import (
	"db"
	"net/url"
	"request/util"
)

type ValidateEmailHandler struct {
	util.RequestInfo
}

func (h *ValidateEmailHandler) Execute() error {
	query := h.R.URL.Query()
	validationCode := url.QueryEscape(string(query["code"][0]))
	db.ValidateEmail(validationCode)

	return nil
}

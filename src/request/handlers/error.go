package handlers

import (
	"request/util"
	"text/template"
)

type ErrorHandler struct {
	util.NoInputValidator
	util.RequestInfo
}

func (h *ErrorHandler) Execute() error {
	msg := h.R.URL.Query()["error"]
	h.W.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../error.html")
	t.Execute(h.W, msg)

	return nil
}

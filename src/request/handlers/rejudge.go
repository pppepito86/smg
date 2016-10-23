package handlers

import (
	"db"
	"net/http"
	"request/util"
	"strconv"
	"submissions"
)

type RejudgeHandler struct {
	util.NoInputValidator
	util.RequestInfo
}

func (h *RejudgeHandler) Execute() error {
	if h.User.RoleName != "admin" {
		return nil
	}

	id, _ := strconv.ParseInt(h.R.URL.Query()["id"][0], 10, 64)
	p, _ := db.GetProblem(id)
	limits := util.LimitsFromString(p.Languages)
	ss, _ := db.ListProblemSubmissions(id)
	for _, s := range ss {
		s.Limit = limits[s.Language]

		// FIXME:
		if s.Language == "nodejs" {
			s.Limit = limits["java"]
		}
		submissions.Push(s)
	}

	http.Redirect(h.W, h.R, "/problems", http.StatusFound)

	return nil
}

package contest

import (
	"db"
	"net/http"
	"request/util"
	"strconv"
	"strings"
)

type ContestRequestInfo struct {
	R    *http.Request
	W    http.ResponseWriter
	User db.User
	Cid  int64
	Args []string
}

func Route(w http.ResponseWriter, r *http.Request, user db.User) {
	path := r.URL.Path

	split := strings.Split(strings.TrimPrefix(path, "/contest/"), "/")
	if len(split) < 2 {
		//assignmentsAdminHtml(w, r)
		//TODO return error page
		return
	}

	contestId, err := strconv.ParseInt(split[0], 10, 64)
	if err != nil {
		//assignmentsAdminHtml(w, r)
		//TODO return error page
		return
	}

	page := split[1]
	info := ContestRequestInfo{r, w, user, contestId, split[2:]}
	var handler util.RequestHandler
	if page == "problems" {
		handler = &ProblemsHandler{ContestRequestInfo: info}
	} else if page == "problem" {
		handler = &ProblemHandler{ContestRequestInfo: info}
	} else if page == "submit" {
		handler = &SubmitHandler{ContestRequestInfo: info}
	} else if page == "mysubmissions" {
		handler = &MySubmissionsHandler{ContestRequestInfo: info}
	} else if page == "submissions" {
		handler = &SubmissionsHandler{ContestRequestInfo: info}
	} else if page == "submission" {
		handler = &SubmissionHandler{ContestRequestInfo: info}
	} else if page == "contestants" {
		handler = &ContestantsHandler{ContestRequestInfo: info}
	} else if page == "standings" {
		handler = &StandingsHandler{ContestRequestInfo: info}
	} else if page == "submitcode" {
		handler = &SubmitCodeHandler{ContestRequestInfo: info}
	} else if page == "edit" {
		handler = &EditHandler{ContestRequestInfo: info}
	} else if page == "editassignment" {
		handler = &EditAssignmentHandler{ContestRequestInfo: info}
	}

	handler.Execute()
}

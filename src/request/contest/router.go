package contest

import (
	"db"
	"net/http"
	"request/util"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type ContestRequestInfo struct {
	R          *http.Request
	W          http.ResponseWriter
	User       db.User
	Cid        int64
	Args       []string
	Assignment db.Assignment
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
	assignment, _ := db.ListAssignment(contestId)
	time := time.Now()
	assignment.IsActive = time.After(assignment.StartTime) && time.Before(assignment.EndTime)
	assignment.HasFinished = time.After(assignment.EndTime)

	info := ContestRequestInfo{r, w, user, contestId, split[2:], assignment}
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

	if handler == nil {
		handler = &ProblemsHandler{ContestRequestInfo: info}
	}

	handler.Execute()
}
func ServeContestHtml(info ContestRequestInfo, html string, data interface{}) {
	info.W.Header().Set("Content-Type", "text/html")

	t, _ := template.ParseFiles("../templates/contest/"+html, "../templates/contest/header.html", "../templates/contest/menu.html", "../templates/contest/footer.html")
	response := util.Response{info.Cid, data, info.User.RoleName, false, info.Assignment}
	if info.User.RoleName == "teacher" {
		a, _ := db.ListAssignment(info.Cid)
		response.Author = a.AuthorId == info.User.Id
	}
	t.Execute(info.W, response)
}

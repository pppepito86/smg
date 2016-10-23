package request

import (
	"fmt"
	"net/http"
	"request/contest"
	"request/handlers"
	"request/util"
	"strings"
)

func Route(w http.ResponseWriter, r *http.Request) {
	user := GetUser(*r)
	path := r.URL.Path
	fmt.Println(r.URL, user.Id, r.RemoteAddr)

	if user.RoleName == "" {
		HandleGuest(w, r)
		return
	}

	if strings.Index(path, "/contest/") == 0 {
		contest.Route(w, r, user)
		return
	}

	info := util.RequestInfo{r, w, user}
	var handler util.RequestHandler

	if path == "/problems.html" {
		handler = &handlers.ProblemsHandler{RequestInfo: info}
	} else if path == "/problem.html" {
		handler = &handlers.ProblemHandler{RequestInfo: info}
	} else if path == "/addproblem.html" || path == "/addproblem" {
		handler = &handlers.AddProblemHandler{RequestInfo: info}
	} else if path == "/editproblem.html" || path == "/editproblem" {
		handler = &handlers.EditProblemHandler{RequestInfo: info}
	}

	if path == "/groups.html" {
		handler = &handlers.GroupsHandler{RequestInfo: info}
	} else if path == "/addgroup.html" || path == "/addgroup" {
		handler = &handlers.AddGroupHandler{RequestInfo: info}
	} else if path == "/joingroup.html" || path == "/joingroup" {
		handler = &handlers.JoinGroupHandler{RequestInfo: info}
	}

	if path == "/assignments.html" {
		handler = &handlers.AssignmentsHandler{RequestInfo: info}
	} else if path == "/addassignment.html" || path == "/addassignment" {
		handler = &handlers.AddAssignmentHandler{RequestInfo: info}
	}

	if path == "/users.html" {
		handler = &handlers.UsersHandler{RequestInfo: info}
	} else if path == "/changeuserrole" {
		handler = &handlers.ChangeUserRoleHandler{RequestInfo: info}
	}

	if path == "/dashboard.html" {
		handler = &handlers.DashboardHandler{RequestInfo: info}
	}

	if path == "/rejudge.html" {
		handler = &handlers.RejudgeHandler{RequestInfo: info}
	}

	if path == "/studentprogress.html" {
		handler = &handlers.StudentProgressHandler{RequestInfo: info}
	}
	if path == "/pointsperweek" {
		handler = &handlers.WeekPointsHandler{RequestInfo: info}
	}
	if path == "/logout" {
		handler = &handlers.LogoutHandler{RequestInfo: info}
	}
	if path == "/error.html" {
		handler = &handlers.ErrorHandler{RequestInfo: info}
	}

	if handler == nil {
		if user.RoleName == "admin" {
			handler = &handlers.ProblemsHandler{RequestInfo: info}
		} else {
			handler = &handlers.DashboardHandler{RequestInfo: info}
		}
	}

	handler.Execute()
}

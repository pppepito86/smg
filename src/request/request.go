package request

import (
	"fmt"
	"net/http"
	"request/api"
	"request/contest"
	"request/handlers"
	"request/util"
	"strings"
)

func Route(w http.ResponseWriter, r *http.Request) {
	user := GetUser(*r)
	path := r.URL.Path
	fmt.Println(r.URL, user.Id, r.RemoteAddr)

	if path == "/noi" {
		noiUrl := "https://docs.google.com/forms/d/e/1FAIpQLSdeodrF0AYlnDkHRs2s1YtD-qWkBACv2GZAsk3WVdB7uFUfrw/viewform?usp=send_form"
		http.Redirect(w, r, noiUrl, http.StatusFound)
		return
	}

	if path == "/error.html" {
		info := util.RequestInfo{r, w, user}
		handler := &handlers.ErrorHandler{RequestInfo: info}
		handler.Execute()
		return
	}

	if strings.Index(path, "/emailvalidation") == 0 {
		info := util.RequestInfo{r, w, user}
		handler := &handlers.ValidateEmailHandler{RequestInfo: info}
		handler.Execute()
		return
	} else if strings.Index(path, "/changepassword") == 0 {
	}

	if user.RoleName == "" {
		HandleGuest(w, r)
		return
	}

	if strings.Index(path, "/contest/") == 0 {
		contest.Route(w, r, user)
		return
	}

	if strings.Index(path, "/api/") == 0 {
		api.Route(w, r, user)
		return
	}

	info := util.RequestInfo{r, w, user}
	var handler util.RequestHandler

	if path == "/problems.html" {
		handler = &handlers.ProblemsHandler{RequestInfo: info}
	} else if path == "/myproblems.html" {
		handler = &handlers.MyProblemsHandler{RequestInfo: info}
	} else if path == "/problem.html" {
		handler = &handlers.ProblemHandler{RequestInfo: info}
	} else if path == "/addproblem.html" || path == "/addproblem" {
		handler = &handlers.AddProblemHandler{RequestInfo: info}
	} else if path == "/editproblem.html" || path == "/editproblem" {
		handler = &handlers.EditProblemHandler{RequestInfo: info}
	} else if path == "/duplicateproblem.html" {
		handler = &handlers.DuplicateProblemHandler{RequestInfo: info}
	}

	if path == "/groups.html" {
		handler = &handlers.GroupsHandler{RequestInfo: info}
	} else if path == "/addgroup.html" || path == "/addgroup" {
		handler = &handlers.AddGroupHandler{RequestInfo: info}
	} else if path == "/joingroup.html" || path == "/joingroup" {
		handler = &handlers.JoinGroupHandler{RequestInfo: info}
	}

	if path == "/assignments.html" {
		handler = &handlers.AllAssignmentsHandler{RequestInfo: info}
	} else if path == "/myassignments.html" {
		handler = &handlers.MyAssignmentsHandler{RequestInfo: info}
	} else if path == "/addassignment.html" || path == "/addassignment" {
		handler = &handlers.AddAssignmentHandler{RequestInfo: info}
	}

	if path == "/users.html" {
		handler = &handlers.UsersHandler{RequestInfo: info}
	} else if path == "/myusers.html" {
		handler = &handlers.MyUsersHandler{RequestInfo: info}
	} else if path == "/usersingroup.html" {
		handler = &handlers.UsersInGroupHtmlHandler{RequestInfo: info}
	} else if path == "/usersingroup" {
		handler = &handlers.UsersInGroupHandler{RequestInfo: info}
	} else if path == "/changeuserrole" {
		handler = &handlers.ChangeUserRoleHandler{RequestInfo: info}
	}

	if path == "/dashboard.html" {
		handler = &handlers.DashboardHandler{RequestInfo: info}
	}

	if path == "/rejudge" {
		handler = &handlers.RejudgeHandler{RequestInfo: info}
	}

	if path == "/studentprogress.html" {
		handler = &handlers.StudentProgressHandler{RequestInfo: info}
	}
	if path == "/groupprogress.html" {
		handler = &handlers.GroupProgressHandler{RequestInfo: info}
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

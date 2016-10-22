package request

import (
	"db"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

func HandleUser(w http.ResponseWriter, r *http.Request, user db.User) {
	path := r.URL.Path
	fmt.Println("path", path)
	if strings.HasPrefix(path, "/contest/") {
		split := strings.Split(strings.TrimPrefix(path, "/contest/"), "/")
		if len(split) < 2 {
			assignmentsUserHtml(w, r, user)
			return
		}
		contestId, err := strconv.ParseInt(split[0], 10, 64)
		if err != nil {
			assignmentsUserHtml(w, r, user)
			return
		}
		page := split[1]
		if page == "problems" {
			problemsHtml(w, r, user, contestId)
		} else if page == "problem" {
			problemHtml(w, r, user, contestId, split[2:])
		} else if page == "submit" {
			submitCodeHtml(w, r, user, contestId)
		} else if page == "submissions" {
			mySubmissionsHtml(w, r, user, contestId)
		} else if page == "submission" {
			submissionHtml(w, r, user, contestId, split[2:])
		} else if page == "submitcode" {
			submitCode(w, r, user, contestId)
		} else {
			fmt.Println("error", page)
		}

	} else if path == "/competition.html" {
		competitionHtml(w, r)
	} else if path == "/groups.html" {
		userGroupsHtml(w, r, user)
	} else if path == "/joingroup.html" {
		joinGroupHtml(w, r)
	} else if path == "/joingroup" {
		joinGroup(w, r, user)
	} else if path == "/logout" {
		logout(w, r)
	} else if path == "/dashboard.html" {
		dashboardUserHtml(w, r, user)
	} else if path == "/pointsperweek" {
		pointPerWeek(w, r, user)
	} else if path == "/assignments.html" {
		assignmentsUserHtml(w, r, user)
	} else {
		dashboardUserHtml(w, r, user)
	}
}

func sendUserResponse(w http.ResponseWriter, page string, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../user/"+page+".html", "../templates/header.html", "../templates/menu.html", "../templates/footer.html")
	t.Execute(w, Resp{data, "user"})
}

func joinGroup(w http.ResponseWriter, r *http.Request, user db.User) {
	r.ParseForm()
	groupName := r.Form["groupname"]
	group, _ := db.FindGroupByName(groupName[0])
	db.CreateUserGroup(user.Id, group.Id)
	http.Redirect(w, r, "/groups.html", http.StatusFound)
}

func assignmentsUserHtml(w http.ResponseWriter, r *http.Request, user db.User) {
	assignments, _ := db.ListAssignmentsForUser(user)
	sendUserResponse(w, "assignments", assignments)
}

func dashboardUserHtml(w http.ResponseWriter, r *http.Request, user db.User) {
	sendUserResponse(w, "dashboard", nil)
}

func competitionHtml(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	response := Response{Id: id}
	t, _ := template.ParseFiles("../user/competition.html")
	t.Execute(w, response)
}

func userGroupsHtml(w http.ResponseWriter, r *http.Request, user db.User) {
	userGroups, _ := db.ListGroupsForUser(user.Id)
	sendUserResponse(w, "groups", userGroups)
}

func joinGroupHtml(w http.ResponseWriter, r *http.Request) {
	sendUserResponse(w, "joingroup", nil)
}

func pointPerWeek(w http.ResponseWriter, r *http.Request, user db.User) {
	w.Header().Set("Content-Type", "application/json")

	Response := db.GetPointsPerWeek(user.Id)

	json, err := json.Marshal(Response)
	fmt.Println("json", json)
	if err != nil {
		fmt.Println("err", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

package request

import (
	"db"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"
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
	} else {
		assignmentsUserHtml(w, r, user)
	}
}

func joinGroup(w http.ResponseWriter, r *http.Request, user db.User) {
	r.ParseForm()
	groupName := r.Form["groupname"]
	group, _ := db.FindGroupByName(groupName[0])
	db.CreateUserGroup(user.Id, group.Id)
	http.Redirect(w, r, "/groups.html", http.StatusFound)
}

func assignmentsUserHtml(w http.ResponseWriter, r *http.Request, user db.User) {
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../user/assignments.html")
	assignments, _ := db.ListAssignmentsForUser(user)
	t.Execute(w, assignments)
}

func dashboardUserHtml(w http.ResponseWriter, r *http.Request, user db.User) {
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../user/dashboard.html")
	t.Execute(w, nil)
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
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../user/groups.html")
	userGroups, _ := db.ListGroupsForUser(user.Id)
	t.Execute(w, userGroups)
}

func joinGroupHtml(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../user/joingroup.html")
	t.Execute(w, nil)
}

func pointPerWeek(w http.ResponseWriter, r *http.Request, user db.User) {
	w.Header().Set("Content-Type", "application/json")
	type Response1 struct {
		Week   string
		Points int
	}
	Response := make([]Response1, 0)
	subs, _ := db.ListMyAllSubmissions(user.Id)

	monday := func(t time.Time) time.Time {
		tt := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		tt = tt.AddDate(0, 0, -(int(tt.Weekday())+6)%7)
		return tt
	}

	currWeek := monday(subs[0].Time)
	nextWeek := currWeek.AddDate(0, 0, 7)

	problemPoints := make(map[int64]int, 0)
	totalPoints := 0
	subIdx := 0
	for subIdx < len(subs) {
		currWeekResponse := Response1{
			currWeek.String()[:10], 0,
		}

		for subIdx < len(subs) && subs[subIdx].Time.After(currWeek) && subs[subIdx].Time.Before(nextWeek) {
			lastPts, _ := problemPoints[subs[subIdx].ApId]
			currPts := subs[subIdx].Points
			if currPts > lastPts {
				totalPoints += currPts - lastPts
			}

			problemPoints[subs[subIdx].ApId] = currPts

			subIdx++
		}
		currWeekResponse.Points = totalPoints
		Response = append(Response, currWeekResponse)
		// add totalPoints for current week
		currWeek = nextWeek
		nextWeek = nextWeek.AddDate(0, 0, 7)
	}
	fmt.Println(Response)

	json, err := json.Marshal(Response)
	fmt.Println("json", json)
	if err != nil {
		fmt.Println("err", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

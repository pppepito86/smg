package handlers

import (
	"db"
	"request/util"
	"strconv"
	"request/contest"
	"net/http"
)

type GroupHandler struct {
	util.NoInputValidator
	util.RequestInfo
}

func (h *GroupHandler) Execute() error {
	h.R.ParseForm()
	groupId, _ := strconv.ParseInt(h.R.Form["groupid"][0], 10, 64)


	if h.User.RoleName != "admin" && h.User.RoleName != "teacher" {
		inGroup, err := db.IsUserInGroup(h.User.Id, groupId)
		if err != nil || !inGroup {
			http.Redirect(h.W, h.R, "/error.html?error=\"You are not allowed to access this group\"",
				http.StatusFound)
		}
	}

	assignments, _ := db.ListAssignmentsForGroup(groupId)
	group, _ := db.ListGroup(groupId)

	totalPoints := 0
	for _, assignment := range assignments {
		points, _ := contest.GetUserPointsInAssignment(h.User.Id, assignment.Id)
		totalPoints += points
	}

	type data struct {
		GroupName	string
		Assignments	 []db.Assignment
		Points  	int
	}
	d := data {
		GroupName:   group.GroupName,
		Assignments: assignments,
		Points: 	 totalPoints,
	}

	util.ServeHtml(h.W, h.User, "group.html", d)

	return nil
}

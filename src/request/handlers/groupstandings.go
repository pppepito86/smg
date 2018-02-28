package handlers

import (
	"request/util"
	"request/contest"
	"strconv"
	"db"
	"sort"
)

type GroupStandingsHandler struct {
	util.NoInputValidator
	util.RequestInfo
}

func (h *GroupStandingsHandler) Execute() error {
	h.R.ParseForm()
	groupStr := h.R.Form["groupid"]
	groupId, _ := strconv.ParseInt(groupStr[0], 10, 64)

	assignments, err := db.ListAssignmentsForGroup(groupId)

	if err != nil {
		return err
	}

	usersMap := make(map[string]int)

	var standings contest.Result;

	for i, _ := range assignments {
		standings, err = contest.GetStandings(assignments[i].Id)
		if err != nil {
			return err
		}

		for _, user := range standings.Info {
			usersMap[user.UserName] += user.Points[0]
		}
	}

	info := contest.UsersInfo{}
	for _, user := range standings.Info {

		totalPoints := usersMap[user.UserName]
		userInfo := contest.UserInfo {UserName: user.UserName,
			FirstName: user.FirstName,
			LastName: user.LastName,
			Points: []int { totalPoints } }
		info = append(info, userInfo)
	}

	sort.Sort(info)
	for i, _ := range info {
		info[i].Place = i + 1
	}

	result := contest.Result{}
	result.Info = info
	result.Problems = []string{}

	if h.User.RoleName != "admin" && h.User.RoleName != "teacher" {
		top := 5 // show only the first few to users
		if top < len(result.Info) {
			result.Info = result.Info[:top]
		}
	}

	util.ServeHtml(h.W, h.User, "contest/standings.html", result)

	return nil
}

package contest

import (
	"db"
	"request/util"
)

type StandingsHandler struct {
	util.NoInputValidator
	ContestRequestInfo
}

func (h *StandingsHandler) Execute() error {
	if h.User.RoleName != "admin" {
		return nil
	}

	problems, _ := db.ListAssignmentProblems(h.Cid)
	users, _ := db.ListUsersForAssignment(h.Cid)
	submissions, _ := db.ListSubmissionsForAssignment(h.Cid)
	type UserInfo struct {
		UserName  string
		FirstName string
		LastName  string
		Points    []int
	}
	type Result struct {
		Problems []string
		Info     []UserInfo
	}
	result := Result{}
	problemsMap := make(map[int64]int)
	for idx, problem := range problems {
		problemsMap[problem.ProblemId] = idx
		result.Problems = append(result.Problems, problem.ProblemName)
	}
	usersMap := make(map[int64]int)
	info := []UserInfo{}
	for idx, user := range users {
		usersMap[user.Id] = idx
		userInfo := UserInfo{user.UserName, user.FirstName, user.LastName, make([]int, len(problems)+1)}
		info = append(info, userInfo)
	}
	for _, submission := range submissions {
		points := submission.Points
		uId, ok := usersMap[submission.UserId]
		if !ok {
			continue
		}
		pId := problemsMap[submission.ProblemId] + 1
		if points > info[uId].Points[pId] {
			diff := points - info[uId].Points[pId]
			info[uId].Points[pId] = points
			info[uId].Points[0] += diff
		}
	}
	result.Info = info
	response := util.Response{h.Cid, result, ""}
	util.ServeContestHtml(h.W, h.R, h.User, "standings.html", response)

	return nil
}

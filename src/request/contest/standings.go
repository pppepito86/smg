package contest

import (
	"db"
	"request/util"
	"sort"
)

type StandingsHandler struct {
	util.NoInputValidator
	ContestRequestInfo
}
type UserInfo struct {
	Place     int
	UserName  string
	FirstName string
	LastName  string
	Points    []int
}

type UsersInfo []UserInfo

func (slice UsersInfo) Len() int {
	return len(slice)
}

func (slice UsersInfo) Less(i, j int) bool {
	return slice[i].Points[0] > slice[j].Points[0]
}

func (slice UsersInfo) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (h *StandingsHandler) Execute() error {
	if h.User.RoleName != "admin" && h.User.RoleName != "teacher" {
		return nil
	}

	problems, _ := db.ListAssignmentProblems(h.Cid)
	users, _ := db.ListUsersForAssignment(h.Cid)
	submissions, _ := db.ListSubmissionsForAssignment(h.Cid)

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
	info := UsersInfo{}
	for idx, user := range users {
		usersMap[user.Id] = idx
		userInfo := UserInfo{0, user.UserName, user.FirstName, user.LastName, make([]int, len(problems)+1)}
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
	sort.Sort(info)
	for i, _ := range info {
		info[i].Place = i + 1
	}

	result.Info = info
	ServeContestHtml(h.ContestRequestInfo, "standings.html", result)

	return nil
}

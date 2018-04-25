package contest

import (
	"db"
	"request/util"
	"sort"
	"strconv"
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

type Result struct {
	Problems []string
	Info     []UserInfo
}

func (slice UsersInfo) Len() int {
	return len(slice)
}

func (slice UsersInfo) Less(i, j int) bool {
	return slice[i].Points[0] > slice[j].Points[0]
}

func (slice UsersInfo) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func GetUserPointsInAssignment(userId int64, contestId int64) (int, error) {
	aps, err := db.ListAssignmentProblems(contestId)
	if err != nil {
		return -1, err
	}

	submissions, err := db.ListMySubmissions(userId, contestId)
	if err != nil {
		return -1, err
	}

	problemIdToPoints := make(map[int64]int)

	for _, problem := range aps {
		problemIdToPoints[problem.ProblemId] = 0
	}
	for _, submission := range submissions {

		recordedPoints := problemIdToPoints[submission.ProblemId];
		if recordedPoints < submission.Points {
			problemIdToPoints[submission.ProblemId] = submission.Points
		}
	}

	totalPoints := 0
	for _, problem := range aps {
		totalPoints += problemIdToPoints[problem.ProblemId]
	}

	return totalPoints, nil
}

func GetStandings(contestId int64) (Result, error) {
	problems, _ := db.ListAssignmentProblems(contestId)
	users, _ := db.ListUsersForAssignment(contestId)
	submissions, _ := db.ListSubmissionsForAssignment(contestId)


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
		userInfo := UserInfo {0,
		user.UserName, user.FirstName, user.LastName,
		make([]int, len(problems)+1)}
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
	return result, nil
}

func (h *StandingsHandler) Execute() error {
	top, _ := strconv.Atoi(h.Assignment.Standings)
	if h.User.RoleName != "admin" && h.User.RoleName != "teacher" && top == 0 {
		return nil
	}

	result, err := GetStandings(h.Cid)
	if err != nil {
		return err
	}

	if h.User.RoleName != "admin" && h.User.RoleName != "teacher" {
		if top < len(result.Info) {
			result.Info = result.Info[:top]
		}
	}

	ServeContestHtml(h.ContestRequestInfo, "standings.html", result)

	return nil
}

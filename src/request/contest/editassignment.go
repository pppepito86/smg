package contest

import (
	"db"
	"net/http"
	"request/util"
	"strconv"
	"strings"
	"time"
)

type EditAssignmentHandler struct {
	util.NoInputValidator
	ContestRequestInfo
}

func (h *EditAssignmentHandler) Execute() error {
	if h.User.RoleName != "admin" && h.User.RoleName != "teacher" {
		return nil
	}

	a, _ := db.ListAssignment(h.Cid)
	if h.User.RoleName != "admin" && a.AuthorId != h.User.Id {
		return nil
	}

	h.R.ParseForm()
	name := h.R.Form["assignmentname"]
	p1 := h.R.Form["problem1"]
	groupId := h.R.Form["groupid"]

	testInfo := "show"
	if len(h.R.Form["test-info"]) > 0 {
		testInfo = "hide"
	}

	standings := ""
	if len(h.R.Form["standings"]) > 0 {
		standings = "5"
	}

	gid, _ := strconv.ParseInt(groupId[0], 10, 64)
	if len(name) != 1 {
		h.W.WriteHeader(http.StatusInternalServerError)
		return nil
	}
	sd := strings.Split(h.R.Form["startdate"][0], "/")
	st := strings.Split(h.R.Form["starttime"][0], ":")
	ed := strings.Split(h.R.Form["enddate"][0], "/")
	et := strings.Split(h.R.Form["endtime"][0], ":")
	y1, _ := strconv.Atoi(sd[2])
	M1, _ := strconv.Atoi(sd[1])
	d1, _ := strconv.Atoi(sd[0])
	h1, _ := strconv.Atoi(st[0])
	m1, _ := strconv.Atoi(st[1])
	y2, _ := strconv.Atoi(ed[2])
	M2, _ := strconv.Atoi(ed[1])
	d2, _ := strconv.Atoi(ed[0])
	h2, _ := strconv.Atoi(et[0])
	m2, _ := strconv.Atoi(et[1])
	location, _ := time.LoadLocation("Europe/Sofia")
	startTime := time.Date(y1, time.Month(M1), d1, h1, m1, 0, 0, location)
	endTime := time.Date(y2, time.Month(M2), d2, h2, m2, 0, 0, location)
	a.AssignmentName = name[0]
	a.StartTime = startTime
	a.EndTime = endTime
	a.GroupId = gid
	a.TestInfo = testInfo
	a.Standings = standings
	db.UpdateAssignment(a)
	if p1[0] != "" {
		db.DeleteAssignmentProblems(h.Cid)
		ppp := strings.Split(p1[0], ",")
		for i, pp := range ppp {
			pp = strings.TrimSpace(pp)
			p1Id, _ := strconv.ParseInt(pp, 10, 64)
			problem, _ := db.GetProblem(p1Id)
			db.AddProblemToAssignment(h.Cid, p1Id, int64(i+1), problem.Points)
		}
	}
	http.Redirect(h.W, h.R, "/contest/"+strconv.FormatInt(h.Cid, 10)+"/problems", http.StatusFound)

	return nil
}

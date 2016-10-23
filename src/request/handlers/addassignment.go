package handlers

import (
	"db"
	"net/http"
	"request/util"
	"strconv"
	"strings"
	"time"
)

type AddAssignmentHandler struct {
	util.NoInputValidator
	util.RequestInfo
}

func (h *AddAssignmentHandler) Execute() error {
	if h.User.RoleName != "admin" {
		return nil
	}

	if h.R.Method == "POST" {
		h.executePost()
	} else {
		h.executeGet()
	}

	return nil
}

func (h *AddAssignmentHandler) executeGet() {
	util.ServeHtml(h.W, h.User, "addassignment.html", nil)
}

func (h *AddAssignmentHandler) executePost() {
	h.R.ParseForm()
	name := h.R.Form["assignmentname"]
	p1 := h.R.Form["problem1"]
	groupId := h.R.Form["groupid"]
	testInfo := "show"
	if len(h.R.Form["test-info"]) > 0 {
		testInfo = "hide"
	}
	gid, _ := strconv.ParseInt(groupId[0], 10, 64)
	if len(name) != 1 {
		h.W.WriteHeader(http.StatusInternalServerError)
		return
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
	a := db.Assignment{
		AssignmentName: name[0],
		AuthorId:       h.User.Id,
		GroupId:        gid,
		StartTime:      startTime,
		EndTime:        endTime,
		TestInfo:       testInfo,
	}
	a, _ = db.CreateAssignment(a)
	if p1[0] != "" {
		ppp := strings.Split(p1[0], ",")
		for i, pp := range ppp {
			pp = strings.TrimSpace(pp)
			p1Id, _ := strconv.ParseInt(pp, 10, 64)
			problem, _ := db.GetProblem(p1Id)
			db.AddProblemToAssignment(a.Id, p1Id, int64(i+1), problem.Points)
		}
	}
	http.Redirect(h.W, h.R, "/assignments.html", http.StatusFound)
}

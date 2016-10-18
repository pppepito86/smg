package request

import (
	"db"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"submissions"
	"text/template"
	"time"
)

func HandleAdmin(w http.ResponseWriter, r *http.Request, user db.User) {
	path := r.URL.Path
	if strings.HasPrefix(path, "/contest/") {
		split := strings.Split(strings.TrimPrefix(path, "/contest/"), "/")
		if len(split) < 2 {
			assignmentsAdminHtml(w, r)
			return
		}
		contestId, err := strconv.ParseInt(split[0], 10, 64)
		if err != nil {
			assignmentsAdminHtml(w, r)
			return
		}
		page := split[1]
		if page == "problems" {
			problemsHtml(w, r, user, contestId)
		} else if page == "problem" {
			problemHtml(w, r, user, contestId, split[2:])
		} else if page == "submit" {
			submitCodeHtml(w, r, user, contestId)
		} else if page == "mysubmissions" {
			mySubmissionsHtml(w, r, user, contestId)
		} else if page == "submissions" {
			allSubmissionsHtml(w, r, user, contestId)
		} else if page == "submission" {
			submissionHtml(w, r, user, contestId, split[2:])
		} else if page == "contestants" {
			contestantsAdminHtml(w, r, contestId)
		} else if page == "standings" {
			standingsHtml(w, r, user, contestId)
		} else if page == "submitcode" {
			submitCode(w, r, user, contestId)
		} else if page == "edit" {
			editHtml(w, r, user, contestId)
		} else if page == "editassignment" {
			editAssignment(w, r, user, contestId)
		} else {
			fmt.Println("error", page)
		}
	} else if path == "/changeuserrole" && r.Method == "POST" {
		changeUserRole(w, r)
	} else if path == "/addgroup" && r.Method == "POST" {
		addGroup(w, r, user)
	} else if path == "/addproblem" && r.Method == "POST" {
		addProblem(w, r, user)
	} else if path == "/editproblem" && r.Method == "POST" {
		editProblem(w, r, user)
	} else if path == "/addassignment" && r.Method == "POST" {
		addAssignment(w, r, user)
	} else if path == "/users.html" {
		usersAdminHtml(w, r)
	} else if path == "/groups.html" {
		groupsAdminHtml(w, r)
	} else if path == "/addgroup.html" {
		addAdminGroupHtml(w, r)
	} else if path == "/problems.html" {
		problemsAdminHtml(w, r)
	} else if path == "/rejudge" {
		rejudge(w, r)
	} else if path == "/problem.html" {
		problemAdminHtml(w, r)
	} else if path == "/addproblem.html" {
		addAdminProblemHtml(w, r)
	} else if path == "/editproblem.html" {
		editAdminProblemHtml(w, r)
	} else if path == "/assignments.html" {
		assignmentsAdminHtml(w, r)
	} else if path == "/addassignment.html" {
		addAdminAssignmentHtml(w, r)
    } else if path == "/studentprogress.html" {
		studentProgressHtml(w, r)
    }else if path == "/pointsperweek" {
        pointPerWeekAdmin(w, r, user)
    } else if path == "/logout" {
		logout(w, r)
	} else {
		assignmentsAdminHtml(w, r)
	}
}

func addGroup(w http.ResponseWriter, r *http.Request, user db.User) {
	r.ParseForm()
	name := r.Form["groupname"]
	description := r.Form["description"]
	if len(name) != 1 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	g := db.Group{
		GroupName:   name[0],
		Description: description[0],
		CreatorId:   user.Id,
	}
	g, _ = db.CreateGroup(g)
	http.Redirect(w, r, "/groups.html", http.StatusFound)
}

func addAssignment(w http.ResponseWriter, r *http.Request, user db.User) {
	r.ParseForm()
	name := r.Form["assignmentname"]
	p1 := r.Form["problem1"]
	groupId := r.Form["groupid"]
	gid, _ := strconv.ParseInt(groupId[0], 10, 64)
	if len(name) != 1 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sd := strings.Split(r.Form["startdate"][0], "/")
	st := strings.Split(r.Form["starttime"][0], ":")
	ed := strings.Split(r.Form["enddate"][0], "/")
	et := strings.Split(r.Form["endtime"][0], ":")
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
		AuthorId:       user.Id,
		GroupId:        gid,
		StartTime:      startTime,
		EndTime:        endTime,
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
	http.Redirect(w, r, "/assignments.html", http.StatusFound)
}

func editAssignment(w http.ResponseWriter, r *http.Request, user db.User, cid int64) {
	r.ParseForm()
	name := r.Form["assignmentname"]
	p1 := r.Form["problem1"]
	groupId := r.Form["groupid"]
	gid, _ := strconv.ParseInt(groupId[0], 10, 64)
	if len(name) != 1 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sd := strings.Split(r.Form["startdate"][0], "/")
	st := strings.Split(r.Form["starttime"][0], ":")
	ed := strings.Split(r.Form["enddate"][0], "/")
	et := strings.Split(r.Form["endtime"][0], ":")
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
	a, _ := db.ListAssignment(cid)
	a.AssignmentName = name[0]
	a.StartTime = startTime
	a.EndTime = endTime
	a.GroupId = gid
	db.UpdateAssignment(a)
	aps, _ := db.ListAssignmentProblems(cid)
	if p1[0] != "" {
		ppp := strings.Split(p1[0], ",")
		for i, pp := range ppp {
			pp = strings.TrimSpace(pp)
			p1Id, _ := strconv.ParseInt(pp, 10, 64)
			problem, _ := db.GetProblem(p1Id)
			if i < len(aps) {
				db.UpdateAssignmentProblem(aps[i].Id, p1Id, problem.Points)
			} else {
				db.AddProblemToAssignment(cid, p1Id, int64(i+1), problem.Points)
			}
		}
		if len(ppp) < len(aps) {
			for _, ap := range aps[len(ppp):len(aps)] {
				db.DeleteAssignmentProblem(ap.Id)
			}
		}
	}
	http.Redirect(w, r, "/contest/"+strconv.FormatInt(cid, 10)+"/problems", http.StatusFound)
}

func changeUserRole(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user := r.Form["userid"]
	role := r.Form["roleid"]
	if len(user) != 1 || len(role) != 1 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userId, _ := strconv.ParseInt(user[0], 10, 64)
	roleId, _ := strconv.ParseInt(role[0], 10, 64)
	db.UpdateUserRole(userId, roleId)
}

func addAdminGroupHtml(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../admin/addgroup.html")
	users, _ := db.ListUsers()
	t.Execute(w, users)
}

func addAdminProblemHtml(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../admin/addproblem.html")
	t.Execute(w, nil)
}

func LimitsFromString(limitsStr string) db.Limits {
	var limits db.Limits
	if err := json.Unmarshal([]byte(limitsStr), &limits); err != nil {
		fmt.Println("error", err.Error())
		limits = make(map[string]db.Limit, 0)
		limits["c++"] = db.Limit{"c++", 1000, 64}
		limits["java"] = db.Limit{"java", 1000, 64}
	}
	return limits
}

func editAdminProblemHtml(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../admin/editproblem.html")
	id, _ := strconv.ParseInt(r.URL.Query()["id"][0], 10, 64)
	problem, _ := db.GetProblem(id)

	problem.LangLimits = LimitsFromString(problem.Languages)
	dir := filepath.Join("workdir", "problems", strconv.FormatInt(id, 10))
	files, _ := ioutil.ReadDir(dir)
	tests := ""
	if len(files)%2 == 0 {
		for i := 1; ; i++ {
			input, err := ioutil.ReadFile(filepath.Join(dir, "input"+strconv.Itoa(i)))
			if err != nil {
				break
			}
			output, _ := ioutil.ReadFile(filepath.Join(dir, "output"+strconv.Itoa(i)))
			tests += string(input)
			tests += "#\n"
			tests += string(output)
			tests += "###\n"
		}
		if len(tests) > 4 {
			tests = tests[0 : len(tests)-4]
		}
	}
	problem.Tests = tests
	t.Execute(w, problem)
}

func addAdminAssignmentHtml(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../admin/addassignment.html")
	t.Execute(w, nil)
}

func studentProgressHtml(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../admin/dashboard.html")
	t.Execute(w, nil)
}

func usersAdminHtml(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../admin/users.html")
	users, _ := db.ListUsers()
	t.Execute(w, users)
}

func contestantsAdminHtml(w http.ResponseWriter, r *http.Request, id int64) {
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../admin/contest/contestants.html")
	users, _ := db.ListUsersForAssignment(id)
	response := Response{id, users}
	t.Execute(w, response)
}

func assignmentsAdminHtml(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../admin/assignments.html")
	assignments, _ := db.ListAssignments()
	t.Execute(w, assignments)
}

func groupsAdminHtml(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../admin/groups.html")
	t.Execute(w, db.ListGroups())
}

func problemsAdminHtml(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../admin/problems.html")
	problems, _ := db.ListProblems()
	t.Execute(w, problems)
}

func problemAdminHtml(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../admin/problem.html")
	id, _ := strconv.ParseInt(r.URL.Query()["id"][0], 10, 64)
	problem, _ := db.GetProblem(id)
	problem.LangLimits = LimitsFromString(problem.Languages)
	t.Execute(w, problem)
}

func allSubmissionsHtml(w http.ResponseWriter, r *http.Request, user db.User, cid int64) {
	submissions, _ := db.ListSubmissionsForAssignment(cid)
	response := Response{cid, submissions}
	serveContestHtml(w, r, user, "submissions.html", response)
}

func rejudge(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.URL.Query()["id"][0], 10, 64)
	p, _ := db.GetProblem(id)
	limits := LimitsFromString(p.Languages)
	ss, _ := db.ListProblemSubmissions(id)
	for _, s := range ss {
		s.Limit = limits[s.Language]
		submissions.Push(s)
	}

	http.Redirect(w, r, "/problems", http.StatusFound)
}

func pointPerWeekAdmin(w http.ResponseWriter, r *http.Request, user db.User) {
	w.Header().Set("Content-Type", "application/json")
	
    r.ParseForm()
	userId, err := strconv.ParseInt(r.Form["id"][0], 10, 64)
    if err != nil {
		fmt.Println("err", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
    
    Response := db.GetPointsPerWeek(userId)

	json, err := json.Marshal(Response)
	fmt.Println("json", json)
	if err != nil {
		fmt.Println("err", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(json)
}


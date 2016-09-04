package request

import (
	"db"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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
			adminProblemsHtml(w, r, user, contestId)
		} else if page == "problem" {
			adminProblemHtml(w, r, contestId, split[2:])
		} else if page == "submit" {
			submitAdminCodeHtml(w, r, contestId)
		} else if page == "submissions" {
			allSubmissionsHtml(w, r, user, contestId)
		} else if page == "submission" {
			submissionHtml(w, r, user, contestId, split[2:])
		} else if page == "contestants" {
			contestantsAdminHtml(w, r, contestId)
		} else if page == "submitcode" {
			submitCode(w, r, user)
		} else {
			fmt.Println("error", page)
		}
	} else if path == "/changeuserrole" && r.Method == "POST" {
		changeUserRole(w, r)
	} else if path == "/addgroup" && r.Method == "POST" {
		addGroup(w, r, user)
	} else if path == "/addproblem" && r.Method == "POST" {
		addProblem(w, r, user)
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
	} else if path == "/addproblem.html" {
		addAdminProblemHtml(w, r)
	} else if path == "/assignments.html" {
		assignmentsAdminHtml(w, r)
	} else if path == "/addassignment.html" {
		addAdminAssignmentHtml(w, r)
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
	p2 := r.Form["problem2"]
	p3 := r.Form["problem3"]
	groupId := r.Form["groupid"]
	gid, _ := strconv.ParseInt(groupId[0], 10, 64)
	if len(name) != 1 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	a := db.Assignment{
		AssignmentName: name[0],
		AuthorId:       user.Id,
		GroupId:        gid,
	}
	a, _ = db.CreateAssignment(a)
	if p1[0] != "" {
		ppp := strings.Split(p1[0], ",")
		for _, pp := range ppp {
			p1Id, _ := strconv.ParseInt(pp, 10, 64)
			db.AddProblemToAssignment(a.Id, p1Id, 1)
		}
	}
	if p2[0] != "" {
		p2Id, _ := strconv.ParseInt(p2[0], 10, 64)
		db.AddProblemToAssignment(a.Id, p2Id, 1)
	}
	if p3[0] != "" {
		p3Id, _ := strconv.ParseInt(p3[0], 10, 64)
		db.AddProblemToAssignment(a.Id, p3Id, 1)
	}
	http.Redirect(w, r, "/assignments.html", http.StatusFound)
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

func addAdminAssignmentHtml(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../admin/addassignment.html")
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

func allSubmissionsHtml(w http.ResponseWriter, r *http.Request, user db.User, cid int64) {
	submissions, _ := db.ListSubmissionsForAssignment(cid)
	response := Response{cid, submissions}
	serveCompetitionHtml(w, r, user, "../admin/contest/submissions.html", response)
}

func submissionHtml(w http.ResponseWriter, r *http.Request, user db.User, cid int64, args []string) {
	id, _ := strconv.ParseInt(args[0], 10, 64)
	mySubmission, _ := db.ListSubmission(id)
	source, _ := ioutil.ReadFile(mySubmission.SourceFile)
	mySubmission.Source = html.EscapeString(string(source))
	details, _ := db.ListSubmissionDetails(id)
	mySubmission.SubmissionDetails = details
	response := Response{cid, mySubmission}
	serveCompetitionHtml(w, r, user, "../admin/contest/submission.html", response)
}

func adminProblemsHtml(w http.ResponseWriter, r *http.Request, user db.User, cid int64) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if cid == 0 {
		cid = id
	}
	aps, _ := db.ListAssignmentProblems(cid)
	type data struct {
		Problems []db.AssignmentProblem
		Status   map[int64]string
	}
	d := data{
		Problems: aps,
		Status:   make(map[int64]string),
	}
	for _, ap := range aps {
		submissions, _ := db.ListMySubmissionsForProblem(ap.Id)
		if len(submissions) > 0 {
			d.Status[ap.Id] = "#ff0000"
			for _, s := range submissions {
				if s.Verdict == "Accepted" {
					d.Status[ap.Id] = "#00ff00"
					break
				}
			}
		} else {
			d.Status[ap.Id] = "ffffff"
		}
	}
	response := Response{cid, d}
	serveCompetitionHtml(w, r, db.User{}, "../admin/contest/problems.html", response)
}

func adminProblemHtml(w http.ResponseWriter, r *http.Request, cid int64, args []string) {
	apId, _ := strconv.ParseInt(args[0], 10, 64)
	ap, _ := db.GetAssignmentProblem(apId)
	problem, _ := db.GetProblem(ap.ProblemId)
	response := Response{cid, problem}
	serveCompetitionHtml(w, r, db.User{}, "../admin/contest/problem.html", response)
}

func submitAdminCodeHtml(w http.ResponseWriter, r *http.Request, cid int64) {
	aps, _ := db.ListAssignmentProblems(cid)
	response := Response{cid, aps}
	serveCompetitionHtml(w, r, db.User{}, "../admin/contest/submitcode.html", response)
}

func submitAdminCode(w http.ResponseWriter, r *http.Request, user db.User) {
	r.ParseForm()
	file, header, _ := r.FormFile("file")
	language := r.Form["language"]
	apIdStr := r.Form["apid"]
	apId, _ := strconv.ParseInt(apIdStr[0], 10, 64)

	t := time.Now().UTC()

	fp := filepath.Join("workdir", "users", strconv.FormatInt(user.Id, 10), strconv.FormatInt(t.UnixNano(), 16), header.Filename)
	os.MkdirAll(filepath.Dir(fp), 0755)
	out, _ := os.Create(fp)
	defer out.Close()
	_, _ = io.Copy(out, file)

	s := db.Submission{
		Id:         -1,
		ApId:       apId,
		UserId:     user.Id,
		Language:   language[0],
		SourceFile: fp,
		Verdict:    "pending",
	}

	s, _ = db.AddSubmission(s)
	submissions.Push(s)
	ap, _ := db.GetAssignmentProblem(apId)
	cId := ap.AssignmentId
	http.Redirect(w, r, "/contest/"+strconv.FormatInt(cId, 10)+"/submissions", http.StatusFound)
}

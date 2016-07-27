package request

import (
	"db"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"submissions"
	"text/template"
	"time"
)

func HandleUser(w http.ResponseWriter, r *http.Request, user db.User) {
	path := r.URL.Path
	if path == "/competition.html" {
		competitionHtml(w, r)
	} else if path == "/groups.html" {
		userGroupsHtml(w, r, user)
	} else if path == "/joingroup.html" {
		joinGroupHtml(w, r)
	} else if path == "/joingroup" {
		joinGroup(w, r, user)
	} else if path == "/problems.html" {
		userProblemsHtml(w, r, user)
	} else if path == "/problem.html" {
		userProblemHtml(w, r)
	} else if path == "/submit.html" {
		submitCodeHtml(w, r)
	} else if path == "/mysubmissions.html" {
		mySubmissionsHtml(w, r, user)
	} else if path == "/submitcode" {
		submitCode(w, r, user)
	} else {
		assignmentsUserHtml(w, r)
	}
}

func submitCode(w http.ResponseWriter, r *http.Request, user db.User) {
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
	http.Redirect(w, r, "/mysubmissions.html?id="+strconv.FormatInt(cId, 10), http.StatusFound)
}

func joinGroup(w http.ResponseWriter, r *http.Request, user db.User) {
	r.ParseForm()
	groupName := r.Form["groupname"]
	group, _ := db.FindGroupByName(groupName[0])
	db.CreateUserGroup(user.Id, group.Id)
	http.Redirect(w, r, "/groups.html", http.StatusFound)
}

func assignmentsUserHtml(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../user/assignments.html")
	assignments, _ := db.ListAssignments()
	t.Execute(w, assignments)
}

func competitionHtml(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	response := Response{Id: id}
	t, _ := template.ParseFiles("../user/competition.html")
	t.Execute(w, response)
}

type Response struct {
	Id   int64
	Data interface{}
}

func submitCodeHtml(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	aps, _ := db.ListAssignmentProblems(id)
	response := Response{id, aps}
	serveCompetitionHtml(w, r, db.User{}, "../user/submitcode.html", response)
}

func serveCompetitionHtml(w http.ResponseWriter, r *http.Request, user db.User, html string, response Response) {
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles(html)
	t.Execute(w, response)
}

func mySubmissionsHtml(w http.ResponseWriter, r *http.Request, user db.User) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	mySubmissions, _ := db.ListMySubmissions(user.Id, id)
	response := Response{id, mySubmissions}
	serveCompetitionHtml(w, r, user, "../user/mysubmissions.html", response)
}

func userProblemsHtml(w http.ResponseWriter, r *http.Request, user db.User) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	ok, err := isUserAssignedToCompetition(r, user)
	if !ok || err != nil {
		return
	}
	aps, _ := db.ListAssignmentProblems(id)
	response := Response{id, aps}
	serveCompetitionHtml(w, r, db.User{}, "../user/problems.html", response)
}

func userProblemHtml(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	apId, _ := db.GetAssignmentProblem(id)
	problem, _ := db.GetProblem(apId.ProblemId)
	response := Response{apId.AssignmentId, problem}
	serveCompetitionHtml(w, r, db.User{}, "../user/problem.html", response)
}

func userGroupsHtml(w http.ResponseWriter, r *http.Request, user db.User) {
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../user/groups.html")
	t.Execute(w, db.ListGroupsForUser(user.Id))
}

func joinGroupHtml(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../user/joingroup.html")
	t.Execute(w, nil)
}

func isUserAssignedToCompetition(r *http.Request, user db.User) (bool, error) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	return db.IsUserAssignedToCompetition(user.Id, id)
}

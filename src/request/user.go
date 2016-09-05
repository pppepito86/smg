package request

import (
	"db"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

func HandleUser(w http.ResponseWriter, r *http.Request, user db.User) {
	path := r.URL.Path
	fmt.Println("path", path)
	if strings.HasPrefix(path, "/contest/") {
		split := strings.Split(strings.TrimPrefix(path, "/contest/"), "/")
		if len(split) < 2 {
			assignmentsUserHtml(w, r, user)
			return
		}
		contestId, err := strconv.ParseInt(split[0], 10, 64)
		if err != nil {
			assignmentsUserHtml(w, r, user)
			return
		}
		page := split[1]
		if page == "problems" {
			problemsHtml(w, r, user, contestId)
		} else if page == "problem" {
			userProblemHtml(w, r, contestId, split[2:])
		} else if page == "submit" {
			submitCodeHtml(w, r, contestId)
		} else if page == "submissions" {
			mySubmissionsHtml(w, r, user, contestId)
		} else if page == "submission" {
			mySubmissionHtml(w, r, user, contestId, split[2:])
		} else if page == "submitcode" {
			submitCode(w, r, user, contestId)
		} else {
			fmt.Println("error", page)
		}

	} else if path == "/competition.html" {
		competitionHtml(w, r)
	} else if path == "/groups.html" {
		userGroupsHtml(w, r, user)
	} else if path == "/joingroup.html" {
		joinGroupHtml(w, r)
	} else if path == "/joingroup" {
		joinGroup(w, r, user)
	} else {
		assignmentsUserHtml(w, r, user)
	}
}

func joinGroup(w http.ResponseWriter, r *http.Request, user db.User) {
	r.ParseForm()
	groupName := r.Form["groupname"]
	group, _ := db.FindGroupByName(groupName[0])
	db.CreateUserGroup(user.Id, group.Id)
	http.Redirect(w, r, "/groups.html", http.StatusFound)
}

func assignmentsUserHtml(w http.ResponseWriter, r *http.Request, user db.User) {
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../user/assignments.html")
	assignments, _ := db.ListAssignmentsForUser(user)
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

func submitCodeHtml(w http.ResponseWriter, r *http.Request, cid int64) {
	aps, _ := db.ListAssignmentProblems(cid)
	response := Response{cid, aps}
	serveContestHtml(w, r, db.User{}, "../user/submitcode.html", response)
}

func serveCompetitionHtml(w http.ResponseWriter, r *http.Request, user db.User, html string, response Response) {
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles(html)
	t.Execute(w, response)
}

func mySubmissionsHtml(w http.ResponseWriter, r *http.Request, user db.User, cid int64) {
	mySubmissions, _ := db.ListMySubmissions(user.Id, cid)
	response := Response{cid, mySubmissions}
	serveCompetitionHtml(w, r, user, "../user/mysubmissions.html", response)
}

func mySubmissionHtml(w http.ResponseWriter, r *http.Request, user db.User, cid int64, args []string) {
	fmt.Println("***", args[0])
	id, _ := strconv.ParseInt(args[0], 10, 64)
	mySubmission, _ := db.ListSubmission(id)
	if mySubmission.UserId != user.Id {
		return
	}
	source, _ := ioutil.ReadFile(mySubmission.SourceFile)
	mySubmission.Source = html.EscapeString(string(source))
	response := Response{cid, mySubmission}
	serveCompetitionHtml(w, r, user, "../user/mysubmission.html", response)
}

func userProblemHtml(w http.ResponseWriter, r *http.Request, cid int64, args []string) {
	apId, _ := strconv.ParseInt(args[0], 10, 64)
	ap, _ := db.GetAssignmentProblem(apId)
	problem, _ := db.GetProblem(ap.ProblemId)
	response := Response{cid, problem}
	serveCompetitionHtml(w, r, db.User{}, "../user/problem.html", response)
}

func userGroupsHtml(w http.ResponseWriter, r *http.Request, user db.User) {
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../user/groups.html")
	userGroups, _ := db.ListGroupsForUser(user.Id)
	t.Execute(w, userGroups)
}

func joinGroupHtml(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../user/joingroup.html")
	t.Execute(w, nil)
}

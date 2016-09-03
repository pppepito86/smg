package request

import (
	"db"
	"fmt"
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
			userProblemsHtml(w, r, user, contestId)
		} else if page == "problem" {
			userProblemHtml(w, r, contestId, split[2:])
		} else if page == "submit" {
			submitCodeHtml(w, r, contestId)
		} else if page == "submissions" {
			mySubmissionsHtml(w, r, user, contestId)
		} else if page == "submission" {
			mySubmissionHtml(w, r, user, contestId, split[2:])
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
	} else if path == "/problems.html" {
		userProblemsHtml(w, r, user, 0)
	} else if path == "/problem.html" {
		userProblemHtml(w, r, 0, []string{})
	} else if path == "/submit.html" {
		submitCodeHtml(w, r, 0)
	} else if path == "/mysubmission.html" {
		mySubmissionHtml(w, r, user, 1, []string{})
	} else if path == "/mysubmissions.html" {
		mySubmissionsHtml(w, r, user, 0)
	} else if path == "/submitcode" {
		submitCode(w, r, user)
	} else {
		assignmentsUserHtml(w, r, user)
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
	http.Redirect(w, r, "/contest/"+strconv.FormatInt(cId, 10)+"/submissions", http.StatusFound)
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

type Response struct {
	Id   int64
	Data interface{}
}

func submitCodeHtml(w http.ResponseWriter, r *http.Request, cid int64) {
	aps, _ := db.ListAssignmentProblems(cid)
	response := Response{cid, aps}
	serveCompetitionHtml(w, r, db.User{}, "../user/submitcode.html", response)
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
	idStr := r.URL.Query().Get("id")
	if len(args) == 1 {
		idStr = args[0]
	}
	id, _ := strconv.ParseInt(idStr, 10, 64)
	mySubmission, _ := db.ListSubmission(id)
	source, _ := ioutil.ReadFile(mySubmission.SourceFile)
	mySubmission.Source = string(source)
	response := Response{cid, mySubmission}
	serveCompetitionHtml(w, r, user, "../user/mysubmission.html", response)
}

func userProblemsHtml(w http.ResponseWriter, r *http.Request, user db.User, cid int64) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if cid == 0 {
		cid = id
	}
	ok, err := isUserAssignedToCompetition(user, cid)
	if !ok || err != nil {
		return
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
			d.Status[ap.Id] = "#ffffff"
		}
	}
	response := Response{cid, d}
	serveCompetitionHtml(w, r, db.User{}, "../user/problems.html", response)
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

func isUserAssignedToCompetition(user db.User, id int64) (bool, error) {
	return db.IsUserAssignedToCompetition(user.Id, id)
}

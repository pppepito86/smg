package request

import (
	"db"
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

type Response struct {
	Id   int64
	Data interface{}
}

func problemsHtml(w http.ResponseWriter, r *http.Request, user db.User, cid int64) {
	if !isUserAssignedToContest(user, cid) {
        http.Redirect(w, r, "/error.html?error=\"You are not allowed to access this assignment\"", http.StatusFound)
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
		submissions, _ := db.ListMySubmissionsForProblem(user.Id, ap.AssignmentId, ap.ProblemId)
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
	serveContestHtml(w, r, user, "problems.html", response)
}

func mySubmissionsHtml(w http.ResponseWriter, r *http.Request, user db.User, cid int64) {
	if !isUserAssignedToContest(user, cid) {
        http.Redirect(w, r, "/error.html?error=\"You are not allowed to access this assignment\"", http.StatusFound)
		return
	}
	mySubmissions, _ := db.ListMySubmissions(user.Id, cid)
	response := Response{cid, mySubmissions}
	serveContestHtml(w, r, user, "submissions.html", response)
}

func submissionHtml(w http.ResponseWriter, r *http.Request, user db.User, cid int64, args []string) {
	id, _ := strconv.ParseInt(args[0], 10, 64)
	submission, _ := db.ListSubmission(id)
	if user.RoleName != "admin" && submission.UserId != user.Id {
		return
	}

	source, _ := ioutil.ReadFile(submission.SourceFile)
	submission.Source = html.EscapeString(string(source))
	details, _ := db.ListSubmissionDetails(id)
	for i := range details {
		d := &details[i]
		if strings.HasPrefix(d.Step, "Test #") {
			testIndex := d.Step[6:len(d.Step)]
			dir := filepath.Dir(submission.SourceFile)
			input, _ := ioutil.ReadFile(filepath.Join(dir, "input"+testIndex))
			output, _ := ioutil.ReadFile(filepath.Join(dir, "output"+testIndex))
			d.Input = string(input)
			d.Output = string(output)
			if len(d.Input) > 1000 {
				d.Input = "input too long"
			}
			d.Input = strings.Replace(d.Input, "\n", "<br>", -1)
			if len(d.Output) > 1000 {
				d.Output = "output too long"
			}
			d.Output = strings.Replace(d.Output, "\n", "<br>", -1)
		}
	}
	submission.SubmissionDetails = details
	response := Response{cid, submission}
	serveContestHtml(w, r, user, "submission.html", response)
}

func editHtml(w http.ResponseWriter, r *http.Request, user db.User, cid int64) {
	if user.RoleName != "admin" {
		return
	}

	assignment, _ := db.ListAssignment(cid)
	aps, _ := db.ListAssignmentProblems(cid)
	problems := ""
	for _, ap := range aps {
		problems += "," + strconv.FormatInt(ap.ProblemId, 10)
	}
	if len(problems) > 1 {
		problems = problems[1:len(problems)]
	}
	assignment.Problems = problems
	response := Response{cid, assignment}
	serveContestHtml(w, r, user, "edit.html", response)
}

func submitCodeHtml(w http.ResponseWriter, r *http.Request, user db.User, cid int64) {
	if !isUserAssignedToContest(user, cid) {
        http.Redirect(w, r, "/error.html?error=\"You are not allowed to access this assignment\"", http.StatusFound)
		return
	}
	aps, _ := db.ListAssignmentProblems(cid)
	response := Response{cid, aps}
	serveContestHtml(w, r, user, "submitcode.html", response)
}

func standingsHtml(w http.ResponseWriter, r *http.Request, user db.User, cid int64) {
	if user.RoleName != "admin" {
		return
	}

	problems, _ := db.ListAssignmentProblems(cid)
	users, _ := db.ListUsersForAssignment(cid)
	submissions, _ := db.ListSubmissionsForAssignment(cid)
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
	response := Response{cid, result}
	serveContestHtml(w, r, user, "standings.html", response)
}

func serveContestHtml(w http.ResponseWriter, r *http.Request, user db.User, html string, response Response) {
	w.Header().Set("Content-Type", "text/html")
	if user.RoleName == "admin" {
		html = "../admin/contest/" + html
	} else {
		html = "../user/" + html
	}
	t, _ := template.ParseFiles(html)
	t.Execute(w, response)
}

func problemHtml(w http.ResponseWriter, r *http.Request, user db.User, cid int64, args []string) {
	apId, _ := strconv.ParseInt(args[0], 10, 64)
	ap, _ := db.GetAssignmentProblem(apId)
	problem, _ := db.GetProblem(ap.ProblemId)
	problem.LangLimits = LimitsFromString(problem.Languages)
	response := Response{cid, problem}
	serveContestHtml(w, r, user, "problem.html", response)
}

func submitCode(w http.ResponseWriter, r *http.Request, user db.User, cid int64) {
	if !isUserAssignedToContest(user, cid) {
        http.Redirect(w, r, "/error.html?error=\"You are not allowed to access this assignment\"", http.StatusFound)
		return
	}

	// set max upload size to 1 MB
	r.Body = http.MaxBytesReader(w, r.Body, 50*1024)

	r.ParseForm()
	file, header, _ := r.FormFile("file")
	language := r.Form["language"]
	apIdStr := r.Form["apid"]
	apId, _ := strconv.ParseInt(apIdStr[0], 10, 64)
	ap, _ := db.GetAssignmentProblem(apId)
	if ap.AssignmentId != cid {
		return
	}

	t := time.Now().UTC()

	fp := filepath.Join("workdir", "users", strconv.FormatInt(user.Id, 10), strconv.FormatInt(t.UnixNano(), 16), header.Filename)
	os.MkdirAll(filepath.Dir(fp), 0755)
	out, _ := os.Create(fp)
	defer out.Close()
	_, _ = io.Copy(out, file)

	limits := LimitsFromString(ap.Languages)
	s := db.Submission{
		Id:            -1,
		AssignmentId:  ap.AssignmentId,
		ProblemId:     ap.ProblemId,
		UserId:        user.Id,
		Language:      language[0],
		SourceFile:    fp,
		Verdict:       "pending",
		ProblemPoints: ap.Points,
		Limit:         limits[language[0]],
	}

	s, _ = db.AddSubmission(s)
	submissions.Push(s)
	http.Redirect(w, r, "/contest/"+strconv.FormatInt(cid, 10)+"/submission/"+strconv.FormatInt(s.Id, 10), http.StatusFound)
}

func isUserAssignedToContest(user db.User, id int64) bool {
	if user.RoleName == "admin" {
		return true
	}
	ok, _ := db.IsUserAssignedToCompetition(user.Id, id)
	if !ok {
		return false
	}
	a, _ := db.ListAssignment(id)
	time := time.Now()
	return time.After(a.StartTime) && time.Before(a.EndTime)
}

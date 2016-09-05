package request

import (
	"db"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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
		submissions, _ := db.ListMySubmissionsForProblem(user.Id, ap.Id)
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

func isUserAssignedToContest(user db.User, id int64) bool {
	if user.RoleName == "admin" {
		return true
	}
	ok, _ := db.IsUserAssignedToCompetition(user.Id, id)
	return ok
}

func submitCode(w http.ResponseWriter, r *http.Request, user db.User, cid int64) {
	if !isUserAssignedToContest(user, cid) {
		return
	}
	r.ParseForm()
	file, header, _ := r.FormFile("file")
	language := r.Form["language"]
	apIdStr := r.Form["apid"]
	apId, _ := strconv.ParseInt(apIdStr[0], 10, 64)
	ap, _ := db.GetAssignmentProblem(apId)
	fmt.Println("cid", cid, "apid", apId, "apcid", ap.AssignmentId)
	if ap.AssignmentId != cid {
		return
	}

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
	http.Redirect(w, r, "/contest/"+strconv.FormatInt(cid, 10)+"/submission/"+strconv.FormatInt(s.Id, 10), http.StatusFound)
}

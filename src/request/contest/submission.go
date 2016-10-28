package contest

import (
	"db"
	"html"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"request/util"
	"strconv"
	"strings"
)

type SubmissionHandler struct {
	util.NoInputValidator
	ContestRequestInfo
}

func (h *SubmissionHandler) Execute() error {
	if !util.IsUserAssignedToContest(h.User, h.Assignment) {
		http.Redirect(h.W, h.R, "/error.html?error=\"You are not allowed to access this assignment\"", http.StatusFound)
		return nil
	}

	id, _ := strconv.ParseInt(h.Args[0], 10, 64)
	submission, _ := db.ListSubmission(id)
	if h.User.RoleName != "admin" && h.User.RoleName != "teacher" && submission.UserId != h.User.Id {
		return nil
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
	ServeContestHtml(h.ContestRequestInfo, "submission.html", submission)

	return nil
}

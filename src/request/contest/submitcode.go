package contest

import (
	"db"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"request/util"
	"strconv"
	"submissions"
	"time"
)

type SubmitCodeHandler struct {
	util.NoInputValidator
	ContestRequestInfo
}

func (h *SubmitCodeHandler) Execute() error {
	if !util.IsUserAssignedToContest(h.User, h.Assignment) {
		http.Redirect(h.W, h.R, "/error.html?error=\"You are not allowed to access this assignment\"", http.StatusFound)
		return nil
	}
	if h.User.RoleName != "admin" && h.User.RoleName != "teacher" && !h.Assignment.IsActive {
		http.Redirect(h.W, h.R, "/error.html?error=\"You are not allowed to submit in this assignment\"", http.StatusFound)
		return nil
	}

	// set max upload size to 1 MB
	h.R.Body = http.MaxBytesReader(h.W, h.R.Body, 50*1024)

	h.R.ParseForm()
	file, header, _ := h.R.FormFile("file")
	language := h.R.Form["language"]
	apIdStr := h.R.Form["apid"]
	apId, _ := strconv.ParseInt(apIdStr[0], 10, 64)
	ap, _ := db.GetAssignmentProblem(apId)
	if ap.AssignmentId != h.Cid {
		return nil
	}

	t := time.Now().UTC()

	fp := filepath.Join("workdir", "users", strconv.FormatInt(h.User.Id, 10), strconv.FormatInt(t.UnixNano(), 16), header.Filename)
	os.MkdirAll(filepath.Dir(fp), 0755)
	out, _ := os.Create(fp)
	defer out.Close()
	_, _ = io.Copy(out, file)

	limits := util.LimitsFromString(ap.Languages)
	limit := limits[language[0]]

	// FIXME:
	if language[0] == "nodejs" {
		limit = limits["java"]
	}
	s := db.Submission{
		Id:            -1,
		AssignmentId:  ap.AssignmentId,
		ProblemId:     ap.ProblemId,
		UserId:        h.User.Id,
		Language:      language[0],
		SourceFile:    fp,
		Verdict:       "pending",
		ProblemPoints: ap.Points,
		Limit:         limit,
	}

	s, _ = db.AddSubmission(s)
	submissions.Push(s)
	http.Redirect(h.W, h.R, "/contest/"+strconv.FormatInt(h.Cid, 10)+"/submission/"+strconv.FormatInt(s.Id, 10), http.StatusFound)

	return nil
}

package contest

import (
	"db"
	"net/http"
	"request/util"
)

type ProblemsHandler struct {
	util.NoInputValidator
	ContestRequestInfo
}

func (h *ProblemsHandler) Execute() error {
	if !util.IsUserAssignedToContest(h.User, h.Assignment) {
		http.Redirect(h.W, h.R, "/error.html?error=\"You are not allowed to access this assignment\"", http.StatusFound)
		return nil
	}

	aps, _ := db.ListAssignmentProblems(h.Cid)
	type data struct {
		Problems []db.AssignmentProblem
		Status   map[int64]string
	}
	d := data{
		Problems: aps,
		Status:   make(map[int64]string),
	}
	for _, ap := range aps {
		submissions, _ := db.ListMySubmissionsForProblem(h.User.Id, ap.AssignmentId, ap.ProblemId)
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
	ServeContestHtml(h.ContestRequestInfo, "problems.html", d)

	return nil
}

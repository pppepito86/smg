package util

import (
	"db"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
	"time"
)

type InputValidator interface {
	Validate() error
}

type RequestHandler interface {
	InputValidator
	Execute() error
}

type NoInputValidator struct{}

func (NoInputValidator) Validate() error {
	return nil
}

type Response struct {
	Id   int64
	Data interface{}
	Role string
}

func ServeContestHtml(w http.ResponseWriter, r *http.Request, user db.User, html string, response Response) {
	w.Header().Set("Content-Type", "text/html")
	if user.RoleName == "admin" {
		response.Role = "admin"
	} else {
		response.Role = "user"
	}
	t, _ := template.ParseFiles("../templates/contest/"+html, "../templates/contest/header.html", "../templates/contest/menu.html", "../templates/contest/footer.html")
	t.Execute(w, response)
}

func IsUserAssignedToContest(user db.User, id int64) bool {
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

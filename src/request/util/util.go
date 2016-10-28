package util

import (
	"db"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

type RequestInfo struct {
	R    *http.Request
	W    http.ResponseWriter
	User db.User
}

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
	Id         int64
	Data       interface{}
	Role       string
	Author     bool
	Assignment db.Assignment
}

func ServeHtml(w http.ResponseWriter, user db.User, html string, data interface{}) {
	ServeHtmlWithAuthor(w, user, html, data, false)
}

func ServeHtmlWithAuthor(w http.ResponseWriter, user db.User, html string, data interface{}, isAuthor bool) {
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../templates/"+html, "../templates/header.html", "../templates/menu.html", "../templates/footer.html")
	response := Response{0, data, user.RoleName, isAuthor, db.Assignment{}}
	t.Execute(w, response)
}

func IsUserAssignedToContest(user db.User, a db.Assignment) bool {
	if user.RoleName == "admin" || user.RoleName == "teacher" {
		return true
	}
	ok, _ := db.IsUserAssignedToCompetition(user.Id, a.Id)
	if !ok {
		return false
	}
	return a.IsActive || a.HasFinished
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

func LimitsAsString(limits db.Limits) string {
	b, _ := json.Marshal(limits)
	return string(b)
}

func LimitsFromRequest(r *http.Request) db.Limits {
	limits := db.Limits{}
	if len(r.Form["cpp"]) > 0 {
		tl, _ := strconv.Atoi(r.Form["cpptime"][0])
		ml, _ := strconv.Atoi(r.Form["cppmemory"][0])
		if tl == 0 {
			tl = 1000
		}
		if ml == 0 {
			ml = 64
		}
		cpplimit := db.Limit{"c++", tl, ml}
		limits["c++"] = cpplimit
	}
	if len(r.Form["java"]) > 0 {
		tl, _ := strconv.Atoi(r.Form["javatime"][0])
		ml, _ := strconv.Atoi(r.Form["javamemory"][0])
		if tl <= 0 {
			tl = 1000
		}
		if tl > 10000 {
			tl = 10000
		}
		if ml <= 0 {
			ml = 64
		}
		if ml > 128 {
			ml = 128
		}
		javalimit := db.Limit{"java", tl, ml}
		limits["java"] = javalimit
	}
	return limits
}

func ParseTagList(tags string) []string {
	var tagList []string
	tagList = []string{}
	tokens := strings.Split(tags, ",")
	for _, token := range tokens {
		token = strings.TrimSpace(token)
		if token != "" {
			tagList = append(tagList, token)
		}
	}
	return tagList
}

func ErrorHtml(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query()["error"]
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../error.html")
	t.Execute(w, msg)
}

func GetSessionIdCookie(r http.Request) *http.Cookie {
	for _, cookie := range r.Cookies() {
		if cookie.Name == "session.id" {
			return cookie
		}
	}
	return new(http.Cookie)
}

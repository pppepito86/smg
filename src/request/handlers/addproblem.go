package handlers

import (
	"db"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"request/util"
	"strconv"
	"strings"
)

type AddProblemHandler struct {
	util.NoInputValidator
	util.RequestInfo
}

func (h *AddProblemHandler) Execute() error {
	if h.User.RoleName != "admin" && h.User.RoleName != "teacher" {
		return nil
	}

	if h.R.Method == "POST" {
		h.executePost()
	} else {
		h.executeGet()
	}

	return nil
}

func (h *AddProblemHandler) executeGet() {
	util.ServeHtml(h.W, h.User, "addproblem.html", nil)
}

func (h *AddProblemHandler) executePost() {
	h.R.ParseForm()
	file, header, _ := h.R.FormFile("file")
	name := h.R.Form["problemname"]
	version := h.R.Form["version"]
	tags := h.R.Form["tags"]
	description := h.R.Form["text"]
	test := h.R.Form["test"][0]
	visibility := h.R.Form["visibility"]
	points, _ := strconv.Atoi(h.R.Form["points"][0])

	if len(name) != 1 || len(visibility) != 1 {
		h.W.WriteHeader(http.StatusInternalServerError)
		return
	}

	var tagList []string
	if len(tags) > 0 && tags[0] != "" {
		tagList = util.ParseTagList(tags[0])
	} else {
		tagList = []string{}
	}

	limits := util.LimitsAsString(util.LimitsFromRequest(h.R))
	p := db.Problem{
		ProblemName: name[0],
		Version:     version[0],
		Tags:        tagList,
		Description: description[0],
		Visibility:  visibility[0],
		Languages:   limits,
		AuthorId:    h.User.Id,
		Points:      points,
	}
	p, _ = db.CreateProblem(p)

	os.MkdirAll(filepath.Join("workdir", "problems", strconv.FormatInt(p.Id, 10)), 0755)
	test = strings.Replace(test, "\r", "", -1)
	if len(test) > 0 {
		fmt.Println("using text field")
		tests := strings.Split(test, "###")
		for i, t := range tests {
			if len(t) > 0 && t[0] == 13 {
				t = t[1:]
			}
			if len(t) > 0 && t[0] == 10 {
				t = t[1:]
			}

			testio := strings.Split(t, "#")
			testin := testio[0]
			testout := testio[1]
			if len(testout) > 0 && testout[0] == 13 {
				testout = testout[1:]
			}
			if len(testout) > 0 && testout[0] == 10 {
				testout = testout[1:]
			}

			filein := filepath.Join("workdir", "problems", strconv.FormatInt(p.Id, 10), fmt.Sprintf("input%d", i+1))
			fileout := filepath.Join("workdir", "problems", strconv.FormatInt(p.Id, 10), fmt.Sprintf("output%d", i+1))
			ioutil.WriteFile(filein, []byte(testin), 0755)
			ioutil.WriteFile(fileout, []byte(testout), 0755)
		}
	} else {
		fmt.Println("using file")
		fp := filepath.Join("workdir", "problems", strconv.FormatInt(p.Id, 10), header.Filename)
		fmt.Println("file is", file)
		out, _ := os.Create(fp)
		defer out.Close()
		_, err := io.Copy(out, file)
		if err != nil {
			fmt.Println("error copying", err)
		}

		err = exec.Command("chmod", "777", fp).Run()
		if err != nil {
			fmt.Println("error chmod", err)
		}

		unzip := exec.Command("unzip", fp, "-d", filepath.Dir(fp))
		err = unzip.Run()
		if err != nil {
			fmt.Println("error unzip", err)
		}

		replaceR := exec.Command("bash", "-c", "sed -i 's/\r$//g' *put*")
		replaceR.Dir = filepath.Dir(fp)
		replaceR.Run()
	}

	http.Redirect(h.W, h.R, "/problems.html", http.StatusFound)

}

type DuplicateProblemHandler struct {
	util.NoInputValidator
	util.RequestInfo
}

func (h *DuplicateProblemHandler) Execute() error {
	if h.User.RoleName != "admin" && h.User.RoleName != "teacher" {
		return nil
	}

	id, _ := strconv.ParseInt(h.R.URL.Query()["id"][0], 10, 64)
	problem, _ := db.GetProblem(id)
	problem.Id = 0
	problem.Version = h.User.UserName

	problem, _ = db.CreateProblem(problem)
	oldProblemDir := filepath.Join("workdir", "problems", strconv.FormatInt(id, 10))
	newProblemDir := filepath.Join("workdir", "problems", strconv.FormatInt(problem.Id, 10))
	exec.Command("cp", "-R", oldProblemDir, newProblemDir).Run()

	http.Redirect(h.W, h.R, "/editproblem.html?id="+strconv.FormatInt(problem.Id, 10), http.StatusFound)

	return nil
}

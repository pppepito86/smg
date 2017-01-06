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

type EditProblemHandler struct {
	util.NoInputValidator
	util.RequestInfo
}

func (h *EditProblemHandler) Execute() error {
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

func (h *EditProblemHandler) executeGet() {
	id, _ := strconv.ParseInt(h.R.URL.Query()["id"][0], 10, 64)
	problem, _ := db.GetProblem(id)

	if h.User.RoleName == "teacher" && h.User.Id != problem.AuthorId {
		return
	}

	problem.LangLimits = util.LimitsFromString(problem.Languages)
	dir := filepath.Join("workdir", "problems", strconv.FormatInt(id, 10))
	files, _ := ioutil.ReadDir(dir)
	tests := ""
	if len(files)%2 == 0 {
		for i := 1; ; i++ {
			input, err := ioutil.ReadFile(filepath.Join(dir, "input"+strconv.Itoa(i)))
			if err != nil {
				break
			}
			output, _ := ioutil.ReadFile(filepath.Join(dir, "output"+strconv.Itoa(i)))
			tests += string(input)
			tests += "#\n"
			tests += string(output)
			tests += "###\n"
		}
		if len(tests) > 4 {
			tests = tests[0 : len(tests)-4]
		}
	}
	problem.Tests = tests
	util.ServeHtml(h.W, h.User, "editproblem.html", problem)
}

func (h *EditProblemHandler) executePost() {
	h.R.ParseForm()
	file, header, fileErr := h.R.FormFile("file")
	id, _ := strconv.ParseInt(h.R.URL.Query()["id"][0], 10, 64)

	problem, _ := db.GetProblem(id)
	if h.User.RoleName == "teacher" && h.User.Id != problem.AuthorId {
		return
	}

	name := h.R.Form["problemname"]
	version := h.R.Form["version"]
	tags := h.R.Form["tags"]
	description := h.R.Form["text"]
	test := h.R.Form["test"][0]
	visibility := h.R.Form["visibility"]
	points, _ := strconv.Atoi(h.R.Form["points"][0])
	limits := util.LimitsAsString(util.LimitsFromRequest(h.R))

	var tagList []string
	if len(tags) > 0 && tags[0] != "" {
		tagList = util.ParseTagList(tags[0])
	} else {
		tagList = []string{}
	}

	p := db.Problem{
		Id:          id,
		ProblemName: name[0],
		Tags:        tagList,
		Version:     version[0],
		Description: description[0],
		Visibility:  visibility[0],
		Languages:   limits,
		AuthorId:    h.User.Id,
		Points:      points,
	}
	db.UpdateProblem(p)

	dir := filepath.Join("workdir", "problems", strconv.FormatInt(id, 10))
	test = strings.Replace(test, "\r", "", -1)
	if len(test) > 0 {
		exec.Command("rm", "-rf", dir).Run()
		os.MkdirAll(dir, 0755)

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

			filein := filepath.Join(dir, fmt.Sprintf("input%d", i+1))
			fileout := filepath.Join(dir, fmt.Sprintf("output%d", i+1))
			ioutil.WriteFile(filein, []byte(testin), 0755)
			ioutil.WriteFile(fileout, []byte(testout), 0755)
		}
	} else if fileErr == nil {
		exec.Command("rm", "-rf", dir).Run()
		os.MkdirAll(dir, 0755)

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

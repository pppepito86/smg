package request

import (
	"db"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

func HandleTeacher(w http.ResponseWriter, r *http.Request, user db.User) {
	path := r.URL.Path
	if path == "/addproblem" && r.Method == "POST" {
		addProblem(w, r, user)
	} else {
		if path == "/addproblem.html" {
			addProblemHtml(w, r)
		} else {
			//problemsHtml(w, r)
		}
	}
}

func addProblem(w http.ResponseWriter, r *http.Request, user db.User) {
	r.ParseForm()
	file, header, _ := r.FormFile("file")
	name := r.Form["problemname"]
	version := r.Form["version"]
	description := r.Form["text"]
	test := r.Form["test"][0]
	visibility := r.Form["visibility"]
	if len(name) != 1 || len(visibility) != 1 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	p := db.Problem{
		ProblemName: name[0],
		Version:     version[0],
		Description: description[0],
		Visibility:  visibility[0],
		Languages:   "java",
		AuthorId:    user.Id,
	}
	p, _ = db.CreateProblem(p)

	os.MkdirAll(filepath.Join("problems", strconv.FormatInt(p.Id, 10)), 0755)
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

			filein := filepath.Join("problems", strconv.FormatInt(p.Id, 10), fmt.Sprintf("input%d", i+1))
			fileout := filepath.Join("problems", strconv.FormatInt(p.Id, 10), fmt.Sprintf("output%d", i+1))
			ioutil.WriteFile(filein, []byte(testin), 0755)
			ioutil.WriteFile(fileout, []byte(testout), 0755)
		}
	} else {
		fmt.Println("using file")
		fp := filepath.Join("problems", strconv.FormatInt(p.Id, 10), header.Filename)
		fmt.Println("file is", file)
		//os.MkdirAll(filepath.Dir(fp), 0755)
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

	}

	http.Redirect(w, r, "/problems.html", http.StatusFound)
}

func editProblem(w http.ResponseWriter, r *http.Request, user db.User) {
	r.ParseForm()
	r.FormFile("file")
	id, _ := strconv.ParseInt(r.URL.Query()["id"][0], 10, 64)
	name := r.Form["problemname"]
	version := r.Form["version"]
	description := r.Form["text"]
	test := r.Form["test"][0]
	p := db.Problem{
		Id:          id,
		ProblemName: name[0],
		Version:     version[0],
		Description: description[0],
		Visibility:  "",
		Languages:   "java",
		AuthorId:    user.Id,
	}
	db.UpdateProblem(p)

	dir := filepath.Join("problems", strconv.FormatInt(id, 10))
	exec.Command("rm", "-rf", dir).Run()
	os.MkdirAll(dir, 0755)
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

			filein := filepath.Join(dir, fmt.Sprintf("input%d", i+1))
			fileout := filepath.Join(dir, fmt.Sprintf("output%d", i+1))
			ioutil.WriteFile(filein, []byte(testin), 0755)
			ioutil.WriteFile(fileout, []byte(testout), 0755)
		}
	}
	http.Redirect(w, r, "/problems.html", http.StatusFound)
}

func addProblemHtml(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../teacher/addproblem.html")
	t.Execute(w, nil)
}

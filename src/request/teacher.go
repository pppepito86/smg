package request

import (
	"db"
	"encoding/json"
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

func limitsFromRequest(r *http.Request) db.Limits {
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

func limitsAsString(limits db.Limits) string {
	b, _ := json.Marshal(limits)
	return string(b)
}

func parseTagList(tags string) []string {
    var tagList []string
    tagList = []string{}
    tokens := strings.Split(tags, ",")
    for _, token := range tokens {
        token = strings.TrimSpace(token)
        if(token != "") {
            tagList = append(tagList, token)
        }
    }
    return tagList
}

func addProblem(w http.ResponseWriter, r *http.Request, user db.User) {
	r.ParseForm()
	file, header, _ := r.FormFile("file")
	name := r.Form["problemname"]
	version := r.Form["version"]
    tags := r.Form["tags"]
	description := r.Form["text"]
	test := r.Form["test"][0]
	visibility := r.Form["visibility"]
	points, _ := strconv.Atoi(r.Form["points"][0])
    
	if len(name) != 1 || len(visibility) != 1 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
    
    var tagList []string
    if len(tags) > 0 && tags[0] != "" {
        tagList = parseTagList(tags[0])
    } else {
        tagList = []string{}
    }
    
	limits := limitsAsString(limitsFromRequest(r))
	p := db.Problem{
		ProblemName: name[0],
		Version:     version[0],
        Tags:        tagList,
		Description: description[0],
		Visibility:  visibility[0],
		Languages:   limits,
		AuthorId:    user.Id,
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

	http.Redirect(w, r, "/problems.html", http.StatusFound)
}

func editProblem(w http.ResponseWriter, r *http.Request, user db.User) {
	r.ParseForm()
	file, header, fileErr := r.FormFile("file")
	id, _ := strconv.ParseInt(r.URL.Query()["id"][0], 10, 64)
	name := r.Form["problemname"]
	version := r.Form["version"]
    tags := r.Form["tags"]
	description := r.Form["text"]
	test := r.Form["test"][0]
	points, _ := strconv.Atoi(r.Form["points"][0])
	limits := limitsAsString(limitsFromRequest(r))
    
    var tagList []string
    if len(tags) > 0 && tags[0] != "" {
        tagList = parseTagList(tags[0])
    } else {
        tagList = []string{}
    }
    
	p := db.Problem{
		Id:          id,
		ProblemName: name[0],
        Tags:        tagList,
		Version:     version[0],
		Description: description[0],
		Visibility:  "",
		Languages:   limits,
		AuthorId:    user.Id,
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

	http.Redirect(w, r, "/problems.html", http.StatusFound)
}

func addProblemHtml(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("../teacher/addproblem.html")
	t.Execute(w, nil)
}
